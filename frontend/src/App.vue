<script setup lang="ts">
import SiteHeader from '@/components/layout/SiteHeader.vue'
import SiteFooter from '@/components/layout/SiteFooter.vue'
</script>

<template>
  <div class="app">
    <SiteHeader />
    <div class="page">
      <main class="page__main">
        <RouterView />
      </main>
    </div>
    <SiteFooter />
  </div>
</template>

<style scoped>
.app {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
}

.page {
  display: flex;
  flex: 1;
  flex-direction: column;
  width: 100%;
  max-width: var(--page-width);
  margin-inline: auto;
  background: var(--color-bg);
}

.page__main {
  display: flex;
  flex: 1;
  flex-direction: column;
}
</style>
