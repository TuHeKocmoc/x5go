package main

import "fmt"

func task1() {
	for number := 1; number <= 30; number++ {
		switch {
		case number%15 == 0:
			fmt.Println("FizzBuzz")
		case number%3 == 0:
			fmt.Println("Fizz")
		case number%5 == 0:
			fmt.Println("Buzz")
		default:
			fmt.Println(number)
		}
	}
}
