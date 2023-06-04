echo "This script requires sudo access"

export HOSTNAME=hostname
export API_PORT=port_for_api_number
export KEYS_PORT=port_for_keys_number

sudo bash -c "$(wget -qO- https://raw.githubusercontent.com/Jigsaw-Code/outline-server/master/src/server_manager/install_scripts/install_server.sh)" "" --hostname "$HOSTNAME" --api-port $API_PORT --keys-port $KEYS_PORT
