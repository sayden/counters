package output

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/fsops"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/vassal"
)

func VassalModule(outputPath string, templatesFiles []string) error {
	os.MkdirAll(outputPath, 0755)
	destDir, err := os.MkdirTemp(outputPath, "vassal")
	if err != nil {
		return errors.Wrap(err, "error creating temporary directory")
	}
	defer os.RemoveAll(destDir)

	os.MkdirAll(path.Join(destDir, "images"), 0755)

	listOfListWidgets := make([]counters.ListWidget, 0, 3)

	moduleName := ""
	mapFilename := ""
	hexGrid := counters.HexGrid{
		Color:        "204,0,204",
		CornersLegal: "false",
		DotsVisible:  "false",
		EdgesLegal:   "false",
		SnapTo:       "true",
		Visible:      "false",
	}

	// create all the counters in their respective folders
	for _, inputPath := range templatesFiles {
		if err := counters.ValidateSchemaAtPath[counters.CounterTemplate](inputPath); err != nil {
			return errors.Wrap(err, "schema validation failed during jsonToAsset")
		}

		counterTemplate, err := input.ReadCounterTemplate(inputPath)
		if err != nil {
			return errors.Wrap(err, "error reading counter template")
		}

		newTemplate, err := counterTemplate.ParsePrototype()
		if err != nil {
			return errors.Wrap(err, "error parsing prototyped template")
		}

		newTemplate.OutputFolder = path.Join(destDir, "images")
		moduleName = newTemplate.Vassal.ModuleName
		mapFilename = newTemplate.Vassal.MapFile

		if newTemplate.Vassal.HexGrid != nil {
			hexGrid = *newTemplate.Vassal.HexGrid
			hexGrid.Color = "204,0,204"
			hexGrid.CornersLegal = "false"
			hexGrid.DotsVisible = "false"
			hexGrid.EdgesLegal = "false"
			hexGrid.SnapTo = "true"
			hexGrid.Visible = "false"
		}

		CountersToPNG(newTemplate)

		// get an array of the vassal pieces
		list := counters.ListWidget{
			EntryName: newTemplate.Vassal.SideName,
			PieceSlot: make([]counters.PieceSlot, 0, len(newTemplate.Counters)),
			Scale:     "1.0",
			Height:    "215",
			Width:     "562",
			Divider:   "194",
		}

		for i, counter := range newTemplate.Counters {
			if counter.VassalPiece != nil {
				list.PieceSlot = append(list.PieceSlot, *newTemplate.Counters[i].VassalPiece)
			}
		}

		listOfListWidgets = append(listOfListWidgets, list)
	}

	// Copy map file to the images folder
	if err = fsops.CopyFile(mapFilename, path.Join(destDir, "images", path.Base(mapFilename))); err != nil {
		return errors.Wrap(err, "error copying map file")
	}

	return writeXMLFiles(destDir, moduleName, mapFilename, listOfListWidgets, outputPath, &hexGrid)
}

func writeXMLFiles(dir, moduleName, mapFilename string, listOfWidgets []counters.ListWidget, outputPath string, hexGrid *counters.HexGrid) error {
	// buildFile.xml
	buildFile := vassal.GetBuildFile()

	buildFile.Name = moduleName
	buildFile.Map.BoardPicker.Board.Image = path.Base(mapFilename)
	buildFile.Map.BoardPicker.Board.Name = moduleName
	if hexGrid != nil {
		buildFile.Map.BoardPicker.Board.HexGrid = *hexGrid
	}
	buildFile.PieceWindow.TabWidget.ListWidget = listOfWidgets

	f, err := os.Create(path.Join(dir, "buildFile.xml"))
	if err != nil {
		return fmt.Errorf("error creating buildFile.xml: %w", err)
	}
	defer f.Close()

	err = xml.NewEncoder(f).Encode(buildFile)
	if err != nil {
		return fmt.Errorf("error encoding buildFile.xml: %w", err)
	}

	// moduledata
	moduleData := vassal.GetModuleData()
	moduleData.Name = moduleName
	f2, err := os.Create(path.Join(dir, "moduledata"))
	if err != nil {
		return fmt.Errorf("error creating buildFile.xml: %w", err)
	}
	defer f2.Close()
	xml.NewEncoder(f2).Encode(moduleData)

	// Compress
	return WriteZipFileWithFolderContent(path.Join("/tmp/test", moduleName+".vmod"), dir)
}
