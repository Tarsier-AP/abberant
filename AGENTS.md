# AGENTS.md

## Stack

- Backend: Go using the Fiber web framework (implemented, `backend/`).
- Database: PostgreSQL via `jackc/pgx/v5`; dev runs an `embedded-postgres` instance (no Docker needed) when `ABBERANT_DATABASE_DSN` is unset.
- Server/reverse proxy: nginx (`nginx/` has Linux and Windows variants plus a README).
- Target OS: Linux and Windows.

## Project shape

- Gaming news site with a Russian-language UI. Keep new user-facing copy in Russian.
- Auth is real: Go/Fiber backend under `/api/v1` with server-side sessions (opaque HttpOnly cookie), Argon2id passwords, single-use emailed codes and 2FA (TOTP + email codes + backup codes; SMS is rejected as «скоро»).
- The root `.gitignore` covers Go artifacts and the embedded-postgres data dir (`backend/.data/`).
- There is no root `package.json`; the root `package-lock.json` is empty. Run all npm commands from `frontend/`.
- Node `^22.18.0 || >=24.12.0`; use npm (despite the old "change from bun to npm" commit).

## Workflow

- After completing a task, commit the changes with a clear, descriptive message. Do not leave completed work uncommitted.

## Commands

Frontend (run from `frontend/`):

```sh
npm ci
npm run dev
npm run type-check
npm run lint
npm run build

# If dev/build exits with "Illegal instruction (core dumped)":
NAPI_RS_FORCE_WASI=true npm run dev
NAPI_RS_FORCE_WASI=true npm run build
```

Backend (run from `backend/`):

```sh
go build ./...
go run ./cmd/server          # starts embedded PostgreSQL on :55432 + API on :8090
ABBERANT_ADDR=":18080" go run ./cmd/server   # when :8090 is taken on this machine
```

- Dev proxy: `vite.config.ts` forwards `/api` → `http://localhost:8090`.
- Without SMTP configured the backend logs every code and returns it to the client as `devCode` (dev mode).

- Both `dev` and `build` currently need `NAPI_RS_FORCE_WASI=true` on this machine because the native Rolldown binding crashes; the README only documents this for `dev`.
- `build` runs `type-check` and `vite build` in parallel via `npm-run-all2`.
- For a focused lint check, run `npx eslint <path>`; there is no per-file typecheck.
- There is no test runner, test script, or test directory. Do not invent a test command; add Vitest only if tests are requested.
- `npm run lint` currently passes; ESLint only warns that `.oxlintrc.json` is missing. The old empty `src/components/Header.vue` has been removed.
- `npm run format` uses Prettier but the repo has no Prettier config, so it currently rewrites every file with Prettier defaults (semicolons, double quotes). Avoid it unless doing a deliberate formatting-only change; match existing style manually.

## Architecture and conventions

- Entrypoint: `src/main.ts` mounts `App.vue` with Pinia and vue-router. `App.vue` owns the global header/footer shell.
- Routes live in `src/router/index.ts`. Most routes, including `/article/:id`, render `PlaceholderView`; replace the placeholder when implementing a section.
- Article types are in `src/types/article.ts`; all content is hard-coded mock data in `src/data/articles.ts`. There is no fetch/API layer yet, so do not assume one.
- `@/` maps to `src/` in both `vite.config.ts` and `tsconfig.app.json`.
- Components use `<script setup lang="ts">`, scoped styles, and BEM class names (`card__title`, `card__meta`).
- Styling relies on CSS custom properties from `src/assets/main.css`. Theme is selected by `data-theme="light|dark"` on `<html>`, persisted under localStorage key `abberant-theme`; the initial value is set inline in `index.html` and toggled in `SiteHeader.vue`. Add new colors to both theme blocks rather than hard-coding them.
- Dates are ISO strings displayed as `dd.mm.yyyy` via `src/utils/date.ts`.
- `src/stores/counter.ts` is unused Vue template scaffolding; do not use it as a state-management pattern to copy.

### Account and auth pages

- Auth is backed by the Go API (`src/api/client.ts`, `src/api/auth.ts`); the session lives in an HttpOnly cookie, `src/stores/auth.ts` holds only the loaded account and pending registration state. `useAuthStore().init()` runs before app mount (`src/main.ts`).
- Flow: `/register` (email + username) → `/register/password` → `/register/activate` (6-digit emailed code) → `/login`. `/activate` is the recovery form (email + code). Login goes to `/2fa` only when the backend answers `2fa_required`.
- Registration is email-only; phone is an editable profile field in settings.
- 2FA: email codes immediately; authenticator apps go through enrollment (`POST /2fa/totp/setup` → QR + secret, `/2fa/totp/confirm` with a live code → 8 one-time backup codes shown once). SMS is listed but disabled («скоро») and rejected server-side.
- Settings also offer «Выйти» and «Выйти со всех устройств» (`POST /logout/all`).
- `.account-settings__card` is a transparent layout wrapper: `max-width: 800px`, `padding: 8px`, `gap: 10px`, `border-radius` only — **no `background` and no `box-shadow`**. Do not restore them; surfaces belong to the inner blocks.
- Visual surfaces: `.account-settings__section` uses `--color-surface` + `--shadow`; `.account-settings__row`, `.account-settings__profile` and `.account-settings__email-form` use `--color-card` + `--shadow`.
- Section order: «Личная информация» (avatar, phone, additional email, email), «Безопасность» (2FA switch + method, enrollment, sessions), «Уведомления» (email, SMS). The old debug role switch is gone — roles come from the server.
- The avatar is a 72px circle with a hover/focus overlay «Изменить»; uploads are cropped to 160×160 WebP and stored as a data URL through `PATCH /me`.

### Grid system

Card grids are rendered by `HomeSection` (`src/components/home/HomeSection.vue`) and its `ArticleSection` data.

- Default grid (`/`, and index pages): 3 columns at `>=1024px`, 2 at `768–1023px`, 1 at `<=767px`. The default caps each section at two rows, hiding cards with `nth-child` rules (6/4/2 cards respectively).
- `show-all`: disables the two-row cap. Use it on filtered/archive pages so every card on the current page is visible. Filtered `/news`, `/reviews`, and `/walkthroughs` pages pass it.
- News pages (`/news`, `/news/:platform`) use the default 3-column grid with six cards per section/page, giving 3×2 on desktop. `PAGE_LIMIT = 6` there.
- `dense`: legacy archive variant (2 columns on mobile/tablet, 3 on desktop, no row cap). No current view uses it; prefer the default grid plus `show-all`.

Pagination limits: filtered news pages slice cards at `PAGE_LIMIT = 6` per page; filtered review and walkthrough feeds use `PAGE_LIMIT = 8`. The `/walkthroughs` index groups whole sections instead and pages by preview-card count.

## TypeScript notes

- TypeScript is the `typescript-native-bridge` (tsgo) package, not upstream `typescript`; `type-check` prints `TNB ACTIVE`. Keep this setup unless explicitly changing toolchains.
- `noUncheckedIndexedAccess` is enabled, so indexed access is `T | undefined`.
