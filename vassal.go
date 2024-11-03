package counters

import "encoding/xml"

const (
	counterTemplate = `+/null/prototype;UnitStep	prototype;RU\	emb2;Activate;128;A;;128;;;128;;;;1;false;0;0;{{ .FilenameFront }},{{ .FilenameBack }};,;false;{{ .CounterName }};;;true;StepValue;1;1;true;65,130;;;;1.0;;true\\	piece;;;{{ .FilenameFront }};{{ .Id }}/	\	1\\	null;0;0;398;0`
	oldTemplate     = `+/null/prototype;BasicPrototype	piece;;;{{ .Filename }};{{ .PieceName}}/	null;0;0;{{ .Id }};0`
)

type VassalGameModule struct {
	XMLName xml.Name `xml:"VASSAL.build.GameModule"`

	ModuleOther1               string  `xml:"ModuleOther1,attr"`
	ModuleOther2               string  `xml:"ModuleOther2,attr"`
	VassalVersion              string  `xml:"VassalVersion,attr"`
	Description                string  `xml:"description,attr"`
	Name                       string  `xml:"name,attr"`
	NextPieceSlotId            string  `xml:"nextPieceSlotId,attr"`
	Version                    string  `xml:"version,attr"`
	BasicCommandEncoder        Capture `xml:"VASSAL.build.module.BasicCommandEncoder"`
	Documentation              Capture `xml:"VASSAL.build.module.Documentation"`
	Chatter                    Capture `xml:"VASSAL.build.module.Chatter"`
	KeyNamer                   Capture `xml:"VASSAL.build.module.KeyNamer"`
	PieceWindow                PieceWindow
	DiceButton                 []DiceButton `xml:"VASSAL.build.module.DiceButton"`
	PlayerRoster               Capture      `xml:"VASSAL.build.module.PlayerRoster"`
	GlobalOptions              Capture      `xml:"VASSAL.build.module.GlobalOptions"`
	GamePieceDefinitions       Capture      `xml:"VASSAL.build.module.gamepieceimage.GamePieceImageDefinitions"`
	GlobalProperties           Capture      `xml:"VASSAL.build.module.properties.GlobalProperties"`
	GlobalTranslatableMessages Capture      `xml:"VASSAL.build.module.properties.GlobalTranslatableMessages"`
	PrototypesContainer        Capture      `xml:"VASSAL.build.module.PrototypesContainer"`
	Language                   Capture      `xml:"VASSAL.i18n.Language"`
	Map                        Capture      `xml:"VASSAL.build.module.Map"`
}

type Capture struct {
	Raw string `xml:",innerxml"`
}

type DiceButton struct {
	Raw          string `xml:",innerxml"`
	AddToTotal   int    `xml:"addToTotal,attr"`
	CanDisable   bool   `xml:"canDisable,attr"`
	DisabledIcon string `xml:"disabledIcon,attr"`
	Hotkey       string `xml:"hotkey,attr"`
	Icon         string `xml:"icon,attr"`
	KeepCount    string `xml:"keepCount,attr"`
	KeepDice     string `xml:"keepDice,attr"`
	KeepOption   string `xml:"keepOption,attr"`
	LockAdd      string `xml:"lockAdd,attr"`
	LockDice     string `xml:"lockDice,attr"`
	LockPlus     string `xml:"lockPlus,attr"`
	LockSides    string `xml:"lockSides,attr"`
	NDice        string `xml:"nDice,attr"`
	NSides       string `xml:"nSides,attr"`
	Name         string `xml:"name,attr"`
	Plus         string `xml:"plus,attr"`
	Prompt       string `xml:"prompt,attr"`
	PropertyGate string `xml:"propertyGate,attr"`
	ReportFormat string `xml:"reportFormat,attr"`
	ReportTotal  string `xml:"reportTotal,attr"`
	SortDice     string `xml:"sortDice,attr"`
	Text         string `xml:"text,attr"`
	Tooltip      string `xml:"tooltip,attr"`
}

type PieceWindow struct {
	XMLName xml.Name `xml:"VASSAL.build.module.PieceWindow"`

	DefaultWidth string `xml:"defaultWidth,attr"`
	Hidden       string `xml:"hidden,attr"`
	Hotkey       string `xml:"hotkey,attr"`
	Icon         string `xml:"icon,attr"`
	Scale        string `xml:"scale,attr"`
	Text         string `xml:"text,attr"`
	ToolTip      string `xml:"tooltip,attr"`
	TabWidget    TabWidget
}

type TabWidget struct {
	XMLName xml.Name `xml:"VASSAL.build.widget.TabWidget"`

	EntryName  string       `xml:"entryName,attr"`
	ListWidget []ListWidget `xml:"VASSAL.build.widget.ListWidget"`
}

type ListWidget struct {
	XMLName xml.Name `xml:"VASSAL.build.widget.ListWidget"`

	Divider   string `xml:"divider,attr"`
	EntryName string `xml:"entryName,attr"`
	Height    string `xml:"height,attr"`
	Scale     string `xml:"scale,attr"`
	Width     string `xml:"width,attr"`
	PieceSlot []PieceSlot
}

type PieceSlot struct {
	XMLName xml.Name `xml:"VASSAL.build.widget.PieceSlot"`

	EntryName string `xml:"entryName,attr"`
	Gpid      string `xml:"gpid,attr"`
	Height    int    `xml:"height,attr"`
	Width     int    `xml:"width,attr"`
	Data      string `xml:",chardata"`
}

type TemplateData struct {
	Filename  string
	PieceName string
	Id        string
}

type PieceTemplateData struct {
	BackFilename  string
	FrontFilename string
	PieceName     string
	Id            string
	FlipName      string
}

type VassalCounterTemplateSettings struct {
	SideName string `json:"side_name,omitempty"`
}
