import type { Article, ArticleSection, Headline } from '@/types/article'

const image = '/img/news-placeholder.svg'

function article(
  id: string,
  title: string,
  comments: number,
  date = '2026-07-01',
): Article {
  return {
    id,
    title,
    image,
    date,
    comments,
    url: `/article/${id}`,
  }
}

export const headlines: Headline[] = [
  {
    id: 'fable-campaign',
    title: 'Прохождение кампании Fable займёт 15-20 часов',
    comments: 12,
    url: '/article/fable-campaign',
  },
  {
    id: 'witcher-3-gameplay',
    title: '12 минут геймплея «Баллад прошлого» для «Ведьмака 3»',
    comments: 24,
    url: '/article/witcher-3-gameplay',
  },
  {
    id: 'warlock-interview',
    title:
      'Основатель студии разработчиков Warlock рассказал, как игра сохраняет дух настольной D&D, и как найденные в тестированиях баги стали фичами',
    comments: 6,
    url: '/article/warlock-interview',
  },
  {
    id: 'dlss-5',
    title: 'Устанавливаем DLSS 5 на любую игру (почти)',
    comments: 32,
    url: '/article/dlss-5',
  },
]

export const sections: ArticleSection[] = [
  {
    id: 'news',
    title: 'Новости',
    href: '/news',
    articles: [
      article('news-1', 'Прохождение кампании Fable займёт 15-20 часов', 12, '2026-07-01'),
      article('news-2', '12 минут геймплея «Баллад прошлого» для «Ведьмака 3»', 24, '2026-06-30'),
      article(
        'news-3',
        'Основатель студии Warlock рассказал, как игра сохраняет дух настольной D&D',
        6,
        '2026-06-29',
      ),
      article('news-4', 'Устанавливаем DLSS 5 на любую игру (почти)', 32, '2026-06-28'),
      article('news-5', 'Анонсирована новая часть культовой RPG — первые подробности', 18, '2026-06-27'),
      article(
        'news-6',
        'Обновление 2.0 для Starfield: что изменилось и стоит ли возвращаться',
        41,
        '2026-06-26',
      ),
    ],
  },
  {
    id: 'reviews',
    title: 'Обзоры',
    href: '/reviews',
    articles: [
      article(
        'review-1',
        'Обзор Elden Ring: Shadow of the Erdtree — достойное продолжение легенды',
        54,
        '2026-06-25',
      ),
      article(
        'review-2',
        'Обзор Dragon’s Dogma 2: свобода, которая иногда пугает',
        37,
        '2026-06-22',
      ),
      article(
        'review-3',
        'Обзор Senua’s Saga: Hellblade II — кино, в которое можно играть',
        29,
        '2026-06-19',
      ),
    ],
  },
  {
    id: 'walkthroughs',
    title: 'Прохождения',
    href: '/walkthroughs',
    articles: [
      article('walkthrough-1', 'Как пройти все испытания в Ghost of Tsushima', 15, '2026-06-18'),
      article('walkthrough-2', 'Полное прохождение Baldur’s Gate 3: советы и секреты', 63, '2026-06-15'),
      article('walkthrough-3', 'Гайд по ресурсам в Fallout 4: где найти всё самое нужное', 21, '2026-06-12'),
    ],
  },
]
