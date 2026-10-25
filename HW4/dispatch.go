package main

import (
	"errors"
	"fmt"
)

func safeSend(notifier Notifier, message string) (err error) {
	defer func() {
		if value := recover(); value != nil {
			if cause, ok := value.(error); ok {
				err = fmt.Errorf("паника при отправке: %w", cause)
			} else {
				err = fmt.Errorf("паника при отправке: %v", value)
			}
		}
	}()

	if notifier == nil {
		return errors.New("канал уведомлений не задан")
	}
	return notifier.Send(message)
}

func Notify(notifier Notifier, message string) error {
	if err := safeSend(notifier, message); err != nil {
		return fmt.Errorf("отправить уведомление: %w", err)
	}
	return nil
}

func Dispatch(all []Notifier, message string) error {
	for index, notifier := range all {
		if err := Notify(notifier, message); err != nil {
			return fmt.Errorf("канал %d: %w", index+1, err)
		}
	}
	return nil
}
