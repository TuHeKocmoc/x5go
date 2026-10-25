package main

import (
	"fmt"
	"strings"
)

type Notifier interface {
	Send(message string) error
}

type BaseNotifier struct {
	Channel string
}

func (b BaseNotifier) Log(message string) {
	fmt.Println(b.Channel, message)
}

type EmailNotifier struct {
	BaseNotifier
	Address string
}

func (n EmailNotifier) Send(message string) error {
	if strings.TrimSpace(message) == "" {
		return ErrEmptyMessage
	}
	if strings.TrimSpace(n.Address) == "" {
		return &SendError{Channel: "email", Reason: "email-адрес не задан"}
	}
	n.Log(fmt.Sprintf("to %s: %s", n.Address, message))
	return nil
}

type PushNotifier struct {
	BaseNotifier
	DeviceID string
}

func (n PushNotifier) Send(message string) error {
	if strings.TrimSpace(message) == "" {
		return ErrEmptyMessage
	}
	if strings.TrimSpace(n.DeviceID) == "" {
		return &SendError{Channel: "push", Reason: "не задан device id"}
	}
	n.Log(fmt.Sprintf("to %s: %s", n.DeviceID, message))
	return nil
}

type BrokenNotifier struct{}

func (BrokenNotifier) Send(string) error {
	panic("драйвер канала повреждён")
}
