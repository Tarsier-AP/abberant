<template>
  <svg class="theme-icon" viewBox="0 0 48 48" fill="none" aria-hidden="true">
    <circle class="theme-icon__shadow" cx="21" cy="26" r="12" />
    <circle class="theme-icon__core" cx="24" cy="24" r="10" />
    <circle class="theme-icon__ring" cx="24" cy="24" r="8" />
    <g class="theme-icon__rays">
      <path d="M24 4v6" />
      <path d="M24 38v6" />
      <path d="M4 24h6" />
      <path d="M38 24h6" />
      <path d="m9.5 9.5 4.2 4.2" />
      <path d="m34.3 34.3 4.2 4.2" />
      <path d="m38.5 9.5-4.2 4.2" />
      <path d="m13.7 34.3-4.2 4.2" />
    </g>
  </svg>
</template>

<style scoped>
.theme-icon {
  width: 32px;
  height: 32px;
}

.theme-icon__shadow {
  fill: var(--theme-icon-shadow);
}

.theme-icon__core {
  fill: var(--color-card);
}

.theme-icon__ring {
  fill: none;
  stroke: var(--theme-icon-ring);
  stroke-width: 4;
}

.theme-icon__rays {
  stroke: #fff39a;
  stroke-width: 2;
  stroke-linecap: round;
}
</style>
