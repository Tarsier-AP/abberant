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

const reviews: Article[] = [
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
]

export { reviews }

/** Reviews feed groups; every group has its own archive page. */
export const reviewGroups: ArticleSection[] = [
  {
    id: 'recent',
    title: 'Недавние',
    href: '/reviews/recent',
    articles: [...reviews].sort((a, b) => b.date.localeCompare(a.date)),
  },
  {
    id: 'popular',
    title: 'Популярные',
    href: '/reviews/popular',
    articles: [...reviews].sort((a, b) => b.comments - a.comments),
  },
]

function newsSection(id: string, title: string, entries: ArticleEntry[]): ArticleSection {
  return {
    id,
    title,
    href: `/news/${id}`,
    articles: entries.map(([title, comments, date], index) =>
      article(`news-${id}-${index + 1}`, title, comments, date),
    ),
  }
}

/**
 * Platform feeds for `/news`. The `all` section is assembled from every
 * platform below, newest first, and is also reused by the home page.
 */
const newsPlatforms: ArticleSection[] = [
  newsSection('pc', 'PC', [
    ['Steam Next Fest: 12 демо, которые стоит попробовать в июле', 14, '2026-07-01'],
    ['Анонсирован ремастер культовой стратегии 2000-х', 27, '2026-06-28'],
    ['Новый патч для Baldur’s Gate 3 исправил сохранения на Steam Deck', 9, '2026-06-24'],
    ['Valve обновила правила модерации Steam Workshop', 13, '2026-06-13'],
    ['Вышло крупное обновление для Cities: Skylines II', 22, '2026-06-08'],
    ['Microsoft прекращает поддержку Windows 10: что делать геймерам', 30, '2026-06-03'],
  ]),
  newsSection('xbox', 'Xbox Series X|S', [
    ['Game Pass пополнится шестью играми в июле', 31, '2026-06-30'],
    ['Obsidian показала новый геймплей Avowed 2', 18, '2026-06-26'],
    ['Xbox Series X|S получит поддержку новых видеокодеков', 7, '2026-06-21'],
    ['Sea of Thieves получила сезон с новыми островами', 12, '2026-06-13'],
    ['Xbox Cloud Gaming выходит на новые устройства', 15, '2026-06-09'],
    ['Halo: Infinite — серверы перевели на новую инфраструктуру', 20, '2026-06-04'],
  ]),
  newsSection('playstation', 'PlayStation 5', [
    ['State of Play: главные анонсы и даты релизов', 42, '2026-06-29'],
    ['Ghost of Yotei получит бесплатное сюжетное дополнение осенью', 36, '2026-06-25'],
    ['PlayStation 5 Pro: первые тесты обновлённого PSSR', 23, '2026-06-20'],
    ['God of War: новая глава получила возрастной рейтинг', 28, '2026-06-12'],
    ['PS Plus: июльская подборка игр', 34, '2026-06-07'],
    ['Astro Bot получит бесплатное дополнение с уровнями', 17, '2026-06-02'],
  ]),
  newsSection('nintendo', 'Nintendo', [
    ['Nintendo Direct: всё, что показали на презентации', 39, '2026-06-27'],
    ['Новая часть Mario Kart получила дату релиза', 25, '2026-06-23'],
    ['Switch 2: подробности обратной совместимости с играми', 16, '2026-06-19'],
    ['Pokémon Legends Z-A: новые подробности боевой системы', 26, '2026-06-12'],
    ['Switch 2: слух о цене и комплектации', 44, '2026-06-06'],
    ['The Legend of Zelda: классика получила обновление', 21, '2026-05-30'],
  ]),
  newsSection('vr', 'VR', [
    ['Meta Quest 3: обновление улучшает трекинг рук', 11, '2026-06-22'],
    ['Half-Life: Alyx получила неофициальный кооперативный мод', 33, '2026-06-18'],
    ['Sony представила новый VR-контроллер для PS5', 8, '2026-06-15'],
    ['Появился мод с поддержкой VR для Cyberpunk 2077', 38, '2026-06-11'],
    ['Meta анонсировала Quest 4 с новыми линзами', 29, '2026-06-05'],
    ['Лучшие VR-игры лета 2026 года', 13, '2026-05-28'],
  ]),
  newsSection('hardware', 'Железо', [
    ['Обзор RTX 5080: стоит ли апгрейдиться в 2026 году', 47, '2026-06-26'],
    ['AMD снижает цены на процессоры Ryzen', 19, '2026-06-17'],
    ['Тесты: SSD PCIe 5.0 в играх почти не даёт прироста', 12, '2026-06-14'],
    ['Вышли драйверы NVIDIA с оптимизацией для новых игр', 16, '2026-06-10'],
    ['Обзор корпуса с необычной системой охлаждения', 9, '2026-06-04'],
    ['DDR5 продолжает дешеветь: стоит ли брать сейчас', 24, '2026-05-29'],
  ]),
]

export const latestNews: Article[] = newsPlatforms
  .flatMap((section) => section.articles)
  .sort((a, b) => b.date.localeCompare(a.date))

export const newsSections: ArticleSection[] = [
  { id: 'all', title: 'Все', href: '/news/all', articles: latestNews },
  ...newsPlatforms,
]

export const sections: ArticleSection[] = [
  {
    id: 'news',
    title: 'Новости',
    href: '/news',
    articles: latestNews.slice(0, 6),
  },
  {
    id: 'reviews',
    title: 'Обзоры',
    href: '/reviews',
    articles: reviews,
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
