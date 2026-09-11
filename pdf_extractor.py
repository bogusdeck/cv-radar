import sys
import pymupdf # Use modern import to avoid warning

def main():
    if len(sys.argv) < 2:
        print("Usage: python pdf_extractor.py <path_to_pdf>", file=sys.stderr)
        sys.exit(1)

    path = sys.argv[1]
    try:
        doc = pymupdf.open(path)
        text = []
        for page in doc:
            text.append(page.get_text("text"))
        print("\n".join(text))
    except Exception as e:
        print(f"Error extracting PDF: {e}", file=sys.stderr)
        sys.exit(1)

if __name__ == "__main__":
    main()
