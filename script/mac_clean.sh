#!/bin/sh

set -e

OS_NAME=$(uname -s)
JB_PRODUCTS="idea clion phpstorm goland pycharm webstorm webide rider datagrip rubymine appcode dataspell gateway jetbrains_client jetbrainsclient studio devecostudio"

KDE_ENV_DIR="${HOME}/.config/plasma-workspace/env"

PROFILE_PATH="${HOME}/.profile"
ZSH_PROFILE_PATH="${HOME}/.zshrc"
PLIST_PATH="${HOME}/Library/LaunchAgents/jetbrains.vmoptions.plist"

if [ $OS_NAME = "Darwin" ]; then
  BASH_PROFILE_PATH="${HOME}/.bash_profile"
else
  BASH_PROFILE_PATH="${HOME}/.bashrc"
fi

touch "${PROFILE_PATH}"
touch "${BASH_PROFILE_PATH}"
touch "${ZSH_PROFILE_PATH}"

MY_VMOPTIONS_SHELL_NAME="jetbrains.vmoptions.sh"
MY_VMOPTIONS_SHELL_FILE="${HOME}/.${MY_VMOPTIONS_SHELL_NAME}"

rm -rf "${MY_VMOPTIONS_SHELL_FILE}"

if [ $OS_NAME = "Darwin" ]; then
  for PRD in $JB_PRODUCTS; do
    ENV_NAME=$(echo $PRD | tr '[a-z]' '[A-Z]')"_VM_OPTIONS"
    
    launchctl unsetenv "${ENV_NAME}"
  done

  rm -rf "${PLIST_PATH}"

  sed -i '' '/___MY_VMOPTIONS_SHELL_FILE="${HOME}\/\.jetbrains\.vmoptions\.sh"; if /d' "${PROFILE_PATH}" >/dev/null 2>&1
  sed -i '' '/___MY_VMOPTIONS_SHELL_FILE="${HOME}\/\.jetbrains\.vmoptions\.sh"; if /d' "${BASH_PROFILE_PATH}" >/dev/null 2>&1
  sed -i '' '/___MY_VMOPTIONS_SHELL_FILE="${HOME}\/\.jetbrains\.vmoptions\.sh"; if /d' "${ZSH_PROFILE_PATH}" >/dev/null 2>&1

  echo 'done.'
else
  sed -i '/___MY_VMOPTIONS_SHELL_FILE="${HOME}\/\.jetbrains\.vmoptions\.sh"; if /d' "${PROFILE_PATH}" >/dev/null 2>&1
  sed -i '/___MY_VMOPTIONS_SHELL_FILE="${HOME}\/\.jetbrains\.vmoptions\.sh"; if /d' "${BASH_PROFILE_PATH}" >/dev/null 2>&1
  sed -i '/___MY_VMOPTIONS_SHELL_FILE="${HOME}\/\.jetbrains\.vmoptions\.sh"; if /d' "${ZSH_PROFILE_PATH}" >/dev/null 2>&1

  rm -rf "${KDE_ENV_DIR}/${MY_VMOPTIONS_SHELL_NAME}"
  echo "done. you'd better log off first!"
fi
