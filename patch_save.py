import re

with open("web/src/App.tsx", "r") as f:
    text = f.read()

old_download = """                    const blob = await compileRes.blob();
                    const url = URL.createObjectURL(blob);
                    const a = document.createElement('a');
                    a.href = url;
                    a.download = `Optimized_Resume_${result.platform}.pdf`;
                    a.click();"""

new_download = """                    const blob = await compileRes.blob();
                    try {
                      if ('showSaveFilePicker' in window) {
                        const handle = await (window as any).showSaveFilePicker({
                          suggestedName: `Optimized_Resume_${result.platform}.pdf`,
                          types: [{
                            description: 'PDF Document',
                            accept: {'application/pdf': ['.pdf']},
                          }],
                        });
                        const writable = await handle.createWritable();
                        await writable.write(blob);
                        await writable.close();
                      } else {
                        throw new Error('Fallback');
                      }
                    } catch (err: any) {
                      if (err.name !== 'AbortError') {
                        const url = URL.createObjectURL(blob);
                        const a = document.createElement('a');
                        a.href = url;
                        a.download = `Optimized_Resume_${result.platform}.pdf`;
                        a.click();
                      }
                    }"""

text = text.replace(old_download, new_download)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

