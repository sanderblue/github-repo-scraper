package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type CodeSample struct {
	RepoURL   string `json:"repo"`
	CommitSHA string `json:"commit"`
	Path      string `json:"path"`
	Text      string `json:"text"`
}

func main() {
	// Flags for extensions (comma-separated) and output file
	extsFlag := flag.String("ext", "go,py,js", "Comma-separated list of file extensions (no dot)")
	outputFlag := flag.String("out", "code_dataset.jsonl", "Output JSONL file path")
	skipTests := flag.Bool("skip-tests", true, "Skip test files matching *_test.* or *.test.*")
	flag.Parse()

	repos := flag.Args()
	if len(repos) == 0 {
		log.Fatal("Usage: github-scraper [flags] <repo1> <repo2> ...")
	}

	// Prepare extensions set
	exts := parseExtensions(*extsFlag)

	// Open output file
	outFile, err := os.OpenFile(*outputFlag, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open output file: %v", err)
	}
	defer outFile.Close()
	writer := bufio.NewWriter(outFile)

	total := 0
	for _, repoURL := range repos {
		fmt.Printf("Processing repo: %s\n", repoURL)

		tempDir, err := os.MkdirTemp("", "repo-*")
		if err != nil {
			log.Printf("Failed to create temp dir for %s: %v", repoURL, err)
			continue
		}
		defer os.RemoveAll(tempDir)

		// Clone repository
		fmt.Println(" Cloning...")
		cloneCmd := exec.Command("git", "clone", "--depth=1", repoURL, tempDir)
		if out, err := cloneCmd.CombinedOutput(); err != nil {
			log.Printf("Git clone failed for %s: %v\n%s", repoURL, err, string(out))
			continue
		}

		// Get commit SHA
		shaCmd := exec.Command("git", "-C", tempDir, "rev-parse", "HEAD")
		shaBytes, err := shaCmd.CombinedOutput()
		if err != nil {
			log.Printf("Failed to get commit SHA for %s: %v", repoURL, err)
		}
		commitSHA := strings.TrimSpace(string(shaBytes))

		// Walk directory
		err = filepath.WalkDir(tempDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			for _, ext := range exts {
				if strings.HasSuffix(path, "."+ext) {
					name := d.Name()
					// Skip test files if requested
					if *skipTests && (strings.Contains(name, "_test.") || strings.Contains(name, ".test.")) {
						return nil
					}

					content, err := os.ReadFile(path)
					if err != nil {
						log.Printf("Failed to read %s: %v", path, err)
						return nil
					}

					relPath, err := filepath.Rel(tempDir, path)
					if err != nil {
						relPath = path
					}

					sample := CodeSample{
						RepoURL:   repoURL,
						CommitSHA: commitSHA,
						Path:      relPath,
						Text:      string(content),
					}
					line, err := json.Marshal(sample)
					if err != nil {
						log.Printf("JSON marshal error for %s: %v", path, err)
						return nil
					}
					writer.Write(line)
					writer.WriteString("\n")
					total++
					break
				}
			}
			return nil
		})
		if err != nil {
			log.Printf("Error walking %s: %v", repoURL, err)
		}
	}

	writer.Flush()
	fmt.Printf("Done. Wrote %d files to %s\n", total, *outputFlag)
}

func parseExtensions(exts string) []string {
	parts := strings.Split(exts, ",")
	var out []string
	for _, p := range parts {
		e := strings.TrimSpace(p)
		if e != "" {
			out = append(out, e)
		}
	}
	return out
}
