#!/bin/bash
set -e

# CV-RADAR Installer
# Usage: curl -fsSL https://raw.githubusercontent.com/bogusdeck/cv-radar/main/install.sh | bash
# Or:    ./install.sh

REPO="https://github.com/bogusdeck/cv-radar.git"
INSTALL_DIR="${CV_RADAR_DIR:-./cv-radar}"
BOLD="\033[1m"
GREEN="\033[0;32m"
YELLOW="\033[0;33m"
RED="\033[0;31m"
RESET="\033[0m"

echo ""
echo -e "${BOLD}  ██████╗██╗   ██╗      ██████╗  █████╗ ██████╗  █████╗ ██████╗${RESET}"
echo -e "${BOLD} ██╔════╝██║   ██║      ██╔══██╗██╔══██╗██╔══██╗██╔══██╗██╔══██╗${RESET}"
echo -e "${BOLD} ██║     ██║   ██║█████╗██████╔╝███████║██║  ██║███████║██████╔╝${RESET}"
echo -e "${BOLD} ██║     ╚██╗ ██╔╝╚════╝██╔══██╗██╔══██║██║  ██║██╔══██║██╔══██╗${RESET}"
echo -e "${BOLD} ╚██████╗ ╚████╔╝       ██║  ██║██║  ██║██████╔╝██║  ██║██║  ██║${RESET}"
echo -e "${BOLD}  ╚═════╝  ╚═══╝        ╚═╝  ╚═╝╚═╝  ╚═╝╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═╝${RESET}"
echo ""
echo -e "${BOLD}  ATS Resume Scanner & CV Optimizer${RESET}"
echo -e "  github.com/bogusdeck/cv-radar"
echo ""

# ─── Check dependencies ───────────────────────────────────────────────────────

check() {
  if ! command -v "$1" &>/dev/null; then
    echo -e "${RED}✗ $1 not found.${RESET} $2"
    exit 1
  else
    echo -e "${GREEN}✓ $1${RESET}"
  fi
}

echo -e "${BOLD}Checking dependencies...${RESET}"
check "git"   "Install git: https://git-scm.com"
check "go"    "Install Go 1.21+: https://go.dev/dl"

# Tectonic (LaTeX compiler) — optional but needed for PDF compilation
if ! command -v tectonic &>/dev/null; then
  echo -e "${YELLOW}⚠ tectonic not found — PDF compilation will be unavailable.${RESET}"
  echo -e "  Install: curl --proto '=https' --tlsv1.2 -fsSL https://drop.tectonic.typesetting.com/install.sh | sh"
  TECTONIC_MISSING=1
else
  echo -e "${GREEN}✓ tectonic${RESET}"
fi

# Node (optional, for web UI)
if ! command -v node &>/dev/null; then
  echo -e "${YELLOW}⚠ node not found — web UI will be unavailable.${RESET}"
  NODE_MISSING=1
else
  echo -e "${GREEN}✓ node$(node -v)${RESET}"
fi

echo ""

# ─── Clone ────────────────────────────────────────────────────────────────────

if [ -d "$INSTALL_DIR/.git" ]; then
  echo -e "${BOLD}Updating existing installation...${RESET}"
  git -C "$INSTALL_DIR" pull --ff-only
else
  echo -e "${BOLD}Cloning cv-radar into $INSTALL_DIR ...${RESET}"
  git clone "$REPO" "$INSTALL_DIR"
fi

cd "$INSTALL_DIR"
echo ""

# ─── Build Go API server ──────────────────────────────────────────────────────

echo -e "${BOLD}Building Go API server...${RESET}"
go build -o server ./cmd/server
echo -e "${GREEN}✓ server binary built${RESET}"

# ─── Build TUI dashboard ──────────────────────────────────────────────────────

echo -e "${BOLD}Building TUI dashboard...${RESET}"
go build -o cv-tui ./cmd/tui
echo -e "${GREEN}✓ cv-tui binary built${RESET}"

# ─── Build web UI ─────────────────────────────────────────────────────────────

if [ -z "$NODE_MISSING" ]; then
  echo -e "${BOLD}Installing web UI dependencies...${RESET}"
  cd web && npm install --silent && npm run build && cd ..
  echo -e "${GREEN}✓ web UI built${RESET}"
fi

# ─── Make scripts executable ──────────────────────────────────────────────────

chmod +x start.sh optimize.sh

# ─── Done ─────────────────────────────────────────────────────────────────────

echo ""
echo -e "${GREEN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo -e "${GREEN}${BOLD}  cv-radar installed successfully!${RESET}"
echo -e "${GREEN}${BOLD}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${RESET}"
echo ""
echo -e "${BOLD}Next steps:${RESET}"
echo ""
echo -e "  1. Add your CV and Job Description:"
echo -e "     cp your_cv.tex ${INSTALL_DIR}/cv.tex"
echo -e "     cp your_jd.txt ${INSTALL_DIR}/jd.txt"
echo ""
echo -e "  2. Launch the TUI dashboard:"
echo -e "     cd ${INSTALL_DIR} && ./cv-tui"
echo ""
echo -e "  3. Start the full web stack:"
echo -e "     cd ${INSTALL_DIR} && ./start.sh"
echo ""
echo -e "  4. Run headless AI optimization (Claude, Antigravity, or OpenCode):"
echo -e "     ./optimize.sh cv.tex jd.txt Workday claude"
echo ""
echo -e "  5. Or open in your AI coding CLI:"
echo -e "     cd ${INSTALL_DIR} && claude   # Claude Code"
echo -e "     cd ${INSTALL_DIR} && agy      # Antigravity"
echo -e "     cd ${INSTALL_DIR} && opencode # OpenCode"
echo -e "     Then say: /cv-optimizer"
echo ""

if [ -n "$TECTONIC_MISSING" ]; then
  echo -e "${YELLOW}  ⚠ Install tectonic for PDF compilation:${RESET}"
  echo -e "    curl --proto '=https' --tlsv1.2 -fsSL https://drop.tectonic.typesetting.com/install.sh | sh"
  echo ""
fi
