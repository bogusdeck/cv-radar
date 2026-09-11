with open("web/src/App.tsx", "r") as f:
    text = f.read()

target = r"""      </div>
    </div>
  )
}"""

# Actually, the last </div> in DetailPanel is:
#        )}
#      </div>
#    </div>
#  )
#}

# Let's just find the `DetailPanel` component and insert it before the last `</div>`.

start_idx = text.find('function DetailPanel')
end_idx = text.find('function PixelSelect', start_idx)

if end_idx == -1:
    end_idx = text.find('function SettingsModal', start_idx)

detail_panel = text[start_idx:end_idx]

# Remove any previously wrongly injected stuff if we did
# (My previous python script matched `    </div>\s*$` but it didn't do anything because it didn't match!)

popup_ui = r"""
        {errorMsg && (
          <div className="fixed inset-0 bg-black/90 flex items-center justify-center z-[100]">
            <div className="pixel-border p-8 max-w-md w-full bg-[#111] text-center border-red-500 shadow-[8px_8px_0_#ef4444]">
              <h2 className="font-PressStart text-red-500 text-[16px] mb-4 drop-shadow-[2px_2px_0_#000]">COMPILER ERROR</h2>
              <p className="font-mono text-gray-300 text-sm mb-8 leading-relaxed">{errorMsg}</p>
              <button 
                onClick={() => setErrorMsg('')}
                className="nes-btn is-error w-full font-PressStart text-[10px]"
              >
                ACKNOWLEDGE
              </button>
            </div>
          </div>
        )}
      </div>
    </div>
  )
}
"""

new_detail_panel = detail_panel.replace("""      </div>\n    </div>\n  )\n}""", popup_ui)

text = text.replace(detail_panel, new_detail_panel)

with open("web/src/App.tsx", "w") as f:
    f.write(text)

