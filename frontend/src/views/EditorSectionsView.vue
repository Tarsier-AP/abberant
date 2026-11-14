<script setup lang="ts">
import { computed, ref } from 'vue'
import LoadState from '@/components/ui/LoadState.vue'
import { reviewGroupSections } from '@/data/sections'
import type { Section } from '@/types/sections'
import type { SectionType } from '@/api/sections'
import { useSectionsStore } from '@/stores/sections'
import { useArticlesStore } from '@/stores/articles'

const sectionsStore = useSectionsStore()
const articlesStore = useArticlesStore()

const errorMessage = ref('')
const busyId = ref<string | null>(null)

const newsTitle = ref('')
const hubTitle = ref('')

/** Id of the section whose title is being edited right now. */
const editingId = ref<string | null>(null)
const editTitle = ref('')

const ready = computed(() => sectionsStore.loadState === 'ready')

const newsSections = computed(() => sectionsStore.sections.filter((s) => s.type === 'news'))
const hubSections = computed(() => sectionsStore.sections.filter((s) => s.type === 'walkthrough'))

/** Published material count, shown next to the archive action. */
function articleCount(id: string): number {
  return articlesStore.all.filter(
    (article) => article.section === id && article.status === 'published',
  ).length
}

function plural(count: number): string {
  const mod10 = count % 10
  const mod100 = count % 100

  if (mod10 === 1 && mod100 !== 11) return 'материал'
  if (mod10 >= 2 && mod10 <= 4 && (mod100 < 12 || mod100 > 14)) return 'материала'

  return 'материалов'
}

async function add(type: SectionType, title: string) {
  if (busyId.value) return

  errorMessage.value = ''

  const trimmed = title.trim()

  if (!trimmed) {
    errorMessage.value = 'Введите название секции.'
    return
  }

  busyId.value = `new-${type}`

  try {
    await sectionsStore.create(type, trimmed)
    if (type === 'news') newsTitle.value = ''
    else hubTitle.value = ''
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось создать секцию.'
  } finally {
    busyId.value = null
  }
}

function startRename(section: Section) {
  editingId.value = section.id
  editTitle.value = section.title
}

async function applyRename(section: Section) {
  if (busyId.value) return

  errorMessage.value = ''

  const trimmed = editTitle.value.trim()

  if (!trimmed) {
    errorMessage.value = 'Введите название секции.'
    return
  }

  busyId.value = section.id

  try {
    await sectionsStore.rename(section.id, trimmed)
    editingId.value = null
  } catch (error) {
    errorMessage.value = error instanceof Error ? error.message : 'Не удалось переименовать секцию.'
  } finally {
    busyId.value = null
  }
}

async function toggleArchived(section: Section) {
  if (busyId.value) return

  errorMessage.value = ''

  const archiving = !section.archived

  if (
    archiving &&
    articleCount(section.id) > 0 &&
    !window.confirm(
      `Отправить «${section.title}» в архив? Ссылки на ${articleCount(section.id)} ${plural(
        articleCount(section.id),
      )} продолжат работать, но раздел исчезнет из всех лент.`,
    )
  ) {
    return
  }

  busyId.value = section.id

  try {
    await sectionsStore.setArchived(section.id, archiving)
  } catch (error) {
    errorMessage.value =
      error instanceof Error ? error.message : 'Не удалось изменить состояние секции.'
  } finally {
    busyId.value = null
  }
}
</script>

<template>
  <div class="sections-manage container">
    <header class="sections-manage__header">
      <h1 class="sections-manage__heading">Секции</h1>
      <RouterLink class="sections-manage__back" :to="{ name: 'editor-articles' }">
        К статьям
      </RouterLink>
    </header>

    <p class="sections-manage__hint">
      Секции — новостные платформы и хабы прохождений. Удаление отправляет секцию в архив:
      статьи сохраняют адреса, но раздел пропадает из лент и редактора. Архив можно
      восстановить.
    </p>

    <p v-if="errorMessage" class="sections-manage__error" role="alert">{{ errorMessage }}</p>

    <LoadState
      v-if="!ready"
      class="sections-manage__state"
      :error="sectionsStore.errorMessage"
      @retry="sectionsStore.reload()"
    />

    <template v-else>
      <section class="sections-manage__block" aria-labelledby="sections-news-title">
        <h2 id="sections-news-title" class="sections-manage__title">Новостные платформы</h2>

        <form class="sections-manage__add" @submit.prevent="add('news', newsTitle)">
          <input
            v-model="newsTitle"
            class="sections-manage__input"
            type="text"
            placeholder="Название платформы, например «Мобилки»"
          />
          <button
            class="sections-manage__button sections-manage__button--primary"
            type="submit"
            :disabled="busyId !== null"
          >
            Добавить
          </button>
        </form>

        <ul class="sections-manage__list">
          <li v-for="section in newsSections" :key="section.id" class="sections-manage__item">
            <div class="sections-manage__info">
              <span class="sections-manage__name">
                {{ section.title }}
                <span v-if="section.archived" class="sections-manage__badge">в архиве</span>
              </span>
              <span class="sections-manage__meta">
                /news/{{ section.id }} · {{ articleCount(section.id) }}
                {{ plural(articleCount(section.id)) }}
              </span>
            </div>

            <div class="sections-manage__actions">
              <template v-if="editingId === section.id">
                <input
                  v-model="editTitle"
                  class="sections-manage__input sections-manage__input--inline"
                  type="text"
                  @keyup.enter="applyRename(section)"
                />
                <button
                  class="sections-manage__button sections-manage__button--primary"
                  type="button"
                  :disabled="busyId !== null"
                  @click="applyRename(section)"
                >
                  Сохранить
                </button>
                <button
                  class="sections-manage__button sections-manage__button--muted"
                  type="button"
                  @click="editingId = null"
                >
                  Отмена
                </button>
              </template>
              <template v-else>
                <button
                  class="sections-manage__button"
                  type="button"
                  :disabled="busyId !== null"
                  @click="startRename(section)"
                >
                  Переименовать
                </button>
                <button
                  v-if="!section.archived"
                  class="sections-manage__button sections-manage__button--danger"
                  type="button"
                  :disabled="busyId !== null"
                  @click="toggleArchived(section)"
                >
                  В архив
                </button>
                <button
                  v-else
                  class="sections-manage__button"
                  type="button"
                  :disabled="busyId !== null"
                  @click="toggleArchived(section)"
                >
                  Восстановить
                </button>
              </template>
            </div>
          </li>
        </ul>
      </section>

      <section class="sections-manage__block" aria-labelledby="sections-hubs-title">
        <h2 id="sections-hubs-title" class="sections-manage__title">Хабы прохождений</h2>

        <form class="sections-manage__add" @submit.prevent="add('walkthrough', hubTitle)">
          <input
            v-model="hubTitle"
            class="sections-manage__input"
            type="text"
            placeholder="Название игры, например «Hollow Knight: Silksong»"
          />
          <button
            class="sections-manage__button sections-manage__button--primary"
            type="submit"
            :disabled="busyId !== null"
          >
            Добавить
          </button>
        </form>

        <ul class="sections-manage__list">
          <li v-for="section in hubSections" :key="section.id" class="sections-manage__item">
            <div class="sections-manage__info">
              <span class="sections-manage__name">
                {{ section.title }}
                <span v-if="section.archived" class="sections-manage__badge">в архиве</span>
              </span>
              <span class="sections-manage__meta">
                /walkthroughs/{{ section.id }} · {{ articleCount(section.id) }}
                {{ plural(articleCount(section.id)) }}
              </span>
            </div>

            <div class="sections-manage__actions">
              <template v-if="editingId === section.id">
                <input
                  v-model="editTitle"
                  class="sections-manage__input sections-manage__input--inline"
                  type="text"
                  @keyup.enter="applyRename(section)"
                />
                <button
                  class="sections-manage__button sections-manage__button--primary"
                  type="button"
                  :disabled="busyId !== null"
                  @click="applyRename(section)"
                >
                  Сохранить
                </button>
                <button
                  class="sections-manage__button sections-manage__button--muted"
                  type="button"
                  @click="editingId = null"
                >
                  Отмена
                </button>
              </template>
              <template v-else>
                <button
                  class="sections-manage__button"
                  type="button"
                  :disabled="busyId !== null"
                  @click="startRename(section)"
                >
                  Переименовать
                </button>
                <button
                  v-if="!section.archived"
                  class="sections-manage__button sections-manage__button--danger"
                  type="button"
                  :disabled="busyId !== null"
                  @click="toggleArchived(section)"
                >
                  В архив
                </button>
                <button
                  v-else
                  class="sections-manage__button"
                  type="button"
                  :disabled="busyId !== null"
                  @click="toggleArchived(section)"
                >
                  Восстановить
                </button>
              </template>
            </div>
          </li>
        </ul>
      </section>

      <section class="sections-manage__block" aria-labelledby="sections-groups-title">
        <h2 id="sections-groups-title" class="sections-manage__title">Группы обзоров</h2>

        <p class="sections-manage__hint sections-manage__hint--tight">
          Группы обзоров управляются сайтом и не редактируются.
        </p>

        <ul class="sections-manage__list">
          <li v-for="group in reviewGroupSections" :key="group.id" class="sections-manage__item">
            <div class="sections-manage__info">
              <span class="sections-manage__name">{{ group.title }}</span>
              <span class="sections-manage__meta">{{ group.href }}</span>
            </div>
          </li>
        </ul>
      </section>
    </template>
  </div>
</template>

<style scoped>
.sections-manage {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  max-width: 1060px;
  padding-block: var(--gap);
  margin-inline: auto;
}

.sections-manage__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 10px;
  padding: 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.sections-manage__heading {
  font: 600 24px/30px Roboto, sans-serif;
}

.sections-manage__back {
  display: inline-flex;
  align-items: center;
  min-height: 36px;
  padding: 6px 14px;
  color: var(--color-text);
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
  font: 500 13px/16px Roboto, sans-serif;
}

.sections-manage__back:hover {
  background: var(--color-card-hover);
}

.sections-manage__hint {
  padding: 8px 12px;
  color: var(--color-muted);
  background: var(--color-surface);
  border-radius: var(--radius);
  font: 400 13px/18px Roboto, sans-serif;
}

.sections-manage__hint--tight {
  padding: 0;
  background: none;
}

.sections-manage__error {
  padding: 8px;
  color: var(--color-danger);
  font: 400 12px/16px Roboto, sans-serif;
}

.sections-manage__state {
  width: 100%;
}

.sections-manage__block {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 12px;
  background: var(--color-surface);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.sections-manage__title {
  font: 600 16px/20px Roboto, sans-serif;
}

.sections-manage__add {
  display: flex;
  gap: 6px;
}

.sections-manage__input {
  flex: 1;
  min-height: 34px;
  width: 100%;
  padding: 4px 8px;
  color: var(--color-text);
  background: var(--color-auth-input);
  border: 1px solid var(--color-input-border);
  border-radius: var(--radius);
  font: 400 14px/18px Roboto, sans-serif;
}

.sections-manage__input:focus-visible {
  outline: 2px solid var(--color-accent);
  outline-offset: 0;
}

.sections-manage__input--inline {
  width: auto;
  min-width: 160px;
}

.sections-manage__list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin: 0;
  padding: 0;
  list-style: none;
}

.sections-manage__item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 12px;
  padding: 8px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.sections-manage__info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.sections-manage__name {
  color: var(--color-text);
  font: 500 14px/18px Roboto, sans-serif;
}

.sections-manage__badge {
  margin-left: 6px;
  padding: 1px 6px;
  color: var(--color-warning);
  border: 1px solid var(--color-warning);
  border-radius: var(--radius);
  font: 500 11px/14px Roboto, sans-serif;
}

.sections-manage__meta {
  color: var(--color-muted);
  font: 400 12px/16px Roboto, sans-serif;
}

.sections-manage__actions {
  display: flex;
  flex: none;
  flex-wrap: wrap;
  justify-content: flex-end;
  align-items: center;
  gap: 6px;
}

.sections-manage__button {
  display: inline-flex;
  align-items: center;
  min-height: 32px;
  padding: 4px 12px;
  color: var(--color-text);
  background: var(--color-surface);
  border: 0;
  border-radius: var(--radius);
  font: 500 12px/16px Roboto, sans-serif;
  cursor: pointer;
}

.sections-manage__button:hover {
  background: var(--color-card-hover);
}

.sections-manage__button:disabled {
  opacity: 0.6;
  cursor: default;
}

.sections-manage__button--primary {
  color: var(--color-accent-text);
  background: var(--color-accent);
}

.sections-manage__button--primary:hover {
  background: var(--color-accent);
  opacity: 0.9;
}

.sections-manage__button--danger {
  color: var(--color-danger);
}

.sections-manage__button--muted {
  color: var(--color-muted);
}

@media (max-width: 767px) {
  .sections-manage__item {
    flex-direction: column;
    align-items: stretch;
  }

  .sections-manage__actions {
    justify-content: flex-start;
  }
}
</style>
