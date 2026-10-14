package main

import (
	"fmt"
	"strconv"
	"strings"
)

func readLine() string {
	line := ""
	for {
		var character rune
		count, _ := fmt.Scanf("%c", &character)
		if count != 1 || character == '\n' {
			break
		}
		line += string(character)
	}
	return line
}

func task2() {
	fields := strings.Fields(readLine())

	amounts := make([]int, len(fields))
	for index, field := range fields {
		amount, _ := strconv.Atoi(field)
		amounts[index] = amount
	}

	minimum := amounts[0]
	maximum := amounts[0]
	total := 0
	for _, amount := range amounts {
		if amount < minimum {
			minimum = amount
		}
		if amount > maximum {
			maximum = amount
		}
		total += amount
	}

	average := float64(total) / float64(len(amounts))
	fmt.Println(minimum)
	fmt.Println(maximum)
	fmt.Println(total)
	fmt.Println(average)
}
