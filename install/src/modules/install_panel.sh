#!/bin/bash
# AURORA Panel installation module

_do_install_panel() {
    local PANEL_DOMAIN="$1"
    local DB_PASSWORD="$2"
    local JWT_ACCESS="$3"
    local JWT_REFRESH="$4"
    local ADMIN_PASSWORD="$5"
    local WEBSERVER="${6:-nginx}"  # nginx or caddy

    mkdir -p /opt/aurora/{configs,data/postgres,data/redis,ssl}

    # Create .env file
    cat > /opt/aurora/.env <<EOF
# AURORA Panel Environment
DB_PASSWORD=${DB_PASSWORD}
AURORA_DATABASE_HOST=postgres
AURORA_DATABASE_PORT=5432
AURORA_DATABASE_USER=aurora
AURORA_DATABASE_PASSWORD=${DB_PASSWORD}
AURORA_DATABASE_DBNAME=aurora
AURORA_REDIS_HOST=redis
AURORA_REDIS_PORT=6379
AURORA_REDIS_PASSWORD=
AURORA_JWT_ACCESSSECRET=${JWT_ACCESS}
AURORA_JWT_REFRESHSECRET=${JWT_REFRESH}
AURORA_PANEL_DOMAIN=${PANEL_DOMAIN}
AURORA_ADMIN_PASSWORD=${ADMIN_PASSWORD}
EOF

    # Create config.yaml
    cat > /opt/aurora/configs/config.yaml <<EOF
server:
  host: "0.0.0.0"
  port: 8080
  readTimeout: "15s"
  writeTimeout: "30s"
  idleTimeout: "60s"

database:
  host: "postgres"
  port: 5432
  user: "aurora"
  password: "${DB_PASSWORD}"
  dbname: "aurora"
  sslmode: "disable"
  maxConns: 25
  minConns: 5

redis:
  host: "redis"
  port: 6379
  password: ""
  db: 0
  poolSize: 20

jwt:
  accessSecret: "${JWT_ACCESS}"
  refreshSecret: "${JWT_REFRESH}"
  accessTTL: "15m"
  refreshTTL: "168h"

xray:
  grpcTimeout: "10s"
  retryMax: 3
  retryWait: "2s"

log:
  level: "info"
  format: "json"
EOF

    # Copy SSL certificates (for Nginx mode)
    if [ "$WEBSERVER" = "nginx" ]; then
        local cert_dir="/etc/letsencrypt/live/${PANEL_DOMAIN}"
        local base_domain=$(echo "$PANEL_DOMAIN" | awk -F'.' '{if (NF > 2) {print $(NF-1)"."$NF} else {print $0}}')
        if [ ! -d "$cert_dir" ]; then
            cert_dir="/etc/letsencrypt/live/${base_domain}"
        fi
        if [ -d "$cert_dir" ]; then
            cp "$cert_dir/fullchain.pem" /opt/aurora/ssl/fullchain.pem
            cp "$cert_dir/privkey.pem" /opt/aurora/ssl/privkey.pem
        fi
    fi

    # ─── Docker-compose ───
    cat > /opt/aurora/docker-compose.yml <<YAML
services:
  postgres:
    image: postgres:16-alpine
    container_name: aurora-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: aurora
      POSTGRES_USER: aurora
      POSTGRES_PASSWORD: \${DB_PASSWORD}
    volumes:
      - ./data/postgres:/var/lib/postgresql/data
    ports:
      - "127.0.0.1:5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U aurora"]
      interval: 5s
      timeout: 5s
      retries: 5
    networks:
      - aurora-net

  redis:
    image: redis:7-alpine
    container_name: aurora-redis
    restart: unless-stopped
    command: redis-server --appendonly yes --maxmemory 256mb --maxmemory-policy allkeys-lru
    volumes:
      - ./data/redis:/data
    ports:
      - "127.0.0.1:6379:6379"
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5
    networks:
      - aurora-net

  backend:
    image: ghcr.io/bychikola/aurora-backend:latest
    container_name: aurora-backend
    restart: unless-stopped
    env_file:
      - .env
    volumes:
      - ./configs:/app/configs:ro
    ports:
      - "127.0.0.1:8080:8080"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - aurora-net

  frontend:
    image: ghcr.io/bychikola/aurora-frontend:latest
    container_name: aurora-frontend
    restart: unless-stopped
    ports:
      - "127.0.0.1:3000:80"
    depends_on:
      - backend
    networks:
      - aurora-net

YAML

    # ─── Reverse proxy: Caddy or Nginx ───
    if [ "$WEBSERVER" = "caddy" ]; then
        cat >> /opt/aurora/docker-compose.yml <<YAML
  caddy:
    image: caddy:2-alpine
    container_name: aurora-caddy
    restart: unless-stopped
    cap_add:
      - NET_ADMIN
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./Caddyfile:/etc/caddy/Caddyfile:ro
      - ./data/caddy:/data
    depends_on:
      - backend
      - frontend
    networks:
      - aurora-net
YAML

        cat > /opt/aurora/Caddyfile <<EOF
${PANEL_DOMAIN} {
    encode gzip zstd
    handle /api/* {
        reverse_proxy backend:8080
    }
    handle {
        reverse_proxy frontend:80
    }
}
EOF
        info "Caddy configured — auto-SSL will obtain certificate on first request"
    else
        # Nginx
        cat >> /opt/aurora/docker-compose.yml <<YAML
  nginx:
    image: nginx:alpine
    container_name: aurora-nginx
    restart: unless-stopped
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    ports:
      - "80:80"
      - "443:443"
    depends_on:
      - backend
      - frontend
    networks:
      - aurora-net
YAML

        cat > /opt/aurora/nginx.conf <<EOF
events { worker_connections 1024; }

http {
    server {
        listen 80;
        server_name ${PANEL_DOMAIN};
        return 301 https://\$host\$request_uri;
    }

    server {
        listen 443 ssl http2;
        server_name ${PANEL_DOMAIN};

        ssl_certificate /etc/nginx/ssl/fullchain.pem;
        ssl_certificate_key /etc/nginx/ssl/privkey.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305;

        location / {
            proxy_pass http://frontend:80;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }

        location /api/ {
            proxy_pass http://backend:8080;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto \$scheme;
        }
    }
}
EOF

        # SSL renewal cron (only for Nginx)
        local base_domain=$(echo "$PANEL_DOMAIN" | awk -F'.' '{if (NF > 2) {print $(NF-1)"."$NF} else {print $0}}')
        cat > /etc/cron.d/aurora-ssl <<EOF
0 3 * * * root certbot renew --quiet --post-hook "cp /etc/letsencrypt/live/${base_domain}/fullchain.pem /opt/aurora/ssl/fullchain.pem && cp /etc/letsencrypt/live/${base_domain}/privkey.pem /opt/aurora/ssl/privkey.pem && cd /opt/aurora && docker compose restart nginx"
EOF
    fi

    cat >> /opt/aurora/docker-compose.yml <<YAML
networks:
  aurora-net:
    driver: bridge
YAML

    success "Panel configuration written to /opt/aurora/ (webserver: ${WEBSERVER})"
}
