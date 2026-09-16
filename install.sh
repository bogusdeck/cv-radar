#!/bin/bash
set -e

# CV-RADAR One-Line Installer & Setup Script
# Usage: curl -fsSL https://raw.githubusercontent.com/bogusdeck/cv-radar/main/install.sh | bash
# Or:    ./install.sh

REPO="https://github.com/bogusdeck/cv-radar.git"
BOLD="\033[1m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
BLUE="\033[0;34m"
RED="\033[0;31m"
RESET="\033[0m"

echo ""
echo -e "${BOLD}${BLUE}  ██████╗██╗   ██╗      ██████╗  █████╗ ██████╗  █████╗ ██████╗${RESET}"
echo -e "${BOLD}${BLUE} ██╔════╝██║   ██║      ██╔══██╗██╔══██╗██╔══██╗██╔══██╗██╔══██╗${RESET}"
echo -e "${BOLD}${BLUE} ██║     ██║   ██║█████╗██████╔╝███████║██║  ██║███████║██████╔╝${RESET}"
echo -e "${BOLD}${BLUE} ██║     ╚██╗ ██╔╝╚════╝██╔══██╗██╔══██║██║  ██║██╔══██║██╔══██╗${RESET}"
echo -e "${BOLD}${BLUE} ╚██████╗ ╚████╔╝       ██║  ██║██║  ██║██████╔╝██║  ██║██║  ██║${RESET}"
echo -e "${BOLD}${BLUE}  ╚═════╝  ╚═══╝        ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝╚═════╝${RESET}"
echo ""
echo -e "${BOLD}  ATS Resume Scanner & AI CV Optimizer Setup${RESET}"
echo -e "  github.com/bogusdeck/cv-radar"
echo ""

# ─── Detect Execution Context ─────────────────────────────────────────────────

if [ -f "./cmd/tui/main.go" ]; then
  INSTALL_DIR="$(pwd)"
  IS_LOCAL_BUILD=1
else
  INSTALL_DIR="${CV_RADAR_DIR:-$HOME/.cv-radar}"
  IS_LOCAL_BUILD=0
fi

# ─── Check Dependencies ───────────────────────────────────────────────────────

check_dep() {
  if ! command -v "$1" &>/dev/null; then
    echo -e "${RED}✗ $1 not found.${RESET} $2"
    exit 1
  else
    echo -e "${GREEN}✓ $1${RESET}"
  fi
}

echo -e "${BOLD}Checking dependencies...${RESET}"
check_dep "git" "Install git: https://git-scm.com"
check_dep "go"  "Install Go 1.21+: https://go.dev/dl"

# Tectonic (LaTeX compiler)
if ! command -v tectonic &>/dev/null; then
  echo -e "${YELLOW}⚠ tectonic not found — PDF compilation will be unavailable.${RESET}"
  echo -e "  Install: curl --proto '=https' --tlsv1.2 -fsSL https://drop.tectonic.typesetting.com/install.sh | sh"
  TECTONIC_MISSING=1
else
  echo -e "${GREEN}✓ tectonic${RESET}"
fi

echo ""

# ─── Clone or Update Repository ──────────────────────────────────────────────

if [ "$IS_LOCAL_BUILD" -eq 0 ]; then
  if [ -d "$INSTALL_DIR/.git" ]; then
    echo -e "${BOLD}Updating existing installation in $INSTALL_DIR ...${RESET}"
    git -C "$INSTALL_DIR" pull --ff-only
  else
    echo -e "${BOLD}Cloning cv-radar into $INSTALL_DIR ...${RESET}"
    git clone "$REPO" "$INSTALL_DIR"
  fi
  cd "$INSTALL_DIR"
fi

# ─── Build Binary ────────────────────────────────────────────────────────────

echo -e "${BOLD}Building TUI dashboard...${RESET}"
go build -o cv-tui ./cmd/tui
echo -e "${GREEN}✓ cv-tui binary built${RESET}"

chmod +x optimize.sh

# ─── Global CLI & Skill Setup ─────────────────────────────────────────────────

echo ""
echo -e "${BOLD}Configuring global CLI & agent skills...${RESET}"

# 1. Symlink binary to ~/.local/bin if available
BIN_DIR="$HOME/.local/bin"
mkdir -p "$BIN_DIR"
ln -sf "$INSTALL_DIR/cv-tui" "$BIN_DIR/cv-tui"
echo -e "${GREEN}✓ Symlinked binary to $BIN_DIR/cv-tui${RESET}"

# 2. Global AI Agent Skill Registration
SKILL_SRC="$INSTALL_DIR/.agents/skills/cv-optimizer"

if [ -d "$SKILL_SRC" ]; then
  # Antigravity / Open Agent standard
  mkdir -p "$HOME/.agents/skills/cv-optimizer"
  cp -r "$SKILL_SRC/"* "$HOME/.agents/skills/cv-optimizer/"
  echo -e "${GREEN}✓ Installed global skill to ~/.agents/skills/cv-optimizer${RESET}"

  # Claude Code
  mkdir -p "$HOME/.claude/skills/cv-optimizer"
  cp -r "$SKILL_SRC/"* "$HOME/.claude/skills/cv-optimizer/"
  echo -e "${GREEN}✓ Installed global skill to ~/.claude/skills/cv-optimizer${RESET}"

  # Codex
  mkdir -p "$HOME/.codex/skills/cv-optimizer"
  cp -r "$SKILL_SRC/"* "$HOME/.codex/skills/cv-optimizer/"
  echo -e "${GREEN}✓ Installed global skill to ~/.codex/skills/cv-optimizer${RESET}"

  # OpenCode
  mkdir -p "$HOME/.opencode/skills/cv-optimizer" "$HOME/.config/opencode/skills/cv-optimizer"
  cp -r "$SKILL_SRC/"* "$HOME/.opencode/skills/cv-optimizer/"
  cp -r "$SKILL_SRC/"* "$HOME/.config/opencode/skills/cv-optimizer/"
  echo -e "${GREEN}✓ Installed global skill to ~/.opencode/skills/cv-optimizer${RESET}"
fi

# ─── Done Output ─────────────────────────────────────────────────────────────

echo ""
echo -e "${GREEN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${GREEN}${BOLD}   🎉 CV-RADAR & CV-OPTIMIZER setup completed!${RESET}"
echo -e "${GREEN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo ""
echo -e "${BOLD}🚀 How to use:${RESET}"
echo ""
echo -e "  ${BOLD}1. Global AI CLI Skill (Any folder!):${RESET}"
echo -e "     Open your AI CLI in any project folder:"
echo -e "       ${BLUE}agy${RESET} | ${BLUE}claude${RESET} | ${BLUE}opencode${RESET} | ${BLUE}codex${RESET}"
echo -e "     Then type:"
echo -e "       ${BOLD}/cv-optimizer${RESET}"
echo ""
echo -e "  ${BOLD}2. Terminal UI (TUI):${RESET}"
echo -e "     Run from anywhere:"
echo -e "       ${BLUE}cv-tui${RESET}"
echo ""

if [ -n "$TECTONIC_MISSING" ]; then
  echo -e "${YELLOW}⚠ Note: Install tectonic for automated PDF rendering:${RESET}"
  echo -e "  curl --proto '=https' --tlsv1.2 -fsSL https://drop.tectonic.typesetting.com/install.sh | sh"
  echo ""
fi
