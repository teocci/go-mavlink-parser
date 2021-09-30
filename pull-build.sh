#!/bin/bash
# Git pull
git pull

# Build the main
go build main.go

# Rename main as a proctel
mv -v main proctel
cp -v proctel /home/rtt/apps/proctel