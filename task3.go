package main

import (
	"fmt"
	"strconv"
	"strings"
)

func task3() {
	orderAmount, _ := strconv.Atoi(strings.TrimSpace(readLine()))
	if orderAmount < 0 {
		fmt.Println("Ошибка: сумма заказа не может быть отрицательной")
		return
	}

	switch {
	case orderAmount < 10_000:
		fmt.Println("small")
	case orderAmount < 50_000:
		fmt.Println("regular")
	default:
		fmt.Println("large")
	}
}
