package y2021

import (
	"fmt"
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
		freq := make([]int, len(input[0]))
		for _, line := range input {
			for j, char := range line {
				if char == '1' {
					freq[j]++
				} else {
					freq[j]--
				}
			}
		}

		var gamma string
		var epsilon string
		for _, v := range freq {
			if v > 0 {
				gamma += "1"
				epsilon += "0"
			} else {
				gamma += "0"
				epsilon += "1"
			}
		}

		g, err := strconv.ParseUint(gamma, 2, len(input[0]))
		if err != nil {
			panic(err)
		}

		e, err := strconv.ParseUint(epsilon, 2, len(input[0]))
		if err != nil {
			panic(err)
		}

		return fmt.Sprint(g * e)
	}

	// 2021D3L2
	oxy, co := input, input
	var i int

	for {
		var newOxy, newCo []string
		if len(oxy) > 1 {
			oxyFreq := 0
			for _, line := range oxy {
				if len(line) == 0 {
					continue
				}
				if line[i] == '1' {
					oxyFreq++
				} else {
					oxyFreq--
				}
			}

			for _, v := range oxy {
				if len(v) == 0 {
					continue
				}
				if (oxyFreq >= 0 && v[i] == '1') || (oxyFreq < 0 && v[i] == '0') {
					if len(v) > 1 {
						newOxy = append(newOxy, v)
					} else {
						newOxy = append(newOxy, v)
					}
				}
			}
			oxy = newOxy
		}

		if len(co) > 1 {
			coFreq := 0
			for _, line := range co {
				if len(line) == 0 {
					continue
				}
				if line[i] == '1' {
					coFreq++
				} else {
					coFreq--
				}
			}

			for _, v := range co {
				if len(v) == 0 {
					continue
				}
				// If 1 and 0 are equally common, keep the ones starting with 0
				if (coFreq >= 0 && v[i] == '0') || (coFreq < 0 && v[i] == '1') {
					if len(v) > 1 {
						newCo = append(newCo, v)
					} else {
						newCo = append(newCo, v)
					}
				}
			}
			co = newCo
		}

		if len(oxy) == 1 && len(co) == 1 {
			break
		}

		i++
	}
	o, err := strconv.ParseUint(oxy[0], 2, len(oxy[0]))
	if err != nil {
		panic(err)
	}

	c, err := strconv.ParseUint(co[0], 2, len(co[0]))
	if err != nil {
		panic(err)
	}

	return fmt.Sprint(o * c)
}
