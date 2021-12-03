package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gen2brain/beeep"
	"github.com/vonr/aoc/y2021"
)

var token string
var years = map[int]func(int, int, []string) string{
	2021: func(day, level int, input []string) string { return y2021.Days[day-1](level, input) },
}

func main() {
	var year int
	var day int
	var level int
	var prompt bool
	var wait int64
	now := time.Now().UTC()

	flag.IntVar(&year, "year", now.Year(), "Year")
	flag.IntVar(&day, "day", now.Day(), "Day")
	flag.IntVar(&level, "level", 1, "Level")
	flag.StringVar(&token, "token", "", "Session Token")
	flag.BoolVar(&prompt, "prompt", false, "Show Prompt")
	flag.Int64Var(&wait, "wait", -1, "Wait for AoC until seconds before start")

	flag.Parse()

	if wait != -1 {
		aocStart := time.Date(year, now.Month(), day, 5, 0, 0, 0, time.UTC)
		fmt.Printf("The time now is %s\n", now.Local())
		fmt.Printf("AoC starts at %s\n", aocStart.Local())
		for {
			if now.Before(aocStart.Add(-time.Duration(wait) * time.Second)) {
				time.Sleep(time.Duration(aocStart.UnixNano() - now.UnixNano() - time.Duration(wait*int64(time.Second)).Nanoseconds()))
			} else {
				break
			}
		}
		beeep.Alert("Advent of Code", fmt.Sprintf("AoC %d Day %d starts in %d seconds!", year, day, wait), "assets/icon.png")
		time.Sleep(time.Duration(wait) * time.Second)
		prompt = true
	}

	if token == "" {
		_, err := os.Stat(".token")
		if !os.IsNotExist(err) {
			// Read token from file
			file, err := os.Open(".token")
			if err != nil {
				panic(err)
			}
			defer file.Close()

			b, err := ioutil.ReadAll(file)
			if err != nil {
				panic(err)
			}

			token = strings.Split(string(b), "\n")[0]
		}
	}

	if prompt {
		fmt.Println(getPrompt(year, day, level))
		return
	}

	if day > 25 {
		fmt.Println("Day must be between 1 and 25")
		flag.Usage()
		return
	}

	fmt.Printf("Running Year %d Day %d\n", year, day)
	fmt.Println()

	answer := years[year](day, level, readInput(year, day))
	fmt.Printf("Output: %s\n", answer)

	if answer != "" {
		if postAnswers(year, day, level, answer) {
			fmt.Printf("Solved day %d %d level %d\n", year, day, level)
		} else {
			fmt.Printf("Failed to solve %d day %d level %d\n", year, day, level)
		}
	}
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

func getPrompt(year, day, level int) string {
	// Prompt is located in https://adventofcode.com/<year>/day/<day>
	fmt.Print("Downloading prompt from adventofcode.com\n\n")

	req, err := http.NewRequest("GET", fmt.Sprintf("https://adventofcode.com/%d/day/%d#part%d", year, day, level), nil)
	if err != nil {
		panic(err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = token
	req.AddCookie(cookie)

	req.Header.Set("Cookie", "session="+token)

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
		panic("Failed to download prompt")
	}

	html := string(body)
	re := regexp.MustCompile("<(.|\n)*?>")
	html = re.ReplaceAllString(html, "")

	//Prompt starts after the line beginning with "window." and ends before the line beginning with "(function"
	lines := strings.Split(html, "\n")
	for i, line := range lines {
		if strings.Contains(line, "window.") {
			lines = lines[i+1:]
			break
		}
	}

	for i, line := range lines {
		if strings.Contains(line, "(function") {
			lines = lines[:i]
			break
		}
	}

	return strings.ReplaceAll(strings.ReplaceAll(strings.Join(lines, "\n"), ".", ".\n"), "\n ", "\n")
}

func postAnswers(year, day, level int, answer string) bool {
	if token == "" {
		fmt.Println("Please provide a session token")
		flag.Usage()
		os.Exit(1)
		return false
	}

	// Answers are posted to https://adventofcode.com/<year>/day/<day>/answer
	fmt.Println("\nPosting answers to adventofcode.com")

	params := url.Values{}
	params.Add("level", fmt.Sprint(level))
	params.Add("answer", answer)
	payload := strings.NewReader(params.Encode())

	req, err := http.NewRequest("POST", fmt.Sprintf("https://adventofcode.com/%d/day/%d/answer", year, day), payload)
	if err != nil {
		panic(err)
	}

	cookie := new(http.Cookie)
	cookie.Name = "session"
	cookie.Value = token
	req.AddCookie(cookie)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

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

	if strings.Contains(string(body), "already complete") {
		fmt.Printf("Already completed %d day %d level %d\n", year, day, level)
		return false
	}

	if strings.Contains(string(body), "too recently") {
		fmt.Println("Slow down!")
		return false
	}

	return !strings.Contains(string(body), "not the right answer")
}
