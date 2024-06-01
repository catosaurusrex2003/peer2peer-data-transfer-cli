package utils

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/c-bata/go-prompt"
)

func Completer(d prompt.Document) []prompt.Suggest {
	path := d.TextBeforeCursor()
	if path == "" {
		path = "."
	}

	dir := filepath.Dir(path)
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}

	suggestions := []prompt.Suggest{}
	for _, f := range files {
		if strings.HasPrefix(f.Name(), filepath.Base(path)) {
			fullPath := filepath.Join(dir, f.Name())
			if f.IsDir() {
				fullPath += string(os.PathSeparator)
			}
			suggestions = append(suggestions, prompt.Suggest{Text: fullPath})
		}
	}
	return suggestions
}
