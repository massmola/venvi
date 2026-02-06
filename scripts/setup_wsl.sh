#!/usr/bin/env bash
set -e

echo "🔧 Configuring WSL and Nix..."

# 1. Enable Systemd in WSL
if ! grep -q "boot" /etc/wsl.conf 2>/dev/null; then
    echo " -> Enabling systemd in /etc/wsl.conf..."
    # Use tee to write with sudo
    echo -e "\n[boot]\nsystemd=true" | sudo tee -a /etc/wsl.conf > /dev/null
elif ! grep -q "systemd=true" /etc/wsl.conf 2>/dev/null; then
    echo " -> Enabling systemd in /etc/wsl.conf..."
    echo "systemd=true" | sudo tee -a /etc/wsl.conf > /dev/null
else
    echo " -> systemd already enabled in /etc/wsl.conf"
fi

# 2. Enable Experimental Features in Nix
# Check if directory exists
if [ ! -d "/etc/nix" ]; then
    echo " -> Creating /etc/nix directory..."
    sudo mkdir -p /etc/nix
fi

if ! grep -q "experimental-features" /etc/nix/nix.conf 2>/dev/null; then
    echo " -> Enabling 'nix-command' and 'flakes' in /etc/nix/nix.conf..."
    echo "experimental-features = nix-command flakes" | sudo tee -a /etc/nix/nix.conf > /dev/null
else
    echo " -> experimental-features check passed (check /etc/nix/nix.conf if manually edited)"
fi

echo "✅ Configuration complete!"
echo "⚠️  IMPORTANT: You must RESTART WSL for these changes to take effect."
echo "   Run this in PowerShell: wsl --shutdown"
