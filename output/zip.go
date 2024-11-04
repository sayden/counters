package output

import (
	"archive/zip"
	"fmt"
	"os"

	"github.com/charmbracelet/log"
)

func WriteZipFileWithFolderContent(destinationZipfilePath, inputFolder string) error {
	// Create the zip/vmod file
	outFile, err := os.Create(destinationZipfilePath)
	if err != nil {
		log.Fatal("could not create destination vassal file", "error", err)
	}
	defer outFile.Close()

	z := zip.NewWriter(outFile)
	defer func() {
		if err = z.Close(); err != nil {
			log.Error("zip file had a problem when closing", "error", err)
		}
	}()

	log.Info("Using", "basepath", inputFolder, "dest_file", destinationZipfilePath)

	return addFiles(z, inputFolder, "")
}

// https://stackoverflow.com/questions/37869793/how-do-i-zip-a-directory-containing-sub-directories-or-files-in-golang
func addFiles(w *zip.Writer, basePath, baseInZip string) error {
	// Open the Directory
	files, err := os.ReadDir(basePath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			dat, err := os.ReadFile(basePath + "/" + file.Name())
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(baseInZip + file.Name())
			if err != nil {
				return err
			}
			_, err = f.Write(dat)
			if err != nil {
				return err
			}
		} else if file.IsDir() {
			// Recurse
			newBase := basePath + "/" + file.Name() + "/"

			if err = addFiles(w, newBase, baseInZip+file.Name()+"/"); err != nil {
				return err
			}
		}
	}

	return nil
}
