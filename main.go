package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
)

var token string
var years = map[int]func(int, int){
	2021: year2021,
}

func main() {
	var year int
	var day int
	var level int
	var prompt bool
	var wait int64
	now := time.Now().UTC().Add(6 * time.Hour)

	flag.IntVar(&year, "year", now.Year(), "Year")
	flag.IntVar(&day, "day", now.Day(), "Day")
	flag.IntVar(&level, "level", 1, "Level")
	flag.StringVar(&token, "token", "", "Session Token")
	flag.BoolVar(&prompt, "prompt", false, "Show Prompt")
	flag.Int64Var(&wait, "wait", -1, "Wait for AoC until seconds before start")

	flag.Parse()

	if wait != -1 {
		aocStart := time.Date(year, now.Month(), day, 0, 0, 0, 0, time.UTC).Add(6 * time.Hour)
		fmt.Printf("AoC starts at %s\n", aocStart.Local())
		for {
			if now.Before(aocStart.Add(-time.Duration(wait) * time.Second)) {
				time.Sleep(time.Duration(aocStart.UnixNano() - now.UnixNano() - time.Duration(wait).Nanoseconds()))
			} else {
				break
			}
		}
		beeep.Alert("Advent of Code", fmt.Sprintf("AoC %d Day %d starts in %d seconds!", year, day, wait), "assets/icon.png")
		prompt = true
	}

	if prompt {
		fmt.Println(getPrompt(year, day))
		return
	}

	if day > 25 {
		fmt.Println("Day must be between 1 and 25")
		flag.Usage()
		return
	}
	fmt.Printf("Running Year %d Day %d\n", year, day)

	years[year](day, level)
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

	if token == "" {
		fmt.Println("Please provide a session token")
		flag.Usage()
		os.Exit(1)
		return nil
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

func getPrompt(year, day int) string {
	// Prompt is located in https://adventofcode.com/<year>/day/<day>
	fmt.Print("Downloading prompt from adventofcode.com\n\n")

	res, err := http.Get(fmt.Sprintf("https://adventofcode.com/%d/day/%d", year, day))
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
		panic("Failed to download prompt")
	}

	html := string(body)
	re := regexp.MustCompile("<(.|\n)*?>")
	html = re.ReplaceAllString(html, "")

	//Prompt starts after the line beginning with "window.addEventListener" and ends before the line beginning with "To play"
	lines := strings.Split(html, "\n")
	for i, line := range lines {
		if strings.Contains(line, "window.addEventListener") {
			lines = lines[i+1:]
			break
		}
	}

	for i, line := range lines {
		if strings.Contains(line, "To play") {
			lines = lines[:i]
			break
		}
	}

	return strings.Join(lines, "\n")
}

func postAnswers(year, day, level int, answer string) bool {
	if token == "" {
		fmt.Println("Please provide a session token")
		flag.Usage()
		os.Exit(1)
		return false
	}

	// Answers are posted to https://adventofcode.com/<year>/day/<day>/answer
	fmt.Println("Posting answers to adventofcode.com")

	req, err := http.NewRequest("POST", fmt.Sprintf("https://adventofcode.com/%d/day/%d/answer", year, day), strings.NewReader(fmt.Sprintf("level=%d&answer=%s", level, answer)))
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
		panic("Failed to post answers")
	}

	return strings.Contains(string(body), "the right answer")
}

func year2021(day, level int) {
	input := readInput(2021, day)
	fmt.Println()

	days := []func() string{
		func() string {
			if level == 1 {
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

			fmt.Printf("Output: %d\n", increases)
			return fmt.Sprint(increases)
		},
		func() string {
			if level == 1 {
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
		},
	}

	answer := days[day-1]()
	if answer != "" {
		if postAnswers(2021, day, level, answer) {
			fmt.Println("Solved day", day, "level", level)
		} else {
			fmt.Println("Failed to solve day", day, "level", level)
		}
	}
}
