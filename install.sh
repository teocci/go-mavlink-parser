#!/bin/bash

## Define package name
PKG_NAME='proctel'

## Define module
MODULE_NAME='rtt'

## Define local path
LOCAL_PATH='/usr/local'

## Define system library path
LIB_PATH='/usr/lib'

## Define system binary path
BIN_PATH='/usr/bin'

## Define installation path
module_path="${LOCAL_PATH}/${MODULE_NAME}"

## Define module lib path
module_lib_path="${LIB_PATH}/${MODULE_NAME}"

## Define module bin path
module_bin_path="${BIN_PATH}/${MODULE_NAME}"

## Define package path
pkg_path="${module_path}/bin/${PKG_NAME}"

# Build the main
go build main.go

# Rename main as a PKG_NAME
mv -v main "${PKG_NAME}"

# Install
sudo mkdir -p "${module_path}/bin"
sudo cp -vf "${PKG_NAME}" "${pkg_path}"

if [ ! -L "${module_lib_path}" ] || [ ! -e "${module_lib_path}" ]; then
  sudo ln -sv "${module_path}" "${LIB_PATH}"
fi

if [ ! -L "${module_bin_path}/${PKG_NAME}" ] || [ ! -e "${module_bin_path}/${PKG_NAME}" ]; then
  sudo ln -sv "${module_lib_path}/bin/${PKG_NAME}" "${module_bin_path}/${PKG_NAME}"
fi
