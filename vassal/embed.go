package vassal

import (
	_ "embed"
	"encoding/xml"
	"fmt"

	"github.com/sayden/counters"
)

//go:embed moduledata
var moduleDataBytes []byte
var moduleData counters.VassalFileModuleData

//go:embed buildFile.xml
var buildFileBytes []byte
var buildFile counters.VassalGameModule

type moduleTemplateData struct {
	Name string
}

func init() {
	err := xml.Unmarshal(moduleDataBytes, &moduleData)
	if err != nil {
		panic(fmt.Errorf("error trying to decode content: %w", err))
	}

	err = xml.Unmarshal(buildFileBytes, &buildFile)
	if err != nil {
		panic(fmt.Errorf("error trying to decode content: %w", err))
	}
}

func GetBuildFile() *counters.VassalGameModule {
	return &buildFile
}

func GetModuleData() *counters.VassalFileModuleData {
	return &moduleData
}
