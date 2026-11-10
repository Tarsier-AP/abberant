<script setup lang="ts">
interface HelpQuestion {
  id: string
  question: string
  answer: string
}

const benefits = [
  'Комментировать материалы и отвечать другим читателям.',
  'Сохранять статьи в закладки и возвращаться к ним позже.',
  'Получать уведомления об ответах и новых материалах.',
  'Оценивать обзоры и прохождения.',
  'Управлять профилем и историей активности.',
]

const questions: HelpQuestion[] = [
  {
    id: 'nickname',
    question: 'Как сменить никнейм?',
    answer:
      'Никнейм меняется в настройках профиля: откройте меню аккаунта, выберите «Профиль» и нажмите «Изменить никнейм». Новое имя должно быть уникальным и содержать от 3 до 20 символов.',
  },
  {
    id: 'delete-account',
    question: 'Как удалить аккаунт',
    answer:
      'Откройте настройки профиля и выберите «Удалить аккаунт». Подтвердите действие — профиль и связанные с ним данные будут удалены без возможности восстановления. Если аккаунт не удаётся удалить самостоятельно, напишите в поддержку.',
  },
  {
    id: 'blocking',
    question: 'За что пользователь может быть заблокирован?',
    answer:
      'За оскорбления и угрозы, спам и рекламу, публикацию незаконного контента, обход блокировок и другие нарушения правил сайта. Срок блокировки зависит от тяжести нарушения.',
  },
]
</script>

<template>
  <div class="help container">
    <header class="help__header">
      <h1 class="help__title">Помощь по сайту</h1>
      <p class="help__lead">Ответы на частые вопросы о регистрации и аккаунте на Abberant.</p>
    </header>

    <section class="help__card" aria-labelledby="help-registration">
      <h2 id="help-registration" class="help__heading">Зачем мне регистрироваться</h2>
      <p class="help__text">
        Регистрироваться на сайте не обязательно. Но регистрация даст вам следующие возможности:
      </p>
      <ol class="help__list">
        <li v-for="benefit in benefits" :key="benefit" class="help__item">{{ benefit }}</li>
      </ol>
    </section>

    <section
      v-for="item in questions"
      :key="item.id"
      class="help__card"
      :aria-labelledby="`help-${item.id}`"
    >
      <h2 :id="`help-${item.id}`" class="help__heading">{{ item.question }}</h2>
      <p class="help__text">{{ item.answer }}</p>
    </section>
  </div>
</template>

<style scoped>
.help {
  display: flex;
  flex-direction: column;
  gap: var(--gap);
  padding-block: var(--gap);
}

.help__header {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.help__title {
  font: 700 28px/34px Inter, sans-serif;
}

.help__lead {
  max-width: 760px;
  color: var(--color-muted);
  font: 400 15px/22px Roboto, sans-serif;
}

.help__card {
  display: flex;
  flex-direction: column;
  gap: 10px;
  padding: 16px;
  background: var(--color-card);
  border-radius: var(--radius);
  box-shadow: var(--shadow);
}

.help__heading {
  font: 700 21px/26px Roboto, sans-serif;
}

.help__text {
  max-width: 820px;
  font: 400 15px/22px Roboto, sans-serif;
}

.help__list {
  display: flex;
  flex-direction: column;
  gap: 8px;
  max-width: 820px;
  counter-reset: benefit;
}

.help__item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font: 400 15px/22px Roboto, sans-serif;
  counter-increment: benefit;
}

.help__item::before {
  content: counter(benefit);
  display: grid;
  flex: none;
  place-items: center;
  width: 22px;
  height: 22px;
  color: var(--color-bg);
  background: var(--color-text);
  border-radius: 50%;
  font: 700 12px/1 Roboto, sans-serif;
}

@media (max-width: 767px) {
  .help__title {
    font: 700 24px/30px Inter, sans-serif;
  }

  .help__card {
    padding: 14px;
  }

  .help__heading {
    font: 700 18px/22px Roboto, sans-serif;
  }
}
</style>
