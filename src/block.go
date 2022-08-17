package main

import (
	"fmt"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
)

// Filter reads HCL and appends only matched blocks at a given address.
// The child address is relative to parent one.
// If a newline flag is true, it also appends a newline before the new block.
func AppendChildBlock(inFile *hclwrite.File, parent string, child string, addNewline bool) error {
	pTypeName, pLabels, err := parseAddress(parent)
	if err != nil {
		return fmt.Errorf("failed to parse parent address: %s", err)
	}

	cTypeName, cLabels, err := parseAddress(child)
	if err != nil {
		return fmt.Errorf("failed to parse child address: %s", err)
	}

	matched := findBlocks(inFile.Body(), pTypeName, pLabels)

	for _, b := range matched {
		if addNewline {
			b.Body().AppendNewline()
		}
		b.Body().AppendNewBlock(cTypeName, cLabels)
	}

	return nil
}

//
func RenameBlocks(inFile *hclwrite.File, from string, to string, body *hclwrite.Body) error {
	if body == nil {
		body = inFile.Body()
	}

	fromAddressList := strings.Split(from, ".")
	toAddressList := strings.Split(to, ".")

	if len(fromAddressList) != len(toAddressList) {
		return fmt.Errorf("'from' and 'to' must have the same argument size")
	}

	for _, block := range body.Blocks() {
		// Check if block has nested blocks
		if len(block.Body().Blocks()) != 0 {
			// Call function recursively
			err := RenameBlocks(nil, from, to, block.Body())
			if err != nil {
				return err
			}
		}

		for fromIndex, fromAddress := range fromAddressList {
			if fromAddress == block.Type() {
				labelNames := block.Labels()

				fromAddressListEnd := []string{}
				if fromIndex+1 != len(fromAddressList) {
					fromAddressListEnd = fromAddressList[fromIndex+1:]
				}

				// Change block name
				block.SetType(toAddressList[fromIndex])

				// Block has no labels
				if len(fromAddressListEnd) < 1 {
					continue
				}

				// Change block labels
				labels := []string{}
				for endIndex, item := range fromAddressListEnd {
					if Contains(labelNames, item) {
						labels = append(labels, toAddressList[fromIndex+endIndex+1])
					}
				}

				block.SetLabels(labels)
			}
		}
	}

	return nil
}

func RemoveParentBlock(inFile *hclwrite.File, address string) error {
	typeName, labels, err := parseAddress(address)
	if err != nil {
		return err
	}

	matched := findBlocks(inFile.Body(), typeName, labels)

	for _, b := range matched {
		inFile.Body().RemoveBlock(b)
	}

	// Formatter
	tokens := inFile.BuildTokens(nil)
	vertical := verticalFormat(tokens)

	inFile = hclwrite.NewEmptyFile()
	inFile.Body().AppendUnstructuredTokens(vertical)

	return nil
}
