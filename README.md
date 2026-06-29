# AURORA VPN Panel

<p align="center">
  <img src="https://raw.githubusercontent.com/bychikola/aurora-vpn-panel/master/install/media/aurora-banner.png" alt="AURORA VPN Panel" width="600">
</p>

<p align="center">
  <b>Масштабируемая панель управления VPN на базе Xray-core</b><br>
  Go · React · TypeScript · PostgreSQL · Redis · Docker
</p>

<p align="center">
  <a href="#быстрая-установка">🚀 Быстрая установка</a> ·
  <a href="#возможности">✨ Возможности</a> ·
  <a href="#архитектура">🏗 Архитектура</a> ·
  <a href="#ручная-установка">🔧 Ручная установка</a> ·
  <a href="#api">📡 API</a>
</p>

---

## Быстрая установка

Одна команда на чистом VPS с Debian/Ubuntu:

```bash
curl -fsSL https://raw.githubusercontent.com/bychikola/aurora-vpn-panel/master/install/install_aurora.sh | bash
```

**Зеркало (если GitHub недоступен):**

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/bychikola/aurora-vpn-panel@master/install/install_aurora.sh | bash
```

После установки управление через меню:

```bash
aurora
```

### Что делает установщик

1. Проверяет ОС (Debian 11/12, Ubuntu 22.04/24.04)
2. Ставит Docker, certbot, UFW и остальные зависимости
3. Настраивает BBR congestion control
4. Выпускает SSL-сертификат через Let's Encrypt
5. Разворачивает все контейнеры через Docker Compose
6. Создаёт админа и выводит пароль

### Варианты установки

| Пункт меню | Описание |
|-----------|----------|
| **1. Panel + Node** | Всё на одном сервере (до 500 пользователей) |
| **2. Panel only** | Только панель управления |
| **3. Node only** | Только сервер с Xray-core |
| **4. Add node** | Подключить новую ноду к существующей панели |

### Требования к серверу

- Debian 11/12 или Ubuntu 22.04/24.04
- Минимум 1 GB RAM (2 GB рекомендуется)
- 10 GB свободного места
- Открытые порты: 22 (SSH), 443 (HTTPS)
- Домен с A-записью на IP сервера

---

## Возможности

### Протоколы и транспорт

| Протоколы | Транспорт | Безопасность |
|-----------|----------|-------------|
| VLESS | TCP | TLS |
| VMess | HTTP | Reality |
| Trojan | WebSocket | XTLS-Vision |
| Shadowsocks | gRPC | |
| Shadowsocks-2022 | QUIC | |
| Hysteria2 | | |
| TUIC v5 | | |

### Управление

- **Пользователи** — CRUD, лимиты по трафику/дате/IP
- **Подписки** — автоматическая генерация v2ray/clash/sing-box конфигов
- **Ноды** — мульти-серверная архитектура с синхронизацией
- **Inbounds** — hot-reload конфигураций Xray без перезагрузки
- **Мониторинг** — графики трафика, CPU/RAM, активные подключения

### Дашборд

<p align="center">
  <i>Трафик за 24 часа · Распределение протоколов · Рост пользователей · Активные соединения</i>
</p>

- Статистика трафика в реальном времени (обновление каждые 15 секунд)
- Графики нагрузки CPU/RAM/диск по каждой ноде
- Логи подключений с геолокацией
- Автоматическое отключение истёкших пользователей

### Подписки

Публичный эндпоинт для клиентов:

```
https://panel.example.com/api/v1/sub/abc123
```

Автоопределение формата по User-Agent:
- **v2rayNG / Nekoray** → base64 ссылки
- **Clash Meta / Stash** → YAML конфиг
- **Sing-box / SFI** → JSON конфиг

Ручной выбор формата:

```
/sub/abc123/base64   — v2ray формат
/sub/abc123/clash    — Clash Meta YAML
/sub/abc123/singbox  — Sing-box JSON
/sub/abc123/qrcode   — QR-код
```

### Безопасность

- JWT авторизация (access 15m + refresh 7d)
- Argon2id хеширование паролей
- AES-256-GCM шифрование чувствительных данных
- Fail2Ban защита от перебора паролей
- Rate limiting на логин и API
- Автоочистка устаревших логов

---

## Архитектура

```
┌─────────────────────────────────────────┐
│              AURORA Panel               │
│  ┌─────────┐ ┌──────────┐ ┌─────────┐  │
│  │Frontend │ │ Backend  │ │   DB    │  │
│  │ React   │ │   Go     │ │PostgreSQL│  │
│  │ Nginx   │ │  Fiber   │ │  Redis  │  │
│  └─────────┘ └────┬─────┘ └─────────┘  │
│                   │ gRPC                │
└───────────────────┼─────────────────────┘
                    │
        ┌───────────┼───────────┐
        │           │           │
   ┌────┴────┐ ┌───┴─────┐ ┌──┴──────┐
   │ Node 1  │ │ Node 2  │ │ Node 3  │
   │ Xray    │ │ Xray    │ │ Xray    │
   │ Agent   │ │ Agent   │ │ Agent   │
   └─────────┘ └─────────┘ └─────────┘
```

### Стек

| Компонент | Технология |
|-----------|-----------|
| **Frontend** | React 19, TypeScript, Tailwind CSS, TanStack Query, Recharts |
| **Backend** | Go, Fiber, pgx, go-redis, gRPC |
| **Core** | Xray-core (gRPC API) |
| **База данных** | PostgreSQL 16 (основные данные), Redis 7 (кэш/статистика) |
| **Развёртывание** | Docker, Docker Compose, Nginx |

### Структура проекта

```
aurora-vpn-panel/
├── aurora-frontend/          # React SPA
│   ├── src/
│   │   ├── pages/            # Dashboard, Users, Nodes, Inbounds, Subscriptions, Settings
│   │   ├── components/       # Sidebar, Layout
│   │   ├── api/              # TanStack Query hooks, API client, mock data
│   │   └── types/            # TypeScript интерфейсы
│   └── ...
├── aurora-backend/           # Go REST API + gRPC
│   ├── cmd/aurora/           # Точка входа
│   ├── internal/
│   │   ├── adapter/          # PostgreSQL, Redis, Xray gRPC
│   │   ├── service/          # Бизнес-логика
│   │   ├── handler/          # HTTP обработчики
│   │   ├── worker/           # Фоновые задачи
│   │   └── domain/           # Сущности
│   └── deployments/          # Dockerfile, docker-compose
├── install/                  # Установщик (однострочный)
│   ├── install_aurora.sh     # Главный скрипт
│   └── src/
│       ├── lang/             # EN/RU переводы
│       └── modules/          # Модули установки
└── plan.md                   # Спецификация проекта
```

---

## Ручная установка

### Фронтенд

```bash
cd aurora-frontend
npm install
npm run dev      # Dev-сервер на :5173
npm run build    # Production-сборка в dist/
```

### Бэкенд

```bash
cd aurora-backend

# PostgreSQL + Redis
docker compose -f deployments/docker-compose.yml up -d postgres redis

# Миграции
make migrate-up

# Запуск
make dev          # :8080
```

### Docker (полный стек)

```bash
cd aurora-backend
docker compose -f deployments/docker-compose.yml up -d
```

### Продакшен

```bash
# Клонируем репозиторий
git clone https://github.com/bychikola/aurora-vpn-panel.git
cd aurora-vpn-panel

# Запускаем установщик (рекомендуется)
sudo bash install/install_aurora.sh
```

---

## API

### Аутентификация

```http
POST /api/v1/auth/login     # { username, password } → { accessToken, refreshToken }
POST /api/v1/auth/refresh   # { refreshToken } → { accessToken }
POST /api/v1/auth/logout
GET  /api/v1/auth/me
```

### Пользователи

```http
GET    /api/v1/users                      # ?search=&status=&protocol=&page=1&pageSize=15
POST   /api/v1/users                      # Создание пользователя
GET    /api/v1/users/:id
PUT    /api/v1/users/:id
DELETE /api/v1/users/:id
POST   /api/v1/users/:id/reset-traffic    # Сброс трафика
POST   /api/v1/users/:id/reset-token      # Новый токен подписки
GET    /api/v1/users/:id/qrcode           # QR-код (PNG)
```

### Ноды

```http
GET    /api/v1/nodes
GET    /api/v1/nodes/:id
POST   /api/v1/nodes
PUT    /api/v1/nodes/:id
DELETE /api/v1/nodes/:id
```

### Inbounds

```http
GET    /api/v1/inbounds                     # ?nodeId=xxx
GET    /api/v1/inbounds/:id
POST   /api/v1/inbounds                     # Создание + отправка в Xray
PUT    /api/v1/inbounds/:id
DELETE /api/v1/inbounds/:id
POST   /api/v1/inbounds/:id/reload          # Hot-reload
```

### Подписки (публичные, без авторизации)

```http
GET /api/v1/sub/:token          # Автоформат по User-Agent
GET /api/v1/sub/:token/base64   # v2ray формат
GET /api/v1/sub/:token/clash    # Clash Meta YAML
GET /api/v1/sub/:token/singbox  # Sing-box JSON
GET /api/v1/sub/:token/qrcode   # QR-код (PNG)
```

### Дашборд и настройки

```http
GET /api/v1/dashboard/stats
GET /api/v1/settings
PUT /api/v1/settings
```

---

## Конкуренты и ориентиры

Проект ориентируется на лучшие решения рынка:

- **Marzban** — самый популярный, Python/FastAPI
- **Remnawave** — современный, TypeScript/NestJS
- **3X-UI** — минималистичный, Go/Vue.js

AURORA отличается:
- **Go-бэкенд** — минимальный Docker-образ (~12 MB), низкое потребление памяти
- **Мульти-нодность из коробки** — Push-синхронизация через gRPC
- **Три формата подписок** — v2ray + Clash Meta + Sing-box
- **Однострочная установка** — curl | bash, полный цикл за 5 минут

---

## Лицензия

MIT License. Свободное использование, модификация и распространение.

---

## Ссылки

- 🌐 [GitHub](https://github.com/bychikola/aurora-vpn-panel)
- 📦 [Docker Hub](https://github.com/bychikola/aurora-vpn-panel/pkgs/container/aurora-backend)
- 📖 [Документация](https://github.com/bychikola/aurora-vpn-panel/wiki)
- 🐛 [Баг-репорты](https://github.com/bychikola/aurora-vpn-panel/issues)
