package main

import (
	"regexp"
	"strings"
)

func ExtractPackageNames(input string) []string {
	// Use a regular expression to find all words within brackets
	re := regexp.MustCompile(`\[(.*?)\]`)
	matches := re.FindAllStringSubmatch(input, -1)

	// Use a map to avoid duplicates
	packageMap := make(map[string]struct{})

	for _, match := range matches {
		if len(match) > 1 {
			// Split the captured string by spaces and trim whitespace
			parts := strings.Fields(match[1])
			for _, part := range parts {
				packageMap[strings.TrimSpace(part)] = struct{}{}
			}
		}
	}

	// Create a slice to hold the unique package names
	var packages []string
	for pkg := range packageMap {
		packages = append(packages, pkg)
	}

	return packages
}
