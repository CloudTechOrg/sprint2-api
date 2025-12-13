#!/bin/bash
set -e

#############################################
# Configuration
#############################################
DB_HOST="<RDSのエンドポイント>"
DB_USER="<マスターユーザー名>"
DB_PASSWORD="<マスターパスワード>"
DB_NAME="meibo"
GITHUB_REPO="https://github.com/CloudTechOrg/sprint2-api.git"
GITHUB_BRANCH="main"
#############################################

# Update system
dnf update -y

# Install Go and Git
dnf install -y golang git

# Set Go environment variables
export HOME=/root
export GOPATH=/root/go
export GOMODCACHE=/root/go/pkg/mod
export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin

# Create app directory
mkdir -p /opt/api
cd /opt/api

# Clone repository from GitHub
git clone -b $GITHUB_BRANCH $GITHUB_REPO .
rm -rf .git

# Download dependencies
go mod tidy

# Build the application
go build -o meibo-api main.go

# Create systemd service
cat > /etc/systemd/system/meibo-api.service << SERVICEEOF
[Unit]
Description=Meibo API Server
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/api
Environment="DB_HOST=$DB_HOST"
Environment="DB_USER=$DB_USER"
Environment="DB_PASSWORD=$DB_PASSWORD"
Environment="DB_NAME=$DB_NAME"
ExecStart=/opt/api/meibo-api
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
SERVICEEOF

# Enable and start service
systemctl daemon-reload
systemctl enable meibo-api
systemctl start meibo-api
