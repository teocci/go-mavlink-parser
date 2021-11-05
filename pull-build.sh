#!/bin/bash
# Git pull
git pull

# Build the main
go build main.go

# Rename main as a proctel
mv -v main proctel

# Install the app in the rtt app directory
sudo cp -vf proctel /usr/local/rtt/bin/proctel
