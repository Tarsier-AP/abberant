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
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--card-gap);
}

@media (max-width: 1023px) {
  .section__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (max-width: 767px) {
  .section__grid {
    grid-template-columns: 1fr;
  }
}
</style>
