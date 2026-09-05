package utils

import (
	"regexp"
	"sort"
)

var fieldRegex = regexp.MustCompile(`\{\{\s*(?:if\s+|range\s+|with\s+)?\.([a-zA-Z0-9_]+)`)

// ExtractTemplateFields extracts all placeholder field names like {{.FieldName}} from the template content.
func ExtractTemplateFields(content string) []string {
	matches := fieldRegex.FindAllStringSubmatch(content, -1)
	if len(matches) == 0 {
		return []string{}
	}

	unique := make(map[string]struct{})
	for _, m := range matches {
		if len(m) > 1 && m[1] != "" {
			unique[m[1]] = struct{}{}
		}
	}

	fields := make([]string, 0, len(unique))
	for f := range unique {
		fields = append(fields, f)
	}
	sort.Strings(fields)
	return fields
}
