<script setup lang="ts">
import type { ArticleSection } from '@/types/article'
import SectionHeader from '@/components/ui/SectionHeader.vue'
import ArticleCard from '@/components/article/ArticleCard.vue'

defineProps<ArticleSection & { dense?: boolean; showAll?: boolean }>()
</script>

<template>
  <section :id="id" class="section">
    <SectionHeader :title="title" :href="href" />

    <div
      class="section__grid"
      :class="{
        'section__grid--dense': dense,
        'section__grid--all': showAll,
      }"
    >
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

/*
 * Home grids always keep two rows: 3x2 on desktop, 2x2 on tablet, 1x2 on mobile.
 * Cards past the second row are hidden at the matching breakpoint.
 */
.section__grid {
  display: grid;
  grid-template-columns: repeat(3, minmax(0, 1fr));
  gap: var(--card-gap);
}

.section__grid:not(.section__grid--dense):not(.section__grid--all) > :nth-child(n + 7) {
  display: none;
}

@media (max-width: 1023px) {
  .section__grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }

  .section__grid:not(.section__grid--dense):not(.section__grid--all) > :nth-child(n + 5) {
    display: none;
  }
}

@media (max-width: 767px) {
  .section__grid {
    grid-template-columns: 1fr;
  }

  .section__grid:not(.section__grid--dense):not(.section__grid--all) > :nth-child(n + 3) {
    display: none;
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
