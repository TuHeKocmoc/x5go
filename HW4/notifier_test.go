package main

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

type notifierFunc func(string) error

func (f notifierFunc) Send(message string) error {
	return f(message)
}

func TestEmptyMessageChain(t *testing.T) {
	for _, notifier := range []Notifier{
		EmailNotifier{Address: "student@example.com"},
		PushNotifier{DeviceID: "device-42"},
	} {
		for _, message := range []string{"", " \t\n"} {
			for _, err := range []error{
				Notify(notifier, message),
				Dispatch([]Notifier{notifier}, message),
			} {
				if err == ErrEmptyMessage || !errors.Is(err, ErrEmptyMessage) {
					t.Fatalf("%T: ожидалась обёрнутая ErrEmptyMessage, получено %v", notifier, err)
				}
			}
		}
	}
}

func TestSendErrorChain(t *testing.T) {
	for _, recipient := range []string{"", " \t\n"} {
		for _, tc := range []struct {
			notifier Notifier
			channel  string
			reason   string
		}{
			{EmailNotifier{Address: recipient}, "email", "email-адрес не задан"},
			{PushNotifier{DeviceID: recipient}, "push", "не задан device id"},
		} {
			err := Dispatch([]Notifier{tc.notifier}, "Заказ готов")
			var sendErr *SendError
			if !errors.As(err, &sendErr) {
				t.Fatalf("%T: в цепочке нет *SendError: %v", tc.notifier, err)
			}
			if err == sendErr || sendErr.Channel != tc.channel || sendErr.Reason != tc.reason {
				t.Fatalf("Неверная ошибка канала: %v", err)
			}
		}
	}
}

func TestDispatchStopsAndRecovers(t *testing.T) {
	cause := errors.New("ошибка отправки")
	for _, tc := range []struct {
		name      string
		failure   notifierFunc
		wantCause error
		wantText  string
	}{
		{"error", func(string) error { return cause }, cause, "ошибка отправки"},
		{"panic(error)", func(string) error { panic(cause) }, cause, "паника при отправке"},
		{"panic(string)", (BrokenNotifier{}).Send, nil, "драйвер канала повреждён"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls []string
			record := func(name string, next notifierFunc) Notifier {
				return notifierFunc(func(message string) error {
					calls = append(calls, name+": "+message)
					if next != nil {
						return next(message)
					}
					return nil
				})
			}
			all := []Notifier{record("first", nil), record("second", tc.failure), record("third", nil)}
			err := Dispatch(all, "Заказ готов")
			if err == nil || !strings.Contains(err.Error(), tc.wantText) ||
				!strings.Contains(err.Error(), "канал 2: отправить уведомление:") {
				t.Fatalf("Потерян контекст ошибки: %v", err)
			}
			if tc.wantCause != nil && !errors.Is(err, tc.wantCause) {
				t.Fatalf("Потеряна исходная ошибка: %v", err)
			}
			if len(calls) != 2 || calls[0] != "first: Заказ готов" || calls[1] != "second: Заказ готов" {
				t.Fatalf("Отправка должна остановиться после второго канала: %v", calls)
			}

			calls = nil
			if err := Dispatch([]Notifier{all[0]}, "Дальше"); err != nil {
				t.Fatalf("После ошибки работа должна продолжаться: %v", err)
			}
			if len(calls) != 1 || calls[0] != "first: Дальше" {
				t.Fatalf("Следующее уведомление не отправлено: %v", calls)
			}
		})
	}
}

func TestNilAndEmpty(t *testing.T) {
	for _, all := range [][]Notifier{nil, {}} {
		if err := Dispatch(all, "Заказ готов"); err != nil {
			t.Fatalf("Пустой список не должен возвращать ошибку: %v", err)
		}
	}
	for _, err := range []error{
		Notify(nil, "Заказ готов"),
		Dispatch([]Notifier{nil}, "Заказ готов"),
	} {
		if err == nil {
			t.Fatal("Ожидалась ошибка для nil-уведомителя")
		}
	}
}

func ExampleDispatch() {
	all := []Notifier{
		EmailNotifier{BaseNotifier: BaseNotifier{Channel: "email"}, Address: "student@example.com"},
		PushNotifier{BaseNotifier: BaseNotifier{Channel: "push"}, DeviceID: "device-42"},
	}
	if err := Dispatch(all, "Заказ готов"); err != nil {
		fmt.Println("ошибка:", err)
	}
}
