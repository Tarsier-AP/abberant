<script setup lang="ts">
interface FooterLink {
  label: string
  to?: string
  href?: string
}

interface FooterGroup {
  id: 'about' | 'social' | 'documents'
  title: string
  links: FooterLink[]
}

const groups: FooterGroup[] = [
  {
    id: 'about',
    title: 'О нас',
    links: [
      { label: 'О издании', to: '/about' },
      { label: 'Карьера в Abberant', to: '/career' },
    ],
  },
  {
    id: 'social',
    title: 'Социальные сети',
    links: [
      { label: 'ВКонтакте', href: 'https://vk.com' },
      { label: 'Телеграмм', href: 'https://t.me' },
    ],
  },
  {
    id: 'documents',
    title: 'Документы',
    links: [
      { label: 'Политика конфиденциальности', to: '/privacy' },
      { label: 'Условия использования сайта', to: '/terms' },
      { label: 'Политика обработки персональных данных', to: '/personal-data' },
    ],
  },
]
</script>

<template>
  <footer class="footer">
    <div class="footer__inner">
      <section
        v-for="group in groups"
        :key="group.id"
        class="footer__card"
        :class="`footer__${group.id}`"
        :aria-labelledby="`footer-${group.id}`"
      >
        <h2 :id="`footer-${group.id}`" class="footer__heading">{{ group.title }}</h2>

        <template v-for="link in group.links" :key="link.label">
          <a
            v-if="link.href"
            class="footer__link"
            :href="link.href"
            target="_blank"
            rel="noopener noreferrer"
          >
            {{ link.label }}
          </a>
          <RouterLink v-else class="footer__link" :to="link.to ?? '/'">
            {{ link.label }}
          </RouterLink>
        </template>
      </section>

      <section class="footer__card footer__copy" aria-labelledby="footer-copy">
        <h2 id="footer-copy" class="footer__copy-title">© 2026 Abberant</h2>
        <p class="footer__copy-text">
          Использование любых материалов с сайта без согласования с администрацией запрещено
        </p>
      </section>
    </div>
  </footer>
</template>

<style scoped>
.footer {
  width: min(1024px, 100%);
  padding: 8px;
  margin-inline: auto;
  background: var(--color-surface);
  border-radius: 8px 8px 0 0;
}

.footer__inner {
  display: grid;
  grid-template-areas: 'copy about social documents';
  grid-template-columns: 329.33fr 159.67fr 159.67fr 329.33fr;
  gap: 10px;
  width: min(1008px, 100%);
  margin-inline: auto;
}

.footer__card {
  display: flex;
  flex-direction: column;
  align-items: flex-start;
  gap: 4px;
  padding: 8px;
  color: var(--color-text);
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.footer__about {
  grid-area: about;
}

.footer__social {
  grid-area: social;
  border-radius: 4px;
}

.footer__documents {
  grid-area: documents;
  border-radius: 4px;
}

.footer__copy {
  grid-area: copy;
  gap: 10px;
  padding: 10px;
}

.footer__heading {
  font: 800 16px/19px Roboto, sans-serif;
}

.footer__link {
  font: 400 14px/16px Roboto, sans-serif;
}

.footer__link:hover {
  text-decoration: underline;
}

.footer__copy-title {
  font: 700 16px/19px Roboto, sans-serif;
}

.footer__copy-text {
  font: 300 14px/16px Roboto, sans-serif;
}

@media (max-width: 1023px) {
  .footer__inner {
    grid-template-areas:
      'about social'
      'documents documents'
      'copy copy';
    grid-template-columns: 1fr 1fr;
  }
}

@media (max-width: 767px) {
  .footer__inner {
    grid-template-areas:
      'about'
      'social'
      'documents'
      'copy';
    grid-template-columns: 1fr;
  }
}
</style>
