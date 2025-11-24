<script setup>
import { onMounted } from 'vue'
import { useQuizStore } from './store/quiz'
import Sidebar from './components/Sidebar.vue'
import QuestionCard from './components/QuestionCard.vue'
import StatsBar from './components/StatsBar.vue'

const store = useQuizStore()

onMounted(() => {
  store.init()
})
</script>

<template>
  <div class="flex h-screen w-screen overflow-hidden bg-gray-100 font-sans">
    <Sidebar />
    <div class="flex-1 flex flex-col h-full relative">
      <div class="flex-1 overflow-hidden relative">
        <QuestionCard />
      </div>
      <StatsBar />
    </div>
  </div>
</template>
