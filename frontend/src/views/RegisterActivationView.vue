<script setup lang="ts">
import { computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import AccountActivationCard from '@/components/auth/AccountActivationCard.vue'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const auth = useAuthStore()

const login = computed(() => auth.account?.login ?? '')
const username = computed(() => auth.account?.username ?? '')

onMounted(() => {
  if (!auth.account) {
    void router.replace({ name: 'register' })
  }
})
</script>

<template>
  <div class="register-activation">
    <AccountActivationCard v-if="login" :login="login" :username="username" />
  </div>
</template>

<style scoped>
.register-activation {
  position: relative;
  display: flex;
  flex: 1;
  align-items: center;
  justify-content: center;
  padding: 32px 16px;
}
</style>
