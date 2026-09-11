#!/bin/bash

# A script that uses headless AI CLI tools (Claude Code, Antigravity, etc.) to optimize a CV.

if [ "$#" -lt 3 ]; then
    echo "Usage: ./optimize.sh <path_to_cv> <path_to_jd> <target_ats> [ai_tool]"
    echo "Example: ./optimize.sh cv.tex jd.txt Workday claude"
    exit 1
fi

CV_FILE=$1
JD_FILE=$2
ATS_TARGET=$3
AI_TOOL=${4:-claude} # Default to claude if not specified

if [ ! -f "$CV_FILE" ]; then
    echo "Error: CV file not found at $CV_FILE"
    exit 1
fi

if [ ! -f "$JD_FILE" ]; then
    echo "Error: JD file not found at $JD_FILE"
    exit 1
fi

CV_TEXT=$(cat "$CV_FILE")
JD_TEXT=$(cat "$JD_FILE")

# Construct the prompt dynamically
PROMPT="You are an expert ATS (Applicant Tracking System) optimizer.
Your goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100% on the ${ATS_TARGET} ATS platform.

Guidelines:
"

if [ "$ATS_TARGET" = "Workday" ]; then
    PROMPT+="- WORKDAY STRATEGY: Workday uses strict exact-keyword matching and heavily penalizes if required JD skills are missing. You MUST use the exact terminology found in the JD without altering the tense or phrasing.
"
elif [ "$ATS_TARGET" = "Taleo" ]; then
    PROMPT+="- TALEO STRATEGY: Taleo relies on high keyword density. You should repeat critical JD keywords 2-3 times across the summary and experience bullet points.
"
elif [ "$ATS_TARGET" = "Greenhouse" ]; then
    PROMPT+="- GREENHOUSE STRATEGY: Greenhouse uses semantic AI and human reviewers. Avoid keyword stuffing. Focus on contextual accomplishments, measurable impact, and natural phrasing.
"
elif [ "$ATS_TARGET" = "iCIMS" ]; then
    PROMPT+="- iCIMS STRATEGY: iCIMS uses a mix of exact and semantic matching. Ensure job titles closely align with the JD requirements and format dates cleanly.
"
fi

PROMPT+="
- CRITICAL KEYWORD INJECTION: You MUST thoroughly analyze the JD for all required hard skills, tools, and keywords (e.g. Kubernetes, Databases, architecture, etc.) and FORCEFULLY inject them into the CV's Skills, Summary, and Experience sections. Do not leave ANY required JD keywords out!
- Incorporate these keywords seamlessly so they read naturally.
- Highlight relevant experience. DO NOT hallucinate fake jobs, degrees, or metrics.
- MATCH ORIGINAL LENGTH: You MUST strictly ensure the final output fits on the exact same number of pages as the original CV. Be concise, consolidate bullet points, and do NOT add fluff.
- CRITICAL: Format the personal details (email, links, location, etc.) in the header on a SINGLE LINE separated by '|' to save vertical space.
- CRITICAL: You must use the EXACT LaTeX structural commands provided in the CV TEXT below. Only modify the actual textual content to optimize for the JD.
- CRITICAL: Do NOT invent or hallucinate new LaTeX commands! Use standard text for the summary and standard \section{} commands.
- CRITICAL: \resumeSubheading is a COMMAND, not an environment! DO NOT use \begin{resumeSubheading}. You MUST provide exactly 4 arguments in curly braces: {Title}{Dates}{Company}{Location}.
- CRITICAL: Do NOT insert random \\ line breaks in normal text or between arguments. LaTeX handles wrapping automatically.
- Ensure all plain-text special characters like % and $ are escaped. Do NOT escape the alignment ampersands (&) inside tabular environments!

OUTPUT FORMAT:
Start your output EXACTLY with \begin{document} and end with \end{document}.
CRITICAL: Wrap the entire output in a single markdown \`\`\`latex code block. Do NOT include the preamble.

--- CV TEXT (USE THIS EXACT LATEX STRUCTURE) ---
${CV_TEXT}

--- JOB DESCRIPTION ---
${JD_TEXT}
"

echo "Running headless CV optimization using $AI_TOOL..."

if [ "$AI_TOOL" = "claude" ]; then
    # Claude Code headless execution
    echo "$PROMPT" | claude -p > optimized_cv.md
elif [ "$AI_TOOL" = "agy" ] || [ "$AI_TOOL" = "antigravity" ]; then
    # Antigravity headless execution
    echo "$PROMPT" | agy -p > optimized_cv.md
elif [ "$AI_TOOL" = "opencode" ]; then
    # OpenCode headless execution
    echo "$PROMPT" | opencode run > optimized_cv.md
else
    echo "Unsupported AI tool: $AI_TOOL"
    exit 1
fi

echo "Done! The optimized CV has been saved to optimized_cv.md"
