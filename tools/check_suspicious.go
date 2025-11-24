package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Question struct {
	ID      uint     `json:"id"`
	Type    string   `json:"type"`
	Content string   `json:"content"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

func main() {
	data, err := os.ReadFile("../questions.json")
	if err != nil {
		panic(err)
	}

	var questions []Question
	if err := json.Unmarshal(data, &questions); err != nil {
		panic(err)
	}

	count := 0
	for _, q := range questions {
		if q.Type == "多选题" {
			// Check if answer is single letter
			if !strings.Contains(q.Answer, ",") && len(q.Answer) == 1 {
				fmt.Printf("ID: %d, Answer: %s, Content: %s\n", q.ID, q.Answer, q.Content)
				count++
			}
		}
	}
	fmt.Printf("Total suspicious questions: %d\n", count)
}
