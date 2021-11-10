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
pkg_lib_path="${module_lib_path}/bin/${PKG_NAME}"
pkg_bin_path="${module_bin_path}/${PKG_NAME}"

## Uninstall
sudo rm -vf "${pkg_bin_path}"
sudo rm -vf "${pkg_lib_path}"
sudo rm -vf "${pkg_path}"