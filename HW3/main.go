package main

import "fmt"

func main() {
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
		return
	}

	for _, notifier := range all {
		if email, ok := notifier.(EmailNotifier); ok {
			fmt.Println("адрес email-канала:", email.Address)
		}
	}
}
