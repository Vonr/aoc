package y2021

import (
	"fmt"
	"math"
	"strconv"
	"strings"
)

var Days = []func(int, []string) string{
	Day1, Day2, Day3,
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
			continue
		}

		next, err := strconv.Atoi(input[i+1])
		if err != nil {
			continue
		}

		nextNext, err := strconv.Atoi(input[i+2])
		if err != nil {
			continue
		}

		sum := curr + next + nextNext
		if sum > prevSum {
			increases++
		}

		prevSum = sum
	}

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
	return fmt.Sprint(depth * distance)

}

func Day3(level int, input []string) string {
	if level == 1 {
		// 2021D3L1
		freq, gamma := make([]int, len(input[0])), ""
		for _, line := range input {
			for j, char := range line {
				freq[j] += int(char)*2 - 97
			}
		}
		for _, v := range freq {
			gamma += string(byte(math.Min(math.Max(float64(v), 0), 1) + 48))
		}
		x, _ := strconv.ParseUint(gamma, 2, len(gamma))
		return fmt.Sprint(x * (^x & (1<<len(gamma) - 1)))
	}

	// 2021D3L2
	oxy, co, i, calc := input, input, 0, func(in []string, a, b byte, i int) []string {
		if len(in) > 1 {
			var out []string
			var freq int
			for _, line := range in {
				if len(line) > i {
					freq += int(line[i])*2 - 97
				}
			}
			for _, v := range in {
				if len(v) > i && ((freq >= 0 && v[i] == a) || (freq < 0 && v[i] == b)) {
					out = append(out, v)
				}
			}
			return out
		}
		return in
	}
	for ; len(oxy) != 1 && len(co) != 1; i++ {
		oxy = calc(oxy, '1', '0', i)
		co = calc(co, '0', '1', i)
	}
	x, _ := strconv.ParseUint(oxy[0]+co[0], 2, len(oxy[0])+len(co[0]))
	return fmt.Sprint((x >> len(co[0])) * (x & (1<<len(co[0]) - 1)))
}
