package main

import (
	"errors"
	"fmt"
	"testing"
)

type recordingNotifier struct {
	name  string
	calls *[]string
	err   error
}

func (n recordingNotifier) Send(message string) error {
	*n.calls = append(*n.calls, n.name+": "+message)
	return n.err
}

func TestDispatchCallsAll(t *testing.T) {
	var calls []string
	all := []Notifier{
		recordingNotifier{name: "email", calls: &calls},
		recordingNotifier{name: "push", calls: &calls},
	}
	if err := Dispatch(all, "Заказ готов"); err != nil {
		t.Fatal(err)
	}
	if len(calls) != 2 || calls[0] != "email: Заказ готов" || calls[1] != "push: Заказ готов" {
		t.Fatalf("Неверные вызовы уведомителей: %v", calls)
	}
}

func TestDispatchStopsOnFirstError(t *testing.T) {
	var calls []string
	wantErr := errors.New("ошибка второго канала")
	all := []Notifier{
		recordingNotifier{name: "first", calls: &calls},
		recordingNotifier{name: "second", calls: &calls, err: wantErr},
		recordingNotifier{name: "third", calls: &calls, err: errors.New("ошибка третьего канала")},
	}
	if err := Dispatch(all, "Сообщение"); err != wantErr {
		t.Fatalf("Ожидалась первая ошибка %v, получено %v", wantErr, err)
	}
	if len(calls) != 2 || calls[0] != "first: Сообщение" || calls[1] != "second: Сообщение" {
		t.Fatalf("Отправка должна остановиться после второго канала: %v", calls)
	}
}

func TestDispatchEmpty(t *testing.T) {
	for _, all := range [][]Notifier{nil, {}} {
		if err := Dispatch(all, "Сообщение"); err != nil {
			t.Fatalf("Пустой список не должен возвращать ошибку: %v", err)
		}
	}
}

func TestDispatchNilNotifier(t *testing.T) {
	var calls []string
	all := []Notifier{
		recordingNotifier{name: "first", calls: &calls},
		nil,
		recordingNotifier{name: "third", calls: &calls},
	}
	if err := Dispatch(all, "Сообщение"); err == nil {
		t.Fatal("Ожидалась ошибка для nil-уведомителя")
	}
	if len(calls) != 1 || calls[0] != "first: Сообщение" {
		t.Fatalf("Отправка должна остановиться на nil-уведомителе: %v", calls)
	}
}

func TestEmptyRecipients(t *testing.T) {
	for _, recipient := range []string{"", " \t\n"} {
		for _, notifier := range []Notifier{
			EmailNotifier{Address: recipient},
			PushNotifier{DeviceID: recipient},
		} {
			if err := notifier.Send("Сообщение"); err == nil {
				t.Fatalf("%T: ожидалась ошибка для получателя %q", notifier, recipient)
			}
		}
	}
}

func ExampleDispatch() {
	all := []Notifier{
		EmailNotifier{
			BaseNotifier: BaseNotifier{Channel: "email"},
			Address:      "student@example.com",
		},
		PushNotifier{
			BaseNotifier: BaseNotifier{Channel: "push"},
			DeviceID:     "device-42",
		},
	}
	if err := Dispatch(all, "Заказ готов"); err != nil {
		fmt.Println("ошибка:", err)
	}
}
