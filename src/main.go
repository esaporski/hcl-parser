package main

import (
	"log"
	"os"

	"github.com/hashicorp/hcl/v2"
	"github.com/hashicorp/hcl/v2/hclwrite"
)

func main() {
	filepathList, err := FindFiles("infrastructure", ".tf")
	if err != nil {
		log.Fatal(err)
	}

	// Iterate through files
	for _, filepath := range filepathList {
		// Read existing file content
		fileContent, err := os.ReadFile(filepath)
		if err != nil {
			log.Fatalf("Could not read file content %s: %v", filepath, err)
		}

		// Parse []byte into *hclwrite.File
		hclFile, diags := hclwrite.ParseConfig(fileContent, filepath, hcl.InitialPos)
		if diags.HasErrors() {
			log.Fatalf("Parsing file %s: %s", filepath, diags.Error())
		}

		// // // //
		// Attribute
		// AppendAttribute
		err = AppendAttribute(hclFile, "resource.aws_instance.nginx.my_new_attribute", "var.my_new_variable", false)
		if err != nil {
			log.Fatalf("Could not append attribute (%s): %v", filepath, err)
		}

		// GetAttributeValue
		attrValue, err := GetAttributeValue(hclFile, "resource.aws_instance.nginx.tags")
		if err != nil {
			log.Fatalf("Could not get attribute value (%s): %v", filepath, err)
		}

		if len(attrValue) > 0 {
			log.Print(attrValue)
		}

		// RemoveAttribute
		err = RemoveAttribute(hclFile, "resource.aws_instance.nginx.my_new_attribute")
		if err != nil {
			log.Fatalf("Could not remove attribute (%s): %v", filepath, err)
		}

		// SetAttributeValue
		err = SetAttributeValue(hclFile, "resource.aws_instance.nginx.instance_type", "var.instance_type")
		if err != nil {
			log.Fatalf("Could not change attribute value (%s): %v", filepath, err)
		}

		// Block
		// AppendChildBlock
		err = AppendChildBlock(hclFile, "resource.aws_instance.nginx", "block", false)
		if err != nil {
			log.Fatalf("Could not append child block (%s): %v", filepath, err)
		}

		// RenameBlocks
		// TODO: should trigger an attribute search for replacing attribute values with new block name
		err = RenameBlocks(hclFile, "resource.aws_instance.nginx.block", "resource.aws_instance.nginx.block_renamed", nil)
		if err != nil {
			log.Fatalf("Could not rename block (%s): %v", filepath, err)
		}

		// Remove Parent Block
		err = RemoveParentBlock(hclFile, "data.aws_ami.amazon_linux_ami")
		if err != nil {
			log.Fatalf("Could not remove parent block (%s): %v", filepath, err)
		}

		// // // //

		// Update file in-place
		err = os.WriteFile(filepath, hclFile.Bytes(), 0644)
		if err != nil {
			log.Fatalf("Failed updating file: %v", err)
		}
	}
}
