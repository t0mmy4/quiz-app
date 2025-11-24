package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sashabaranov/go-openai"
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
	apiKey := "sk-jisVIO44N30kkyusAKPc5b0YCmMX44CltY2546mQy2H693QA"
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://api.moonshot.cn/v1"
	client := openai.NewClientWithConfig(config)

	data, err := os.ReadFile("../questions.json")
	if err != nil {
		panic(err)
	}

	var questions []Question
	if err := json.Unmarshal(data, &questions); err != nil {
		panic(err)
	}

	targetIDs := map[uint]bool{
		52: true, 53: true, 54: true, 55: true, 56: true, 243: true,
	}

	for i := range questions {
		if targetIDs[questions[i].ID] {
			fmt.Printf("Fixing Question ID %d: %s\n", questions[i].ID, questions[i].Content)

			prompt := fmt.Sprintf(`
题目：%s
选项：%s
这是一道多选题。请给出正确答案的选项字母，如果有多个选项，请用逗号分隔（例如：A,B,C）。
只返回答案字母，不要包含其他文字。
`, questions[i].Content, strings.Join(questions[i].Options, "\n"))

			resp, err := client.CreateChatCompletion(
				context.Background(),
				openai.ChatCompletionRequest{
					Model: "kimi-k2-turbo-preview",
					Messages: []openai.ChatCompletionMessage{
						{
							Role:    openai.ChatMessageRoleSystem,
							Content: "你是一个专业的政治课助教。请只返回选项字母，用逗号分隔。",
						},
						{
							Role:    openai.ChatMessageRoleUser,
							Content: prompt,
						},
					},
				},
			)

			if err != nil {
				fmt.Printf("Error fetching AI for ID %d: %v\n", questions[i].ID, err)
				continue
			}

			newAnswer := strings.TrimSpace(resp.Choices[0].Message.Content)
			// Clean up answer (remove "答案：" etc if present, though system prompt should prevent it)
			newAnswer = strings.ReplaceAll(newAnswer, "答案：", "")
			newAnswer = strings.ReplaceAll(newAnswer, " ", "")
			newAnswer = strings.ReplaceAll(newAnswer, "，", ",") // Replace Chinese comma

			fmt.Printf("Old Answer: %s, New Answer: %s\n", questions[i].Answer, newAnswer)
			questions[i].Answer = newAnswer
		}
	}

	newData, err := json.MarshalIndent(questions, "", "  ")
	if err != nil {
		panic(err)
	}

	if err := os.WriteFile("../questions.json", newData, 0644); err != nil {
		panic(err)
	}

	fmt.Println("Done fixing answers.")
}
