<script setup>
import { useQuizStore } from '../store/quiz'
const store = useQuizStore()

const getBgClass = (item) => {
    if (item.id === store.currentQuestion?.id) return 'ring-2 ring-blue-500 z-10'
    if (item.status === 1) return 'bg-green-500 text-white border-green-600'
    if (item.status === 2) return 'bg-red-500 text-white border-red-600'
    return 'bg-white hover:bg-gray-100 text-gray-700'
}
</script>

<template>
  <div class="w-72 h-full bg-gray-50 border-r flex flex-col shadow-lg z-20">
    <div class="p-4 border-b bg-white shadow-sm">
        <h2 class="font-bold text-lg text-gray-800">题目列表</h2>
        <div class="flex items-center mt-3 space-x-2 text-sm">
            <button 
                @click="store.setMode(false)"
                :class="!store.mistakeMode ? 'bg-blue-600 text-white shadow-md' : 'bg-gray-200 text-gray-600 hover:bg-gray-300'"
                class="flex-1 py-2 rounded transition-all font-medium">
                全部
            </button>
            <button 
                @click="store.setMode(true)"
                :class="store.mistakeMode ? 'bg-red-600 text-white shadow-md' : 'bg-gray-200 text-gray-600 hover:bg-gray-300'"
                class="flex-1 py-2 rounded transition-all font-medium">
                错题本
            </button>
        </div>
    </div>
    <div class="flex-1 overflow-y-auto p-3 custom-scrollbar">
        <div v-if="store.grid.length === 0" class="text-center text-gray-500 mt-10">
            {{ store.mistakeMode ? '暂无错题' : '暂无题目' }}
        </div>
        <div class="grid grid-cols-5 gap-2">
            <div 
                v-for="item in store.grid" 
                :key="item.id"
                @click="store.fetchQuestion(item.id)"
                :class="[getBgClass(item), 'cursor-pointer rounded flex items-center justify-center h-10 text-sm font-medium border relative transition-all duration-200']"
            >
                {{ item.id }}
                <span v-if="item.is_marked" class="absolute -top-1 -right-1 w-3 h-3 bg-yellow-400 rounded-full border-2 border-white"></span>
            </div>
        </div>
    </div>
  </div>
</template>

<style scoped>
.custom-scrollbar::-webkit-scrollbar {
  width: 6px;
}
.custom-scrollbar::-webkit-scrollbar-track {
  background: #f1f1f1; 
}
.custom-scrollbar::-webkit-scrollbar-thumb {
  background: #c1c1c1; 
  border-radius: 3px;
}
.custom-scrollbar::-webkit-scrollbar-thumb:hover {
  background: #a8a8a8; 
}
</style>
