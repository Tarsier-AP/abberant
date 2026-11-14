# abberant

Gaming news site

Стек: Vue 3 + Vite (фронтенд), Go + Fiber (бэкенд), PostgreSQL, nginx.

## Структура

```
backend/    — API-сервер на Go/Fiber (/api/v1), PostgreSQL, SMTP-почта
frontend/   — сайт на Vue 3 + Vite
nginx/      — конфиги nginx для Linux и Windows
```

## Требования

- Node `^22.18.0 || >=24.12.0` и npm — для фронтенда
- Go 1.26+ — для бэкенда
- Docker не нужен: в dev-режиме бэкенд сам поднимает встроенный PostgreSQL
  (embedded-postgres, порт `55432`). Для продакшена укажите `ABBERANT_DATABASE_DSN`.

## Фронтенд

Команды выполняются из каталога `frontend/`:

```sh
npm ci
npm run dev       # дев-сервер с прокси /api → http://localhost:8090
npm run build     # production-сборка в dist/ (включая type-check)
npm run lint
```

Если `dev`/`build` падает с «Illegal instruction (core dumped)» (специфика
отдельных машин из-за нативной сборки Rolldown):

```sh
NAPI_RS_FORCE_WASI=true npm run dev
NAPI_RS_FORCE_WASI=true npm run build
```

## Бэкенд

Команды выполняются из каталога `backend/`:

```sh
go run ./cmd/server   # запуск: API на :8090 + встроенный PostgreSQL
go build ./...
go vet ./...
```

Первый запуск скачает бинарники PostgreSQL и создаст базу в `backend/.data/`
(каталог можно переопределить через `ABBERANT_DB_DATA_DIR`).

Без настроенного SMTP сервер работает в dev-режиме: коды активации и 2FA
печатаются в лог сервера и возвращаются клиенту в поле `devCode`.

### Переменные окружения

| Переменная | По умолчанию | Назначение |
| --- | --- | --- |
| `ABBERANT_ADDR` | `:8090` | адрес HTTP-сервера |
| `ABBERANT_DATABASE_DSN` | — (встроенный PG) | строка подключения к внешнему PostgreSQL |
| `ABBERANT_DB_DATA_DIR` | `.data/postgres` | каталог данных встроенного PostgreSQL |
| `ABBERANT_DB_PORT` | `55432` | порт встроенного PostgreSQL |
| `ABBERANT_MEDIA_DIR` | `.data/media` | каталог загруженных изображений (`/media/{hash}/…`) |
| `ABBERANT_COOKIE_SECURE` | `false` | флаг Secure у cookie сессии (включайте за HTTPS) |
| `ABBERANT_PUBLIC_ORIGIN` | `http://localhost:5173` | origin сайта для писем |
| `ABBERANT_DEV_MODE` | `true` | dev-режим (видимость кодов) |
| `SMTP_HOST`, `SMTP_PORT`, `SMTP_USERNAME`, `SMTP_PASSWORD`, `SMTP_FROM` | — | отправка писем; без `SMTP_HOST` — dev-режим |

### Изображения (медиа-хранилище)

Загруженные обложки, картинки в тексте статей и аватары сохраняются на диск
(`ABBERANT_MEDIA_DIR`): оригинал — без перекодирования, плюс варианты ширины
320/480/828/1200 (аватары — 160×160). Адреса вида `/media/{hash}/{файл}`
неизменяемы и кэшируются навсегда.

Формат вариантов зависит от платформы бэкенда:

- **Linux** — WebP (качество 80) через cgo (`libwebp`); достаточно
  установленной runtime-библиотеки `libwebp7` (soname `libwebp.so.7`),
  заголовки не нужны.
- **Windows** — PNG (сборка остаётся чистым Go без cgo: `CGO_ENABLED=0`,
  в том числе кросс-компиляция с Linux).

Фронтенд определяет расширение по URL, поэтому один и тот же код работает
с обоими форматами.

## Сборка и запуск на Linux

1. Фронтенд:

   ```sh
   cd frontend && npm ci && npm run build
   ```

2. Бэкенд (нужна установленная библиотека libwebp — пакет `libwebp7` /
   `libwebp`, заголовки не требуются):

   ```sh
   cd backend && go build -o server ./cmd/server
   ```

3. Запуск: `./server` (в фоне — `setsid ./server &`).

4. nginx: скопировать `nginx/abberant.conf` в `/etc/nginx/sites-available/abberant`,
   создать symlink в `sites-enabled`, в конфиге поправить `root` на путь к
   `frontend/dist`, затем `nginx -t && systemctl reload nginx`.

Итог: статика и API на одном origin (`/api/` проксируется на `:8090`),
cookie сессии работает без CORS.

## Сборка и запуск на Windows

1. Фронтенд (сборка платформо-независимая, можно сделать на Linux и
   скопировать `frontend/dist/` на Windows-машину):

   ```sh
   cd frontend && npm ci && npm run build
   ```

2. Бэкенд — чистый Go, можно собрать нативно или кросс-компиляцией с Linux
   (варианты изображений на Windows — PNG, без cgo):

   ```bat
   cd backend
   go build -o server.exe .\cmd\server
   ```

   ```sh
   cd backend
   GOOS=windows GOARCH=amd64 CGO_ENABLED=0 go build -o server.exe ./cmd/server
   ```

   Разместите каталог медиа (по умолчанию `.data\media`, см.
   `ABBERANT_MEDIA_DIR`) там, куда указывает `alias` в nginx-конфиге.

3. Запуск бэкенда (например, из `C:\abberant\`; каталог должен допускать
   запись — там появится `.data\postgres`):

   ```bat
   set ABBERANT_ADDR=:8090
   server.exe
   ```

   Embedded-postgres сам развернёт PostgreSQL при первом запуске.

4. nginx: распаковать официальный zip с nginx.org в `C:\nginx`, скопировать
   `nginx/abberant-windows.conf` в `C:\nginx\conf\`, в конфиге поправить
   `root` (путь к `dist`, слэши прямые: `C:/abberant/dist`). Подключить одним
   из способов:

   - добавить в `C:\nginx\conf\nginx.conf` внутрь `http {}`:
     `include conf/abberant.conf;`
   - либо запустить `nginx.exe -c conf/abberant.conf` из `C:\nginx`.

   Управление: `nginx.exe` (старт), `nginx.exe -s stop`, `nginx.exe -s reload`.

   Порядок запуска: сначала `server.exe`, затем nginx — иначе `/api/` будет
   отдавать 502.

5. Проверка: `http://localhost/` — сайт; `http://localhost/api/v1/me` — JSON
   от бэкенда (401 без сессии — норма).

Для постоянной работы как службы под Windows используйте обёртки NSSM или
WinSW (nginx/Windows не поддерживает штатный режим сервиса).

## Разработка авторизации/2FA

- Флоу: `/register` → `/register/password` → `/register/activate` (код из
  письма) → `/login`; при включённой 2FA вход продолжается на `/2fa`.
- 2FA: TOTP (приложение-аутентификатор, QR + 8 одноразовых резервных кодов)
  и email-коды; SMS помечен «скоро» и отклоняется сервером.
- Подробности по API — `backend/internal/api/` и `plan.md`.
