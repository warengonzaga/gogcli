package cmd

import "strings"

func resolveLabelIDs(labels []string, nameToID map[string]string) []string {
	if len(labels) == 0 {
		return nil
	}
	out := make([]string, 0, len(labels))
	for _, label := range labels {
		trimmed := strings.TrimSpace(label)
		if trimmed == "" {
			continue
		}
		if nameToID != nil {
			if id, ok := nameToID[strings.ToLower(trimmed)]; ok {
				out = append(out, id)
				continue
			}
		}
		out = append(out, trimmed)
	}
	return out
}
