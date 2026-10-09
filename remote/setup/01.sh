#!/bin/bash
set -eu

# ==================================================================================== #
# VARIABLES
# ==================================================================================== #
TIMEZONE=Asia/Kolkata
USERNAME=greenlight

# Prompt for your Neon Database connection string
read -p "Enter your full Neon DB connection string (DSN): " NEON_DSN
read -p "Enter the Resend Api Key (for sending mails): " RESEND_API_KEY

# ==================================================================================== #
# SCRIPT LOGIC
# ==================================================================================== #

echo "Updating system packages..."
apt update
apt --yes -o Dpkg::Options::="--force-confnew" upgrade

echo "Setting timezone and locales..."
timedatectl set-timezone ${TIMEZONE}
apt --yes install locales-all
# Export the locale AFTER the packages are installed to prevent the warning
export LC_ALL=en_US.UTF-8

echo "Creating the ${USERNAME} user..."
useradd --create-home --shell "/bin/bash" --groups sudo "${USERNAME}"
passwd "${USERNAME}"

echo "Creating secure .env file for the greenlight user..."
cat > /home/greenlight/.env << EOF
DATABASE_URL=${NEON_DSN}
RESEND_API_KEY=${RESEND_API_KEY}
EOF

# Lock down the permissions so ONLY the greenlight user can read it
chown ${USERNAME}:${USERNAME} /home/${USERNAME}/.env
chmod 600 /home/${USERNAME}/.env

echo "Installing Caddy..."
apt install -y debian-keyring debian-archive-keyring apt-transport-https curl
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/gpg.key' | gpg --dearmor -o /usr/share/keyrings/caddy-stable-archive-keyring.gpg
curl -1sLf 'https://dl.cloudsmith.io/public/caddy/stable/debian.deb.txt' | tee /etc/apt/sources.list.d/caddy-stable.list
apt update
apt install caddy -y

echo "Script complete! Rebooting..."
reboot