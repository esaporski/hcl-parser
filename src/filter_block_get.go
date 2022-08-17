package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

func parseAddress(address string) (string, []string, error) {
	if len(address) == 0 {
		return "", []string{}, fmt.Errorf("failed to parse address: %s", address)
	}

	a := strings.Split(address, ".")
	typeName := a[0]
	labels := []string{}
	if len(a) > 1 {
		labels = a[1:]
	}
	return typeName, labels, nil
}

// findBlocks returns matching blocks from the body that have the given name
// and labels or returns an empty list if there is currently no matching block.
// The labels can be wildcard (*), but numbers of label must be equal.
func findBlocks(b *hclwrite.Body, typeName string, labels []string) []*hclwrite.Block {
	var matched []*hclwrite.Block
	for _, block := range b.Blocks() {
		if typeName == block.Type() {
			labelNames := block.Labels()
			if len(labels) == 0 && len(labelNames) == 0 {
				matched = append(matched, block)
				continue
			}
			if matchLabels(labels, labelNames) {
				matched = append(matched, block)
			}
		}
	}

	return matched
}

// matchLabels returns true only if the matched and false otherwise.
// The labels can be wildcard (*), but numbers of label must be equal.
func matchLabels(lhs []string, rhs []string) bool {
	if len(lhs) != len(rhs) {
		return false
	}

	for i := range lhs {
		if !(lhs[i] == rhs[i] || lhs[i] == "*" || rhs[i] == "*") {
			return false
		}
	}

	return true
}
