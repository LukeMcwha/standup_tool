package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"
)

type Config struct {
	Members   []string `json:"members"`
	Questions []string `json:"questions"`
	SlackHook string   `json:"slackHook"`
}

type Question struct {
	Title  string `json:"title"`
	Answer string `json:"answer"`
}

func Shuffle(a []string) []string {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	rand.Shuffle(len(a), func(i, j int) { a[i], a[j] = a[j], a[i] })
	return a
}

func FormatSlackMessage(answers map[string][]Question) string {
	var sb strings.Builder

	sb.WriteString("```\n")

	for name, questions := range answers {
		sb.WriteString(fmt.Sprintf("%s:\n", strings.Title(name)))

		for _, question := range questions {
			sb.WriteString(fmt.Sprintf("%s: %s:\n", question.Title, question.Answer))
		}

		sb.WriteString("\n")
	}

	sb.WriteString("```")

	return sb.String()

}

func main() {
	configFileLocation := "config.json"
	saveFileLocation := "save.json"

	// Get Config and read it into struct
	jsonFile, err := os.Open(configFileLocation)
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		panic(err)
	}
	config := Config{}
	err = json.Unmarshal(byteValue, &config)
	if err != nil {
		panic(err)
	}

	// Yesterdays Information
	standupJ, err := os.Open(saveFileLocation)
	if err != nil {
		fmt.Println(err)
	}
	b, err := ioutil.ReadAll(standupJ)
	if err != nil {
		panic(err)
	}
	yesterday := map[string][]Question{}
	err = json.Unmarshal(b, &yesterday)
	if err != nil {
		panic(err)
	}
	standupJ.Close()

	// Create a reader
	reader := bufio.NewReader(os.Stdin)

	randomOrder := Shuffle(config.Members)

	fmt.Printf("The random order for the day is: %v\n\n", strings.Join(randomOrder, ", "))

	answers := make(map[string][]Question, len(randomOrder))

	for i := 0; i < len(randomOrder); i++ {
		name := randomOrder[i]
		questions := make([]Question, len(config.Questions))

		fmt.Printf("-- %s\n\n", name)

		for j, question := range config.Questions {
			info := yesterday[name][j].Answer

			fmt.Printf("%s? (%s)\n", question, info)
			text, _ := reader.ReadString('\n')

			// convert CRLF to LF
			text = strings.Replace(text, "\n", "", -1)

			if text == "" {
				text = info
			}

			ans := Question{
				Title:  question,
				Answer: text,
			}

			questions[j] = ans
		}
		fmt.Print("\n")

		answers[name] = questions
	}

	// Write to file
	file, _ := json.MarshalIndent(answers, "", " ")
	_ = ioutil.WriteFile(saveFileLocation, file, 0644)

	// Print out to the terminal. This can then be copy pasted into slack.
	slackFormated := FormatSlackMessage(answers)
	fmt.Println(slackFormated)
}
