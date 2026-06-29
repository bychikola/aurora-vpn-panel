#!/bin/bash
# AURORA Node installation module

_do_install_node() {
    local NODE_NAME="$1"
    local XRAY_API_PORT="$2"
    local PANEL_URL="$3"

    mkdir -p /opt/aurora-node/{configs,data/xray,ssl}

    # docker-compose for node
    cat > /opt/aurora-node/docker-compose.yml <<'YAML'
services:
  xray-core:
    image: ghcr.io/xtls/xray-core:latest
    container_name: aurora-xray
    restart: unless-stopped
    volumes:
      - ./configs/xray_config.json:/etc/xray/config.json:ro
      - ./ssl:/etc/xray/ssl:ro
    ports:
      - "443:443"
      - "8080:8080"
      - "8388:8388"
      - "8443:8443"
      - "127.0.0.1:${XRAY_API_PORT}:10085"
    command: xray run -config /etc/xray/config.json
    networks:
      - aurora-node-net

  node-agent:
    image: ghcr.io/bychikola/aurora-node-agent:latest
    container_name: aurora-node-agent
    restart: unless-stopped
    environment:
      NODE_NAME: ${NODE_NAME}
      PANEL_URL: ${PANEL_URL}
      XRAY_GRPC_ADDR: xray-core:10085
    volumes:
      - ./configs:/app/configs:ro
      - /var/run/docker.sock:/var/run/docker.sock:ro
    depends_on:
      - xray-core
    networks:
      - aurora-node-net

  beszel-agent:
    image: henrygd/beszel-agent:latest
    container_name: beszel-agent
    restart: unless-stopped
    network_mode: host
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
    environment:
      PORT: 45876
      KEY: ${BESZEL_KEY:-}

networks:
  aurora-node-net:
    driver: bridge
YAML

    # Initial Xray config
    cat > /opt/aurora-node/configs/xray_config.json <<'JSON'
{
  "log": { "loglevel": "warning" },
  "inbounds": [],
  "outbounds": [
    {
      "protocol": "freedom",
      "tag": "direct"
    }
  ],
  "routing": {
    "domainStrategy": "AsIs",
    "rules": []
  },
  "stats": {},
  "policy": {
    "system": {
      "statsInboundUplink": true,
      "statsInboundDownlink": true
    }
  },
  "api": {
    "tag": "api",
    "services": ["HandlerService", "StatsService"]
  }
}
JSON

    # Env file
    cat > /opt/aurora-node/.env <<EOF
NODE_NAME=${NODE_NAME}
PANEL_URL=${PANEL_URL}
XRAY_API_PORT=${XRAY_API_PORT}
EOF

    success "Node configuration written to /opt/aurora-node/"
}
