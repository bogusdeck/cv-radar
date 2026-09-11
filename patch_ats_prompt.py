import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

# I will insert ATS specific rules into the prompt
old_prompt_start = r"""  const promptText = `You are an expert ATS (Applicant Tracking System) optimizer.
Your goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100% on the ${result.platform} ATS platform.

Guidelines:
- CRITICAL KEYWORD INJECTION: You MUST thoroughly analyze the JD for all required hard skills, tools, and keywords (e.g. Kubernetes, Databases, architecture, etc.) and FORCEFULLY inject them into the CV's Skills, Summary, and Experience sections. Do not leave ANY required JD keywords out!"""

new_prompt_start = r"""  let atsSpecificRules = "";
  if (result.platform === "Workday") {
    atsSpecificRules = "- WORKDAY STRATEGY: Workday uses strict exact-keyword matching and heavily penalizes if required JD skills are missing. You MUST use the exact terminology found in the JD without altering the tense or phrasing.";
  } else if (result.platform === "Taleo") {
    atsSpecificRules = "- TALEO STRATEGY: Taleo relies on high keyword density. You should repeat critical JD keywords 2-3 times across the summary and experience bullet points.";
  } else if (result.platform === "Greenhouse") {
    atsSpecificRules = "- GREENHOUSE STRATEGY: Greenhouse uses semantic AI and human reviewers. Avoid keyword stuffing. Focus on contextual accomplishments, measurable impact, and natural phrasing.";
  } else if (result.platform === "iCIMS") {
    atsSpecificRules = "- iCIMS STRATEGY: iCIMS uses a mix of exact and semantic matching. Ensure job titles closely align with the JD requirements and format dates cleanly.";
  }

  const promptText = `You are an expert ATS (Applicant Tracking System) optimizer.
Your goal is to rewrite the provided CV so that it perfectly matches the provided Job Description (JD) and scores 100% on the ${result.platform} ATS platform.

Guidelines:
${atsSpecificRules}
- CRITICAL KEYWORD INJECTION: You MUST thoroughly analyze the JD for all required hard skills, tools, and keywords (e.g. Kubernetes, Databases, architecture, etc.) and FORCEFULLY inject them into the CV's Skills, Summary, and Experience sections. Do not leave ANY required JD keywords out!"""

text = text.replace(old_prompt_start, new_prompt_start)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

