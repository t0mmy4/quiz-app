<script setup>
import { computed, ref, watch } from 'vue'
import { useQuizStore } from '../store/quiz'
import ExplanationCard from './ExplanationCard.vue'

const store = useQuizStore()
const selected = ref([]) // For multi-select
const singleSelected = ref('') // For single select

// Reset selection when question changes
watch(() => store.currentQuestion, (newQ) => {
    selected.value = []
    singleSelected.value = ''
    if (newQ && newQ.user_answer && newQ.status > 0) {
        // If already answered, show what they selected?
        // Since we don't have the full result object from GetQuestion, we can't show the full explanation card immediately
        // unless we fetch it or store it.
        // But we can show the user's answer.
        if (newQ.type === '多选题') {
            selected.value = newQ.user_answer.split(',')
        } else {
            singleSelected.value = newQ.user_answer
        }
    }
})

const isMulti = computed(() => store.currentQuestion?.type === '多选题')

const toggleSelection = (opt) => {
    if (store.showExplanation) return // Disable interaction if result shown
    // Also disable if status > 0 (already answered) unless we allow re-answering?
    // Let's allow re-answering only if not in "showExplanation" mode (which is transient).
    // But if status > 0, we might want to block.
    // For now, let's assume if showExplanation is false, they can try.

    const letter = opt.split('、')[0]
    if (isMulti.value) {
        if (selected.value.includes(letter)) {
            selected.value = selected.value.filter(s => s !== letter)
        } else {
            selected.value.push(letter)
        }
    } else {
        singleSelected.value = letter
        store.submit(letter)
    }
}

const submitMulti = () => {
    if (selected.value.length === 0) return
    const ans = selected.value.sort().join(',')
    store.submit(ans)
}

const getOptionClass = (opt) => {
    const letter = opt.split('、')[0]
    const isSelected = isMulti.value ? selected.value.includes(letter) : singleSelected.value === letter
    
    // If explanation is shown, highlight correct/wrong
    if (store.showExplanation && store.lastResult) {
        const correct = store.lastResult.correct_answer.includes(letter)
        // Logic for highlighting:
        // Green if it is part of the correct answer.
        // Red if selected but NOT part of correct answer.
        // Also, if it IS part of correct answer but NOT selected, maybe yellow?
        // Let's keep it simple:
        // Correct Answer -> Green
        // Selected Wrong -> Red
        
        if (correct) return 'bg-green-100 border-green-500 text-green-900 font-bold'
        if (isSelected && !correct) return 'bg-red-100 border-red-500 text-red-900'
    }
    
    if (isSelected) return 'bg-blue-100 border-blue-500 text-blue-900 font-bold shadow-sm'
    return 'bg-white border-gray-200 hover:bg-gray-50 text-gray-700'
}
</script>

<template>
  <div v-if="store.currentQuestion" class="max-w-4xl mx-auto p-8 h-full overflow-y-auto custom-scrollbar">
    <div class="mb-8">
        <div class="flex items-center mb-4">
            <span class="inline-block px-3 py-1 text-sm font-bold text-white bg-blue-600 rounded-full shadow-sm mr-3">
                {{ store.currentQuestion.type }}
            </span>
            <span class="text-gray-500 font-mono text-sm">ID: {{ store.currentQuestion.id }}</span>
        </div>
        <h1 class="text-2xl font-bold leading-relaxed text-gray-800 tracking-wide">
            {{ store.currentQuestion.content }}
        </h1>
    </div>

    <div class="space-y-4">
        <div 
            v-for="opt in store.currentQuestion.options" 
            :key="opt"
            @click="toggleSelection(opt)"
            :class="[getOptionClass(opt), 'p-5 border-2 rounded-xl cursor-pointer transition-all duration-200 flex items-center text-lg']"
        >
            <span class="w-8 h-8 flex items-center justify-center rounded-full border-2 mr-4 text-sm font-bold"
                :class="getOptionClass(opt).includes('bg-blue') ? 'border-blue-500 bg-blue-200' : 'border-gray-300 bg-gray-100'">
                {{ opt.split('、')[0] }}
            </span>
            <span>{{ opt.substring(opt.indexOf('、') + 1) }}</span>
        </div>
    </div>

    <div v-if="isMulti && !store.showExplanation" class="mt-8">
        <button 
            @click="submitMulti"
            :disabled="selected.length === 0"
            :class="selected.length === 0 ? 'bg-gray-300 cursor-not-allowed' : 'bg-blue-600 hover:bg-blue-700 shadow-lg hover:shadow-xl'"
            class="w-full py-4 text-white rounded-xl font-bold text-lg transition-all transform active:scale-95">
            提交答案
        </button>
    </div>

    <ExplanationCard />
    
    <!-- Spacer for bottom bar -->
    <div class="h-20"></div>
  </div>
  <div v-else class="flex flex-col items-center justify-center h-full text-gray-500">
    <div class="text-6xl mb-4">📚</div>
    <p class="text-xl">请从左侧选择一道题目开始刷题</p>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 8px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: transparent; 
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #d1d5db; 
  border-radius: 4px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #9ca3af; 
}
</style>
