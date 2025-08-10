# WireGuard VPN Server

This directory contains the WireGuard VPN server setup for the VPS infrastructure.

## Quick Start

### Deploy to Production
```bash
export SSHPASS=your_ssh_password
make deploy
```

### Get Client Configurations
```bash
make get-config
make qr-codes  # Show QR codes for mobile clients
```

### Management Commands
```bash
make status    # Check server status
make logs      # View server logs
make restart   # Restart server
make stop      # Stop server
```

## Configuration

The WireGuard server is configured with:
- **Port**: 51820/UDP
- **Network**: 10.13.13.0/24
- **Peers**: 10 client configurations pre-generated
- **DNS**: Auto-configured
- **Allowed IPs**: 0.0.0.0/0 (full tunnel)

## Client Setup

1. Deploy the server: `make deploy`
2. Wait 2-3 minutes for initialization
3. Download configs: `make get-config`
4. Import the `.conf` file into your WireGuard client

### Mobile Clients (iOS/Android)
Use QR codes for easy setup:
```bash
make qr-codes
```

### Desktop Clients
Use the `.conf` files from the `./configs/` directory.

## Security Notes

- All client configs are automatically generated with unique keys
- Server uses persistent storage for configuration data
- Traffic is routed through the VPS server
- DNS queries are handled by the VPN
- **Important**: Client configs contain private keys and are automatically ignored by git
- Never commit `.conf` files to version control - they contain sensitive cryptographic keys

## Troubleshooting

### Check Server Status
```bash
make status
```

### View Logs
```bash
make logs
```

### Regenerate Configs
If you need to regenerate client configs, restart the container:
```bash
make restart
```

## Network Configuration

The server creates a VPN network on `10.13.13.0/24`:
- Server: `10.13.13.1`
- Clients: `10.13.13.2` - `10.13.13.11`

## Files

- `docker-compose.wireguard.yml`: WireGuard server configuration
- `Makefile`: Deployment and management commands
- `configs/`: Client configuration files (generated after deployment)