package output

import (
	"encoding/xml"
	"fmt"
	"os"
	"path"
	"sort"

	"github.com/pkg/errors"
	"github.com/sayden/counters"
	"github.com/sayden/counters/fsops"
	"github.com/sayden/counters/input"
	"github.com/sayden/counters/vassal"
)

func VassalModule(outputPath string, templatesFiles []string) error {
	// TODO: Side effects
	os.MkdirAll(outputPath, 0755)
	tempDir, err := os.MkdirTemp(outputPath, "vassal")
	if err != nil {
		return errors.Wrap(err, "error creating temporary directory")
	}
	defer os.RemoveAll(tempDir)

	_ = os.MkdirAll(path.Join(tempDir, "images"), 0755)

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

		newTemplate.OutputFolder = path.Join(tempDir, "images")
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

		// TODO: Side effects
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
		sort.Sort(list.PieceSlot)

		listOfListWidgets = append(listOfListWidgets, list)
	}

	// TODO: Side effects
	// Copy map file to the images folder
	if err = fsops.CopyFile(mapFilename, path.Join(tempDir, "images", path.Base(mapFilename))); err != nil {
		return errors.Wrap(err, "error copying map file")
	}

	return writeXMLFiles(tempDir, moduleName, mapFilename, outputPath, listOfListWidgets, &hexGrid)
}

func writeXMLFiles(tempDir, moduleName, mapFilename, outputPath string,
	listOfWidgets []counters.ListWidget, hexGrid *counters.HexGrid) error {
	// buildFile.xml
	buildFile := vassal.GetBuildFile()

	buildFile.Name = moduleName
	buildFile.Version = "0.2"
	buildFile.Map.BoardPicker.Board.Image = path.Base(mapFilename)
	buildFile.Map.BoardPicker.Board.Name = moduleName
	if hexGrid != nil {
		buildFile.Map.BoardPicker.Board.HexGrid = *hexGrid
	}
	buildFile.PieceWindow.TabWidget.ListWidget = listOfWidgets

	// TODO: Side effects
	f, err := os.Create(path.Join(tempDir, "buildFile.xml"))
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
	// TODO: Side effects
	f2, err := os.Create(path.Join(tempDir, "moduledata"))
	if err != nil {
		return fmt.Errorf("error creating buildFile.xml: %w", err)
	}
	defer f2.Close()
	xml.NewEncoder(f2).Encode(moduleData)

	// Compress
	return WriteZipFileWithFolderContent(path.Join(outputPath, moduleName+".vmod"), tempDir)
}
