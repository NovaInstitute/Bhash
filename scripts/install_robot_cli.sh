#!/usr/bin/env bash
# Installs the ROBOT command-line interface along with its dependencies.
# Designed for Debian/Ubuntu environments with apt package manager.

set -euo pipefail

ROBOT_VERSION="${ROBOT_VERSION:-1.9.5}"
ROBOT_JAR_URL="https://github.com/ontodev/robot/releases/download/v${ROBOT_VERSION}/robot.jar"
INSTALL_PREFIX="${INSTALL_PREFIX:-$HOME/opt/robot}"
BIN_DIR="${BIN_DIR:-$HOME/bin}"
ROBOT_JAR_PATH="${INSTALL_PREFIX}/robot.jar"
ROBOT_WRAPPER_PATH="${BIN_DIR}/robot"

log() {
  printf '[ROBOT-SETUP] %s\n' "$1"
}

require_command() {
  if ! command -v "$1" >/dev/null 2>&1; then
    return 1
  fi
}

ensure_apt_package() {
  local package="$1"
  if dpkg -s "$package" >/dev/null 2>&1; then
    return 0
  fi

  if ! command -v apt-get >/dev/null 2>&1; then
    log "apt-get not available. Please install '$package' manually."
    exit 1
  fi

  if [[ $EUID -ne 0 ]]; then
    if command -v sudo >/dev/null 2>&1; then
      log "Installing '${package}' via sudo apt-get."
      sudo apt-get update
      sudo apt-get install -y "$package"
    else
      log "Please rerun this script as root or install '${package}' manually."
      exit 1
    fi
  else
    log "Installing '${package}' via apt-get."
    apt-get update
    apt-get install -y "$package"
  fi
}

ensure_java() {
  if require_command java; then
    return
  fi
  ensure_apt_package openjdk-17-jre-headless
}

ensure_curl() {
  if require_command curl; then
    return
  fi
  ensure_apt_package curl
}

install_robot() {
  mkdir -p "$INSTALL_PREFIX"
  mkdir -p "$BIN_DIR"

  log "Downloading ROBOT v${ROBOT_VERSION} from ${ROBOT_JAR_URL}."
  curl -L -o "$ROBOT_JAR_PATH" "$ROBOT_JAR_URL"

  cat > "$ROBOT_WRAPPER_PATH" <<SCRIPT
#!/usr/bin/env bash
exec java -jar "$ROBOT_JAR_PATH" "\$@"
SCRIPT
  chmod +x "$ROBOT_WRAPPER_PATH"
}

verify_installation() {
  if ! "$ROBOT_WRAPPER_PATH" --version; then
    log "ROBOT installation failed."
    exit 1
  fi
}

log "Ensuring Java runtime is available."
ensure_java

log "Ensuring curl is available for downloads."
ensure_curl

log "Installing ROBOT CLI."
install_robot

log "Verifying installation."
verify_installation

log "ROBOT CLI installation completed. Ensure '${BIN_DIR}' is on your PATH."
