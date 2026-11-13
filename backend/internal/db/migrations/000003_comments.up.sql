CREATE TABLE comments (
    id         bigint GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    article_id text   NOT NULL REFERENCES articles (id) ON DELETE CASCADE,
    user_id    bigint NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    parent_id  bigint REFERENCES comments (id) ON DELETE CASCADE,
    body       text   NOT NULL CHECK (length(trim(body)) BETWEEN 1 AND 500),
    created_at timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX comments_article_created_idx ON comments (article_id, created_at, id);
CREATE INDEX comments_parent_idx ON comments (parent_id);

-- Example accounts. They are deactivated (no password hash), so nobody can log
-- in as them: they exist to make the comment sections look alive.
INSERT INTO users (email, username, activated, role, avatar_data_url) VALUES
    ('ivan.sokolov@example.abberant', 'Иван Соколов', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%23a53b6e'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%98</text></svg>$$),
    ('alina.vetrova@example.abberant', 'Алина Ветрова', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%233b6ea5'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%90</text></svg>$$),
    ('demid.krylov@example.abberant', 'Демид Крылов', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%236ea53b'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%94</text></svg>$$),
    ('marina.zaytseva@example.abberant', 'Марина Зайцева', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%23a56e3b'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%9C</text></svg>$$),
    ('oleg.tarasov@example.abberant', 'Олег Тарасов', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%236e3ba5'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%9E</text></svg>$$),
    ('polina.morozova@example.abberant', 'Полина Морозова', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%233ba56e'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%9F</text></svg>$$),
    ('artem.gusev@example.abberant', 'Артём Гусев', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%23a53b3b'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%90</text></svg>$$),
    ('vera.lebedeva@example.abberant', 'Вера Лебедева', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%233ba5a5'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%92</text></svg>$$),
    ('kirill.ezhov@example.abberant', 'Кирилл Ежов', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%23843ba5'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%9A</text></svg>$$),
    ('sonya.romanova@example.abberant', 'Соня Романова', false, 'reader',
     $$data:image/svg+xml;utf8,<svg%20xmlns='http://www.w3.org/2000/svg'%20width='160'%20height='160'><rect%20width='160'%20height='160'%20fill='%23a5843b'/><text%20x='80'%20y='108'%20font-family='Arial,sans-serif'%20font-size='80'%20font-weight='600'%20fill='%23ffffff'%20text-anchor='middle'>%D0%A1</text></svg>$$);

-- One example comment on every published article, rotating through the demo
-- accounts and a pool of generic reader reactions.
DO $seed$
DECLARE
    target       record;
    author_ids   bigint[];
    author_count int;
    texts        text[] := ARRAY[
        'Отличный материал, спасибо! Ждём продолжения.',
        'Не согласен по паре пунктов, но в целом написано хорошо.',
        'Как раз то, что искал. Закладку поставил.',
        'Прочитал на одном дыхании, у автора отлично получается.',
        'Аналогичная ситуация у меня, подтверждаю слова из статьи.',
        'Спасибо за подробный разбор, всё по делу.',
        'Жаль, что об этом мало где пишут. Респект редакции.',
        'Интересная точка зрения, но есть нюансы, о которых стоит упомянуть.',
        'Добавил в избранное, буду советовать друзьям.',
        'Полезно, даже для тех, кто далёк от темы.',
        'Жду обновлений по этой теме, надеюсь, будет продолжение.',
        'Согласен целиком и полностью, всё проверено на практике.',
        'Коротко и ясно, без воды. Так держать!',
        'Пара советов из статьи уже пригодилась. Благодарю.',
        'Не думал, что всё настолько запущено. Есть о чём подумать.',
        'Прочитал, проверил — работает. Рекомендую.'
    ];
    text_count   int;
    counter      int := 0;
BEGIN
    SELECT coalesce(array_agg(id ORDER BY id), '{}') INTO author_ids
    FROM users WHERE email LIKE '%@example.abberant';

    author_count := array_length(author_ids, 1);
    text_count := array_length(texts, 1);

    IF author_count IS NULL OR author_count = 0 THEN
        RETURN;
    END IF;

    FOR target IN
        SELECT id, date FROM articles WHERE status = 'published' ORDER BY date DESC, id
    LOOP
        counter := counter + 1;

        INSERT INTO comments (article_id, user_id, body, created_at)
        VALUES (
            target.id,
            author_ids[1 + (counter % author_count)],
            texts[1 + (counter % text_count)],
            target.date + interval '13 hours' + ((counter % 9) * interval '5 hours')
        );
    END LOOP;
END
$seed$;

-- Richer example threads with replies for the featured articles.
INSERT INTO comments (article_id, user_id, parent_id, body, created_at) VALUES
    ('news-pc-1', (SELECT id FROM users WHERE email = 'demid.krylov@example.abberant'), NULL,
     'Скачал три демоверсии из списка — тактическая оказалась неожиданно глубокой.',
     (SELECT max(created_at) + interval '26 hours' FROM comments WHERE article_id = 'news-pc-1')),
    ('news-pc-1', (SELECT id FROM users WHERE email = 'marina.zaytseva@example.abberant'), NULL,
     'Присоединяюсь: сюжетная адвенчура из списка — главная находка фестиваля.',
     (SELECT max(created_at) + interval '30 hours' FROM comments WHERE article_id = 'news-pc-1')),
    ('news-pc-1', (SELECT id FROM users WHERE email = 'oleg.tarasov@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'news-pc-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'А как вам кооперативная? У меня на процедурных картах вылеты.',
     (SELECT max(created_at) + interval '27 hours' FROM comments WHERE article_id = 'news-pc-1')),
    ('news-pc-1', (SELECT id FROM users WHERE email = 'demid.krylov@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'news-pc-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'Без вылетов, но кооп пока сыроват. Ждём патч до релиза.',
     (SELECT max(created_at) + interval '31 hours' FROM comments WHERE article_id = 'news-pc-1')),

    ('news-xbox-1', (SELECT id FROM users WHERE email = 'artem.gusev@example.abberant'), NULL,
     'Три игры из шести уже в бэклоге, отлично подобрали.',
     (SELECT max(created_at) + interval '24 hours' FROM comments WHERE article_id = 'news-xbox-1')),
    ('news-xbox-1', (SELECT id FROM users WHERE email = 'vera.lebedeva@example.abberant'), NULL,
     'Осталось только найти время на всё это.',
     (SELECT max(created_at) + interval '28 hours' FROM comments WHERE article_id = 'news-xbox-1')),
    ('news-xbox-1', (SELECT id FROM users WHERE email = 'kirill.ezhov@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'news-xbox-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'Тот же вопрос. Game Pass копится быстрее, чем успеваешь играть.',
     (SELECT max(created_at) + interval '29 hours' FROM comments WHERE article_id = 'news-xbox-1')),

    ('news-playstation-1', (SELECT id FROM users WHERE email = 'sonya.romanova@example.abberant'), NULL,
     'State of Play в этот раз порадовал, даты релизов наконец-то назвали.',
     (SELECT max(created_at) + interval '25 hours' FROM comments WHERE article_id = 'news-playstation-1')),
    ('news-playstation-1', (SELECT id FROM users WHERE email = 'ivan.sokolov@example.abberant'), NULL,
     'Единственный минус — половина анонсов без геймплея.',
     (SELECT max(created_at) + interval '33 hours' FROM comments WHERE article_id = 'news-playstation-1')),
    ('news-playstation-1', (SELECT id FROM users WHERE email = 'alina.vetrova@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'news-playstation-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'Зато у показанных проектов графика на уровне.',
     (SELECT max(created_at) + interval '34 hours' FROM comments WHERE article_id = 'news-playstation-1')),

    ('news-nintendo-1', (SELECT id FROM users WHERE email = 'polina.morozova@example.abberant'), NULL,
     'Презентация была на редкость плотной, ни одной пустой минуты.',
     (SELECT max(created_at) + interval '24 hours' FROM comments WHERE article_id = 'news-nintendo-1')),
    ('news-nintendo-1', (SELECT id FROM users WHERE email = 'demid.krylov@example.abberant'), NULL,
     'Дата релиза новой Mario Kart — лучшая новость месяца.',
     (SELECT max(created_at) + interval '32 hours' FROM comments WHERE article_id = 'news-nintendo-1')),
    ('news-nintendo-1', (SELECT id FROM users WHERE email = 'sonya.romanova@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'news-nintendo-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'Согласна, сразу предзаказ оформила.',
     (SELECT max(created_at) + interval '35 hours' FROM comments WHERE article_id = 'news-nintendo-1')),

    ('news-vr-1', (SELECT id FROM users WHERE email = 'oleg.tarasov@example.abberant'), NULL,
     'Трекинг рук реально стал лучше после обновления, проверил вчера.',
     (SELECT max(created_at) + interval '26 hours' FROM comments WHERE article_id = 'news-vr-1')),
    ('news-vr-1', (SELECT id FROM users WHERE email = 'artem.gusev@example.abberant'), NULL,
     'Осталось дождаться нормальных игр под это обновление.',
     (SELECT max(created_at) + interval '29 hours' FROM comments WHERE article_id = 'news-vr-1')),

    ('news-hardware-1', (SELECT id FROM users WHERE email = 'kirill.ezhov@example.abberant'), NULL,
     'Для апгрейда с 30-й серии смысла мало, а вот с 20-й — самое время.',
     (SELECT max(created_at) + interval '27 hours' FROM comments WHERE article_id = 'news-hardware-1')),
    ('news-hardware-1', (SELECT id FROM users WHERE email = 'marina.zaytseva@example.abberant'), NULL,
     'Цены всё равно кусаются, подожду следующего поколения.',
     (SELECT max(created_at) + interval '31 hours' FROM comments WHERE article_id = 'news-hardware-1')),
    ('news-hardware-1', (SELECT id FROM users WHERE email = 'ivan.sokolov@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'news-hardware-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'В тестах прирост честный, но блок питания на 850 Вт уже не хватает.',
     (SELECT max(created_at) + interval '36 hours' FROM comments WHERE article_id = 'news-hardware-1')),

    ('review-1', (SELECT id FROM users WHERE email = 'vera.lebedeva@example.abberant'), NULL,
     'Полностью согласна с оценкой, дополнение выше всяких похвал.',
     (SELECT max(created_at) + interval '40 hours' FROM comments WHERE article_id = 'review-1')),
    ('review-1', (SELECT id FROM users WHERE email = 'oleg.tarasov@example.abberant'), NULL,
     'Боссы в новом регионе — лучшие в серии, имхо.',
     (SELECT max(created_at) + interval '44 hours' FROM comments WHERE article_id = 'review-1')),
    ('review-1', (SELECT id FROM users WHERE email = 'polina.morozova@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'review-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'Только вот баланс поздних этапов всё спорный, как и отмечено.',
     (SELECT max(created_at) + interval '41 hours' FROM comments WHERE article_id = 'review-1')),
    ('review-1', (SELECT id FROM users WHERE email = 'artem.gusev@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'review-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'Патчи это подтянули, сейчас играется отлично.',
     (SELECT max(created_at) + interval '45 hours' FROM comments WHERE article_id = 'review-1')),

    ('review-2', (SELECT id FROM users WHERE email = 'sonya.romanova@example.abberant'), NULL,
     'Пешки действительно выручают, без них я бы заблудилась.',
     (SELECT max(created_at) + interval '38 hours' FROM comments WHERE article_id = 'review-2')),
    ('review-2', (SELECT id FROM users WHERE email = 'kirill.ezhov@example.abberant'), NULL,
     'Свобода иногда вредит: побочные квесты сильно слабее основной линии.',
     (SELECT max(created_at) + interval '42 hours' FROM comments WHERE article_id = 'review-2')),
    ('review-2', (SELECT id FROM users WHERE email = 'vera.lebedeva@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'review-2' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'Зато ни одно прохождение не повторяется, в этом и прелесть.',
     (SELECT max(created_at) + interval '43 hours' FROM comments WHERE article_id = 'review-2')),

    ('review-3', (SELECT id FROM users WHERE email = 'alina.vetrova@example.abberant'), NULL,
     'Визуал действительно топ, но четыре часа кампании — это мало.',
     (SELECT max(created_at) + interval '39 hours' FROM comments WHERE article_id = 'review-3')),
    ('review-3', (SELECT id FROM users WHERE email = 'demid.krylov@example.abberant'), NULL,
     'Согласен: впечатляюще, но цена не оправдана.',
     (SELECT max(created_at) + interval '46 hours' FROM comments WHERE article_id = 'review-3')),

    ('control-resonant-1', (SELECT id FROM users WHERE email = 'ivan.sokolov@example.abberant'), NULL,
     'Прохождение выручает: без него застрял бы на втором акте.',
     (SELECT max(created_at) + interval '30 hours' FROM comments WHERE article_id = 'control-resonant-1')),
    ('control-resonant-1', (SELECT id FROM users WHERE email = 'polina.morozova@example.abberant'), NULL,
     'Способности Джесси описаны точнее, чем в официальной вики.',
     (SELECT max(created_at) + interval '34 hours' FROM comments WHERE article_id = 'control-resonant-1')),
    ('control-resonant-1', (SELECT id FROM users WHERE email = 'kirill.ezhov@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'control-resonant-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'Подтверждаю, сам так же прошёл этот момент.',
     (SELECT max(created_at) + interval '37 hours' FROM comments WHERE article_id = 'control-resonant-1')),

    ('ghost-of-tsushima-1', (SELECT id FROM users WHERE email = 'marina.zaytseva@example.abberant'), NULL,
     'Испытания на золото прошли со второй попытки, награды того стоят.',
     (SELECT max(created_at) + interval '28 hours' FROM comments WHERE article_id = 'ghost-of-tsushima-1')),
    ('ghost-of-tsushima-1', (SELECT id FROM users WHERE email = 'artem.gusev@example.abberant'), NULL,
     'Подскажите, все ли испытания доступны до финала?',
     (SELECT max(created_at) + interval '31 hours' FROM comments WHERE article_id = 'ghost-of-tsushima-1')),
    ('ghost-of-tsushima-1', (SELECT id FROM users WHERE email = 'marina.zaytseva@example.abberant'),
     (SELECT id FROM comments WHERE article_id = 'ghost-of-tsushima-1' AND parent_id IS NULL ORDER BY id LIMIT 1),
     'Все, кроме одного — оно открывается уже после финальных титров.',
     (SELECT max(created_at) + interval '32 hours' FROM comments WHERE article_id = 'ghost-of-tsushima-1')),

    ('silent-hill-townfall-1', (SELECT id FROM users WHERE email = 'sonya.romanova@example.abberant'), NULL,
     'Прошла без смертей с первого раза, гайд написан понятно.',
     (SELECT max(created_at) + interval '29 hours' FROM comments WHERE article_id = 'silent-hill-townfall-1')),
    ('silent-hill-townfall-1', (SELECT id FROM users WHERE email = 'oleg.tarasov@example.abberant'), NULL,
     'Секция с туманом — самое нервное место, тут гайд незаменим.',
     (SELECT max(created_at) + interval '33 hours' FROM comments WHERE article_id = 'silent-hill-townfall-1')),

    ('walkthrough-1', (SELECT id FROM users WHERE email = 'demid.krylov@example.abberant'), NULL,
     'Карты с отметками очень выручают, спасибо редакции.',
     (SELECT max(created_at) + interval '27 hours' FROM comments WHERE article_id = 'walkthrough-1')),
    ('walkthrough-1', (SELECT id FROM users WHERE email = 'alina.vetrova@example.abberant'), NULL,
     'Боссфайты разобраны лучше, чем в любых видео.',
     (SELECT max(created_at) + interval '35 hours' FROM comments WHERE article_id = 'walkthrough-1'));
