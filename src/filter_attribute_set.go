package main

import (
	"fmt"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

// buildExpression returns a new expressions for a given name and value of attribute.
// At the time of wrting this, there is no way to parse expression from string.
// So we generate a temporarily config on memory and parse it, and extract a generated expression.
func buildExpression(name string, value string) (*hclwrite.Expression, error) {
	src := name + " = " + value
	f, err := safeParseConfig([]byte(src), "generated_by_buildExpression", hcl.Pos{Line: 1, Column: 1})
	if err != nil {
		return nil, fmt.Errorf("failed to build expression at the parse phase: %s", err)
	}

	attr := f.Body().GetAttribute(name)
	if attr == nil {
		return nil, fmt.Errorf("failed to build expression at the get phase. name = %s, value = %s", name, value)
	}

	return attr.Expr(), nil
}
