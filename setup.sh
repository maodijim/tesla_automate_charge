#!/usr/bin/env bash

function create_dir() {
    mkdir -p /opt/tesla_charge_control
    cp tesla_automated_charge_control /opt/tesla_charge_control/tesla_automated_charge_control
    cp configs.yml.template /opt/tesla_charge_control/configs.yml
}

function create_systemd() {
    # Generate Systemd service file
    cat <<EOF > /etc/systemd/system/tesla_charge_control.service
[Unit]
Description=Tesla Automated Charge Control
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/tesla_charge_control/
ExecStart=/opt/tesla_charge_control/tesla_automated_charge_control
Restart=always
RestartSec=10

[Install]
WantedBy=multi-user.target
EOF
}

function reload_service() {
    systemctl daemon-reload
}

function enable_service() {
    systemctl enable tesla_charge_control.service
}

function start_service() {
    systemctl start tesla_charge_control.service
}

create_dir
create_systemd
reload_service
enable_service
start_service
