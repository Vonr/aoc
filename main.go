package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

var token string
var years = map[int]func(int){
	2021: year2021,
}

func main() {
	var year int
	var day int
	now := time.Now().UTC().Add(6 * time.Hour)

	flag.IntVar(&year, "year", now.Year(), "Year")
	flag.IntVar(&day, "day", now.Day(), "Day")
	flag.StringVar(&token, "token", "", "Session Token")
	flag.Parse()

	if token == "" && !inputExists(year, day) {
		fmt.Println("Please provide a session token")
		os.Exit(1)
		flag.Usage()
		return
	}

	if day > 25 {
		fmt.Println("Day must be between 1 and 25")
		flag.Usage()
		return
	}
	fmt.Printf("Running Year %d Day %d\n", year, day)

	years[year](day)
}

func inputExists(year, day int) bool {
	_, err := os.Stat(fmt.Sprintf("inputs/%d/%d", year, day))
	return !os.IsNotExist(err)
}

// Function inspired by https://github.com/GreenLightning/advent-of-code-downloader/blob/c020acc44d8bfa409252f71aaa8705e9a80c73e2/aocdl/main.go#L254
func readInput(year, day int) []string {
	if inputExists(year, day) {
		fmt.Println("Reading input from file")

		file, err := os.Open(fmt.Sprintf("inputs/%d/%d", year, day))
		if err != nil {
			panic(err)
		}
		defer file.Close()

		b, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		return strings.Split(string(b), "\n")
	}

	// Inputs are located in https://adventofcode.com/<year>/day/<day>/input
	fmt.Println("Downloading input from adventofcode.com")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", year, day), nil)
	if err != nil {
		panic(err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = token
	req.AddCookie(cookie)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		fmt.Println(string(body))
		panic("Failed to download input")
	}

	// Save to file
	err = os.MkdirAll(fmt.Sprintf("inputs/%d", year), 0755)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fmt.Sprintf("inputs/%d/%d", year, day), body, 0644)
	if err != nil {
		panic(err)
	}

	return strings.Split(string(body), "\n")
}

func year2021(day int) {
	input := readInput(2021, day)
	fmt.Println()

	days := []func(){
		func() {
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

			fmt.Printf("Part 1 Output: %d\n", increases)

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

			fmt.Printf("Part 2 Output: %d\n", increases)
		},
		func() {
			depth, distance := 0, 0
			for _, line := range input {
				if strings.HasPrefix(line, "up ") {
					depthChange, err := strconv.Atoi(strings.TrimPrefix(line, "up "))
					if err != nil {
						continue
					}
					depth -= depthChange
				} else if strings.HasPrefix(line, "down ") {
					depthChange, err := strconv.Atoi(strings.TrimPrefix(line, "down "))
					if err != nil {
						continue
					}
					depth += depthChange
				} else if strings.HasPrefix(line, "forward ") {
					distanceChange, err := strconv.Atoi(strings.TrimPrefix(line, "forward "))
					if err != nil {
						continue
					}
					distance += distanceChange
				}

			}
			fmt.Printf("Part 1 Output: %d\n", depth*distance)

			depth, distance, aim := 0, 0, 0
			for _, line := range input {
				if strings.HasPrefix(line, "up ") {
					aimChange, err := strconv.Atoi(strings.TrimPrefix(line, "up "))
					if err != nil {
						continue
					}
					aim -= aimChange
				} else if strings.HasPrefix(line, "down ") {
					aimChange, err := strconv.Atoi(strings.TrimPrefix(line, "down "))
					if err != nil {
						continue
					}
					aim += aimChange
				} else if strings.HasPrefix(line, "forward ") {
					distanceChange, err := strconv.Atoi(strings.TrimPrefix(line, "forward "))
					if err != nil {
						continue
					}
					distance += distanceChange
					depth += aim * distanceChange
				}
			}
			fmt.Printf("Part 2 Output: %d\n", depth*distance)
		},
	}

	days[day-1]()
}
