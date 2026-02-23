package provider

import (
	"strings"
)

// Extracts the date string from a description using the group date prefix.
func ExtractDate(description string) string {
	groupDatePrefix := "group date: "
	descLower := strings.ToLower(description)
	dateIndex := strings.Index(descLower, groupDatePrefix)
	if dateIndex > -1 {
		extracted := SliceUntilWhitespace(description, dateIndex+len(groupDatePrefix))
		return AddLeadingZerosToDate(extracted)
	}
	return "unknown"
}

// Returns the substring from startIndex until the next whitespace.
func SliceUntilWhitespace(str string, startIndex int) string {
	if startIndex >= len(str) {
		return ""
	}
	substring := str[startIndex:]
	for i, c := range substring {
		if c == ' ' || c == '\t' || c == '\n' || c == '\r' {
			return strings.TrimSpace(substring[:i])
		}
	}
	return strings.TrimSpace(substring)
}

// Adds leading zeros to MM/DD/YYYY date strings.
func AddLeadingZerosToDate(date string) string {
	if len(date) != 10 {
		parts := strings.Split(date, "/")
		if len(parts) != 3 {
			return date
		}
		mm := parts[0]
		if len(mm) == 1 {
			mm = "0" + mm
		}
		dd := parts[1]
		if len(dd) == 1 {
			dd = "0" + dd
		}
		return mm + "/" + dd + "/" + parts[2]
	}
	return date
}
