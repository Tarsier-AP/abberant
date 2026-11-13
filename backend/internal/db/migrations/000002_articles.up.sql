CREATE TABLE articles (
    id         text PRIMARY KEY,
    title      text      NOT NULL,
    image      text      NOT NULL,
    date       date      NOT NULL,
    comments   int       NOT NULL DEFAULT 0,
    kind       text      NOT NULL CHECK (kind IN ('news', 'review')),
    content    text      NOT NULL DEFAULT '',
    author     text      NOT NULL,
    status     text      NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'published')),
    rating     smallint  CHECK (rating BETWEEN 1 AND 5),
    platform   text,
    verdict    text,
    section    text,
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX articles_date_idx ON articles (date DESC);

INSERT INTO articles (id, title, image, date, comments, kind, content, author, status, rating, platform, verdict, section) VALUES
    ('news-pc-1', 'Steam Next Fest: 12 демо, которые стоит попробовать в июле', '/img/news-placeholder.svg', '2026-07-01', 14, 'news', '## Итоги презентации

![Демоверсии Steam Next Fest](/img/news-placeholder.svg "Одна из демоверсий фестиваля в действии")

Фестиваль собрал больше сотни бесплатных демоверсий. Мы отобрали двенадцать проектов, которые стоит попробовать до конца недели.

### На что обратить внимание

- Тактические стратегии с необычной системой приказов
- Сюжетные адвенчуры от независимых студий
- Кооперативные выживалки с процедурными картами

Полный список демоверсий опубликован [в новостях Steam](https://store.steampowered.com).
', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'pc'),
    ('news-pc-2', 'Анонсирован ремастер культовой стратегии 2000-х', '/img/news-placeholder.svg', '2026-06-28', 27, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'pc'),
    ('news-pc-3', 'Новый патч для Baldur’s Gate 3 исправил сохранения на Steam Deck', '/img/news-placeholder.svg', '2026-06-24', 9, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'pc'),
    ('news-pc-4', 'Valve обновила правила модерации Steam Workshop', '/img/news-placeholder.svg', '2026-06-13', 13, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'pc'),
    ('news-pc-5', 'Вышло крупное обновление для Cities: Skylines II', '/img/news-placeholder.svg', '2026-06-08', 22, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'pc'),
    ('news-pc-6', 'Microsoft прекращает поддержку Windows 10: что делать геймерам', '/img/news-placeholder.svg', '2026-06-03', 30, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'pc'),
    ('news-xbox-1', 'Game Pass пополнится шестью играми в июле', '/img/news-placeholder.svg', '2026-06-30', 31, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'xbox'),
    ('news-xbox-2', 'Obsidian показала новый геймплей Avowed 2', '/img/news-placeholder.svg', '2026-06-26', 18, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'xbox'),
    ('news-xbox-3', 'Xbox Series X|S получит поддержку новых видеокодеков', '/img/news-placeholder.svg', '2026-06-21', 7, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'xbox'),
    ('news-xbox-4', 'Sea of Thieves получила сезон с новыми островами', '/img/news-placeholder.svg', '2026-06-13', 12, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'xbox'),
    ('news-xbox-5', 'Xbox Cloud Gaming выходит на новые устройства', '/img/news-placeholder.svg', '2026-06-09', 15, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'xbox'),
    ('news-xbox-6', 'Halo: Infinite — серверы перевели на новую инфраструктуру', '/img/news-placeholder.svg', '2026-06-04', 20, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'xbox'),
    ('news-playstation-1', 'State of Play: главные анонсы и даты релизов', '/img/news-placeholder.svg', '2026-06-29', 42, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'playstation'),
    ('news-playstation-2', 'Ghost of Yotei получит бесплатное сюжетное дополнение осенью', '/img/news-placeholder.svg', '2026-06-25', 36, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'playstation'),
    ('news-playstation-3', 'PlayStation 5 Pro: первые тесты обновлённого PSSR', '/img/news-placeholder.svg', '2026-06-20', 23, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'playstation'),
    ('news-playstation-4', 'God of War: новая глава получила возрастной рейтинг', '/img/news-placeholder.svg', '2026-06-12', 28, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'playstation'),
    ('news-playstation-5', 'PS Plus: июльская подборка игр', '/img/news-placeholder.svg', '2026-06-07', 34, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'playstation'),
    ('news-playstation-6', 'Astro Bot получит бесплатное дополнение с уровнями', '/img/news-placeholder.svg', '2026-06-02', 17, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'playstation'),
    ('news-nintendo-1', 'Nintendo Direct: всё, что показали на презентации', '/img/news-placeholder.svg', '2026-06-27', 39, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'nintendo'),
    ('news-nintendo-2', 'Новая часть Mario Kart получила дату релиза', '/img/news-placeholder.svg', '2026-06-23', 25, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'nintendo'),
    ('news-nintendo-3', 'Switch 2: подробности обратной совместимости с играми', '/img/news-placeholder.svg', '2026-06-19', 16, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'nintendo'),
    ('news-nintendo-4', 'Pokémon Legends Z-A: новые подробности боевой системы', '/img/news-placeholder.svg', '2026-06-12', 26, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'nintendo'),
    ('news-nintendo-5', 'Switch 2: слух о цене и комплектации', '/img/news-placeholder.svg', '2026-06-06', 44, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'nintendo'),
    ('news-nintendo-6', 'The Legend of Zelda: классика получила обновление', '/img/news-placeholder.svg', '2026-05-30', 21, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'nintendo'),
    ('news-vr-1', 'Meta Quest 3: обновление улучшает трекинг рук', '/img/news-placeholder.svg', '2026-06-22', 11, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'vr'),
    ('news-vr-2', 'Half-Life: Alyx получила неофициальный кооперативный мод', '/img/news-placeholder.svg', '2026-06-18', 33, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'vr'),
    ('news-vr-3', 'Sony представила новый VR-контроллер для PS5', '/img/news-placeholder.svg', '2026-06-15', 8, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'vr'),
    ('news-vr-4', 'Появился мод с поддержкой VR для Cyberpunk 2077', '/img/news-placeholder.svg', '2026-06-11', 38, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'vr'),
    ('news-vr-5', 'Meta анонсировала Quest 4 с новыми линзами', '/img/news-placeholder.svg', '2026-06-05', 29, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'vr'),
    ('news-vr-6', 'Лучшие VR-игры лета 2026 года', '/img/news-placeholder.svg', '2026-05-28', 13, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'vr'),
    ('news-hardware-1', 'Обзор RTX 5080: стоит ли апгрейдиться в 2026 году', '/img/news-placeholder.svg', '2026-06-26', 47, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'hardware'),
    ('news-hardware-2', 'AMD снижает цены на процессоры Ryzen', '/img/news-placeholder.svg', '2026-06-17', 19, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'hardware'),
    ('news-hardware-3', 'Тесты: SSD PCIe 5.0 в играх почти не даёт прироста', '/img/news-placeholder.svg', '2026-06-14', 12, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'hardware'),
    ('news-hardware-4', 'Вышли драйверы NVIDIA с оптимизацией для новых игр', '/img/news-placeholder.svg', '2026-06-10', 16, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'hardware'),
    ('news-hardware-5', 'Обзор корпуса с необычной системой охлаждения', '/img/news-placeholder.svg', '2026-06-04', 9, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'hardware'),
    ('news-hardware-6', 'DDR5 продолжает дешеветь: стоит ли брать сейчас', '/img/news-placeholder.svg', '2026-05-29', 24, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'hardware'),
    ('review-1', 'Обзор Elden Ring: Shadow of the Erdtree — достойное продолжение легенды', '/img/news-placeholder.svg', '2026-06-25', 54, 'review', '## Коротко о главном

![Кадр из игры Elden Ring: Shadow of the Erdtree](/img/news-placeholder.svg "Кадр из игры: Elden Ring: Shadow of the Erdtree")

Междуземье снова зовёт: дополнение предлагает десятки часов контента, новый регион и одни из самых запоминающихся боссов в серии.

### Что понравилось

- Атмосфера, от которой не хочется отрываться
- Продуманный боевой и прогрессионный системы
- Отличное музыкальное сопровождение

### Что не понравилось

- Спорные решения в балансе поздних этапов
- Технические недочёты на старте продаж

> Elden Ring: Shadow of the Erdtree — тот случай, когда продолжение оправдывает ожидания.

Подробнее о рейтинге игр можно прочитать в разделе [«Помощь по сайту»](/help).
', 'Редакция Abberant', 'published', 5, 'PC, PlayStation 5, Xbox Series X|S', 'Обязательное дополнение для всех, кому понравилась базовая игра.', NULL),
    ('review-2', 'Обзор Dragon’s Dogma 2: свобода, которая иногда пугает', '/img/news-placeholder.svg', '2026-06-22', 37, 'review', '## Коротко о главном

![Кадр из игры Dragon’s Dogma 2](/img/news-placeholder.svg "Кадр из игры: Dragon’s Dogma 2")

Капком не стал чинить то, что не ломалось: хаотичные путешествия, непредсказуемые встречи и пешки, которых реально хочется хвалить.

### Что понравилось

- Атмосфера, от которой не хочется отрываться
- Продуманный боевой и прогрессионный системы
- Отличное музыкальное сопровождение

### Что не понравилось

- Спорные решения в балансе поздних этапов
- Технические недочёты на старте продаж

> Dragon’s Dogma 2 — тот случай, когда продолжение оправдывает ожидания.

Подробнее о рейтинге игр можно прочитать в разделе [«Помощь по сайту»](/help).
', 'Редакция Abberant', 'published', 4, 'PC, PlayStation 5, Xbox Series X|S', 'Смелая и живая RPG, которую прощают её шероховатости.', NULL),
    ('review-3', 'Обзор Senua’s Saga: Hellblade II — кино, в которое можно играть', '/img/news-placeholder.svg', '2026-06-19', 29, 'review', '', 'Редакция Abberant', 'published', 4, 'PC, Xbox Series X|S', 'Визуальное великолепие и сильная история ценой короткой кампании.', NULL),
    ('recent-1', 'Прохождение Control Resonant: как открыть все способности', '/img/news-placeholder.svg', '2026-07-01', 18, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'recent'),
    ('recent-2', 'Silent Hill: Townfall — полное прохождение без смертей', '/img/news-placeholder.svg', '2026-06-30', 27, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'recent'),
    ('recent-3', 'Ghost of Tsushima: все испытания и награды', '/img/news-placeholder.svg', '2026-06-29', 22, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'recent'),
    ('control-resonant-1', 'Control Resonant: полное прохождение кампании', '/img/news-placeholder.svg', '2026-06-28', 33, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'control-resonant'),
    ('control-resonant-2', 'Все способности Джесси в Control Resonant', '/img/news-placeholder.svg', '2026-06-25', 21, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'control-resonant'),
    ('control-resonant-3', 'Control Resonant: где найти каждое секретное оружие', '/img/news-placeholder.svg', '2026-06-22', 17, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'control-resonant'),
    ('control-resonant-4', 'Боссы Control Resonant: тактика и слабые места', '/img/news-placeholder.svg', '2026-06-19', 29, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'control-resonant'),
    ('control-resonant-5', 'Control Resonant: все побочные задания и награды', '/img/news-placeholder.svg', '2026-06-16', 12, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'control-resonant'),
    ('control-resonant-6', 'Как открыть истинную концовку в Control Resonant', '/img/news-placeholder.svg', '2026-06-13', 46, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'control-resonant'),
    ('control-resonant-7', 'Control Resonant: гайд по ресурсам и прокачке', '/img/news-placeholder.svg', '2026-06-10', 15, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'control-resonant'),
    ('control-resonant-8', 'Секретные комнаты Control Resonant: полный список', '/img/news-placeholder.svg', '2026-06-07', 24, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'control-resonant'),
    ('control-resonant-9', 'Control Resonant: все достижения и трофеи', '/img/news-placeholder.svg', '2026-06-04', 8, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'control-resonant'),
    ('silent-hill-townfall-1', 'Silent Hill: Townfall — полное прохождение', '/img/news-placeholder.svg', '2026-06-27', 52, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'silent-hill-townfall'),
    ('silent-hill-townfall-2', 'Все концовки Silent Hill: Townfall и условия их получения', '/img/news-placeholder.svg', '2026-06-24', 38, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'silent-hill-townfall'),
    ('silent-hill-townfall-3', 'Silent Hill: Townfall: где найти все записки', '/img/news-placeholder.svg', '2026-06-21', 19, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'silent-hill-townfall'),
    ('silent-hill-townfall-4', 'Головоломки Silent Hill: Townfall — все решения', '/img/news-placeholder.svg', '2026-06-18', 27, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'silent-hill-townfall'),
    ('silent-hill-townfall-5', 'Silent Hill: Townfall: как выжить в тумане', '/img/news-placeholder.svg', '2026-06-15', 14, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'silent-hill-townfall'),
    ('silent-hill-townfall-6', 'Silent Hill: Townfall: карта города и тайники', '/img/news-placeholder.svg', '2026-06-12', 11, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'silent-hill-townfall'),
    ('silent-hill-townfall-7', 'Все секреты Silent Hill: Townfall', '/img/news-placeholder.svg', '2026-06-09', 31, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'silent-hill-townfall'),
    ('silent-hill-townfall-8', 'Silent Hill: Townfall: как одолеть финального босса', '/img/news-placeholder.svg', '2026-06-06', 23, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'silent-hill-townfall'),
    ('silent-hill-townfall-9', 'Silent Hill: Townfall: все достижения', '/img/news-placeholder.svg', '2026-06-03', 7, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'silent-hill-townfall'),
    ('ghost-of-tsushima-1', 'Ghost of Tsushima: прохождение всех актов', '/img/news-placeholder.svg', '2026-06-26', 40, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'ghost-of-tsushima'),
    ('ghost-of-tsushima-2', 'Все испытания Ghost of Tsushima: карты и награды', '/img/news-placeholder.svg', '2026-06-23', 35, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'ghost-of-tsushima'),
    ('ghost-of-tsushima-3', 'Ghost of Tsushima: где найти всю экипировку', '/img/news-placeholder.svg', '2026-06-20', 20, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'ghost-of-tsushima'),
    ('ghost-of-tsushima-4', 'Ghost of Tsushima: стойки и приёмы — полный гайд', '/img/news-placeholder.svg', '2026-06-17', 25, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'ghost-of-tsushima'),
    ('ghost-of-tsushima-5', 'Ghost of Tsushima: все святилища и лисицы', '/img/news-placeholder.svg', '2026-06-14', 16, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'ghost-of-tsushima'),
    ('ghost-of-tsushima-6', 'Ghost of Tsushima: как открыть все способности', '/img/news-placeholder.svg', '2026-06-11', 13, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'ghost-of-tsushima'),
    ('ghost-of-tsushima-7', 'Ghost of Tsushima: секретные шляпы и маски', '/img/news-placeholder.svg', '2026-06-08', 9, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'ghost-of-tsushima'),
    ('ghost-of-tsushima-8', 'Ghost of Tsushima: все дуэли и тактика', '/img/news-placeholder.svg', '2026-06-05', 28, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'ghost-of-tsushima'),
    ('ghost-of-tsushima-9', 'Ghost of Tsushima: достижения и трофеи', '/img/news-placeholder.svg', '2026-06-02', 10, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'ghost-of-tsushima'),
    ('baldurs-gate-3-1', 'Baldur’s Gate 3: полное прохождение всех актов', '/img/news-placeholder.svg', '2026-06-23', 95, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'baldurs-gate-3'),
    ('baldurs-gate-3-2', 'Baldur’s Gate 3: все концовки и их условия', '/img/news-placeholder.svg', '2026-06-20', 71, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'baldurs-gate-3'),
    ('baldurs-gate-3-3', 'Baldur’s Gate 3: лучшие билды для каждого класса', '/img/news-placeholder.svg', '2026-06-17', 64, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'baldurs-gate-3'),
    ('baldurs-gate-3-4', 'Baldur’s Gate 3: где найти всех спутников', '/img/news-placeholder.svg', '2026-06-14', 39, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'baldurs-gate-3'),
    ('baldurs-gate-3-5', 'Baldur’s Gate 3: все секреты первого акта', '/img/news-placeholder.svg', '2026-06-11', 30, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'baldurs-gate-3'),
    ('baldurs-gate-3-6', 'Baldur’s Gate 3: как пройти бой с Кетериком', '/img/news-placeholder.svg', '2026-06-08', 22, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'baldurs-gate-3'),
    ('baldurs-gate-3-7', 'Baldur’s Gate 3: гайд по крафту и алхимии', '/img/news-placeholder.svg', '2026-06-05', 18, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'baldurs-gate-3'),
    ('baldurs-gate-3-8', 'Baldur’s Gate 3: скрытые квесты и награды', '/img/news-placeholder.svg', '2026-06-02', 26, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'baldurs-gate-3'),
    ('baldurs-gate-3-9', 'Baldur’s Gate 3: все достижения', '/img/news-placeholder.svg', '2026-05-30', 12, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, 'baldurs-gate-3'),
    ('walkthrough-1', 'Как пройти все испытания в Ghost of Tsushima', '/img/news-placeholder.svg', '2026-06-18', 15, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, NULL),
    ('walkthrough-2', 'Полное прохождение Baldur’s Gate 3: советы и секреты', '/img/news-placeholder.svg', '2026-06-15', 63, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, NULL),
    ('walkthrough-3', 'Гайд по ресурсам в Fallout 4: где найти всё самое нужное', '/img/news-placeholder.svg', '2026-06-12', 21, 'news', '', 'Редакция Abberant', 'published', NULL, NULL, NULL, NULL);
