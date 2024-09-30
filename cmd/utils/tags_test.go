package utils

import (
	"path/filepath"
	"testing"

	"github.com/kkwon1/apod-forum-backend/cmd/models"
)


func TestExtractTags(t *testing.T) {
	absPath, _ := filepath.Abs("../../internal/const/test_astro_terms.txt")
	InitTagsWithFilePath(absPath)
	tests := []struct {
		name     string
		apod     models.Apod
		expected []string
	}{
		{
			name: "Single match",
			apod: models.Apod{Explanation: "The galaxy is vast and beautiful."},
			expected: []string{"galaxy"},
		},
		{
			name: "Multiple matches",
			apod: models.Apod{Explanation: "The galaxy and nebula are amazing."},
			expected: []string{"galaxy", "nebula"},
		},
		{
			name: "No matches",
			apod: models.Apod{Explanation: "The ocean is deep and mysterious."},
			expected: []string{},
		},
		{
			name: "Case insensitive match",
			apod: models.Apod{Explanation: "The GALAXY is vast."},
			expected: []string{"galaxy"},
		},
		{
			name: "Duplicate words",
			apod: models.Apod{Explanation: "The galaxy galaxy is vast."},
			expected: []string{"galaxy"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTags(tt.apod)
			if !equal(result, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func equal(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	m := make(map[string]struct{}, len(a))
	for _, v := range a {
		m[v] = struct{}{}
	}
	for _, v := range b {
		if _, ok := m[v]; !ok {
			return false
		}
	}
	return true
}
