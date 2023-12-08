package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Println("Можете вводить арифметическое выражение")
	text, _ := reader.ReadString('\n')
	array := strings.Split(strings.TrimSpace(text), " ")

	if len(array) != 3 {
		panic("В выражении должно быть ровно 3 элемента: два числа и одна арифметическая операция")
	}

	var operations []string = []string{"+", "-", "/", "*"}
	if !slices.Contains[[]string](operations, array[1]) {
		panic("На месте второго элемента можно напечатать только +, -, * или /")
	}

	num1, err1 := strconv.Atoi(array[0])
	num2, err2 := strconv.Atoi(array[2])

	var isRome bool = false
	if err1 != nil || err2 != nil {

		var romeIntegers []string = []string{"I", "II", "III", "IV", "V", "VI", "VII", "VIII", "IX", "X"}
		if slices.Contains[[]string](romeIntegers, array[0]) {

			isRome = true
			num1 = RomeNumberToArabic(array[0])

		} else {
			panic("Оба аргумента обязаны быть либо римскими, либо арабскими числами1")
		}

		if slices.Contains[[]string](romeIntegers, array[2]) {

			if !isRome {
				panic("Второе число должно быть римским!")
			}

			num2 = RomeNumberToArabic(array[2])

		} else {
			panic("Оба аргумента обязаны быть либо римскими, либо арабскими числами")
		}
	}

	if num1 < 1 || num1 > 10 {
		panic("Первое число не должно быть меньше 1 или больше 10!")
	}
	if num2 < 1 || num2 > 10 {
		panic("Второе число не должно быть меньше 1 или больше 10!")
	}

	calculate := CreateStandartCalc()
	result, err := calculate.Calc(num1, num2, array[1])
	if err != nil {
		panic(err.Error())
	}
	if isRome {

		if result < 1 {
			panic("Римское число не может быть меньше единицы!")
		}
		fmt.Println(ArabicNumberToRome(result))

	} else {
		fmt.Println(result)
	}
}

func RomeNumberToArabic(number string) int {
	if len(number) < 1 {
		return 0
	}

	var arabic []int = []int{100, 50, 10, 5, 1}
	var rome []string = []string{"C", "L", "X", "V", "I"}

	prefix := string(number[0])
	i := slices.Index[[]string](rome, prefix)

	if slices.Contains[[]string](rome, number) {
		return arabic[i]
	}

	if i%2 == 0 && rome[i-2] == string(number[1]) {
		return (arabic[i-2] - arabic[i]) + RomeNumberToArabic(number[2:])
	}

	return arabic[i] + RomeNumberToArabic(number[1:])
}

func ArabicNumberToRome(number int) string {
	if number < 1 {
		return ""
	}

	var arabic []int = []int{100, 50, 10, 5, 1}
	var rome []string = []string{"C", "L", "X", "V", "I"}

	if slices.Contains[[]int](arabic, number) {
		return rome[slices.Index[[]int](arabic, number)]
	}

	i := slices.IndexFunc[[]int](arabic, func(i int) bool { return i <= number })

	if arabic[i-1] == number+arabic[i] {
		return rome[i] + rome[i-1]
	}

	return rome[i] + ArabicNumberToRome(number-arabic[i])
}

func CreateStandartCalc() Calculate[string, int] {
	return MakeCalculator[string, int]().AddOperation(
		"/",
		func(a, b int) int { return a / b },
		func(a, b int) error {
			if b == 0 {
				return errors.New("Вы попытались поделить на ноль")
			}
			return nil
		},
	).AddOperation(
		"+",
		func(a, b int) int { return a + b },
		nil,
	).AddOperation(
		"-",
		func(a, b int) int { return a - b },
		nil,
	).AddOperation(
		"*",
		func(a, b int) int { return a * b },
		nil,
	)
}

type Calculate[E comparable, T int] struct {
	Checkers   map[E]func(a, b T) error
	Operations map[E]func(a, b T) T
}

func MakeCalculator[E comparable, T int]() Calculate[E, T] {
	return Calculate[E, T]{Operations: make(map[E]func(a, b T) T), Checkers: make(map[E]func(a, b T) error)}
}

func (c Calculate[E, T]) AddOperation(name E, operation func(a, b T) T, checker func(a, b T) error) Calculate[E, T] {
	c.Operations[name] = operation
	if checker != nil {
		c.Checkers[name] = checker
	}
	return c
}

func (c Calculate[E, T]) Calc(a, b T, op E) (T, error) {
	for i, operation := range c.Operations {
		if i == op {
			for j, checker := range c.Checkers {
				if j == op {
					e := checker(a, b)
					if e != nil {
						return a, e
					}
				}
			}
			return operation(a, b), nil
		}
	}
	return a, errors.New("Операция не найдена")
}
