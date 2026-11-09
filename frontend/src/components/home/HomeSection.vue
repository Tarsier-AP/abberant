<script setup lang="ts">
import type { ArticleSection } from '@/types/article'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import ArticleCard from '@/components/article/ArticleCard.vue'

defineProps<ArticleSection>()
</script>

<template>
  <section :id="id" class="section">
    <SectionHeader :title="title" :href="href" />

    <div class="section__grid">
      <ArticleCard v-for="item in articles" :key="item.id" :article="item" />
    </div>
  </section>
</template>

<style scoped>
.section {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
}

.section__grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: var(--card-gap);
}
</style>
