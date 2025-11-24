package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Question struct {
	ID          int      `json:"id"`
	Type        string   `json:"type"`
	Content     string   `json:"content"`
	Options     []string `json:"options"`
	Answer      string   `json:"answer"`
	Explanation string   `json:"explanation"`
}

type AnswerInfo struct {
	Answer      string
	Explanation string
}

func main() {
	// Adjust paths based on where this script is run from
	// Assuming run from quiz-app/tools/
	questions := parseQuestions("../../题库.txt")
	answers := parseAnswers("../../答案.txt")

	// Merge
	mergedCount := 0
	for i := range questions {
		q := &questions[i]
		if ans, ok := answers[q.ID]; ok {
			q.Answer = ans.Answer
			q.Explanation = ans.Explanation
			mergedCount++
		}
	}

	fmt.Printf("Parsed %d questions, merged %d answers.\n", len(questions), mergedCount)

	// Write to JSON in the root of quiz-app
	file, err := os.Create("../questions.json")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(questions)
	fmt.Println("Successfully generated questions.json")
}

func parseQuestions(path string) []Question {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening question file:", err)
		return []Question{}
	}
	defer file.Close()

	var questions []Question
	var currentQ *Question

	scanner := bufio.NewScanner(file)
	// Regex for question start: "1、..."
	reQStart := regexp.MustCompile(`^(\d+)、(.*)`)
	// Regex for option: "A、..."
	reOption := regexp.MustCompile(`^([A-Z])、(.*)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		// Check for Question Start
		if matches := reQStart.FindStringSubmatch(line); matches != nil {
			id, _ := strconv.Atoi(matches[1])
			content := matches[2]

			qType := "单选题" // Default
			if strings.Contains(line, "(多选题)") {
				qType = "多选题"
			} else if strings.Contains(line, "(判断题)") {
				qType = "判断题"
			} else if strings.Contains(line, "(单选题)") {
				qType = "单选题"
			}

			q := Question{
				ID:      id,
				Type:    qType,
				Content: content,
				Options: []string{},
			}
			questions = append(questions, q)
			currentQ = &questions[len(questions)-1]
			continue
		}

		// Check for Option
		if matches := reOption.FindStringSubmatch(line); matches != nil {
			if currentQ != nil {
				currentQ.Options = append(currentQ.Options, line)
			}
			continue
		}

		// Append to content if it's a continuation of the question and not an option
		// But be careful not to append garbage.
		// Assuming if it's not an option and we have a current question, it might be part of the question text.
		if currentQ != nil && len(currentQ.Options) == 0 {
			currentQ.Content += " " + line
		}
	}
	return questions
}

func parseAnswers(path string) map[int]AnswerInfo {
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error opening answer file:", err)
		return nil
	}
	defer file.Close()

	answers := make(map[int]AnswerInfo)
	scanner := bufio.NewScanner(file)

	// Regex to match: "1. B（脱贫攻坚精神）" or "61. 正确" or "20. D（...）- ..."
	// Group 1: ID
	// Group 2: Answer (A-Z, or A,B,C or 正确/错误)
	// Group 3: Explanation (content inside first parenthesis)
	reAnswer := regexp.MustCompile(`^(\d+)\.\s*([A-Z,]+|正确|错误)(?:[（(](.*?)[）)])?`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		if matches := reAnswer.FindStringSubmatch(line); matches != nil {
			id, _ := strconv.Atoi(matches[1])
			ans := matches[2]
			expl := ""
			if len(matches) > 3 {
				expl = matches[3]
			}

			// Clean up answer (remove spaces)
			ans = strings.ReplaceAll(ans, " ", "")

			answers[id] = AnswerInfo{
				Answer:      ans,
				Explanation: expl,
			}
		}
	}
	return answers
}
