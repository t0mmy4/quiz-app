import { defineStore } from 'pinia'
// These imports will work after 'wails dev' generates the bindings
import { GetGrid, GetQuestion, SubmitAnswer, ToggleMark, SetMistakeMode, GetStats, RemoveFromMistakeBook, GenerateAIExplanation, GetCorrectMistakesCount, ClearCorrectMistakes } from '../../wailsjs/go/main/App'

export const useQuizStore = defineStore('quiz', {
  state: () => ({
    currentQuestion: null,
    grid: [],
    stats: { total: 0, done: 0, correct: 0, accuracy: "0%" },
    mistakeMode: false,
    showExplanation: false,
    lastResult: null, // { correct, explanation, correct_answer, ai_explanation }
    loading: false,
    aiThinking: false
  }),
  actions: {
    async init() {
      await this.fetchGrid()
      await this.fetchStats()
      // Load first question if grid is not empty and no current question
      if (this.grid.length > 0 && !this.currentQuestion) {
        await this.fetchQuestion(this.grid[0].id)
      }
    },
    async fetchGrid() {
      try {
        this.grid = await GetGrid() || []
      } catch (e) {
        console.error(e)
        this.grid = []
      }
    },
    async fetchStats() {
      try {
        this.stats = await GetStats()
      } catch (e) {
        console.error(e)
      }
    },
    async fetchQuestion(id) {
      this.loading = true
      this.showExplanation = false
      this.lastResult = null
      try {
        this.currentQuestion = await GetQuestion(id)
        
        if (this.currentQuestion.status > 0) {
            this.showExplanation = true
            this.lastResult = {
                correct: this.currentQuestion.status === 1,
                explanation: this.currentQuestion.explanation,
                correct_answer: this.currentQuestion.correct_answer,
                ai_explanation: this.currentQuestion.ai_explanation
            }
        }
      } catch (e) {
        console.error(e)
      }
      this.loading = false
    },
    async submit(answer) {
      if (!this.currentQuestion) return
      
      try {
        const result = await SubmitAnswer(this.currentQuestion.id, answer)
        this.lastResult = result
        this.showExplanation = true
        
        // Update local state
        this.currentQuestion.status = result.correct ? 1 : 2
        this.currentQuestion.user_answer = answer
        
        await this.fetchGrid()
        await this.fetchStats()

        // If in mistake mode and correct, ask to remove
        // Logic changed: Don't ask immediately. Just auto jump.
        // Cleanup happens when leaving mistake mode.
        if (result.correct) {
            // Auto jump to next question if correct
            setTimeout(() => {
                this.nextQuestion()
            }, 800)
        }
      } catch (e) {
        console.error(e)
      }
    },
    async toggleMark() {
      if (!this.currentQuestion) return
      try {
        const marked = await ToggleMark(this.currentQuestion.id)
        this.currentQuestion.is_marked = marked
        await this.fetchGrid()
      } catch (e) {
        console.error(e)
      }
    },
    async setMode(mode) {
      // If switching FROM mistake mode TO normal mode
      if (this.mistakeMode && !mode) {
        try {
            const count = await GetCorrectMistakesCount()
            if (count > 0) {
                if (confirm(`错题本中有 ${count} 道题已答对，是否将它们移出错题本？`)) {
                    await ClearCorrectMistakes()
                }
            }
        } catch (e) {
            console.error(e)
        }
      }

      this.mistakeMode = mode
      this.currentQuestion = null
      try {
        await SetMistakeMode(mode)
        await this.init()
      } catch (e) {
        console.error(e)
      }
    },
    async generateAI(force = false) {
        if (!this.currentQuestion) return
        this.aiThinking = true
        try {
            const explanation = await GenerateAIExplanation(this.currentQuestion.id, force)
            // Update current question's AI explanation
            if (this.currentQuestion) {
                this.currentQuestion.ai_explanation = explanation
            }
            // Also update lastResult if it exists
            if (this.lastResult) {
                this.lastResult.ai_explanation = explanation
            }
        } catch (e) {
            console.error(e)
        }
        this.aiThinking = false
    },
    nextQuestion() {
        if (!this.currentQuestion || this.grid.length === 0) return
        const idx = this.grid.findIndex(g => g.id === this.currentQuestion.id)
        if (idx !== -1 && idx < this.grid.length - 1) {
            this.fetchQuestion(this.grid[idx + 1].id)
        }
    },
    prevQuestion() {
        if (!this.currentQuestion || this.grid.length === 0) return
        const idx = this.grid.findIndex(g => g.id === this.currentQuestion.id)
        if (idx > 0) {
            this.fetchQuestion(this.grid[idx - 1].id)
        }
    }
  }
})
