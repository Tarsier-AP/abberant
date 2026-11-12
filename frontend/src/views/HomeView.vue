<script setup lang="ts">
import PopularNow from '@/components/home/PopularNow.vue'
import HomeSection from '@/components/home/HomeSection.vue'
import { useArticlesStore } from '@/stores/articles'
import { headlines } from '@/data/articles'

const store = useArticlesStore()
</script>

<template>
  <div class="home container">
    <PopularNow :items="headlines" />
    <HomeSection v-for="section in store.homeSections" :key="section.id" v-bind="section" />
  </div>
</template>

<style scoped>
.home {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  padding-block: var(--gap);
}
</style>
