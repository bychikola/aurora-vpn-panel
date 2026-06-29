# AURORA VPN Panel — Unified Installer

Однострочная установка панели управления VPN на базе Xray-core.

## Быстрый старт

```bash
curl -fsSL https://raw.githubusercontent.com/bychikola/aurora-vpn-panel/main/install/install_aurora.sh | bash
```

Или с зеркалом:

```bash
curl -fsSL https://cdn.jsdelivr.net/gh/bychikola/aurora-vpn-panel@main/install/install_aurora.sh | bash
```

После установки скрипт доступен по команде:

```bash
aurora
```

## Меню

```
╔══════════════════════════════════╗
║     AURORA VPN Panel             ║
║     Unified Installer v1.0.0     ║
╚══════════════════════════════════╝

1. Install AURORA (Panel + Node on one server)
2. Install AURORA Panel only
3. Install AURORA Node only
4. Add node to existing panel
5. Manage panel/node
6. SSL Certificates
7. Backup & Restore
8. Update AURORA
9. Remove AURORA
```

## Что делает установщик

1. **Проверяет ОС** — Debian 11/12, Ubuntu 22.04/24.04
2. **Устанавливает зависимости** — Docker, curl, certbot, UFW
3. **Настраивает BBR** — congestion control для лучшей производительности
4. **Выпускает SSL** — через Let's Encrypt (Cloudflare DNS или HTTP)
5. **Разворачивает контейнеры** — PostgreSQL, Redis, Backend, Frontend, Nginx
6. **Создаёт администратора** — генерирует надёжный пароль
7. **Настраивает автообновление** — cron для SSL и unattended-upgrades

## Структура

```
install/
├── install_aurora.sh         # Главный скрипт (entry point)
├── README.md
├── src/
│   ├── lang/
│   │   ├── en.sh             # Английские переводы
│   │   └── ru.sh             # Русские переводы
│   └── modules/
│       ├── install_panel.sh  # Модуль установки панели
│       ├── install_node.sh   # Модуль установки ноды
│       └── add_node.sh       # Модуль добавления ноды
```

## Установка на VPS

### Панель + Нода (один сервер)

Подходит для небольших инсталляций до 500 пользователей.

```bash
# Установка за 5 минут
bash <(curl -fsSL https://raw.githubusercontent.com/bychikola/aurora-vpn-panel/main/install/install_aurora.sh)
# Выбрать пункт 1
```

### Панель отдельно, Нода отдельно (рекомендуется)

**На сервере панели:**
```bash
bash <(curl -fsSL ...)
# Выбрать пункт 2 — "Install Panel only"
```

**На сервере ноды:**
```bash
bash <(curl -fsSL ...)
# Выбрать пункт 4 — "Add node to existing panel"
# Указать URL панели и API ключ
```

## Требования

- Debian 11/12 или Ubuntu 22.04/24.04
- Минимум 1 GB RAM (2 GB рекомендуется)
- 10 GB свободного места
- Открытый порт 443 (HTTPS) + 22 (SSH)
- Домен, направленный A-записью на IP сервера

## Самообновление

Скрипт проверяет наличие обновлений при запуске и предлагает обновиться:

```
[!] New version available: v1.1.0 (current: v1.0.0)
[?] Update? (y/n):
```

Все модули (языки, модули установки) загружаются с GitHub при первом использовании и кэшируются локально.
