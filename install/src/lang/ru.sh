#!/bin/bash
# AURORA Installer — Russian translations

declare -gA LANG

LANG[CHOOSE_LANG]="Выберите язык / Select language:"
LANG[LANG_EN]="English"
LANG[LANG_RU]="Русский"
LANG[WELCOME]="AURORA VPN Panel — Установщик"
LANG[MENU_TITLE]="AURORA VPN PANEL"
LANG[VERSION_LABEL]="Версия: %s"
LANG[AVAILABLE_UPDATE]="доступно обновление скрипта"
LANG[SELECT_ACTION]="Выберите действие (0-9):"
LANG[EXIT]="Выход"
LANG[INVALID_CHOICE]="Неверный выбор. Выберите 0-9."
LANG[BACK]="« Назад"
LANG[WAITING]="Пожалуйста, подождите..."
LANG[DONE]="Готово!"
LANG[ERROR_ROOT]="Скрипт должен запускаться от root"
LANG[ERROR_OS]="Поддерживаются только: Debian 11/12, Ubuntu 22.04/24.04"
LANG[ERROR_DOCKER]="Docker не установлен"
LANG[CONTINUE_PROMPT]="Продолжить? (y/n):"
LANG[CONFIRM_YES]="y"

# Пункты меню
LANG[MENU_1]="Установить AURORA (Панель + Нода на одном сервере)"
LANG[MENU_2]="Установить только Панель AURORA"
LANG[MENU_3]="Установить только Ноду AURORA"
LANG[MENU_4]="Добавить ноду к существующей панели"
LANG[MENU_5]="Управление панелью/нодой"
LANG[MENU_6]="SSL Сертификаты"
LANG[MENU_7]="Бэкап и восстановление"
LANG[MENU_8]="Обновить AURORA"
LANG[MENU_9]="Удалить AURORA"

# Меню установки
LANG[INSTALL_MENU_TITLE]="Установка компонентов AURORA"
LANG[INSTALL_PANEL_NODE]="Установить Панель + Ноду (один сервер)"
LANG[INSTALL_PANEL]="Установить только Панель"
LANG[INSTALL_NODE]="Установить только Ноду"
LANG[INSTALL_ADD_NODE]="Добавить ноду к панели"
LANG[INSTALL_PROMPT]="Выберите тип установки (0-4):"

# Запросы
LANG[ENTER_DOMAIN]="Введите домен для панели (например, panel.example.com): "
LANG[ENTER_EMAIL]="Введите email для Let's Encrypt: "
LANG[ENTER_DB_PASSWORD]="Введите пароль PostgreSQL (пусто = сгенерировать): "
LANG[ENTER_JWT_SECRET]="Введите JWT секрет (пусто = сгенерировать): "
LANG[ENTER_NODE_HOST]="Введите IP-адрес или хостнейм ноды: "
LANG[ENTER_NODE_NAME]="Введите отображаемое имя ноды: "
LANG[ENTER_PANEL_URL]="Введите URL панели для подключения ноды: "
LANG[ENTER_API_KEY]="Введите API ключ ноды (из настроек ноды в панели): "
LANG[GENERATING_CONFIG]="Генерация конфигурации..."
LANG[STARTING_SERVICES]="Запуск сервисов AURORA..."
LANG[INSTALL_COMPLETE]="╔══════════════════════════════════════╗
║     Установка завершена!               ║
╚══════════════════════════════════════╝"
LANG[PANEL_ACCESS]="Панель доступна по адресу: https://%s"
LANG[ADMIN_CREDENTIALS]="Логин администратора: %s / пароль: %s"
LANG[NODE_REGISTERED]="Нода зарегистрирована. Добавьте её в интерфейсе панели в разделе Nodes."

# Управление
LANG[MANAGE_TITLE]="Управление AURORA"
LANG[MANAGE_LOGS]="Просмотр логов"
LANG[MANAGE_RESTART]="Перезапустить сервисы"
LANG[MANAGE_STATUS]="Статус сервисов"
LANG[MANAGE_MIGRATE]="Запустить миграции БД"

# SSL
LANG[SSL_TITLE]="SSL Сертификаты"
LANG[SSL_RENEW]="Обновить все сертификаты"
LANG[SSL_CHECK]="Проверить срок действия сертификатов"
LANG[SSL_NEW]="Выпустить новый сертификат"

# Бэкап
LANG[BACKUP_TITLE]="Бэкап и восстановление"
LANG[BACKUP_CREATE]="Создать бэкап"
LANG[BACKUP_RESTORE]="Восстановить из бэкапа"
LANG[BACKUP_FILE_PATH]="Введите путь к файлу бэкапа: "

# Обновление
LANG[UPDATE_TITLE]="Обновление AURORA"
LANG[UPDATE_PULL]="Загрузка актуальных образов..."
LANG[UPDATE_DONE]="AURORA обновлена до последней версии"

# Удаление
LANG[REMOVE_TITLE]="Удаление AURORA"
LANG[REMOVE_WARNING]="Будут удалены ВСЕ данные AURORA. Отменить будет невозможно."
LANG[REMOVE_CONFIRM]="Введите 'yes' для подтверждения: "
LANG[REMOVE_CANCELLED]="Удаление отменено."
LANG[REMOVE_DONE]="AURORA удалена."

# Алиас
LANG[ALIAS_ADDED]="Алиас 'aurora' добавлен. Откройте новый терминал или выполните 'source %s' для активации."
