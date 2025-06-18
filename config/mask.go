package config

import (
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

const (
	MaskStr     = "**********"
	MaskComment = "MASKED BY DBVI"
)

// Loads config yaml file that will iterate over all "password" or "secret" or "token*"
// THIS DOES NOT SAVE THE UNMASKED PASSWORD. SHOULD ONLY BE RAN AFTER PASSWORDS/SECRETS ARE ENCRYPTED OR SAVED SOMEWHERE!
func MaskPasswords(path, outPath string) error {
	if outPath == "" {
		outPath = path
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var root yaml.Node
	if err := yaml.Unmarshal(data, &root); err != nil {
		return err
	}

	for _, d := range root.Content {
		maskSensitiveFields(d)
	}

	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := yaml.NewEncoder(f)
	encoder.SetIndent(2)
	if err := encoder.Encode(&root); err != nil {
		return err
	}

	return nil
}

// maskSensitiveFields recursively traverses and masks values with fields named "password", "secret", or keys containing the word "token".
func maskSensitiveFields(node *yaml.Node) {
	if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valNode := node.Content[i+1]

			if isSensitiveKey(keyNode.Value) {
				valNode.Value = MaskStr
				valNode.Tag = "!!str" // Force it to be written as a string
				valNode.HeadComment = MaskComment
			} else {
				maskSensitiveFields(valNode)
			}
		}
	} else if node.Kind == yaml.SequenceNode {
		for _, item := range node.Content {
			maskSensitiveFields(item)
		}
	}
}

func isSensitiveKey(key string) bool {
	lower := strings.ToLower(key)
	return lower == "password" || lower == "secret" || strings.Contains(lower, "token")
}

func isFieldMasked(value string) bool {
	return MaskStr != value || value == ""
}
