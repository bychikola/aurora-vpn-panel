#!/bin/bash
#
# AURORA VPN Panel — Unified Installer
# https://github.com/bychikola/aurora-vpn-panel
#

SCRIPT_VERSION="1.0.0"
INSTALL_DIR="/usr/local/aurora"
LANG_FILE="${INSTALL_DIR}/selected_language"
SCRIPT_URL="https://raw.githubusercontent.com/bychikola/aurora-vpn-panel/master/install/install_aurora.sh"
LANG_BASE_URL="https://raw.githubusercontent.com/bychikola/aurora-vpn-panel/master/install/src/lang"
MODULE_BASE_URL="https://raw.githubusercontent.com/bychikola/aurora-vpn-panel/master/install/src"

COLOR_RESET="\033[0m"
COLOR_GREEN="\033[1;32m"
COLOR_YELLOW="\033[1;33m"
COLOR_WHITE="\033[1;37m"
COLOR_RED="\033[1;31m"
COLOR_GRAY='\033[0;90m'
COLOR_CYAN='\033[1;36m'

# ═══════════════════════════════════════════════
# Language defaults (fallback English)
# ═══════════════════════════════════════════════
declare -gA LANG=(
    [CHOOSE_LANG]="Select language / Выберите язык:"
    [LANG_EN]="English"
    [LANG_RU]="Русский"
    [WELCOME]="AURORA VPN Panel — Installer"
    [SELECT_ACTION]="Select action:"
    [EXIT]="Exit"
    [INVALID_CHOICE]="Invalid choice"
    [BACK]="Back"
    [CONTINUE_PROMPT]="Continue? (y/n):"
    [CONFIRM_YES]="y"
    [WAITING]="Please wait..."
    [DONE]="Done!"
    [ERROR_ROOT]="This script must be run as root"
    [ERROR_OS]="Supported only: Debian 11/12, Ubuntu 22.04/24.04"
    [ERROR_DOCKER]="Docker is required but not installed"
)

# ═══════════════════════════════════════════════
# Helpers
# ═══════════════════════════════════════════════

question() { echo -e "${COLOR_GREEN}[?]${COLOR_RESET} ${COLOR_YELLOW}$*${COLOR_RESET}"; }
info() { echo -e "${COLOR_CYAN}[i]${COLOR_RESET} ${COLOR_WHITE}$*${COLOR_RESET}"; }
success() { echo -e "${COLOR_GREEN}[✓]${COLOR_RESET} ${COLOR_GREEN}$*${COLOR_RESET}"; }
warning() { echo -e "${COLOR_YELLOW}[!]${COLOR_RESET} ${COLOR_YELLOW}$*${COLOR_RESET}"; }
error() { echo -e "${COLOR_RED}[✗]${COLOR_RESET} ${COLOR_RED}$*${COLOR_RESET}"; exit 1; }
reading() { read -rp "$(question "$1")" "$2"; }

spinner() {
    local pid=$1; local text=$2
    local spinstr='⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏'
    while kill -0 "$pid" 2>/dev/null; do
        for ((i=0; i<${#spinstr}; i++)); do
            printf "\r${COLOR_GREEN}[%s]${COLOR_RESET} %s" "${spinstr:$i:1}" "$text"
            sleep 0.1
        done
    done
    printf "\r\033[K"
}

download_file() {
    local url="$1"; local dest="$2"; local label="${3:-file}"
    info "Downloading $label..."
    if command -v curl &>/dev/null; then
        curl -fsSL "$url" -o "$dest" 2>/dev/null && return 0
        # jsDelivr fallback
        local gh_path="${url#https://raw.githubusercontent.com/*/}"
        local jsdelivr_url="https://cdn.jsdelivr.net/gh/bychikola/aurora-vpn-panel@master/${gh_path#master/}"
        curl -fsSL "$jsdelivr_url" -o "$dest" 2>/dev/null && return 0
    elif command -v wget &>/dev/null; then
        wget -q "$url" -O "$dest" 2>/dev/null && return 0
    fi
    return 1
}

# ═══════════════════════════════════════════════
# Language loading
# ═══════════════════════════════════════════════

load_language() {
    if [ -f "$LANG_FILE" ]; then
        local saved=$(cat "$LANG_FILE")
        case $saved in
            1) set_language en; return 0 ;;
            2) set_language ru; return 0 ;;
        esac
    fi
    return 1
}

set_language() {
    local lang="$1"; local force="${2:-false}"
    local lang_file="${INSTALL_DIR}/lang/${lang}.sh"

    unset LANG; declare -gA LANG

    if [ "$force" = true ] || [ ! -f "$lang_file" ]; then
        mkdir -p "${INSTALL_DIR}/lang"
        download_file "${LANG_BASE_URL}/${lang}.sh" "$lang_file" "language pack" || {
            warning "Failed to download language pack, using English fallback"
            _set_fallback_en
        }
    fi

    if [ -f "$lang_file" ]; then
        source "$lang_file"
    else
        _set_fallback_en
    fi
}

_set_fallback_en() {
    declare -gA LANG=(
        [CHOOSE_LANG]="Select language / Выберите язык:"
        [LANG_EN]="English"
        [LANG_RU]="Русский"
        [WELCOME]="AURORA VPN Panel — Installer"
        [MENU_TITLE]="AURORA VPN PANEL"
        [VERSION_LABEL]="Version: %s"
        [AVAILABLE_UPDATE]="update available"
        [SELECT_ACTION]="Select action (0-9):"
        [EXIT]="Exit"
        [INVALID_CHOICE]="Invalid choice"
        [BACK]="« Back"
        [WAITING]="Please wait..."
        [DONE]="Done!"
        [ERROR_ROOT]="This script must be run as root"
        [ERROR_OS]="Supported only: Debian 11/12, Ubuntu 22.04/24.04"
        [ERROR_DOCKER]="Docker is required but not installed"
        [CONTINUE_PROMPT]="Continue? (y/n):"
        [CONFIRM_YES]="y"
        [MENU_1]="Install AURORA (Panel + Node on one server)"
        [MENU_2]="Install AURORA Panel only"
        [MENU_3]="Install AURORA Node only"
        [MENU_4]="Add node to existing panel"
        [MENU_5]="Manage panel/node"
        [MENU_6]="SSL certificates"
        [MENU_7]="Backup & Restore"
        [MENU_8]="Update AURORA"
        [MENU_9]="Remove AURORA"
        [INSTALL_MENU_TITLE]="Install AURORA Components"
        [INSTALL_PANEL_NODE]="Install Panel + Node (single server)"
        [INSTALL_PANEL]="Install Panel only"
        [INSTALL_NODE]="Install Node only"
        [INSTALL_ADD_NODE]="Add node to existing panel"
        [INSTALL_PROMPT]="Select installation type (0-4):"
        [ENTER_DOMAIN]="Enter domain for panel (e.g. panel.example.com):"
        [ENTER_EMAIL]="Enter email for Let's Encrypt:"
        [ENTER_DB_PASSWORD]="Enter PostgreSQL password (leave empty to generate):"
        [ENTER_JWT_SECRET]="Enter JWT secret (leave empty to generate):"
        [ENTER_NODE_HOST]="Enter node IP address or hostname:"
        [ENTER_NODE_NAME]="Enter node display name:"
        [ENTER_PANEL_URL]="Enter panel URL for node to connect to:"
        [ENTER_API_KEY]="Enter node API key (from panel node settings):"
        [GENERATING_CONFIG]="Generating configuration..."
        [STARTING_SERVICES]="Starting AURORA services..."
        [INSTALL_COMPLETE]="Installation complete!"
        [PANEL_ACCESS]="Panel URL: https://%s"
        [ADMIN_CREDENTIALS]="Admin login: %s / password: %s"
        [NODE_REGISTERED]="Node registered. Add it in panel UI."
        [SELECT_WEBSERVER]="Select reverse proxy:"
        [WEBSERVER_NGINX]="1. Nginx + Let's Encrypt"
        [WEBSERVER_CADDY]="2. Caddy (auto-SSL)"
    )
}

show_language_menu() {
    clear
    echo -e ""
    echo -e "${COLOR_CYAN}    ╔══════════════════════════════════╗${COLOR_RESET}"
    echo -e "${COLOR_CYAN}    ║     ${COLOR_GREEN}AURORA VPN Panel${COLOR_CYAN}           ║${COLOR_RESET}"
    echo -e "${COLOR_CYAN}    ║     ${COLOR_WHITE}Unified Installer${COLOR_CYAN}           ║${COLOR_RESET}"
    echo -e "${COLOR_CYAN}    ╚══════════════════════════════════╝${COLOR_RESET}"
    echo -e ""
    echo -e "${COLOR_GREEN}${LANG[CHOOSE_LANG]}${COLOR_RESET}"
    echo -e ""
    echo -e "    ${COLOR_YELLOW}1. ${COLOR_WHITE}English${COLOR_RESET}"
    echo -e "    ${COLOR_YELLOW}2. ${COLOR_WHITE}Русский${COLOR_RESET}"
    echo -e ""
    reading "Select (1-2): " LANG_CHOICE

    case $LANG_CHOICE in
        1) echo "1" > "$LANG_FILE"; set_language en ;;
        2) echo "2" > "$LANG_FILE"; set_language ru ;;
        *) echo "1" > "$LANG_FILE"; set_language en ;;
    esac
}

# ═══════════════════════════════════════════════
# System checks
# ═══════════════════════════════════════════════

check_root() {
    if [[ $EUID -ne 0 ]]; then
        echo -e "${COLOR_RED}${LANG[ERROR_ROOT]}${COLOR_RESET}"
        exit 1
    fi
}

check_os() {
    if grep -qE "bullseye|bookworm|trixie" /etc/os-release 2>/dev/null; then
        OS="debian"; return 0
    elif grep -qE "jammy|noble" /etc/os-release 2>/dev/null; then
        OS="ubuntu"; return 0
    fi
    error "${LANG[ERROR_OS]}"
}

install_dependencies() {
    if [ -f "${INSTALL_DIR}/deps_installed" ]; then
        return 0
    fi

    info "Installing system dependencies..."
    apt-get update -y -qq >/dev/null 2>&1
    apt-get install -y -qq ca-certificates curl jq wget gnupg ufw \
        dnsutils git cron certbot python3-certbot-dns-cloudflare \
        unattended-upgrades >/dev/null 2>&1

    # Docker
    if ! command -v docker &>/dev/null; then
        info "Installing Docker..."
        curl -fsSL https://get.docker.com | sh >/dev/null 2>&1
        systemctl start docker >/dev/null 2>&1
        systemctl enable docker >/dev/null 2>&1
    fi

    # UFW
    ufw allow 22/tcp comment 'SSH' >/dev/null 2>&1
    ufw allow 443/tcp comment 'HTTPS' >/dev/null 2>&1
    ufw --force enable >/dev/null 2>&1

    # BBR
    if ! grep -q "net.core.default_qdisc = fq" /etc/sysctl.conf 2>/dev/null; then
        echo "net.core.default_qdisc = fq" >> /etc/sysctl.conf
        echo "net.ipv4.tcp_congestion_control = bbr" >> /etc/sysctl.conf
        sysctl -p >/dev/null 2>&1
    fi

    touch "${INSTALL_DIR}/deps_installed"
    success "Dependencies installed"
}

# ═══════════════════════════════════════════════
# Password / secret generation
# ═══════════════════════════════════════════════

generate_password() {
    local length="${1:-24}"
    tr -dc 'A-Za-z0-9!@#%^&*()_+' < /dev/urandom | head -c "$length"
}

generate_hex32() {
    tr -dc 'a-f0-9' < /dev/urandom | head -c 64
}

# ═══════════════════════════════════════════════
# Domain + SSL helpers
# ═══════════════════════════════════════════════

extract_domain() {
    echo "$1" | awk -F'.' '{if (NF > 2) {print $(NF-1)"."$NF} else {print $0}}'
}

check_domain() {
    local domain="$1"
    local domain_ip=$(dig +short A "$domain" 2>/dev/null | grep -E '^[0-9]+\.[0-9]+\.[0-9]+\.[0-9]+$' | head -1)
    local server_ip=$(curl -s4 ifconfig.me 2>/dev/null || curl -s4 ipinfo.io/ip 2>/dev/null)

    if [ -z "$domain_ip" ] || [ -z "$server_ip" ]; then
        warning "Cannot resolve domain IP. Make sure DNS A record points to $server_ip"
        reading "${LANG[CONTINUE_PROMPT]}" confirm
        [[ "$confirm" != "${LANG[CONFIRM_YES]}" ]] && return 2
        return 0
    fi

    if [ "$domain_ip" = "$server_ip" ]; then
        success "Domain $domain resolves to this server ($server_ip)"
        return 0
    else
        warning "Domain $domain resolves to $domain_ip, but server IP is $server_ip"
        warning "This may be Cloudflare proxy — make sure SSL mode is Full"
        reading "${LANG[CONTINUE_PROMPT]}" confirm
        [[ "$confirm" != "${LANG[CONFIRM_YES]}" ]] && return 2
        return 0
    fi
}

obtain_ssl() {
    local domain="$1"; local email="$2"

    info "Obtaining SSL certificate for $domain..."
    reading "Use Cloudflare DNS API? (y/n): " use_cf

    if [[ "$use_cf" == "y" || "$use_cf" == "Y" ]]; then
        reading "Enter Cloudflare API token (or global key): " CF_API_KEY
        reading "Enter Cloudflare email: " CF_EMAIL

        mkdir -p ~/.secrets/certbot
        if [[ "$CF_API_KEY" =~ [A-Z] ]]; then
            cat > ~/.secrets/certbot/cloudflare.ini <<EOF
dns_cloudflare_api_token = $CF_API_KEY
EOF
        else
            cat > ~/.secrets/certbot/cloudflare.ini <<EOF
dns_cloudflare_email = $CF_EMAIL
dns_cloudflare_api_key = $CF_API_KEY
EOF
        fi
        chmod 600 ~/.secrets/certbot/cloudflare.ini

        local base_domain=$(extract_domain "$domain")
        certbot certonly --dns-cloudflare \
            --dns-cloudflare-credentials ~/.secrets/certbot/cloudflare.ini \
            --dns-cloudflare-propagation-seconds 60 \
            -d "$base_domain" -d "*.$base_domain" \
            --email "$email" --agree-tos --non-interactive \
            --key-type ecdsa --elliptic-curve secp384r1 2>&1 &
    else
        ufw allow 80/tcp comment 'HTTP ACME' >/dev/null 2>&1
        certbot certonly --standalone \
            -d "$domain" \
            --email "$email" --agree-tos --non-interactive \
            --http-01-port 80 \
            --key-type ecdsa --elliptic-curve secp384r1 2>&1 &
        ufw delete allow 80/tcp >/dev/null 2>&1
    fi

    spinner $! "Obtaining SSL certificate..."
    wait $!

    if [ -d "/etc/letsencrypt/live/$domain" ]; then
        success "SSL certificate obtained for $domain"
        return 0
    else
        error "Failed to obtain SSL certificate"
    fi
}

# ═══════════════════════════════════════════════
# Module loading from remote
# ═══════════════════════════════════════════════

load_module() {
    local module="$1"; local category="${2:-modules}"
    local module_file="${INSTALL_DIR}/${category}/${module}.sh"
    local module_url="${MODULE_BASE_URL}/${category}/${module}.sh"

    if [ -f "$module_file" ]; then
        source "$module_file"
        return 0
    fi

    mkdir -p "${INSTALL_DIR}/${category}"
    if download_file "$module_url" "$module_file" "${category}/${module}.sh"; then
        source "$module_file"
        return 0
    fi
    error "Failed to load module: $module"
}

# ═══════════════════════════════════════════════
# Menus
# ═══════════════════════════════════════════════

show_menu() {
    clear
    echo -e ""
    echo -e "${COLOR_CYAN}    ╔══════════════════════════════════╗${COLOR_RESET}"
    echo -e "${COLOR_CYAN}    ║     ${COLOR_GREEN}AURORA VPN Panel${COLOR_CYAN}           ║${COLOR_RESET}"
    echo -e "${COLOR_CYAN}    ║     ${COLOR_WHITE}Unified Installer v${SCRIPT_VERSION}${COLOR_CYAN}   ║${COLOR_RESET}"
    echo -e "${COLOR_CYAN}    ╚══════════════════════════════════╝${COLOR_RESET}"
    echo -e ""
    echo -e "${COLOR_GRAY}github.com/bychikola/aurora-vpn-panel${COLOR_RESET}"
    echo -e ""
    echo -e "  ${COLOR_YELLOW}1.${COLOR_RESET} ${LANG[MENU_1]}"
    echo -e "  ${COLOR_YELLOW}2.${COLOR_RESET} ${LANG[MENU_2]}"
    echo -e "  ${COLOR_YELLOW}3.${COLOR_RESET} ${LANG[MENU_3]}"
    echo -e "  ${COLOR_YELLOW}4.${COLOR_RESET} ${LANG[MENU_4]}"
    echo -e "  ${COLOR_YELLOW}5.${COLOR_RESET} ${LANG[MENU_5]}"
    echo -e "  ${COLOR_YELLOW}6.${COLOR_RESET} ${LANG[MENU_6]}"
    echo -e "  ${COLOR_YELLOW}7.${COLOR_RESET} ${LANG[MENU_7]}"
    echo -e "  ${COLOR_YELLOW}8.${COLOR_RESET} ${LANG[MENU_8]}"
    echo -e "  ${COLOR_YELLOW}9.${COLOR_RESET} ${LANG[MENU_9]}"
    echo -e ""
    echo -e "  ${COLOR_YELLOW}0.${COLOR_RESET} ${LANG[EXIT]}"
    echo -e ""
}

install_aurora_panel_node() {
    clear
    echo -e ""
    echo -e "${COLOR_GREEN}═══ ${LANG[INSTALL_PANEL_NODE]} ═══${COLOR_RESET}"
    echo -e ""
    warning "This installs both Panel and Node on a single server."
    reading "${LANG[CONTINUE_PROMPT]}" confirm
    [[ "$confirm" != "${LANG[CONFIRM_YES]}" ]] && return 0

    load_module "install_panel" "modules"
    load_module "install_node" "modules"
    install_dependencies

    reading "${LANG[ENTER_DOMAIN]}" PANEL_DOMAIN
    check_domain "$PANEL_DOMAIN" || return 0
    reading "${LANG[ENTER_EMAIL]}" LETSENCRYPT_EMAIL
    obtain_ssl "$PANEL_DOMAIN" "$LETSENCRYPT_EMAIL"

    DB_PASSWORD=$(generate_password 24)
    JWT_ACCESS=$(generate_hex32)
    JWT_REFRESH=$(generate_hex32)
    ADMIN_PASSWORD=$(generate_password 16)

    info "${LANG[GENERATING_CONFIG]}"
    # Call panel installation function from loaded module
    _do_install_panel "$PANEL_DOMAIN" "$DB_PASSWORD" "$JWT_ACCESS" "$JWT_REFRESH" "$ADMIN_PASSWORD"

    info "Installing node alongside panel..."
    _do_install_node "localhost" "10085" "localhost"

    info "${LANG[STARTING_SERVICES]}"
    cd /opt/aurora && docker compose up -d 2>&1 &
    spinner $! "${LANG[WAITING]}"
    wait $!

    echo -e ""
    success "${LANG[INSTALL_COMPLETE]}"
    printf "${LANG[PANEL_ACCESS]}\n" "$PANEL_DOMAIN"
    printf "${LANG[ADMIN_CREDENTIALS]}\n" "admin" "$ADMIN_PASSWORD"
    echo -e ""
    reading "Press Enter to return to menu..." _
}

install_aurora_panel() {
    clear
    echo -e ""
    echo -e "${COLOR_GREEN}═══ ${LANG[INSTALL_PANEL]} ═══${COLOR_RESET}"
    echo -e ""

    load_module "install_panel" "modules"
    install_dependencies

    reading "${LANG[ENTER_DOMAIN]}" PANEL_DOMAIN
    check_domain "$PANEL_DOMAIN" || return 0
    reading "${LANG[ENTER_EMAIL]}" LETSENCRYPT_EMAIL
    obtain_ssl "$PANEL_DOMAIN" "$LETSENCRYPT_EMAIL"

    DB_PASSWORD=$(generate_password 24)
    JWT_ACCESS=$(generate_hex32)
    JWT_REFRESH=$(generate_hex32)
    ADMIN_PASSWORD=$(generate_password 16)

    info "${LANG[GENERATING_CONFIG]}"
    _do_install_panel "$PANEL_DOMAIN" "$DB_PASSWORD" "$JWT_ACCESS" "$JWT_REFRESH" "$ADMIN_PASSWORD"

    info "${LANG[STARTING_SERVICES]}"
    cd /opt/aurora && docker compose up -d 2>&1 &
    spinner $! "${LANG[WAITING]}"
    wait $!

    echo -e ""
    success "${LANG[INSTALL_COMPLETE]}"
    printf "${LANG[PANEL_ACCESS]}\n" "$PANEL_DOMAIN"
    printf "${LANG[ADMIN_CREDENTIALS]}\n" "admin" "$ADMIN_PASSWORD"
    echo -e ""
    reading "Press Enter to return to menu..." _
}

install_aurora_node() {
    clear
    echo -e ""
    echo -e "${COLOR_GREEN}═══ ${LANG[INSTALL_NODE]} ═══${COLOR_RESET}"
    echo -e ""

    load_module "install_node" "modules"
    install_dependencies

    reading "${LANG[ENTER_NODE_NAME]}" NODE_NAME
    reading "${LANG[ENTER_PANEL_URL]}" PANEL_URL

    info "${LANG[GENERATING_CONFIG]}"
    _do_install_node "$NODE_NAME" "10085" "$PANEL_URL"

    info "${LANG[STARTING_SERVICES]}"
    cd /opt/aurora-node && docker compose up -d 2>&1 &
    spinner $! "${LANG[WAITING]}"
    wait $!

    echo -e ""
    success "${LANG[INSTALL_COMPLETE]}"
    success "${LANG[NODE_REGISTERED]}"
    echo -e ""
    reading "Press Enter to return to menu..." _
}

add_node_to_panel() {
    clear
    echo -e ""
    echo -e "${COLOR_GREEN}═══ ${LANG[INSTALL_ADD_NODE]} ═══${COLOR_RESET}"
    echo -e ""
    warning "Run this on the node server, not the panel server!"

    load_module "add_node" "modules"
    reading "${LANG[ENTER_NODE_NAME]}" NODE_NAME
    reading "${LANG[ENTER_PANEL_URL]}" PANEL_URL
    reading "${LANG[ENTER_API_KEY]}" API_KEY

    _do_add_node "$NODE_NAME" "$PANEL_URL" "$API_KEY"

    success "Node added. Complete setup in panel UI."
    reading "Press Enter to return to menu..." _
}

manage_panel() {
    clear
    echo -e ""
    echo -e "${COLOR_GREEN}═══ ${LANG[MENU_5]} ═══${COLOR_RESET}"
    echo -e ""
    echo -e "  ${COLOR_YELLOW}1.${COLOR_RESET} View logs"
    echo -e "  ${COLOR_YELLOW}2.${COLOR_RESET} Restart services"
    echo -e "  ${COLOR_YELLOW}3.${COLOR_RESET} View status"
    echo -e "  ${COLOR_YELLOW}4.${COLOR_RESET} Run migrations"
    echo -e ""
    echo -e "  ${COLOR_YELLOW}0.${COLOR_RESET} ${LANG[BACK]}"
    echo -e ""
    reading "${LANG[SELECT_ACTION]}" MGMT_OPTION

    case $MGMT_OPTION in
        1)
            if [ -d "/opt/aurora" ]; then
                cd /opt/aurora && docker compose logs --tail=100 -f
            else
                warning "Panel not found in /opt/aurora"
            fi
            ;;
        2)
            if [ -d "/opt/aurora" ]; then
                cd /opt/aurora && docker compose restart
                success "Services restarted"
            fi
            ;;
        3)
            if [ -d "/opt/aurora" ]; then
                cd /opt/aurora && docker compose ps
            fi
            ;;
        4)
            if [ -d "/opt/aurora" ]; then
                cd /opt/aurora && docker compose exec backend ./aurora migrate up
                success "Migrations applied"
            fi
            ;;
        0) return 0 ;;
    esac
    reading "Press Enter to return to menu..." _
}

manage_ssl() {
    clear
    echo -e ""
    echo -e "${COLOR_GREEN}═══ ${LANG[MENU_6]} ═══${COLOR_RESET}"
    echo -e ""
    echo -e "  ${COLOR_YELLOW}1.${COLOR_RESET} Renew all certificates"
    echo -e "  ${COLOR_YELLOW}2.${COLOR_RESET} Check certificate expiry"
    echo -e "  ${COLOR_YELLOW}3.${COLOR_RESET} Issue new certificate"
    echo -e ""
    echo -e "  ${COLOR_YELLOW}0.${COLOR_RESET} ${LANG[BACK]}"
    echo -e ""
    reading "${LANG[SELECT_ACTION]}" SSL_OPTION

    case $SSL_OPTION in
        1)
            certbot renew --no-random-sleep-on-renew 2>&1 &
            spinner $! "Renewing certificates..."
            wait $!
            success "Renewal complete"
            ;;
        2)
            for d in /etc/letsencrypt/live/*; do
                [ -d "$d" ] && echo -e "$(basename "$d"): $(openssl x509 -enddate -noout -in "$d/fullchain.pem" 2>/dev/null)"
            done
            ;;
        3)
            reading "${LANG[ENTER_DOMAIN]}" CERT_DOMAIN
            reading "${LANG[ENTER_EMAIL]}" CERT_EMAIL
            obtain_ssl "$CERT_DOMAIN" "$CERT_EMAIL"
            ;;
        0) return 0 ;;
    esac
    reading "Press Enter to return to menu..." _
}

backup_restore() {
    clear
    echo -e ""
    echo -e "${COLOR_GREEN}═══ ${LANG[MENU_7]} ═══${COLOR_RESET}"
    echo -e ""
    echo -e "  ${COLOR_YELLOW}1.${COLOR_RESET} Backup panel data"
    echo -e "  ${COLOR_YELLOW}2.${COLOR_RESET} Restore from backup"
    echo -e ""
    echo -e "  ${COLOR_YELLOW}0.${COLOR_RESET} ${LANG[BACK]}"
    echo -e ""
    reading "${LANG[SELECT_ACTION]}" BAK_OPTION

    case $BAK_OPTION in
        1)
            local bak_file="/root/aurora_backup_$(date +%Y%m%d_%H%M%S).tar.gz"
            info "Creating backup: $bak_file"
            if [ -d "/opt/aurora" ]; then
                tar czf "$bak_file" -C /opt aurora 2>/dev/null
                success "Backup saved to $bak_file"
            else
                warning "No panel data found"
            fi
            ;;
        2)
            reading "Enter backup file path: " RESTORE_FILE
            if [ -f "$RESTORE_FILE" ]; then
                cd /opt/aurora && docker compose down 2>/dev/null
                tar xzf "$RESTORE_FILE" -C /opt 2>/dev/null
                cd /opt/aurora && docker compose up -d 2>&1 &
                spinner $! "Restoring..."
                wait $!
                success "Restored from $RESTORE_FILE"
            else
                warning "Backup file not found"
            fi
            ;;
        0) return 0 ;;
    esac
    reading "Press Enter to return to menu..." _
}

update_aurora() {
    clear
    echo -e ""
    echo -e "${COLOR_GREEN}═══ ${LANG[MENU_8]} ═══${COLOR_RESET}"
    echo -e ""

    if [ -d "/opt/aurora" ]; then
        info "Pulling latest Docker images..."
        cd /opt/aurora && docker compose pull 2>&1 &
        spinner $! "Updating containers..."
        wait $!
        cd /opt/aurora && docker compose up -d 2>&1 &
        spinner $! "Starting..."
        wait $!
        success "AURORA updated to latest version"
    else
        warning "Panel not installed"
    fi
    reading "Press Enter to return to menu..." _
}

remove_aurora() {
    clear
    echo -e ""
    echo -e "${COLOR_RED}═══ ${LANG[MENU_9]} ═══${COLOR_RESET}"
    echo -e ""
    warning "This will remove ALL AURORA data. This cannot be undone."
    reading "Type 'yes' to confirm: " confirm
    if [[ "$confirm" != "yes" ]]; then
        info "Cancelled."
        return 0
    fi

    for dir in /opt/aurora /opt/aurora-node; do
        if [ -d "$dir" ]; then
            cd "$dir" && docker compose down -v --rmi all --remove-orphans 2>/dev/null
            rm -rf "$dir"
        fi
    done
    docker system prune -a --volumes -f 2>/dev/null &
    spinner $! "Cleaning..."
    wait $!
    rm -rf "$INSTALL_DIR" 2>/dev/null
    rm -f /usr/local/bin/aurora 2>/dev/null

    success "AURORA removed."
    reading "Press Enter to exit..." _
    exit 0
}

# ═══════════════════════════════════════════════
# Self-update
# ═══════════════════════════════════════════════

self_update() {
    local remote_version=$(curl -fsSL "$SCRIPT_URL" 2>/dev/null | grep -m1 "SCRIPT_VERSION=" | cut -d'"' -f2)
    if [ -z "$remote_version" ]; then
        warning "Cannot check for updates"
        return 1
    fi
    if [ "$remote_version" = "$SCRIPT_VERSION" ]; then
        success "Already up to date (v$SCRIPT_VERSION)"
        return 0
    fi

    info "New version available: v$remote_version (current: v$SCRIPT_VERSION)"
    reading "Update? (y/n): " upd_confirm
    [[ "$upd_confirm" != "y" ]] && return 0

    local tmp_script="${INSTALL_DIR}/aurora_installer.tmp"
    download_file "$SCRIPT_URL" "$tmp_script" "new installer version" || {
        error "Failed to download update"
    }
    mv "$tmp_script" "${INSTALL_DIR}/aurora_installer.sh"
    chmod +x "${INSTALL_DIR}/aurora_installer.sh"
    ln -sf "${INSTALL_DIR}/aurora_installer.sh" /usr/local/bin/aurora

    success "Updated to v$remote_version. Restarting..."
    exec "${INSTALL_DIR}/aurora_installer.sh"
}

# ═══════════════════════════════════════════════
# Install self (first run)
# ═══════════════════════════════════════════════

install_self() {
    if [ ! -f "${INSTALL_DIR}/aurora_installer.sh" ]; then
        mkdir -p "$INSTALL_DIR"
        cp "$0" "${INSTALL_DIR}/aurora_installer.sh"
        chmod +x "${INSTALL_DIR}/aurora_installer.sh"
        ln -sf "${INSTALL_DIR}/aurora_installer.sh" /usr/local/bin/aurora

        # Alias
        local bashrc="/etc/bash.bashrc"
        grep -q "alias aurora=" "$bashrc" 2>/dev/null || \
            echo "alias aurora='aurora'" >> "$bashrc"

        info "Installer saved to /usr/local/bin/aurora"
        info "Run 'aurora' to open this menu again"
    fi
}

# ═══════════════════════════════════════════════
# Main entry point
# ═══════════════════════════════════════════════

main() {
    check_root
    check_os
    install_self

    if ! load_language; then
        show_language_menu
    fi

    while true; do
        show_menu
        reading "${LANG[SELECT_ACTION]}" OPTION
        case $OPTION in
            1) install_aurora_panel_node ;;
            2) install_aurora_panel ;;
            3) install_aurora_node ;;
            4) add_node_to_panel ;;
            5) manage_panel ;;
            6) manage_ssl ;;
            7) backup_restore ;;
            8) update_aurora ;;
            9) remove_aurora ;;
            0) echo -e "${COLOR_GREEN}Goodbye!${COLOR_RESET}"; exit 0 ;;
            *) warning "${LANG[INVALID_CHOICE]}"; sleep 1 ;;
        esac
    done
}

# ─── Run ───
main "$@"
