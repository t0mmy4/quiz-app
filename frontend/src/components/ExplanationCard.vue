<script setup>
import { useQuizStore } from '../store/quiz'
import { marked } from 'marked'
import { computed } from 'vue'

const store = useQuizStore()

const aiExplanationHtml = computed(() => {
    if (store.lastResult && store.lastResult.ai_explanation) {
        return marked.parse(store.lastResult.ai_explanation)
    }
    return ''
})
</script>

<template>
  <div v-if="store.showExplanation && store.lastResult" class="mt-6 p-5 rounded-xl border-l-4 shadow-md transition-all duration-500 ease-in-out transform translate-y-0 opacity-100"
    :class="store.lastResult.correct ? 'bg-green-50 border-green-500' : 'bg-red-50 border-red-500'">
    <div class="flex items-start">
        <div class="flex-1">
            <h3 class="font-bold text-lg mb-3 flex items-center" :class="store.lastResult.correct ? 'text-green-700' : 'text-red-700'">
                <span class="text-2xl mr-2">{{ store.lastResult.correct ? '✓' : '✗' }}</span>
                {{ store.lastResult.correct ? '回答正确' : '回答错误' }}
            </h3>
            <div class="space-y-2">
                <p class="text-gray-800">
                    <span class="font-bold text-gray-900">正确答案：</span> 
                    <span class="font-mono text-lg font-bold text-blue-600">{{ store.lastResult.correct_answer }}</span>
                </p>
                <div class="bg-white bg-opacity-60 p-3 rounded-lg mt-2">
                    <p v-if="store.lastResult.explanation" class="text-gray-800 leading-relaxed">
                        <span class="font-bold text-gray-900 block mb-1">解析：</span> 
                        {{ store.lastResult.explanation }}
                    </p>
                    <p v-else class="text-gray-500 italic text-sm">暂无解析</p>
                </div>

                <!-- AI Explanation -->
                <div class="mt-4 pt-4 border-t border-gray-200">
                    <div class="flex justify-between items-center mb-2">
                        <h4 class="font-bold text-purple-700 flex items-center">
                            <span class="mr-2">🤖</span> AI 深度解析 (Kimi)
                        </h4>
                        <button 
                            v-if="store.lastResult.ai_explanation && !store.aiThinking"
                            @click="store.generateAI(true)"
                            class="text-xs text-purple-600 hover:text-purple-800 underline"
                        >
                            重新思考
                        </button>
                    </div>

                    <div v-if="store.aiThinking" class="p-4 bg-purple-50 rounded-lg flex items-center justify-center text-purple-600">
                        <svg class="animate-spin h-5 w-5 mr-3" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                        </svg>
                        Kimi 正在思考中...
                    </div>

                    <div v-else-if="store.lastResult.ai_explanation" 
                         class="bg-purple-50 p-4 rounded-lg text-gray-800 text-sm leading-relaxed border border-purple-100 markdown-body"
                         v-html="aiExplanationHtml">
                    </div>

                    <div v-else>
                        <button 
                            @click="store.generateAI(false)"
                            class="w-full py-2 bg-purple-100 hover:bg-purple-200 text-purple-700 rounded-lg transition-colors text-sm font-medium flex items-center justify-center"
                        >
                            <span>✨ 生成 AI 解析</span>
                        </button>
                    </div>
                </div>
            </div>
        </div>
    </div>
  </div>
</template>
