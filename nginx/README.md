# nginx configs

- `abberant.conf` — Linux variant: server block for `sites-available` /
  `include` from the main `http {}` context. Serves the vite build from
  `/var/www/abberant/dist` and proxies `/api/` to `127.0.0.1:8090`.
- `abberant-windows.conf` — standalone config for nginx/Windows: relative
  paths, single worker, `pid logs/nginx.pid`. Put it into `conf/` of the
  nginx install and include it from `conf/nginx.conf`, or run
  `nginx.exe -c conf/abberant.conf` from the install directory. Adjust
  `root` to the deployed `dist` folder.

Both variants serve the frontend statics with `try_files ... /index.html`
(SPA history mode) and reverse-proxy `/api/` to the Fiber backend on port
8090. Cookies are same-origin (`SameSite=Lax`), so no extra CORS setup is
needed.
