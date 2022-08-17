package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Filter reads HCL and appends a new attribute to a given address.
// If a matched block not found, nothing happens.
// If the given attribute already exists, it returns an error.
// If a newline flag is true, it also appends a newline before the new attribute.
func AppendAttribute(inFile *hclwrite.File, address string, value string, addNewline bool) error {
	body := inFile.Body()

	addressList := strings.Split(address, ".")
	if len(addressList) > 1 {
		// if address contains dots, the last element is an attribute name,
		// and the rest is the address of the block.
		address = addressList[len(addressList)-1]
		blockAddr := strings.Join(addressList[:len(addressList)-1], ".")
		blocks, err := findLongestMatchingBlocks(body, blockAddr)
		if err != nil {
			return err
		}

		if len(blocks) == 0 {
			// attribute not found does not return an error
			return nil
		}

		// Use first matching one.
		body = blocks[0].Body()
		if body.GetAttribute(address) != nil {
			return fmt.Errorf("attribute already exists ~> %s", address)
		}
	}

	// To delegate expression parsing to the hclwrite parser,
	// We build a new expression and set back to the attribute by tokens.
	expr, err := buildExpression(address, value)
	if err != nil {
		return err
	}

	if addNewline {
		body.AppendNewline()
	}
	body.SetAttributeRaw(address, expr.BuildTokens(nil))

	return nil
}

// Returns a value of Attribute as string.
// There is no way to get value as string directly,
// so we parse tokens of Attribute and build string representation.
func GetAttributeValue(inFile *hclwrite.File, address string) (string, error) {
	attr, _, err := findAttribute(inFile.Body(), address)
	if err != nil {
		return "", err
	}

	// not found
	if attr == nil {
		return "", nil
	}

	// treat expr as a string without interpreting its meaning.
	out, err := getAttributeValueAsString(attr)

	if err != nil {
		return "", err
	}

	return (out + "\n"), nil
}

// Filter reads HCL and remove a matched attribute at a given address.
func RemoveAttribute(inFile *hclwrite.File, address string) error {
	attr, body, err := findAttribute(inFile.Body(), address)
	if err != nil {
		return err
	}

	if attr != nil {
		a := strings.Split(address, ".")
		attrName := a[len(a)-1]
		body.RemoveAttribute(attrName)
	}

	return nil
}

// Filter reads HCL and updates a value of matched an attribute at a given address.
func SetAttributeValue(inFile *hclwrite.File, address string, value string) error {
	attr, body, err := findAttribute(inFile.Body(), address)
	if err != nil {
		return err
	}

	if attr != nil {
		a := strings.Split(address, ".")
		attrName := a[len(a)-1]

		// To delegate expression parsing to the hclwrite parser,
		// We build a new expression and set back to the attribute by tokens.
		expr, err := buildExpression(attrName, value)
		if err != nil {
			return err
		}
		body.SetAttributeRaw(attrName, expr.BuildTokens(nil))
	}

	return nil
}
