import re

with open("cmd/tui/main.go", "r") as f:
    text = f.read()

text = text.replace('result     *ats.Result', 'result     *models.ATSResult')
text = text.replace('var res ats.Result', 'var res models.ATSResult')
text = text.replace('case ats.Result:', 'case models.ATSResult:')

# Also add the models import
import_block = """import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/bogusdeck/ats-scanner/internal/ats"
	"github.com/bogusdeck/ats-scanner/internal/models"
	"github.com/bogusdeck/ats-scanner/internal/parser"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)"""

text = re.sub(r"import \([^)]+\)", import_block, text)

with open("cmd/tui/main.go", "w") as f:
    f.write(text)

