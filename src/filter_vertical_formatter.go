package main

import (
	"github.com/hashicorp/hcl/v2/hclsyntax"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// verticalFormat formats token in vertical.
func verticalFormat(tokens hclwrite.Tokens) hclwrite.Tokens {
	trimmed := trimLeadingNewLine(tokens)
	removed := removeDuplicatedNewLine(trimmed)
	return removed
}

// trimLeadingNewLine trims leading newlines from tokens.
// We don't need trim trailing newlines because the last newline should be
// kept and others are removed removeDuplicatedNewLine.
func trimLeadingNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
	begin := 0
	for ; begin < len(tokens); begin++ {
		if tokens[begin].Type != hclsyntax.TokenNewline {
			break
		}
	}

	return tokens[begin:]
}

// removeDuplicatedNewLine removes duplicated newlines
// Two consecutive blank lines should be removed.
// In other words, if there are three consecutive TokenNewline tokens,
// the third and subsequent TokenNewline tokens are removed.
func removeDuplicatedNewLine(tokens hclwrite.Tokens) hclwrite.Tokens {
	var removed hclwrite.Tokens
	beforeBefore := false
	before := false

	for _, token := range tokens {
		if token.Type != hclsyntax.TokenNewline {
			removed = append(removed, token)
			// reset
			beforeBefore = false
			before = false
			continue
		}
		// TokenNewLine
		if before && beforeBefore {
			// skip duplicated newlines
			continue
		}
		removed = append(removed, token)
		beforeBefore = before
		before = true
	}

	return removed
}
