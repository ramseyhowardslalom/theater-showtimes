package cinemagic_theater

import (
	"regexp"
	"strings"

	"github.com/gocolly/colly/v2"
)

// normalizeTitleForTMDB cleans movie titles to improve TMDB API match success
// Removes parenthetical suffixes, dash suffixes, and year indicators
//
// Examples:
//   - "Blade Runner (Director's Cut)" -> "Blade Runner"
//   - "Movie Title - 25th Anniversary" -> "Movie Title"
//   - "The Matrix (1999)" -> "The Matrix"
func normalizeTitleForTMDB(rawTitle string) string {
	title := strings.TrimSpace(rawTitle)

	// Remove parenthetical suffixes
	if idx := strings.Index(title, "("); idx > 0 {
		base := strings.TrimSpace(title[:idx])
		if len(base) > 0 {
			title = base
		}
	}

	// Remove dash suffixes (e.g., " - 25th Anniversary")
	if idx := strings.Index(title, " - "); idx > 0 {
		base := strings.TrimSpace(title[:idx])
		if len(base) > 3 { // Ensure we don't strip actual title
			title = base
		}
	}

	// Remove year indicators (e.g., " (2023)")
	yearPattern := regexp.MustCompile(`\s*\(\d{4}\)`)
	title = yearPattern.ReplaceAllString(title, "")

	return strings.TrimSpace(title)
}

// extractFilmFormat detects film format from movie page badges
// Tries multiple selectors in order of likelihood, returns "digital" as default
//
// Supported formats: digital, 35mm, 70mm, IMAX
func extractFilmFormat(e *colly.HTMLElement) string {
	// Try multiple selectors in order of likelihood
	selectors := []string{
		".format-badge",
		".film-format",
		"[data-format]",
		".badge:contains('mm')",
		".badge:contains('igital')",
	}

	for _, selector := range selectors {
		selection := e.DOM.Find(selector).First()
		if selection.Length() > 0 {
			format := selection.Text()
			if format != "" {
				return normalizeFormat(format)
			}
		}
	}

	// Default to digital if no format found
	return "digital"
}

// normalizeFormat standardizes film format strings
// Converts various format representations to canonical forms
//
// Examples:
//   - "35MM", "35 mm", "35mm" -> "35mm"
//   - "Digital", "DIGITAL" -> "digital"
//   - "IMAX" -> "IMAX"
func normalizeFormat(raw string) string {
	raw = strings.ToLower(strings.TrimSpace(raw))

	if strings.Contains(raw, "35") || strings.Contains(raw, "35mm") {
		return "35mm"
	}
	if strings.Contains(raw, "70") || strings.Contains(raw, "70mm") {
		return "70mm"
	}
	if strings.Contains(raw, "imax") {
		return "IMAX"
	}

	return "digital"
}

// extractMovieSlugFromURL extracts the movie slug from a Cinemagic movie URL
// Example: "https://tickets.thecinemagictheater.com/movie/sentimental-value" -> "sentimental-value"
func extractMovieSlugFromURL(url string) string {
	parts := strings.Split(url, "/movie/")
	if len(parts) == 2 {
		// Remove trailing slash if present
		slug := strings.TrimSuffix(parts[1], "/")
		return slug
	}
	return ""
}
