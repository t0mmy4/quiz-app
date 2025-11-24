<script setup>
import { useQuizStore } from '../store/quiz'
const store = useQuizStore()
</script>

<template>
  <div class="h-16 bg-white border-t flex items-center justify-between px-8 text-sm text-gray-600 shadow-[0_-4px_6px_-1px_rgba(0,0,0,0.05)] z-30">
    <div class="flex space-x-8">
        <div class="flex flex-col">
            <span class="text-xs text-gray-400">总题数</span>
            <b class="text-lg text-gray-800">{{ store.stats.total }}</b>
        </div>
        <div class="flex flex-col">
            <span class="text-xs text-gray-400">已完成</span>
            <b class="text-lg text-blue-600">{{ store.stats.done }}</b>
        </div>
        <div class="flex flex-col">
            <span class="text-xs text-gray-400">正确数</span>
            <b class="text-lg text-green-600">{{ store.stats.correct }}</b>
        </div>
        <div class="flex flex-col">
            <span class="text-xs text-gray-400">正确率</span>
            <b class="text-lg text-purple-600">{{ store.stats.accuracy }}</b>
        </div>
    </div>
    <div class="flex space-x-4">
        <button @click="store.prevQuestion()" class="px-4 py-2 rounded-lg border hover:bg-gray-50 hover:text-blue-600 transition-colors flex items-center">
            ← 上一题
        </button>
        <button @click="store.toggleMark()" class="px-4 py-2 rounded-lg border hover:bg-yellow-50 hover:border-yellow-200 hover:text-yellow-600 transition-colors flex items-center min-w-[100px] justify-center">
            <span v-if="store.currentQuestion?.is_marked" class="font-bold text-yellow-500">★ 已标记</span>
            <span v-else>☆ 标记</span>
        </button>
        <button @click="store.nextQuestion()" class="px-4 py-2 rounded-lg bg-blue-600 text-white hover:bg-blue-700 shadow-md hover:shadow-lg transition-all flex items-center">
            下一题 →
        </button>
    </div>
  </div>
</template>
