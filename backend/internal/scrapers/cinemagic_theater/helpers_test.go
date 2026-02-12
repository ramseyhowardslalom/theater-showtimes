package cinemagic_theater

import "testing"

func TestNormalizeTitleForTMDB(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard title unchanged",
			input:    "The Matrix",
			expected: "The Matrix",
		},
		{
			name:     "remove parenthetical suffix",
			input:    "Blade Runner (Director's Cut)",
			expected: "Blade Runner",
		},
		{
			name:     "remove dash suffix",
			input:    "Jaws - 45th Anniversary",
			expected: "Jaws",
		},
		{
			name:     "remove year indicator",
			input:    "The Matrix (1999)",
			expected: "The Matrix",
		},
		{
			name:     "remove multiple suffixes",
			input:    "Star Wars (1977) - Special Edition",
			expected: "Star Wars",
		},
		{
			name:     "trim whitespace",
			input:    "  The Godfather  ",
			expected: "The Godfather",
		},
		{
			name:     "preserve short dash titles",
			input:    "X - Y",
			expected: "X - Y",
		},
		{
			name:     "empty parentheses",
			input:    "Movie Title ()",
			expected: "Movie Title",
		},
		{
			name:     "parentheses at start preserved",
			input:    "(500) Days of Summer",
			expected: "(500) Days of Summer",
		},
		{
			name:     "complex case",
			input:    "The Lord of the Rings: The Fellowship of the Ring (2001) - Extended Edition",
			expected: "The Lord of the Rings: The Fellowship of the Ring",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeTitleForTMDB(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeTitleForTMDB(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestNormalizeFormat(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "digital lowercase",
			input:    "digital",
			expected: "digital",
		},
		{
			name:     "digital uppercase",
			input:    "DIGITAL",
			expected: "digital",
		},
		{
			name:     "digital with spaces",
			input:    "  Digital  ",
			expected: "digital",
		},
		{
			name:     "35mm lowercase",
			input:    "35mm",
			expected: "35mm",
		},
		{
			name:     "35mm uppercase",
			input:    "35MM",
			expected: "35mm",
		},
		{
			name:     "35 with space",
			input:    "35 mm",
			expected: "35mm",
		},
		{
			name:     "70mm lowercase",
			input:    "70mm",
			expected: "70mm",
		},
		{
			name:     "70mm uppercase",
			input:    "70MM",
			expected: "70mm",
		},
		{
			name:     "IMAX uppercase",
			input:    "IMAX",
			expected: "IMAX",
		},
		{
			name:     "IMAX lowercase",
			input:    "imax",
			expected: "IMAX",
		},
		{
			name:     "IMAX mixed case",
			input:    "IMax",
			expected: "IMAX",
		},
		{
			name:     "unknown format defaults to digital",
			input:    "4K Ultra HD",
			expected: "digital",
		},
		{
			name:     "empty string defaults to digital",
			input:    "",
			expected: "digital",
		},
		{
			name:     "35mm in text",
			input:    "Projected in 35mm film",
			expected: "35mm",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeFormat(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeFormat(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractMovieSlugFromURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "standard movie URL",
			input:    "https://tickets.thecinemagictheater.com/movie/sentimental-value",
			expected: "sentimental-value",
		},
		{
			name:     "movie URL with trailing slash",
			input:    "https://tickets.thecinemagictheater.com/movie/blade-runner/",
			expected: "blade-runner",
		},
		{
			name:     "relative URL",
			input:    "/movie/the-matrix",
			expected: "the-matrix",
		},
		{
			name:     "slug with numbers",
			input:    "https://tickets.thecinemagictheater.com/movie/2001-a-space-odyssey",
			expected: "2001-a-space-odyssey",
		},
		{
			name:     "complex slug",
			input:    "https://tickets.thecinemagictheater.com/movie/star-wars-a-new-hope-1977",
			expected: "star-wars-a-new-hope-1977",
		},
		{
			name:     "URL without movie path returns empty",
			input:    "https://tickets.thecinemagictheater.com/now-showing",
			expected: "",
		},
		{
			name:     "empty URL returns empty",
			input:    "",
			expected: "",
		},
		{
			name:     "URL with query parameters",
			input:    "https://tickets.thecinemagictheater.com/movie/inception?date=2026-02-12",
			expected: "inception?date=2026-02-12",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractMovieSlugFromURL(tt.input)
			if result != tt.expected {
				t.Errorf("extractMovieSlugFromURL(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
