package main

import (
	"errors"
	"strings"
)

type Item struct {
	Name  string
	Price int
}

type Order struct {
	ID     int
	Status string
	Items  []Item
}

func NewOrder(id int) (*Order, error) {
	if id <= 0 {
		return nil, errors.New("ID заказа должен быть положительным")
	}
	return &Order{ID: id, Status: "new"}, nil
}

func AddItem(order *Order, name string, price int) error {
	if order == nil {
		return errors.New("заказ не создан")
	}
	if strings.TrimSpace(name) == "" {
		return errors.New("имя позиции не должно быть пустым")
	}
	if price < 0 {
		return errors.New("цена не может быть отрицательной")
	}

	order.Items = append(order.Items, Item{Name: name, Price: price})
	return nil
}

func Total(order Order) int {
	total := 0
	for _, item := range order.Items {
		total += item.Price
	}
	return total
}

func CopyOrder(order Order) Order {
	if order.Items != nil {
		items := make([]Item, len(order.Items))
		copy(items, order.Items)
		order.Items = items
	}
	return order
}
