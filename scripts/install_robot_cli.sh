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
CONFIG_DIR="${CONFIG_DIR:-$HOME/.config/bhash}"
SECRETS_FILE="${SECRETS_FILE:-$CONFIG_DIR/secrets.env}"
PROFILE_FILES=(${PROFILE_FILES:-$HOME/.bashrc $HOME/.profile $HOME/.zshrc})
INTERACTIVE="${INTERACTIVE:-1}"
SECRETS_HEADER="# >>> bhash secrets >>>"
SECRETS_FOOTER="# <<< bhash secrets <<<"

log() {
  printf '[ROBOT-SETUP] %s\n' "$1"
}

persist_env_var() {
  local key="$1"
  local value="$2"
  mkdir -p "$CONFIG_DIR"
  touch "$SECRETS_FILE"
  chmod 600 "$SECRETS_FILE"

  local tmp
  tmp=$(mktemp)
  if ! grep -v "^export ${key}=" "$SECRETS_FILE" >"$tmp" 2>/dev/null; then
    :
  fi
  printf 'export %s=%q\n' "$key" "$value" >>"$tmp"
  mv "$tmp" "$SECRETS_FILE"
  export "$key=$value"
  log "Persisted ${key} to ${SECRETS_FILE}."
}

set_default_env_var() {
  local key="$1"
  local default_value="$2"
  local value="${!key:-$default_value}"
  persist_env_var "$key" "$value"
}

prompt_for_secret() {
  local key="$1"
  local prompt="$2"
  local secret_input="${3:-0}"
  local value="${!key:-}"

  if [[ -n "$value" ]]; then
    log "Using existing value for ${key}."
  elif [[ "$INTERACTIVE" == "0" ]]; then
    log "Environment variable ${key} must be set when running non-interactively."
    exit 1
  else
    if [[ "$secret_input" == "1" ]]; then
      read -r -s -p "$prompt" value || true
      printf '\n'
    else
      read -r -p "$prompt" value || true
    fi
  fi

  if [[ -z "$value" ]]; then
    log "No value provided for ${key}."
    exit 1
  fi

  persist_env_var "$key" "$value"
}

ensure_profile_hook() {
  local hook="[ -f \"$SECRETS_FILE\" ] && source \"$SECRETS_FILE\""
  for profile in "${PROFILE_FILES[@]}"; do
    [[ -z "$profile" ]] && continue
    if [[ ! -e "$profile" ]]; then
      if [[ "$profile" == "$HOME/.bashrc" ]]; then
        touch "$profile"
      else
        continue
      fi
    fi

    if ! grep -Fq "$SECRETS_HEADER" "$profile" 2>/dev/null; then
      {
        printf '\n%s\n' "$SECRETS_HEADER"
        printf '%s\n' "$hook"
        printf '%s\n' "$SECRETS_FOOTER"
      } >>"$profile"
      log "Added secrets sourcing hook to ${profile}."
    fi
  done
}

configure_secrets() {
  log "Configuring Hedera and Fluree secrets."
  prompt_for_secret "HEDERA_OPERATOR_ID" "Enter Hedera operator ID (e.g. 0.0.x): "
  prompt_for_secret "HEDERA_OPERATOR_KEY" "Enter Hedera operator private key: " 1
  set_default_env_var "HEDERA_NETWORK" "testnet"
  prompt_for_secret "FLUREE_API_TOKEN" "Enter Fluree API token: " 1
  prompt_for_secret "FLUREE_HANDLE" "Enter Fluree tenant handle (e.g. team/ledger): "
  set_default_env_var "FLUREE_BASE_URL" "https://data.flur.ee"
  ensure_profile_hook
  log "Secrets stored in ${SECRETS_FILE}."
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

configure_secrets

log "Installing ROBOT CLI."
install_robot

log "Verifying installation."
verify_installation

log "ROBOT CLI installation completed. Ensure '${BIN_DIR}' is on your PATH."
