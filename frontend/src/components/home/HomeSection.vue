<script setup lang="ts">
import type { ArticleSection } from '@/types/article'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import ArticleCard from '@/components/article/ArticleCard.vue'

defineProps<ArticleSection & { dense?: boolean }>()
</script>

<template>
  <section :id="id" class="section">
    <SectionHeader :title="title" :href="href" />

    <div class="section__grid" :class="{ 'section__grid--dense': dense }">
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

/* Archive grids: two cards per row on mobile and tablet, three on laptops. */
.section__grid--dense {
  grid-template-columns: repeat(2, minmax(0, 1fr));
}

@media (min-width: 1024px) {
  .section__grid--dense {
    grid-template-columns: repeat(3, minmax(0, 1fr));
  }
}
</style>
