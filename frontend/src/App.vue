<script setup lang="ts">
import SiteHeader from '@/components/layout/SiteHeader.vue'
import SiteFooter from '@/components/layout/SiteFooter.vue'
</script>

<template>
  <SiteHeader />
  <div class="page">
    <main class="page__main">
      <RouterView />
    </main>
    <SiteFooter />
  </div>
</template>

<style scoped>
.page {
  display: flex;
  flex-direction: column;
  width: 100%;
  max-width: var(--page-width);
  min-height: calc(100vh - var(--header-height));
  margin-inline: auto;
  background: var(--color-bg);
}

.page__main {
  flex: 1;
}
</style>
