package utils

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/kkwon1/apod-forum-backend/cmd/models"
)

var astroTerms map[string]struct{}

func InitTags() {
	absPath, _ := filepath.Abs("internal/const/astro_terms.txt")
	astroTerms = loadAstroTerms(absPath)
}

func InitTagsWithFilePath(filePath string) {
	astroTerms = loadAstroTerms(filePath)
}

func loadAstroTerms(filePath string) (map[string]struct{}) {
	file, err := os.Open(filePath)
    if err != nil {
        log.Fatalf("failed to open file: %s", err)
    }
    defer file.Close()

    set := make(map[string]struct{})
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        set[strings.ToLower(scanner.Text())] = struct{}{}
    }

    if err := scanner.Err(); err != nil {
        log.Fatalf("failed to scan file: %s", err)
    }

    // Print the set
    return set
}

func ExtractTags(apod models.Apod) []string {
	words := strings.Fields(strings.ToLower(apod.Explanation))
	matches := make(map[string]struct{})
	for _, word := range words {
		if _, ok := astroTerms[word]; ok {
			matches[word] = struct{}{}
		}
	}
	return setToList(matches)
}

func setToList(set map[string]struct{}) []string {
	list := make([]string, 0, len(set))
	for key := range set {
		list = append(list, key)
	}
	return list
}