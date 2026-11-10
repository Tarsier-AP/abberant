# AGENTS.md

## Stack

- Backend: Go using the Fiber web framework.
- Database: PostgreSQL.
- Server/reverse proxy: nginx.
- Target OS: Linux and Windows.
- Fiber, PostgreSQL, and nginx are the intended stack; `backend/` currently has no code or Go module.

## Project shape

- Gaming news site with a Russian-language UI. Keep new user-facing copy in Russian.
- `frontend/` is the only implemented part. `backend/` is still empty; the planned backend is Go/Fiber with PostgreSQL behind nginx.
- The root `.gitignore` is a generic Go template; it does not mean a Go backend exists.
- There is no root `package.json`; the root `package-lock.json` is empty. Run all npm commands from `frontend/`.
- Node `^22.18.0 || >=24.12.0`; use npm (despite the old "change from bun to npm" commit).

## Commands

Run from `frontend/`:

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

- Both `dev` and `build` currently need `NAPI_RS_FORCE_WASI=true` on this machine because the native Rolldown binding crashes; the README only documents this for `dev`.
- `build` runs `type-check` and `vite build` in parallel via `npm-run-all2`.
- For a focused lint check, run `npx eslint <path>`; there is no per-file typecheck.
- There is no test runner, test script, or test directory. Do not invent a test command; add Vitest only if tests are requested.
- `npm run lint` currently fails on the empty leftover `src/components/Header.vue` (`vue/multi-word-component-names`). Prefer deleting that unused file. ESLint also warns that `.oxlintrc.json` is missing; that is not the failure.
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

## TypeScript notes

- TypeScript is the `typescript-native-bridge` (tsgo) package, not upstream `typescript`; `type-check` prints `TNB ACTIVE`. Keep this setup unless explicitly changing toolchains.
- `noUncheckedIndexedAccess` is enabled, so indexed access is `T | undefined`.
