package fsops

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/thehivecorporation/log"
)

func ReadMarkupFile(markupFilepath string, destination interface{}) error {
	extension := filepath.Ext(markupFilepath)

	data, err := os.ReadFile(markupFilepath)
	if err != nil {
		return errors.Wrap(err, "could not read file content")
	}

	switch extension {
	case ".xml":
		if err = xml.Unmarshal(data, &destination); err != nil {
			return errors.Wrapf(err, "the file in '%s' has syntax errors", markupFilepath)
		}
		return nil
	case ".json":
		if err = json.Unmarshal(data, &destination); err != nil {
			return errors.Wrapf(err, "the file in '%s' has syntax errors", markupFilepath)
		}
		return nil
	}

	return fmt.Errorf("file extension '%s' not recognized. Use .json or .xml files only", extension)
}

func FilenameExistsInFolder(filename, folder string) bool {
	fs, err := os.ReadDir(folder)
	if err != nil {
		log.WithError(err).Fatal("could not read images folder")
	}

	for _, file := range fs {
		_, existingFilename := filepath.Split(file.Name())
		_, gameImageName := filepath.Split(filename)
		if strings.Contains(gameImageName, existingFilename) {
			return true
		}
	}

	return false
}

// GetFilenamesForPath returns every path+filename found in `path`
func GetFilenamesForPath(path string) ([]string, error) {
	rootPath, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	images := make([]string, 0)

	fullPath := rootPath + path
	err = filepath.Walk(fullPath, func(imagePath string, info os.FileInfo, err error) error {
		if imagePath == fullPath {
			return nil
		}

		if err != nil {
			return err
		}

		images = append(images, imagePath)

		return nil
	})

	if err != nil {
		return nil, errors.Wrap(err, "could not finish reading files from folder")
	}

	return images, nil
}
