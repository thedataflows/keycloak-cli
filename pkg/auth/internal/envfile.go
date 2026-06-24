package internal

import (
	"os"
	"regexp"
	"strings"
)

func writeKeyValue(filePath, key, value, separator string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return os.WriteFile(filePath, formatKeyValue(key, value, separator, "\n"), 0o644)
		}
		return err
	}

	contentString := string(content)
	newline := "\n"
	if strings.Contains(contentString, "\r\n") {
		newline = "\r\n"
	}

	matcher := regexp.MustCompile("(?m)^" + regexp.QuoteMeta(key) + "\\s*" + regexp.QuoteMeta(separator) + ".*$")
	if matcher.MatchString(contentString) {
		contentString = matcher.ReplaceAllString(contentString, key+separator+value)
		return os.WriteFile(filePath, []byte(contentString), 0o644)
	}

	var builder strings.Builder
	builder.WriteString(contentString)
	if contentString != "" && !strings.HasSuffix(contentString, newline) {
		builder.WriteString(newline)
	}
	builder.Write(formatKeyValue(key, value, separator, newline))

	return os.WriteFile(filePath, []byte(builder.String()), 0o644)
}

func formatKeyValue(key, value, separator, newline string) []byte {
	var builder strings.Builder
	builder.WriteString(key)
	builder.WriteString(separator)
	builder.WriteString(value)
	builder.WriteString(newline)
	return []byte(builder.String())
}
