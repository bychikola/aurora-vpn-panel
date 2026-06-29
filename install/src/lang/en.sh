#!/bin/bash
# AURORA Installer — English translations

declare -gA LANG

LANG[CHOOSE_LANG]="Select language / Выберите язык:"
LANG[LANG_EN]="English"
LANG[LANG_RU]="Русский"
LANG[WELCOME]="AURORA VPN Panel — Installer"
LANG[MENU_TITLE]="AURORA VPN PANEL"
LANG[VERSION_LABEL]="Version: %s"
LANG[AVAILABLE_UPDATE]="update available"
LANG[SELECT_ACTION]="Select action (0-9):"
LANG[EXIT]="Exit"
LANG[INVALID_CHOICE]="Invalid choice. Please select 0-9."
LANG[BACK]="« Back"
LANG[WAITING]="Please wait..."
LANG[DONE]="Done!"
LANG[ERROR_ROOT]="This script must be run as root"
LANG[ERROR_OS]="Supported only: Debian 11/12, Ubuntu 22.04/24.04"
LANG[ERROR_DOCKER]="Docker is required but not installed"
LANG[CONTINUE_PROMPT]="Continue? (y/n):"
LANG[CONFIRM_YES]="y"

# Menu items
LANG[MENU_1]="Install AURORA (Panel + Node on one server)"
LANG[MENU_2]="Install AURORA Panel only"
LANG[MENU_3]="Install AURORA Node only"
LANG[MENU_4]="Add node to existing panel"
LANG[MENU_5]="Manage panel/node"
LANG[MENU_6]="SSL Certificates"
LANG[MENU_7]="Backup & Restore"
LANG[MENU_8]="Update AURORA"
LANG[MENU_9]="Remove AURORA"

# Install menu
LANG[INSTALL_MENU_TITLE]="Install AURORA Components"
LANG[INSTALL_PANEL_NODE]="Install Panel + Node (single server)"
LANG[INSTALL_PANEL]="Install Panel only"
LANG[INSTALL_NODE]="Install Node only"
LANG[INSTALL_ADD_NODE]="Add node to existing panel"
LANG[INSTALL_PROMPT]="Select installation type (0-4):"

# Prompts
LANG[ENTER_DOMAIN]="Enter domain for panel (e.g. panel.example.com): "
LANG[ENTER_EMAIL]="Enter email for Let's Encrypt: "
LANG[ENTER_DB_PASSWORD]="Enter PostgreSQL password (empty = auto-generate): "
LANG[ENTER_JWT_SECRET]="Enter JWT secret (empty = auto-generate): "
LANG[ENTER_NODE_HOST]="Enter node IP address or hostname: "
LANG[ENTER_NODE_NAME]="Enter node display name: "
LANG[ENTER_PANEL_URL]="Enter panel URL for node to connect to: "
LANG[ENTER_API_KEY]="Enter node API key (from panel node settings): "
LANG[GENERATING_CONFIG]="Generating configuration..."
LANG[STARTING_SERVICES]="Starting AURORA services..."
LANG[INSTALL_COMPLETE]="╔══════════════════════════════════════╗
║     Installation complete!            ║
╚══════════════════════════════════════╝"
LANG[PANEL_ACCESS]="Panel URL: https://%s"
LANG[ADMIN_CREDENTIALS]="Admin login: %s / password: %s"
LANG[NODE_REGISTERED]="Node registered. Add it in panel UI under Nodes section."

# Manage
LANG[MANAGE_TITLE]="Manage AURORA"
LANG[MANAGE_LOGS]="View logs"
LANG[MANAGE_RESTART]="Restart services"
LANG[MANAGE_STATUS]="View status"
LANG[MANAGE_MIGRATE]="Run database migrations"

# SSL
LANG[SSL_TITLE]="SSL Certificates"
LANG[SSL_RENEW]="Renew all certificates"
LANG[SSL_CHECK]="Check certificate expiry"
LANG[SSL_NEW]="Issue new certificate"

# Backup
LANG[BACKUP_TITLE]="Backup & Restore"
LANG[BACKUP_CREATE]="Create backup"
LANG[BACKUP_RESTORE]="Restore from backup"
LANG[BACKUP_FILE_PATH]="Enter backup file path: "

# Update
LANG[UPDATE_TITLE]="Update AURORA"
LANG[UPDATE_PULL]="Pulling latest images..."
LANG[UPDATE_DONE]="AURORA updated to latest version"

# Remove
LANG[REMOVE_TITLE]="Remove AURORA"
LANG[REMOVE_WARNING]="This will remove ALL AURORA data. This cannot be undone."
LANG[REMOVE_CONFIRM]="Type 'yes' to confirm: "
LANG[REMOVE_CANCELLED]="Removal cancelled."
LANG[REMOVE_DONE]="AURORA has been removed."

# Alias
LANG[ALIAS_ADDED]="Alias 'aurora' added. Run new terminal or 'source %s' to activate."
