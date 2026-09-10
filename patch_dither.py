import re

with open("web/src/App.tsx", "r") as f:
    content = f.read()

# Remove the DitherBackground from the outer container
content = content.replace("<DitherBackground />\n        \n        <div className=\"absolute inset-0", "<div className=\"absolute inset-0")

# And from the previous prompt it might have been different, let's just use regex to remove it
content = re.sub(r'\s*<DitherBackground />\s*', '\n        ', content)

# Now add it inside the p-4 md:p-8 container
target = '<div className="p-4 md:p-8 dither-bg flex-1 flex flex-col gap-4 md:gap-6">'
replacement = '<div className="p-4 md:p-8 flex-1 flex flex-col gap-4 md:gap-6 relative overflow-hidden">\n              <DitherBackground />\n              <div className="absolute inset-0 z-0 bg-[#0b0b0c]/80 pointer-events-none"></div>\n              <div className="relative z-10 flex flex-col h-full">'

content = content.replace(target, replacement)

# Don't forget to close the relative z-10 div
# We need to close it right before the end of the p-4 container
target2 = '</div>\n        </div>\n\n      </div>'
replacement2 = '</div>\n            </div>\n          </div>\n        </div>\n\n      </div>'
# Instead of fragile replacement for closing tags, let's just do it manually

with open("web/src/App.tsx", "w") as f:
    f.write(content)
