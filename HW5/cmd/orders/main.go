package main

import (
	"fmt"

	"example.com/layered-orders/internal/config"
	"example.com/layered-orders/internal/handler"
	"example.com/layered-orders/internal/repository"
	"example.com/layered-orders/internal/service"
)

func main() {
	repo := repository.NewMemory()
	orders := service.NewService(repo)
	h := handler.New(orders)

	order, err := h.Create(handler.CreateRequest{Customer: "Анна"})
	if err != nil {
		fmt.Println("ошибка:", err)

		return
	}

	fmt.Printf("сервис слушал бы %s, создан заказ: %+v\n", config.HTTPAddr(), order)
}
