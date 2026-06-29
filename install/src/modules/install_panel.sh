#!/bin/bash
# AURORA Panel installation module

_do_install_panel() {
    local PANEL_DOMAIN="$1"
    local DB_PASSWORD="$2"
    local JWT_ACCESS="$3"
    local JWT_REFRESH="$4"
    local ADMIN_PASSWORD="$5"

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

    # Docker-compose for panel
    cat > /opt/aurora/docker-compose.yml <<'YAML'
services:
  postgres:
    image: postgres:16-alpine
    container_name: aurora-postgres
    restart: unless-stopped
    environment:
      POSTGRES_DB: aurora
      POSTGRES_USER: aurora
      POSTGRES_PASSWORD: ${DB_PASSWORD}
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

  nginx:
    image: nginx:alpine
    container_name: aurora-nginx
    restart: unless-stopped
    volumes:
      - ./nginx.conf:/etc/nginx/nginx.conf:ro
      - ./ssl:/etc/nginx/ssl:ro
    ports:
      - "443:443"
    depends_on:
      - backend
      - frontend
    networks:
      - aurora-net

networks:
  aurora-net:
    driver: bridge
YAML

    # Nginx config
    cat > /opt/aurora/nginx.conf <<EOF
events { worker_connections 1024; }

http {
    server {
        listen 443 ssl http2;
        server_name ${PANEL_DOMAIN};

        ssl_certificate /etc/nginx/ssl/fullchain.pem;
        ssl_certificate_key /etc/nginx/ssl/privkey.pem;
        ssl_protocols TLSv1.2 TLSv1.3;
        ssl_ciphers ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305;

        # Frontend
        location / {
            proxy_pass http://frontend:80;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
        }

        # Backend API
        location /api/ {
            proxy_pass http://backend:8080;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        }
    }
}
EOF

    # Copy SSL certificates
    if [ -d "/etc/letsencrypt/live/${PANEL_DOMAIN}" ]; then
        cp /etc/letsencrypt/live/${PANEL_DOMAIN}/fullchain.pem /opt/aurora/ssl/fullchain.pem
        cp /etc/letsencrypt/live/${PANEL_DOMAIN}/privkey.pem /opt/aurora/ssl/privkey.pem
    fi

    # SSL renewal cron
    cat > /etc/cron.d/aurora-ssl <<EOF
0 3 * * * root certbot renew --quiet --post-hook "cp /etc/letsencrypt/live/${PANEL_DOMAIN}/fullchain.pem /opt/aurora/ssl/fullchain.pem && cp /etc/letsencrypt/live/${PANEL_DOMAIN}/privkey.pem /opt/aurora/ssl/privkey.pem && cd /opt/aurora && docker compose restart nginx"
EOF

    success "Panel configuration written to /opt/aurora/"
}
