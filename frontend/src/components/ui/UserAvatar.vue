<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    /** Account name used for the fallback letter and color. */
    name: string
    /** Uploaded avatar as a data URL; empty means the letter fallback. */
    avatar?: string
    /** Side of the square avatar in px. */
    size?: number
  }>(),
  { avatar: '', size: 36 },
)

const initial = computed(() => {
  const letter = props.name.trim().charAt(0)

  return letter ? letter.toUpperCase() : '?'
})

const hue = computed(() => {
  let hash = 0

  for (const char of props.name) {
    hash = (hash * 31 + (char.codePointAt(0) ?? 0)) % 360
  }

  return hash
})

const sizeStyle = computed(() => ({
  width: `${props.size}px`,
  height: `${props.size}px`,
  fontSize: `${Math.round(props.size * 0.44)}px`,
}))

const letterStyle = computed(() => ({
  ...sizeStyle.value,
  background: `hsl(${hue.value} 42% 42%)`,
}))
</script>

<template>
  <img
    v-if="avatar"
    class="user-avatar"
    :src="avatar"
    :alt="`Аватар: ${name}`"
    :style="sizeStyle"
  />
  <span
    v-else
    class="user-avatar user-avatar--letter"
    :style="letterStyle"
    aria-hidden="true"
    >{{ initial }}</span
  >
</template>

<style scoped>
.user-avatar {
  display: inline-flex;
  flex: none;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  object-fit: cover;
  color: #fff;
  font-family: Roboto, sans-serif;
  font-weight: 500;
  line-height: 1;
  user-select: none;
}
</style>
