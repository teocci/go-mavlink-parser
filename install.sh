#!/bin/bash

## Define local path
LOCAL_PATH='/usr/local'

## Define system library path
LIB_PATH='/usr/lib'

## Define system binary path
BIN_PATH='/usr/bin'

## Define package name
PKG_NAME='proctel'

## Define module
MODULE_NAME='rtt'

## Define installation path
module_path="${LOCAL_PATH}/${MODULE_NAME}/bin"

## Define module lib path
module_lib_path="${LIB_PATH}/${MODULE_NAME}/bin"

## Define module bin path
module_bin_path="${BIN_PATH}/${MODULE_NAME}"

## Define package path
pkg_path="${module_path}/${PKG_NAME}"

# Build the main
go build main.go

# Rename main as a PKG_NAME
mv -v main "${PKG_NAME}"

# Install
sudo mkdir -p "${module_path}"
sudo mkdir -p "${module_lib_path}"
sudo mkdir -p "${module_bin_path}"
sudo cp -v "${PKG_NAME}" "${pkg_path}"

if [ ! -L "${pkg_path}" ] || [ ! -e "${pkg_path}" ]; then
  sudo ln -sv "${module_path}" "${module_lib_path}"
  sudo ln -sv "${module_lib_path}/${PKG_NAME}" "${module_bin_path}/${PKG_NAME}"
fi