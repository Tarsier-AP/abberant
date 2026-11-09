<script setup lang="ts">
import PopularNow from '@/components/home/PopularNow.vue'
import HomeSection from '@/components/home/HomeSection.vue'
import { headlines, sections } from '@/data/articles'
</script>

<template>
  <div class="home container">
    <PopularNow :items="headlines" />
    <HomeSection v-for="section in sections" :key="section.id" v-bind="section" />
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
