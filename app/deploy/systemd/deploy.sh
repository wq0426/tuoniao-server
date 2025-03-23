#!/bin/bash

sudo touch /etc/systemd/system/flashbear.service

sudo systemctl daemon-reload
sudo systemctl enable flashbear
sudo systemctl start flashbear

sudo systemctl status flashbear