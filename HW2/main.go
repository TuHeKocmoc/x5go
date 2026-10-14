package main

import "fmt"

func main() {
	statusNames := make(map[string]string)
	statusNames["new"] = "Новый"
	statusNames["paid"] = "Оплачен"
	statusNames["completed"] = "Завершён"
	statusNames["cancelled"] = "Отменён"

	order, err := NewOrder(1)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(Total(*order))

	err = AddItem(order, "Кофе", 12000)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = AddItem(order, "Чай", 7500)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(Total(*order))
	fmt.Println(statusNames[order.Status])

	orderCopy := CopyOrder(*order)
	orderCopy.ID = 2
	orderCopy.Status = "paid"
	orderCopy.Items[0].Price = 1000

	fmt.Println("Исходный заказ:", order.ID, statusNames[order.Status], Total(*order))
	fmt.Println("Копия заказа:", orderCopy.ID, statusNames[orderCopy.Status], Total(orderCopy))

	err = AddItem(nil, "Кофе", 12000)
	if err != nil {
		fmt.Println(err)
	}
}
