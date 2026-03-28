// Package templates provides utilities for loading and parsing HTML templates.
package templates

import (
	"embed"
	"fmt"
	"html/template"
	"strings"
)

// FuncMap returns the default template functions used across all templates.
func FuncMap() template.FuncMap {
	return template.FuncMap{
		"lower":     strings.ToLower,
		"upper":     strings.ToUpper,
		"hasSuffix": strings.HasSuffix,
		"hasPrefix": strings.HasPrefix,
		"trimSpace": strings.TrimSpace,
		"dict":      dictFunc,
		"deref":     derefBool,
		"sub":       func(a, b int) int { return a - b },
	}
}

// derefBool dereferences a *bool pointer, returning false if nil.
func derefBool(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

// dictFunc creates a map from alternating key/value pairs.
// Example: dict "key1" value1 "key2" value2
func dictFunc(values ...any) (map[string]any, error) {
	if len(values)%2 != 0 {
		return nil, fmt.Errorf("dict expects an even number of arguments, got %d", len(values))
	}
	m := make(map[string]any, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, fmt.Errorf("dict expects string keys, got %T at index %d", values[i], i)
		}
		m[key] = values[i+1]
	}
	return m, nil
}

// ParseFS parses templates from an embedded filesystem.
// It returns a template with the given name and all specified patterns parsed.
func ParseFS(fs embed.FS, name string, patterns ...string) (*template.Template, error) {
	tmpl := template.New(name).Funcs(FuncMap())
	tmpl, err := tmpl.ParseFS(fs, patterns...)
	if err != nil {
		return nil, fmt.Errorf("parse templates (%s): %w", name, err)
	}
	return tmpl, nil
}
