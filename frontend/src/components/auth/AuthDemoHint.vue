<script setup lang="ts">
import { TWO_FACTOR_DEMO_HOTKEY } from '@/composables/useTwoFactorDemoShortcut'
</script>

<template>
  <p class="auth-demo-hint">Демо 2FA: {{ TWO_FACTOR_DEMO_HOTKEY }}</p>
</template>

<style scoped>
.auth-demo-hint {
  position: absolute;
  right: 16px;
  bottom: 8px;
  color: var(--color-danger);
  font: 700 12px/16px Roboto, sans-serif;
  text-align: right;
}
</style>
