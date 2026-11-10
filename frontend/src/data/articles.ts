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

type ArticleEntry = [title: string, comments: number, date: string]

function archiveSection(
  id: string,
  sectionTitle: string,
  entries: ArticleEntry[],
): ArticleSection {
  return {
    id,
    title: sectionTitle,
    href: `/walkthroughs/${id}`,
    articles: entries.map(([title, comments, date], index) =>
      article(`${id}-${index + 1}`, title, comments, date),
    ),
  }
}

/**
 * Archive feed for `/walkthroughs`. Every section is a game hub; the first one
 * collects the latest walkthroughs. `href` points to the game's own page.
 */
export const walkthroughs: ArticleSection[] = [
  archiveSection('recent', 'Недавние', [
    ['Прохождение Control Resonant: как открыть все способности', 18, '2026-07-01'],
    ['Silent Hill: Townfall — полное прохождение без смертей', 27, '2026-06-30'],
    ['Ghost of Tsushima: все испытания и награды', 22, '2026-06-29'],
  ]),
  archiveSection('control-resonant', 'CONTROL RESONANT', [
    ['Control Resonant: полное прохождение кампании', 33, '2026-06-28'],
    ['Все способности Джесси в Control Resonant', 21, '2026-06-25'],
    ['Control Resonant: где найти каждое секретное оружие', 17, '2026-06-22'],
    ['Боссы Control Resonant: тактика и слабые места', 29, '2026-06-19'],
    ['Control Resonant: все побочные задания и награды', 12, '2026-06-16'],
    ['Как открыть истинную концовку в Control Resonant', 46, '2026-06-13'],
    ['Control Resonant: гайд по ресурсам и прокачке', 15, '2026-06-10'],
    ['Секретные комнаты Control Resonant: полный список', 24, '2026-06-07'],
    ['Control Resonant: все достижения и трофеи', 8, '2026-06-04'],
  ]),
  archiveSection('silent-hill-townfall', 'SILENT HILL: TOWNFALL', [
    ['Silent Hill: Townfall — полное прохождение', 52, '2026-06-27'],
    ['Все концовки Silent Hill: Townfall и условия их получения', 38, '2026-06-24'],
    ['Silent Hill: Townfall: где найти все записки', 19, '2026-06-21'],
    ['Головоломки Silent Hill: Townfall — все решения', 27, '2026-06-18'],
    ['Silent Hill: Townfall: как выжить в тумане', 14, '2026-06-15'],
    ['Silent Hill: Townfall: карта города и тайники', 11, '2026-06-12'],
    ['Все секреты Silent Hill: Townfall', 31, '2026-06-09'],
    ['Silent Hill: Townfall: как одолеть финального босса', 23, '2026-06-06'],
    ['Silent Hill: Townfall: все достижения', 7, '2026-06-03'],
  ]),
  archiveSection('ghost-of-tsushima', 'GHOST OF TSUSHIMA', [
    ['Ghost of Tsushima: прохождение всех актов', 40, '2026-06-26'],
    ['Все испытания Ghost of Tsushima: карты и награды', 35, '2026-06-23'],
    ['Ghost of Tsushima: где найти всю экипировку', 20, '2026-06-20'],
    ['Ghost of Tsushima: стойки и приёмы — полный гайд', 25, '2026-06-17'],
    ['Ghost of Tsushima: все святилища и лисицы', 16, '2026-06-14'],
    ['Ghost of Tsushima: как открыть все способности', 13, '2026-06-11'],
    ['Ghost of Tsushima: секретные шляпы и маски', 9, '2026-06-08'],
    ['Ghost of Tsushima: все дуэли и тактика', 28, '2026-06-05'],
    ['Ghost of Tsushima: достижения и трофеи', 10, '2026-06-02'],
  ]),
  archiveSection('baldurs-gate-3', 'BALDUR’S GATE 3', [
    ['Baldur’s Gate 3: полное прохождение всех актов', 95, '2026-06-23'],
    ['Baldur’s Gate 3: все концовки и их условия', 71, '2026-06-20'],
    ['Baldur’s Gate 3: лучшие билды для каждого класса', 64, '2026-06-17'],
    ['Baldur’s Gate 3: где найти всех спутников', 39, '2026-06-14'],
    ['Baldur’s Gate 3: все секреты первого акта', 30, '2026-06-11'],
    ['Baldur’s Gate 3: как пройти бой с Кетериком', 22, '2026-06-08'],
    ['Baldur’s Gate 3: гайд по крафту и алхимии', 18, '2026-06-05'],
    ['Baldur’s Gate 3: скрытые квесты и награды', 26, '2026-06-02'],
    ['Baldur’s Gate 3: все достижения', 12, '2026-05-30'],
  ]),
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
        'Основатель студии разработчиков Warlock рассказал, как игра сохраняет дух настольной D&D',
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
