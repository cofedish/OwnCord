# Форк OwnCord

Linux-first форк OwnCord для self-host использования: с Docker/Compose-деплоем, более безопасными дефолтами для публикации в интернет и первым слоем security hardening.

## Что это за форк

Этот репозиторий сохраняет кодовую базу OwnCord, но меняет эксплуатационный фокус:

- сервер ориентирован на запуск под Linux
- Docker/Compose является основным способом деплоя
- внешний доступ предполагается через reverse proxy, например Caddy или Nginx
- исправлена часть проблем безопасности, найденных при локальном аудите
- desktop-клиент остаётся отдельным Tauri-приложением и пока по-прежнему ориентирован на Windows

## Текущее состояние

Проект всё ещё остаётся ранним ПО. Этот форк безопаснее и удобнее для self-host, чем исходное состояние, которое я аудировал, но это всё ещё не зрелая защищённая messaging-платформа.

Подходит для:

- лабораторных стендов
- домашнего self-host
- небольших доверенных групп при нормальной настройке TLS и reverse proxy

Не подходит для:

- чувствительной переписки
- враждебной multi-tenant среды
- серьёзной high-scale production-нагрузки

## Что изменено в форке

- добавлены Dockerfile и Compose-деплой для сервера
- добавлена Linux-friendly документация по деплою
- усилена модель доступа к `/admin`
- админка переведена с хранения токена в `localStorage` на `HttpOnly` cookie-backed session flow
- приватные вложения больше не отдаются как публично кэшируемые
- desktop-клиент больше не отключает проверку TLS для REST/media
- production devtools убраны из дефолтной клиентской сборки
- tightened Tauri defaults и снижена опасная production surface
- для Linux/container deployment обновления предполагаются через rebuild/redeploy образа, а не через built-in updater

## Быстрый старт

### Сервер

1. Склонировать репозиторий.
2. Скопировать `.env.example` в `.env`.
3. Установить `OWNCORD_DOMAIN`.
4. Проверить `deploy/config/owncord.container.yaml`.
5. Поднять стек:

```bash
docker compose --profile proxy up -d --build
```

6. Открыть админку через локальный туннель:

```bash
ssh -L 8080:127.0.0.1:8080 your-server
```

7. Открыть `http://127.0.0.1:8080/admin`, создать owner-аккаунт и затем выпустить invite-коды для пользователей.

### Клиент

Desktop-клиент устанавливается отдельно от сервера.

Сборка из исходников:

```bash
cd Client/tauri-client
npm install
npm run tauri build
```

Windows-инсталлятор появится здесь:

`Client/tauri-client/src-tauri/target/release/bundle/nsis/`

Важно:

- в этом форке больше нет небезопасного TLS bypass для REST/media
- клиент должен подключаться к нормальному HTTPS-адресу за reverse proxy
- локальный `127.0.0.1:8080` нужен для админки, а не как обычный пользовательский endpoint

## Состояние безопасности

В форке есть реальные исправления по итогам аудита, но “стало безопаснее” не означает “проект полностью доверенный”.

Что уже усилено:

- desktop TLS trust model
- admin session handling
- caching/privacy для вложений
- Linux deployment path
- дефолты для reverse proxy и admin exposure

Что всё ещё ограничивает проект:

- SQLite остаётся базой и узким местом по масштабированию
- проект всё ещё alpha-grade
- desktop-клиенту всё ещё нужен более широкий runtime hardening и дополнительная end-to-end валидация
- voice/LiveKit остаётся чувствительным с эксплуатационной точки зрения компонентом

## Рекомендуемая схема деплоя

Нормальная practical-схема для self-host:

- OwnCord server в Docker
- Caddy или Nginx для публичного TLS
- `/admin` закрыт снаружи и доступен только через localhost/VPN/SSH tunnel
- регулярные бэкапы `deploy/runtime/data`
- обновления через rebuild/redeploy контейнера

Смотри также:

- [README.md](README.md)
- [README.en.md](README.en.md)
- [docs/docker-deployment.md](docs/docker-deployment.md)
- [docs/production-deployment.md](docs/production-deployment.md)

## Структура репозитория

- `Server/` — Go-сервер
- `Client/tauri-client/` — Tauri desktop client
- `deploy/` — Compose, Caddy, конфиги, примеры systemd
- `docs/` — документация по деплою и проекту

## Лицензия

AGPL-3.0
