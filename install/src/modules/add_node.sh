#!/bin/bash
# AURORA — Add node to existing panel module

_do_add_node() {
    local NODE_NAME="$1"
    local PANEL_URL="$2"
    local API_KEY="$3"

    # This runs on the NEW node server.
    # It registers the node with the central panel via API.

    mkdir -p /opt/aurora-node/{configs,data,ssl}

    # Fetch node join token from panel
    local REGISTER_RESPONSE
    REGISTER_RESPONSE=$(curl -s -X POST "${PANEL_URL}/api/v1/nodes" \
        -H "Authorization: Bearer ${API_KEY}" \
        -H "Content-Type: application/json" \
        -d "{\"name\":\"${NODE_NAME}\",\"host\":\"$(curl -s4 ifconfig.me)\",\"port\":443,\"apiPort\":10085,\"weight\":1}")

    local NODE_ID
    NODE_ID=$(echo "$REGISTER_RESPONSE" | jq -r '.id // empty')

    if [ -z "$NODE_ID" ]; then
        warning "Failed to register node with panel."
        warning "Response: $REGISTER_RESPONSE"
        reading "Continue with manual setup? (y/n): " manual
        [[ "$manual" != "${LANG[CONFIRM_YES]}" ]] && return 1
    else
        success "Node registered with ID: $NODE_ID"
    fi

    # Install node using the same module
    source "${INSTALL_DIR}/modules/install_node.sh" 2>/dev/null || return 1
    _do_install_node "$NODE_NAME" "10085" "$PANEL_URL"

    # Generate Xray API key and encrypt it
    local XRAY_API_KEY
    XRAY_API_KEY=$(tr -dc 'A-Za-z0-9' < /dev/urandom | head -c 32)

    # Store the API key for node-agent
    cat > /opt/aurora-node/.env <<EOF
NODE_NAME=${NODE_NAME}
PANEL_URL=${PANEL_URL}
XRAY_API_PORT=10085
XRAY_API_KEY=${XRAY_API_KEY}
NODE_ID=${NODE_ID}
EOF

    success "Node ready. Start with: cd /opt/aurora-node && docker compose up -d"
}
