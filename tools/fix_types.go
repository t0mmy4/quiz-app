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
		content := questions[i].Content
		if strings.Contains(content, "(多选题)") {
			if questions[i].Type != "多选题" {
				questions[i].Type = "多选题"
				count++
			}
		} else if strings.Contains(content, "(单选题)") {
			if questions[i].Type != "单选题" {
				questions[i].Type = "单选题"
				count++
			}
		}
	}

	newData, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("../questions.json", newData, 0644); err != nil {
		panic(err)
	}

	fmt.Printf("Fixed %d questions\n", count)
}
