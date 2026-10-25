package main

import (
	"errors"
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
	if strings.TrimSpace(n.Address) == "" {
		return errors.New("email-адрес не задан")
	}
	n.Log(fmt.Sprintf("to %s: %s", n.Address, message))
	return nil
}

type PushNotifier struct {
	BaseNotifier
	DeviceID string
}

func (n PushNotifier) Send(message string) error {
	if strings.TrimSpace(n.DeviceID) == "" {
		return errors.New("устройство не задано")
	}
	n.Log(fmt.Sprintf("to %s: %s", n.DeviceID, message))
	return nil
}

func Dispatch(all []Notifier, message string) error {
	for _, notifier := range all {
		if notifier == nil {
			return errors.New("канал уведомлений не задан")
		}
		if err := notifier.Send(message); err != nil {
			return err
		}
	}
	return nil
}
