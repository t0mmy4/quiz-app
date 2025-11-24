package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Question struct {
	ID            uint     `json:"id"`
	Type          string   `json:"type"`
	Content       string   `json:"content"`
	Options       []string `json:"options"`
	Answer        string   `json:"answer"`
	Explanation   string   `json:"explanation"`
	AIExplanation string   `json:"ai_explanation"`
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
	for i := range questions {
		// Check for empty options or explicit "判断题" type/content
		isTF := false
		if len(questions[i].Options) == 0 {
			isTF = true
		} else if questions[i].Type == "判断题" {
			isTF = true
		} else if strings.Contains(questions[i].Content, "(判断题)") {
			isTF = true
		}

		if isTF {
			// Fix Type
			questions[i].Type = "判断题"

			// Fix Options if empty
			if len(questions[i].Options) == 0 {
				questions[i].Options = []string{"A、正确", "B、错误"}
			}

			// Fix Answer
			ans := strings.TrimSpace(questions[i].Answer)
			if ans == "正确" || ans == "对" {
				questions[i].Answer = "A"
			} else if ans == "错误" || ans == "错" {
				questions[i].Answer = "B"
			}

			count++
		}
	}

	newData, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("../questions.json", newData, 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Fixed %d True/False questions\n", count)
}
