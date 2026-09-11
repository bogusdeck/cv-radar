import re

with open("cmd/server/fix_resume.go", "r") as f:
    text = f.read()

# We need to prepend the preamble string
preamble_str = r"""const Preamble = `%-------------------------
% Resume in Latex
%------------------------
\documentclass[letterpaper,10.8pt]{article}
\usepackage{latexsym}
\usepackage[empty]{fullpage}
\usepackage{titlesec}
\usepackage{marvosym}
\usepackage[usenames,dvipsnames]{color}
\usepackage{verbatim}
\usepackage{enumitem}
\usepackage[hidelinks]{hyperref}
\usepackage{fancyhdr}
\usepackage[english]{babel}
\usepackage{tabularx}
\input{glyphtounicode}
\pagestyle{fancy}
\fancyhf{}
\fancyfoot{}
\renewcommand{\headrulewidth}{0pt}
\renewcommand{\footrulewidth}{0pt}
\addtolength{\oddsidemargin}{-0.5in}
\addtolength{\evensidemargin}{-0.5in}
\addtolength{\textwidth}{1in}
\addtolength{\topmargin}{-.5in}
\addtolength{\textheight}{1.1in}
\urlstyle{same}
\raggedbottom
\raggedright
\setlength{\tabcolsep}{0in}
\titleformat{\section}{
  \vspace{-6pt}\bfseries\scshape\raggedright\large
}{}{0em}{}[\color{black}\titlerule \vspace{-5pt}]
\pdfgentounicode=1

%-------------------------
% Custom commands
\newcommand{\resumeItem}[1]{
  \item\small{
    {#1 \vspace{-2pt}}
  }
}
\newcommand{\resumeSubheading}[4]{
  \vspace{-2pt}\item
    \begin{tabular*}{0.97\textwidth}[t]{l@{\extracolsep{\fill}}r}
      \textbf{#1} & #2 \\
      \textit{\small#3} & \textit{\small #4} \\
    \end{tabular*}\vspace{-7pt}
}
\newcommand{\resumeSubSubheading}[2]{
    \item
    \begin{tabular*}{0.97\textwidth}{l@{\extracolsep{\fill}}r}
      \textit{\small#1} & \textit{\small #2} \\
    \end{tabular*}\vspace{-7pt}
}
\newcommand{\resumeProjectHeading}[2]{
    \item
    \begin{tabular*}{0.97\textwidth}{l@{\extracolsep{\fill}}r}
      \small#1 & #2 \\
    \end{tabular*}\vspace{-7pt}
}
\newcommand{\resumeSubItem}[1]{\resumeItem{#1}\vspace{-4pt}}
\renewcommand\labelitemii{$\vcenter{\hbox{\tiny$\bullet$}}$}
\newcommand{\resumeSubHeadingListStart}{\begin{itemize}[leftmargin=0.15in, label={}]}
\newcommand{\resumeSubHeadingListEnd}{\end{itemize}}
\newcommand{\resumeItemListStart}{\begin{itemize}[topsep=2pt, itemsep=0pt, parsep=0pt]}
\newcommand{\resumeItemListEnd}{\end{itemize}\vspace{-5pt}}

`

func extractLatex(output string) string {
	start := strings.Index(output, "\\begin{document}")
	if start != -1 {
		endStr := "\\end{document}"
		end := strings.LastIndex(output, endStr)
		if end != -1 {
			output = output[start : end+len(endStr)]
		} else {
			output = output[start:]
		}
	} else {
	    // If LLM still outputs \documentclass, find it
	    startDoc := strings.Index(output, "\\documentclass")
	    if startDoc != -1 {
	        return output[startDoc:]
	    }
	}
	
	// Remove common problematic unicode characters
	output = strings.ReplaceAll(output, "≈", "~")
	output = strings.ReplaceAll(output, "–", "--") // en-dash
	output = strings.ReplaceAll(output, "—", "---") // em-dash

	return Preamble + output
}
"""

old_extract = r"func extractLatex.*?return output\n\}"
text = re.sub(old_extract, preamble_str, text, flags=re.DOTALL)

with open("cmd/server/fix_resume.go", "w") as f:
    f.write(text)

