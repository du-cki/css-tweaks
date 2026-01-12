package main

import (
	"fmt"
	"hash/crc32"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/flosch/pongo2/v6"
)

var LookUp SnippetLookUp

var (
	Branch     = "dev"
	CommitHash = "unknown"
	Origin     = "unknown"
)

const UniverseSize = 4096

func init() {
	LookUp = SnippetLookUp{
		Flat:    make(map[int]Snippet),
		Grouped: make(map[string][]Snippet),
	}

	Unwrap(struct{}{}, filepath.WalkDir("snippets", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() || filepath.Ext(path) != ".css" {
			return nil
		}

		relPath := Unwrap(filepath.Rel("snippets", path))
		relPath = filepath.ToSlash(relPath)

		hash := crc32.ChecksumIEEE([]byte(relPath))
		id := int(hash % UniverseSize)

		if existing, exists := LookUp.Flat[id]; exists {
			log.Fatalf("Collision Detected!\nFile 1: %s\nFile 2: %s\nBoth map to ID: %d. Rename one file to resolve.", existing.Name, relPath, id)
		}

		content := Unwrap(os.ReadFile(path))

		parts := strings.Split(relPath, "/")
		category := "misc"
		if len(parts) > 1 {
			category = parts[0]
		}

		snippet := Snippet{
			ID:       id,
			Name:     d.Name(),
			Category: category,
			Content:  content,
		}

		LookUp.Flat[id] = snippet
		LookUp.Grouped[category] = append(LookUp.Grouped[category], snippet)

		return nil
	}))

	for cat := range LookUp.Grouped {
		sort.Slice(LookUp.Grouped[cat], func(i, j int) bool {
			return LookUp.Grouped[cat][i].Name < LookUp.Grouped[cat][j].Name
		})
	}

	log.Printf("loaded %d snippets", len(LookUp.Flat))
}

func main() {
	tpl := Unwrap(pongo2.FromFile("templates/home.html.jinja"))
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Because http.ServeMux catches any fall
		// through routes with this route, kinda dumb.
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		err := tpl.ExecuteWriter(
			pongo2.Context{
				"snippets":   LookUp.Grouped,
				"repository": Repository{Origin: Origin, Commit: CommitHash, Branch: Branch},
			}, w,
		)

		if err != nil {
			http.Error(w, "error rendering template", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/s/{slug}", func(w http.ResponseWriter, r *http.Request) {
		slug := r.PathValue("slug")
		parts := strings.Split(slug, ",")

		w.Header().Set("Content-Type", "text/css")

		for _, p := range parts {
			id, err := strconv.Atoi(p)
			if err != nil {
				continue
			}

			if snippet, ok := LookUp.Flat[id]; ok {
				fmt.Fprintf(w, "/* %s */\n", snippet.Name)
				w.Write(snippet.Content)
				w.Write([]byte("\n"))
			}
		}
	})

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	log.Printf("server running on http://localhost:8080")
	Unwrap(struct{}{}, http.ListenAndServe(":8080", nil))
}
