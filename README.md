# GitHub Repo Scraper

A lightweight command-line tool written in Go that clones one or more public GitHub repositories, extracts source files by extension (e.g., `.go`, `.py`, `.js`), enriches each entry with repo URL, commit SHA, and file path, and serializes them into a JSON-lines dataset for AI training tasks.

---

## Features

- **Batch Processing**: Accept multiple repositories in one invocation.
- **Multi-Language Support**: Specify file extensions via `-ext` flag (default: `go,py,js`).
- **Metadata Enrichment**: Each JSONL entry includes `repo`, `commit`, and `path`.
- **Test-File Filtering**: Skip test files by default; toggle via `-skip-tests`.
- **Custom Output**: Define output file path with `-out` flag.

---

## Installation

1. **Prerequisites**
   - Go **1.18+** installed and in your `PATH`.
   - Git CLI available.

2. **Build the binary**
   ```bash
   git clone https://github.com/yourusername/github-repo-scraper.git
   cd github-repo-scraper
   go build -o github-repo-scraper main.go
   ```

## Usage
Process a single repository with default settings
```bash
./github-repo-scraper https://github.com/<user>/<repo>.git
```

# Batch-process multiple repositories, custom extensions, and specify output
```bash
./github-repo-scraper \
  -ext go,py,js \
  -out code_dataset.jsonl \
  -skip-tests=false \
  https://github.com/user/repo1.git \
  https://github.com/user/repo2.git
```
