package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// NavItems preserves ordered explicit sidebar configuration.
type NavItems []NavItemConfig

// NavItemConfig is a single explicit navigation entry.
type NavItemConfig struct {
	Label string
	Page  string
	Items []NavItemConfig
}

func (n *NavItems) UnmarshalYAML(node *yaml.Node) error {
	if node.Kind != yaml.SequenceNode {
		return fmt.Errorf("nav.items must be a list")
	}

	items := make([]NavItemConfig, 0, len(node.Content))
	for _, child := range node.Content {
		var item NavItemConfig
		if err := item.unmarshalYAMLNode(child); err != nil {
			return err
		}
		items = append(items, item)
	}

	*n = items
	return nil
}

func (n *NavItems) UnmarshalJSON(data []byte) error {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("nav.items must be an array: %w", err)
	}

	items := make([]NavItemConfig, 0, len(raw))
	for _, entry := range raw {
		var item NavItemConfig
		if err := item.unmarshalJSONBytes(entry); err != nil {
			return err
		}
		items = append(items, item)
	}

	*n = items
	return nil
}

func (i *NavItemConfig) unmarshalYAMLNode(node *yaml.Node) error {
	if node.Kind != yaml.MappingNode || len(node.Content) != 2 {
		return fmt.Errorf("each nav item must be a single label-to-value mapping")
	}

	labelNode := node.Content[0]
	valueNode := node.Content[1]
	i.Label = strings.TrimSpace(labelNode.Value)
	if i.Label == "" {
		return fmt.Errorf("nav item label cannot be empty")
	}

	return i.decodeValueYAML(valueNode)
}

func (i *NavItemConfig) decodeValueYAML(node *yaml.Node) error {
	switch node.Kind {
	case yaml.ScalarNode:
		i.Page = strings.TrimSpace(node.Value)
		return nil
	case yaml.MappingNode:
		var payload struct {
			Page  string   `yaml:"page"`
			Items NavItems `yaml:"items"`
		}
		if err := node.Decode(&payload); err != nil {
			return err
		}
		i.Page = strings.TrimSpace(payload.Page)
		i.Items = payload.Items
		return nil
	default:
		return fmt.Errorf("nav item %q must be a path string or an object with page/items", i.Label)
	}
}

func (i *NavItemConfig) unmarshalJSONBytes(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()

	var raw map[string]json.RawMessage
	if err := dec.Decode(&raw); err != nil {
		return fmt.Errorf("nav item must be an object: %w", err)
	}
	if len(raw) != 1 {
		return fmt.Errorf("each nav item must contain exactly one label")
	}

	for label, value := range raw {
		i.Label = strings.TrimSpace(label)
		if i.Label == "" {
			return fmt.Errorf("nav item label cannot be empty")
		}
		return i.decodeValueJSON(value)
	}

	return nil
}

func (i *NavItemConfig) decodeValueJSON(data []byte) error {
	var page string
	if err := json.Unmarshal(data, &page); err == nil {
		i.Page = strings.TrimSpace(page)
		return nil
	}

	var payload struct {
		Page  string   `json:"page"`
		Items NavItems `json:"items"`
	}
	if err := json.Unmarshal(data, &payload); err != nil {
		return fmt.Errorf("nav item %q must be a path string or an object with page/items: %w", i.Label, err)
	}
	i.Page = strings.TrimSpace(payload.Page)
	i.Items = payload.Items
	return nil
}

func (c *Config) Validate() error {
	for _, item := range c.Nav.Items {
		if err := validateNavItem(item, "nav.items"); err != nil {
			return err
		}
	}
	return nil
}

func validateNavItem(item NavItemConfig, path string) error {
	current := path + "." + item.Label
	if strings.TrimSpace(item.Page) == "" && len(item.Items) == 0 {
		return fmt.Errorf("%s must define page, items, or both", current)
	}
	if len(item.Items) == 0 {
		return nil
	}
	for _, child := range item.Items {
		if err := validateNavItem(child, current); err != nil {
			return err
		}
	}
	return nil
}
