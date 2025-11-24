package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/glebarez/sqlite"
	"github.com/sashabaranov/go-openai"
	"gorm.io/gorm"
)

//go:embed questions.json
var questionsJSON []byte

// Models
type Question struct {
	ID            uint   `gorm:"primaryKey" json:"id"`
	Type          string `json:"type"`
	Content       string `json:"content"`
	Options       string `json:"options"` // JSON string
	Answer        string `json:"-"`
	Explanation   string `json:"-"`
	AIExplanation string `json:"ai_explanation"`
}

type UserProgress struct {
	QuestionID uint   `gorm:"primaryKey" json:"question_id"`
	Status     int    `json:"status"` // 0: Unanswered, 1: Correct, 2: Wrong
	UserAnswer string `json:"user_answer"`
	IsMarked   bool   `json:"is_marked"`
}

type MistakeBook struct {
	QuestionID uint `gorm:"primaryKey" json:"question_id"`
	Count      int  `json:"count"`
}

// App struct
type App struct {
	ctx            context.Context
	db             *gorm.DB
	MistakeMode    bool
	aiClient       *openai.Client
	mistakeSession map[uint]int
}

// NewApp creates a new App application struct
func NewApp() *App {
	config := openai.DefaultConfig("")
	config.BaseURL = "https://api.moonshot.cn/v1"
	client := openai.NewClientWithConfig(config)
	return &App{
		aiClient:       client,
		mistakeSession: make(map[uint]int),
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.initDB()
}

func (a *App) initDB() {
	// Initialize SQLite DB
	cwd, _ := os.Getwd()
	dbPath := filepath.Join(cwd, "quiz.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	a.db = db

	// Migrate
	db.AutoMigrate(&Question{}, &UserProgress{}, &MistakeBook{})

	// Check if questions exist
	var count int64
	db.Model(&Question{}).Count(&count)
	if count == 0 {
		a.loadQuestions()
	} else {
		// Sync questions (e.g. fix types)
		a.syncQuestions()
	}
}

func (a *App) syncQuestions() {
	if len(questionsJSON) == 0 {
		return
	}

	var rawQuestions []struct {
		ID            uint     `json:"id"`
		Type          string   `json:"type"`
		Content       string   `json:"content"`
		Options       []string `json:"options"`
		Answer        string   `json:"answer"`
		Explanation   string   `json:"explanation"`
		AIExplanation string   `json:"ai_explanation"`
	}

	json.Unmarshal(questionsJSON, &rawQuestions)

	// Batch update might be complex with GORM for different values, so loop is fine for <10k items
	// Or we can just check and update if needed
	for _, rq := range rawQuestions {
		var q Question
		if err := a.db.First(&q, rq.ID).Error; err == nil {
			updated := false
			if q.Type != rq.Type {
				q.Type = rq.Type
				updated = true
			}

			// Check options
			opts, _ := json.Marshal(rq.Options)
			if q.Options != string(opts) {
				q.Options = string(opts)
				updated = true
			}

			// Check answer (for TF fix)
			if q.Answer != rq.Answer {
				q.Answer = rq.Answer
				updated = true
			}

			if updated {
				a.db.Save(&q)
			}
		}
	}
}

func (a *App) loadQuestions() {
	// Use embedded file
	if len(questionsJSON) == 0 {
		fmt.Println("questions.json is empty")
		return
	}

	var rawQuestions []struct {
		ID            uint     `json:"id"`
		Type          string   `json:"type"`
		Content       string   `json:"content"`
		Options       []string `json:"options"`
		Answer        string   `json:"answer"`
		Explanation   string   `json:"explanation"`
		AIExplanation string   `json:"ai_explanation"`
	}

	json.Unmarshal(questionsJSON, &rawQuestions)

	for _, rq := range rawQuestions {
		opts, _ := json.Marshal(rq.Options)
		q := Question{
			ID:            rq.ID,
			Type:          rq.Type,
			Content:       rq.Content,
			Options:       string(opts),
			Answer:        rq.Answer,
			Explanation:   rq.Explanation,
			AIExplanation: rq.AIExplanation,
		}
		a.db.Create(&q)
		// Init progress
		a.db.Create(&UserProgress{QuestionID: rq.ID, Status: 0})
	}
}

// API Methods

type QuestionView struct {
	ID            uint     `json:"id"`
	Type          string   `json:"type"`
	Content       string   `json:"content"`
	Options       []string `json:"options"`
	UserAnswer    string   `json:"user_answer"`
	Status        int      `json:"status"`
	IsMarked      bool     `json:"is_marked"`
	AIExplanation string   `json:"ai_explanation"`
	Explanation   string   `json:"explanation,omitempty"`
	CorrectAnswer string   `json:"correct_answer,omitempty"`
}

func (a *App) GetQuestion(id uint) QuestionView {
	var q Question
	var p UserProgress
	a.db.First(&q, id)
	a.db.First(&p, id)

	var opts []string
	json.Unmarshal([]byte(q.Options), &opts)

	// In Mistake Mode, use session status
	status := p.Status
	userAnswer := p.UserAnswer

	if a.MistakeMode {
		if s, ok := a.mistakeSession[id]; ok {
			status = s
			// If status is 0, user answer should be empty
			if s == 0 {
				userAnswer = ""
			}
		} else {
			status = 0
			userAnswer = ""
		}
	}

	qv := QuestionView{
		ID:            q.ID,
		Type:          q.Type,
		Content:       q.Content,
		Options:       opts,
		UserAnswer:    userAnswer,
		Status:        status,
		IsMarked:      p.IsMarked,
		AIExplanation: q.AIExplanation,
	}

	if status > 0 {
		qv.Explanation = q.Explanation
		qv.CorrectAnswer = q.Answer
	}

	return qv
}

type SubmitResult struct {
	Correct       bool   `json:"correct"`
	Explanation   string `json:"explanation"`
	CorrectAnswer string `json:"correct_answer"`
	AIExplanation string `json:"ai_explanation"`
}

func (a *App) SubmitAnswer(id uint, answer string) SubmitResult {
	var q Question
	var p UserProgress
	a.db.First(&q, id)
	a.db.First(&p, id)

	correct := q.Answer == answer

	status := 2
	if correct {
		status = 1
	}

	p.Status = status
	p.UserAnswer = answer
	a.db.Save(&p)

	if !correct {
		// Add to mistake book
		var mb MistakeBook
		res := a.db.First(&mb, id)
		if res.Error != nil {
			a.db.Create(&MistakeBook{QuestionID: id, Count: 1})
		} else {
			mb.Count++
			a.db.Save(&mb)
		}
	}

	// Update session if in mistake mode
	if a.MistakeMode {
		a.mistakeSession[id] = status
	}

	return SubmitResult{
		Correct:       correct,
		Explanation:   q.Explanation,
		CorrectAnswer: q.Answer,
		AIExplanation: q.AIExplanation,
	}
}

func (a *App) ToggleMark(id uint) bool {
	var p UserProgress
	a.db.First(&p, id)
	p.IsMarked = !p.IsMarked
	a.db.Save(&p)
	return p.IsMarked
}

type GridItem struct {
	ID       uint `json:"id"`
	Status   int  `json:"status"`
	IsMarked bool `json:"is_marked"`
}

func (a *App) GetGrid() []GridItem {
	var progress []UserProgress

	if a.MistakeMode {
		// Only return mistakes
		var mistakes []MistakeBook
		a.db.Find(&mistakes)
		ids := make([]uint, len(mistakes))
		for i, m := range mistakes {
			ids[i] = m.QuestionID
		}
		if len(ids) > 0 {
			a.db.Where("question_id IN ?", ids).Find(&progress)
		} else {
			progress = []UserProgress{}
		}
	} else {
		a.db.Find(&progress)
	}

	// Sort by ID
	sort.Slice(progress, func(i, j int) bool {
		return progress[i].QuestionID < progress[j].QuestionID
	})

	var grid []GridItem
	for _, p := range progress {
		status := p.Status
		if a.MistakeMode {
			if s, ok := a.mistakeSession[p.QuestionID]; ok {
				status = s
			} else {
				status = 0
			}
		}
		grid = append(grid, GridItem{
			ID:       p.QuestionID,
			Status:   status,
			IsMarked: p.IsMarked,
		})
	}
	return grid
}

func (a *App) SetMistakeMode(enable bool) {
	a.MistakeMode = enable
	if enable {
		// Clear session when entering mistake mode
		a.mistakeSession = make(map[uint]int)
	}
}

func (a *App) GetCorrectMistakesCount() int64 {
	var count int64
	// Find questions in mistake book that have status = 1 (Correct) in UserProgress
	// Note: We check the GLOBAL UserProgress, because that's what we want to clean up
	// based on the user's latest attempt (which is updated in SubmitAnswer)
	a.db.Table("mistake_books").
		Joins("JOIN user_progresses ON user_progresses.question_id = mistake_books.question_id").
		Where("user_progresses.status = ?", 1).
		Count(&count)
	return count
}

func (a *App) ClearCorrectMistakes() {
	// Delete from mistake_books where corresponding user_progress status is 1
	// SQLite doesn't support JOIN in DELETE easily, so do it in two steps or subquery
	subQuery := a.db.Table("user_progresses").Select("question_id").Where("status = ?", 1)
	a.db.Where("question_id IN (?)", subQuery).Delete(&MistakeBook{})
}

func (a *App) RemoveFromMistakeBook(id uint) {
	a.db.Delete(&MistakeBook{}, id)
}

type Stats struct {
	Total    int64  `json:"total"`
	Done     int64  `json:"done"`
	Correct  int64  `json:"correct"`
	Accuracy string `json:"accuracy"`
}

func (a *App) GetStats() Stats {
	var total, done, correct int64

	if a.MistakeMode {
		a.db.Model(&MistakeBook{}).Count(&total)
	} else {
		a.db.Model(&Question{}).Count(&total)
	}

	a.db.Model(&UserProgress{}).Where("status > 0").Count(&done)
	a.db.Model(&UserProgress{}).Where("status = 1").Count(&correct)

	acc := "0%"
	if done > 0 {
		acc = fmt.Sprintf("%.1f%%", float64(correct)/float64(done)*100)
	}

	return Stats{
		Total:    total,
		Done:     done,
		Correct:  correct,
		Accuracy: acc,
	}
}

// AI Generation
type AIResponse struct {
	Answer   string `json:"answer"`
	Analysis string `json:"analysis"`
}

func (a *App) GenerateAIExplanation(id uint, force bool) string {
	var q Question
	if err := a.db.First(&q, id).Error; err != nil {
		return "题目不存在"
	}

	// If not force and already has explanation, return it
	if !force && q.AIExplanation != "" {
		return q.AIExplanation
	}

	// Construct Prompt
	prompt := fmt.Sprintf(`
你是一个专业的政治课助教。请对以下题目进行解析。
题目：%s
选项：%s

要求：
1. 使用联网搜索功能查找相关背景知识。
2. 必须返回合法的 JSON 格式，包含 "answer" (你的答案，例如 "A") 和 "analysis" (详细解析内容) 两个字段。
3. 解析内容要深入浅出，逻辑清晰。
`, q.Content, q.Options)

	// Call API
	resp, err := a.aiClient.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: "kimi-k2-turbo-preview", // Updated to faster model
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "你是一个帮助学生学习的AI助手。请以JSON格式输出。",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			ResponseFormat: &openai.ChatCompletionResponseFormat{
				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
			},
		},
	)

	if err != nil {
		// Fallback to standard model if custom model fails (though API usually returns error)
		// For now just return error message
		return fmt.Sprintf("AI 生成失败: %v", err)
	}

	content := resp.Choices[0].Message.Content

	// Parse JSON to ensure format and maybe reformat for display if needed
	// But user wants specific format in the JSON.
	// Let's just save the raw content or parsed content?
	// The prompt asks for specific style inside the JSON fields.
	// Let's try to parse it to verify.
	var aiResp AIResponse
	if err := json.Unmarshal([]byte(content), &aiResp); err == nil {
		// Reformat to string for storage/display
		finalText := fmt.Sprintf("kimi说：\n答案：%s\n解析：%s", aiResp.Answer, aiResp.Analysis)
		q.AIExplanation = finalText
		a.db.Save(&q)
		return finalText
	}

	// If JSON parse fails, just save the raw content if it looks like text
	q.AIExplanation = content
	a.db.Save(&q)
	return content
}
