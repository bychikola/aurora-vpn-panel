#!/bin/bash
# AURORA Panel installation module

_do_install_panel() {
    local PANEL_DOMAIN="$1"
    local DB_PASSWORD="$2"
    local JWT_ACCESS="$3"
    local JWT_REFRESH="$4"
    local ADMIN_PASSWORD="$5"
    local WEBSERVER="${6:-nginx}"

    local AURORA_DIR="/opt/aurora"
    local REPO_URL="https://github.com/bychikola/aurora-vpn-panel.git"

    mkdir -p "$AURORA_DIR"/{configs,data/postgres,data/redis,ssl,build}

    # ─── Clone repo if not present ───
    if [ ! -d "$AURORA_DIR/repo/.git" ]; then
        info "Cloning AURORA repository..."
        git clone --depth 1 "$REPO_URL" "$AURORA_DIR/repo" 2>&1 || {
            warning "Git clone failed. Trying with https + GIT_SSL_NO_VERIFY..."
            GIT_SSL_NO_VERIFY=true git clone --depth 1 "$REPO_URL" "$AURORA_DIR/repo" 2>&1 || {
                error "Cannot clone repository. Check internet connection."
            }
        }
        success "Repository cloned"
    fi

    # ─── .env ───
    cat > "$AURORA_DIR/.env" <<EOF
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

    # ─── config.yaml ───
    cat > "$AURORA_DIR/configs/config.yaml" <<EOF
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

    # ─── SSL certs for Nginx ───
    if [ "$WEBSERVER" = "nginx" ]; then
        local cert_dir="/etc/letsencrypt/live/${PANEL_DOMAIN}"
        local base_domain
        base_domain=$(echo "$PANEL_DOMAIN" | awk -F'.' '{if (NF > 2) {print $(NF-1)"."$NF} else {print $0}}')
        [ ! -d "$cert_dir" ] && cert_dir="/etc/letsencrypt/live/${base_domain}"
        if [ -d "$cert_dir" ]; then
            cp "$cert_dir/fullchain.pem" "$AURORA_DIR/ssl/fullchain.pem"
            cp "$cert_dir/privkey.pem" "$AURORA_DIR/ssl/privkey.pem"
        fi
    fi

    # ─── docker-compose.yml ───
    cat > "$AURORA_DIR/docker-compose.yml" <<YAML
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
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 3s
      retries: 5
    networks:
      - aurora-net

  backend:
    build:
      context: ./repo/aurora-backend
      dockerfile: deployments/Dockerfile
    image: aurora-backend:local
    container_name: aurora-backend
    restart: unless-stopped
    env_file:
      - .env
    volumes:
      - ./configs:/app/configs:ro
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - aurora-net

  frontend:
    build:
      context: ./repo/aurora-frontend
      dockerfile_inline: |
        FROM node:22-alpine AS builder
        WORKDIR /app
        COPY package.json package-lock.json ./
        RUN npm ci
        COPY . .
        RUN npm run build

        FROM nginx:alpine
        COPY --from=builder /app/dist /usr/share/nginx/html
        COPY nginx-frontend.conf /etc/nginx/conf.d/default.conf
        EXPOSE 80
    image: aurora-frontend:local
    container_name: aurora-frontend
    restart: unless-stopped
    depends_on:
      - backend
    networks:
      - aurora-net

YAML

    # Frontend nginx config
    cat > "$AURORA_DIR/repo/aurora-frontend/nginx-frontend.conf" <<'NGX'
server {
    listen 80;
    server_name _;
    root /usr/share/nginx/html;
    index index.html;
    location / {
        try_files $uri $uri/ /index.html;
    }
    location /api/ {
        proxy_pass http://backend:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }
}
NGX

    # ─── Reverse proxy ───
    if [ "$WEBSERVER" = "caddy" ]; then
        cat >> "$AURORA_DIR/docker-compose.yml" <<YAML
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
        cat > "$AURORA_DIR/Caddyfile" <<EOF
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
    else
        cat >> "$AURORA_DIR/docker-compose.yml" <<YAML
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
        cat > "$AURORA_DIR/nginx.conf" <<EOF
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
        }
        location /api/ {
            proxy_pass http://backend:8080;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        }
    }
}
EOF

        # SSL renewal cron
        local base_domain
        base_domain=$(echo "$PANEL_DOMAIN" | awk -F'.' '{if (NF > 2) {print $(NF-1)"."$NF} else {print $0}}')
        cat > /etc/cron.d/aurora-ssl <<EOF
0 3 * * * root certbot renew --quiet --post-hook "cp /etc/letsencrypt/live/${base_domain}/fullchain.pem $AURORA_DIR/ssl/fullchain.pem && cp /etc/letsencrypt/live/${base_domain}/privkey.pem $AURORA_DIR/ssl/privkey.pem && cd $AURORA_DIR && docker compose restart nginx"
EOF
    fi

    cat >> "$AURORA_DIR/docker-compose.yml" <<YAML
networks:
  aurora-net:
    driver: bridge
YAML

    # ─── Build images (takes a few minutes first time) ───
    info "Building AURORA Docker images (this may take 3-5 minutes on first run)..."
    info "--- Build output below ---"
    cd "$AURORA_DIR" && docker compose build --progress=plain 2>&1
    local build_exit=$?
    info "--- End of build output ---"

    if [ $build_exit -ne 0 ]; then
        echo -e ""
        error "Docker build FAILED (exit code: $build_exit). Check output above for details."
    fi
    success "All images built successfully"

    success "Panel configuration written to $AURORA_DIR/ (webserver: $WEBSERVER)"
}
