package main

import (
	"io/fs"
	"path/filepath"
)

// Scan directory looking for files with specified extension
func FindFiles(directory string, extension string) ([]string, error) {
	var filepathList []string

	err := filepath.WalkDir(directory, func(path string, dirEntry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(dirEntry.Name()) == extension {
			filepathList = append(filepathList, path)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return filepathList, nil
}
