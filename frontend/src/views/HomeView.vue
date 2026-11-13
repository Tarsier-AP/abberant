<script setup lang="ts">
import PopularNow from '@/components/home/PopularNow.vue'
import HomeSection from '@/components/home/HomeSection.vue'
import LoadState from '@/components/ui/LoadState.vue'
import { useArticlesStore } from '@/stores/articles'

const store = useArticlesStore()
</script>

<template>
  <div class="home container">
    <LoadState
      v-if="store.loadState !== 'ready'"
      :error="store.errorMessage"
      @retry="store.reload()"
    />
    <template v-else>
      <PopularNow :items="store.popularNow" />
      <HomeSection v-for="section in store.homeSections" :key="section.id" v-bind="section" />
    </template>
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
