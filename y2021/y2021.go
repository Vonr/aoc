package y2021

import (
	"fmt"
	"strconv"
	"strings"
)

var Days = []func(int, []string) string{
	Day1,
	Day2,
}

func Day1(level int, input []string) string {
	if level == 1 {
		// 2021D1L1
		lastNum, increases := -1, -1

		for _, line := range input {
			num, err := strconv.Atoi(line)

			if err != nil {
				continue
			}

			if num > lastNum {
				increases++
			}

			lastNum = num

		}

		fmt.Printf("Output: %d\n", increases)
		return fmt.Sprint(increases)
	}

	// 2021D1L2
	length := len(input)
	prevSum, increases := -1, -1
	for i, line := range input {
		if i+2 >= length {
			break
		}

		curr, err := strconv.Atoi(line)
		if err != nil {
			fmt.Println("failed to convert curr")
			continue
		}

		next, err := strconv.Atoi(input[i+1])
		if err != nil {
			fmt.Println("failed to convert next")
			continue
		}

		nextNext, err := strconv.Atoi(input[i+2])
		if err != nil {
			fmt.Println("failed to convert nextNext")
			continue
		}

		sum := curr + next + nextNext
		if sum > prevSum {
			increases++
		}

		prevSum = sum
	}

	fmt.Printf("Output: %d\n", increases)
	return fmt.Sprint(increases)
}

func Day2(level int, input []string) string {
	if level == 1 {
		// 2021D2L1
		depth, distance := 0, 0
		for _, line := range input {
			split := strings.Split(line, " ")
			if len(split) != 2 {
				continue
			}
			amount, err := strconv.Atoi(split[1])
			if err != nil {
				continue
			}

			switch split[0] {
			case "up":
				depth -= amount
			case "down":
				depth += amount
			case "forward":
				distance += amount
			}

		}
		fmt.Printf("Output: %d\n", depth*distance)
		return fmt.Sprint(depth * distance)
	}

	// 2021D2L2
	depth, distance, aim := 0, 0, 0
	for _, line := range input {
		split := strings.Split(line, " ")
		if len(split) != 2 {
			continue
		}
		amount, err := strconv.Atoi(split[1])
		if err != nil {
			continue
		}

		switch split[0] {
		case "up":
			aim -= amount
		case "down":
			aim += amount
		case "forward":
			distance += amount
			depth += aim * amount
		}
	}
	fmt.Printf("Output: %d\n", depth*distance)
	return fmt.Sprint(depth * distance)

}
