import re

with open("cmd/server/llm.go", "r") as f:
    text = f.read()

target = r"""		reqBody, _ := json.Marshal(map[string]interface{}{
			"model":  model,
			"prompt": prompt,
			"stream": false,
		})"""

new_code = r"""		reqBody, _ := json.Marshal(map[string]interface{}{
			"model":  model,
			"prompt": prompt,
			"stream": false,
			"options": map[string]interface{}{
				"num_ctx": 8192,
				"num_predict": 4096,
			},
		})"""

text = text.replace(target, new_code)

with open("cmd/server/llm.go", "w") as f:
    f.write(text)

