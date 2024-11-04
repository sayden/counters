package counters

import "encoding/xml"

const (
	Template_VassalPiece           = `+/null/prototype;UnitStep	prototype;RU\	emb2;Activate;128;A;;128;;;128;;;;1;false;0;0;{{ .FrontFilename }},{{ .BackFilename }};,;false;{{ .PieceName }};;;true;StepValue;1;1;true;65,130;;;;1.0;;true\\	piece;;;{{ .FrontFilename }};{{ .Id }}/	\	1\\	null;0;0;398;0`
	Template_NewVassalPiece        = `+/null/prototype;Prototype	emb2;Next;128;A;;128;;;128;;;;1;false;0;0;{{.BackFilename }};;true;{{.FrontFilename}};;;false;StepValue;1;1;false;65,130;;;;1.0;;true\	piece;;;{{.FrontFilename}};{{.PieceName}}/	-1\	null;117;107;89;1;ppScale;1.0`
	Template_Reference_VassalPiece = `+/null/prototype;UnitStep	prototype;RU\	emb2;Activate;128;A;;128;;;128;;;;1;false;0;0;{{ .FilenameFront }},{{ .FilenameBack }};,;false;{{ .CounterName }};;;true;StepValue;1;1;true;65,130;;;;1.0;;true\\	piece;;;{{ .FilenameFront }};{{ .Id }}/	\	1\\	null;0;0;398;0`
	Template_OldPiece              = `+/null/prototype;BasicPrototype	piece;;;{{ .Filename }};{{ .PieceName}}/	null;0;0;{{ .Id }};0`
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
	BasicCommandEncoder        capture `xml:"VASSAL.build.module.BasicCommandEncoder"`
	Documentation              capture `xml:"VASSAL.build.module.Documentation"`
	Chatter                    capture `xml:"VASSAL.build.module.Chatter"`
	KeyNamer                   capture `xml:"VASSAL.build.module.KeyNamer"`
	PieceWindow                PieceWindow
	DiceButton                 []DiceButton `xml:"VASSAL.build.module.DiceButton"`
	PlayerRoster               capture      `xml:"VASSAL.build.module.PlayerRoster"`
	GlobalOptions              capture      `xml:"VASSAL.build.module.GlobalOptions"`
	GamePieceDefinitions       capture      `xml:"VASSAL.build.module.gamepieceimage.GamePieceImageDefinitions"`
	GlobalProperties           capture      `xml:"VASSAL.build.module.properties.GlobalProperties"`
	GlobalTranslatableMessages capture      `xml:"VASSAL.build.module.properties.GlobalTranslatableMessages"`
	PrototypesContainer        capture      `xml:"VASSAL.build.module.PrototypesContainer"`
	Language                   capture      `xml:"VASSAL.i18n.Language"`
	Map                        Map          `xml:"VASSAL.build.module.Map"`
}

type Map struct {
	XMLName xml.Name `xml:"VASSAL.build.module.Map"`

	BoardPicker               BoardPicker `xml:"VASSAL.build.module.map.BoardPicker"`
	Text                      string      `xml:",chardata"`
	AllowMultiple             string      `xml:"allowMultiple,attr"`
	Backgroundcolor           string      `xml:"backgroundcolor,attr"`
	ButtonName                string      `xml:"buttonName,attr"`
	ChangeFormat              string      `xml:"changeFormat,attr"`
	Color                     string      `xml:"color,attr"`
	CreateFormat              string      `xml:"createFormat,attr"`
	EdgeHeight                string      `xml:"edgeHeight,attr"`
	EdgeWidth                 string      `xml:"edgeWidth,attr"`
	HideKey                   string      `xml:"hideKey,attr"`
	Hotkey                    string      `xml:"hotkey,attr"`
	Icon                      string      `xml:"icon,attr"`
	Launch                    string      `xml:"launch,attr"`
	MapName                   string      `xml:"mapName,attr"`
	MarkMoved                 string      `xml:"markMoved,attr"`
	MarkUnmovedHotkey         string      `xml:"markUnmovedHotkey,attr"`
	MarkUnmovedIcon           string      `xml:"markUnmovedIcon,attr"`
	MarkUnmovedReport         string      `xml:"markUnmovedReport,attr"`
	MarkUnmovedText           string      `xml:"markUnmovedText,attr"`
	MarkUnmovedTooltip        string      `xml:"markUnmovedTooltip,attr"`
	MoveKey                   string      `xml:"moveKey,attr"`
	MoveToFormat              string      `xml:"moveToFormat,attr"`
	MoveWithinFormat          string      `xml:"moveWithinFormat,attr"`
	OnlyReportChangedLocation string      `xml:"onlyReportChangedLocation,attr"`
	ShowKey                   string      `xml:"showKey,attr"`
	Thickness                 string      `xml:"thickness,attr"`
	StackMetrics              struct {
		Text     string `xml:",chardata"`
		Bottom   string `xml:"bottom,attr"`
		Disabled string `xml:"disabled,attr"`
		Down     string `xml:"down,attr"`
		ExSepX   string `xml:"exSepX,attr"`
		ExSepY   string `xml:"exSepY,attr"`
		Top      string `xml:"top,attr"`
		UnexSepX string `xml:"unexSepX,attr"`
		UnexSepY string `xml:"unexSepY,attr"`
		Up       string `xml:"up,attr"`
	} `xml:"VASSAL.build.module.map.StackMetrics"`
	ForwardToKeyBuffer string `xml:"VASSAL.build.module.map.ForwardToKeyBuffer"`
	Scroller           string `xml:"VASSAL.build.module.map.Scroller"`
	ForwardToChatter   string `xml:"VASSAL.build.module.map.ForwardToChatter"`
	MenuDisplayer      string `xml:"VASSAL.build.module.map.MenuDisplayer"`
	MapCenterer        string `xml:"VASSAL.build.module.map.MapCenterer"`
	StackExpander      string `xml:"VASSAL.build.module.map.StackExpander"`
	PieceMover         string `xml:"VASSAL.build.module.map.PieceMover"`
	KeyBufferer        string `xml:"VASSAL.build.module.map.KeyBufferer"`
	ImageSaver         struct {
		Text             string `xml:",chardata"`
		ButtonText       string `xml:"buttonText,attr"`
		CanDisable       string `xml:"canDisable,attr"`
		DisabledIcon     string `xml:"disabledIcon,attr"`
		HideWhenDisabled string `xml:"hideWhenDisabled,attr"`
		Hotkey           string `xml:"hotkey,attr"`
		Icon             string `xml:"icon,attr"`
		PropertyGate     string `xml:"propertyGate,attr"`
		Tooltip          string `xml:"tooltip,attr"`
	} `xml:"VASSAL.build.module.map.ImageSaver"`
	CounterDetailViewer struct {
		Text                   string `xml:",chardata"`
		BgColor                string `xml:"bgColor,attr"`
		BorderColor            string `xml:"borderColor,attr"`
		BorderInnerThickness   string `xml:"borderInnerThickness,attr"`
		BorderThickness        string `xml:"borderThickness,attr"`
		BorderWidth            string `xml:"borderWidth,attr"`
		CenterAll              string `xml:"centerAll,attr"`
		CenterPiecesVertically string `xml:"centerPiecesVertically,attr"`
		CenterText             string `xml:"centerText,attr"`
		CombineCounterSummary  string `xml:"combineCounterSummary,attr"`
		CounterReportFormat    string `xml:"counterReportFormat,attr"`
		Delay                  string `xml:"delay,attr"`
		Description            string `xml:"description,attr"`
		Display                string `xml:"display,attr"`
		EmptyHexReportForma    string `xml:"emptyHexReportForma,attr"`
		EnableHTML             string `xml:"enableHTML,attr"`
		ExtraTextPadding       string `xml:"extraTextPadding,attr"`
		FgColor                string `xml:"fgColor,attr"`
		FontSize               string `xml:"fontSize,attr"`
		GraphicsZoom           string `xml:"graphicsZoom,attr"`
		Hotkey                 string `xml:"hotkey,attr"`
		LayerList              string `xml:"layerList,attr"`
		MinDisplayPieces       string `xml:"minDisplayPieces,attr"`
		OnlyShowFirstSummary   string `xml:"onlyShowFirstSummary,attr"`
		PropertyFilter         string `xml:"propertyFilter,attr"`
		ShowDeck               string `xml:"showDeck,attr"`
		ShowDeckDepth          string `xml:"showDeckDepth,attr"`
		ShowDeckMasked         string `xml:"showDeckMasked,attr"`
		ShowMoveSelectde       string `xml:"showMoveSelectde,attr"`
		ShowNoStack            string `xml:"showNoStack,attr"`
		ShowNonMovable         string `xml:"showNonMovable,attr"`
		ShowOnlyTopOfStack     string `xml:"showOnlyTopOfStack,attr"`
		ShowOverlap            string `xml:"showOverlap,attr"`
		ShowTerrainBeneath     string `xml:"showTerrainBeneath,attr"`
		ShowTerrainHeight      string `xml:"showTerrainHeight,attr"`
		ShowTerrainSnappy      string `xml:"showTerrainSnappy,attr"`
		ShowTerrainText        string `xml:"showTerrainText,attr"`
		ShowTerrainWidth       string `xml:"showTerrainWidth,attr"`
		ShowTerrainZoom        string `xml:"showTerrainZoom,attr"`
		Showgraph              string `xml:"showgraph,attr"`
		Showgraphsingle        string `xml:"showgraphsingle,attr"`
		Showtext               string `xml:"showtext,attr"`
		Showtextsingle         string `xml:"showtextsingle,attr"`
		StopAfterShowing       string `xml:"stopAfterShowing,attr"`
		StretchWidthPieces     string `xml:"stretchWidthPieces,attr"`
		StretchWidthSummary    string `xml:"stretchWidthSummary,attr"`
		SummaryReportFormat    string `xml:"summaryReportFormat,attr"`
		UnrotatePieces         string `xml:"unrotatePieces,attr"`
		Version                string `xml:"version,attr"`
		VerticalBottomText     string `xml:"verticalBottomText,attr"`
		VerticalOffset         string `xml:"verticalOffset,attr"`
		VerticalTopText        string `xml:"verticalTopText,attr"`
		Zoomlevel              string `xml:"zoomlevel,attr"`
	} `xml:"VASSAL.build.module.map.CounterDetailViewer"`
	Flare struct {
		Text              string `xml:",chardata"`
		CircleColor       string `xml:"circleColor,attr"`
		CircleScale       string `xml:"circleScale,attr"`
		CircleSize        string `xml:"circleSize,attr"`
		FlareKey          string `xml:"flareKey,attr"`
		FlareName         string `xml:"flareName,attr"`
		FlarePulses       string `xml:"flarePulses,attr"`
		FlarePulsesPerSec string `xml:"flarePulsesPerSec,attr"`
		ReportFormat      string `xml:"reportFormat,attr"`
	} `xml:"VASSAL.build.module.map.Flare"`
	Zoomer struct {
		Text           string `xml:",chardata"`
		InButtonText   string `xml:"inButtonText,attr"`
		InIconName     string `xml:"inIconName,attr"`
		InTooltip      string `xml:"inTooltip,attr"`
		OutButtonText  string `xml:"outButtonText,attr"`
		OutIconName    string `xml:"outIconName,attr"`
		OutTooltip     string `xml:"outTooltip,attr"`
		PickButtonText string `xml:"pickButtonText,attr"`
		PickIconName   string `xml:"pickIconName,attr"`
		PickTooltip    string `xml:"pickTooltip,attr"`
		ZoomInKey      string `xml:"zoomInKey,attr"`
		ZoomLevels     string `xml:"zoomLevels,attr"`
		ZoomOutKey     string `xml:"zoomOutKey,attr"`
		ZoomPickKey    string `xml:"zoomPickKey,attr"`
		ZoomStart      string `xml:"zoomStart,attr"`
	} `xml:"VASSAL.build.module.map.Zoomer"`
	VASSALBuildModulePropertiesGlobalProperties string `xml:"VASSAL.build.module.properties.GlobalProperties"`
	VASSALBuildModuleMapSelectionHighlighters   string `xml:"VASSAL.build.module.map.SelectionHighlighters"`
	HighlightLastMoved                          struct {
		Text      string `xml:",chardata"`
		Color     string `xml:"color,attr"`
		Enabled   string `xml:"enabled,attr"`
		Thickness string `xml:"thickness,attr"`
	} `xml:"VASSAL.build.module.map.HighlightLastMoved"`
	HidePiecesButton struct {
		Text        string `xml:",chardata"`
		ButtonText  string `xml:"buttonText,attr"`
		HiddenIcon  string `xml:"hiddenIcon,attr"`
		Hotkey      string `xml:"hotkey,attr"`
		ShowingIcon string `xml:"showingIcon,attr"`
		Tooltip     string `xml:"tooltip,attr"`
	} `xml:"VASSAL.build.module.map.HidePiecesButton"`
}

type capture struct {
	Raw string `xml:",innerxml"`
}

type DiceButton struct {
	XMLName xml.Name `xml:"VASSAL.build.module.DiceButton"`

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

	EntryName string `xml:"entryName,attr"`
	// PanelWidget PanelWidget  `xml:"VASSAL.build.widget.PanelWidget"`
	ListWidget []ListWidget `xml:"VASSAL.build.widget.ListWidget"`
}

type PanelWidget struct {
	XMLName xml.Name `xml:"VASSAL.build.widget.PanelWidget"`

	Text       string       `xml:",chardata"`
	EntryName  string       `xml:"entryName,attr"`
	Fixed      string       `xml:"fixed,attr"`
	NColumns   string       `xml:"nColumns,attr"`
	Scale      string       `xml:"scale,attr"`
	Vert       string       `xml:"vert,attr"`
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

type VassalFileModuleData struct {
	XMLName xml.Name `xml:"data"`

	Text          string `xml:",chardata"`
	AttrVersion   string `xml:"version,attr"`
	Version       string `xml:"version"`
	Extra1        string `xml:"extra1"`
	Extra2        string `xml:"extra2"`
	VassalVersion string `xml:"VassalVersion"`
	DateSaved     string `xml:"dateSaved"`
	Description   string `xml:"description"`
	Name          string `xml:"name"`
}

type BoardPicker struct {
	XMLName       xml.Name `xml:"VASSAL.build.module.map.BoardPicker"`
	Text          string   `xml:",chardata"`
	AddColumnText string   `xml:"addColumnText,attr"`
	AddRowText    string   `xml:"addRowText,attr"`
	BoardPrompt   string   `xml:"boardPrompt,attr"`
	SlotHeight    string   `xml:"slotHeight,attr"`
	SlotScale     string   `xml:"slotScale,attr"`
	SlotWidth     string   `xml:"slotWidth,attr"`
	Title         string   `xml:"title,attr"`
	Board         Board    `xml:"VASSAL.build.module.map.boardPicker.Board"`
}

type Board struct {
	Text       string  `xml:",chardata"`
	Image      string  `xml:"image,attr"`
	Name       string  `xml:"name,attr"`
	Reversible string  `xml:"reversible,attr"`
	HexGrid    HexGrid `xml:"VASSAL.build.module.map.boardPicker.board.HexGrid"`
}

type HexGrid struct {
	XMLName xml.Name `xml:"VASSAL.build.module.map.boardPicker.board.HexGrid" json:"-"`

	Text         string `xml:",chardata" json:"text,omitempty"`
	Color        string `xml:"color,attr" json:"color,omitempty"`
	CornersLegal string `xml:"cornersLegal,attr" json:"cornersLegal,omitempty"`
	DotsVisible  string `xml:"dotsVisible,attr" json:"dotsVisible,omitempty"`
	Dx           string `xml:"dx,attr" json:"dx,omitempty"`
	Dy           string `xml:"dy,attr" json:"dy,omitempty"`
	EdgesLegal   string `xml:"edgesLegal,attr" json:"edgesLegal,omitempty"`
	Sideways     string `xml:"sideways,attr" json:"sideways,omitempty"`
	SnapTo       string `xml:"snapTo,attr" json:"snapTo,omitempty"`
	Visible      string `xml:"visible,attr" json:"visible,omitempty"`
	X0           string `xml:"x0,attr" json:"x0,omitempty"`
	Y0           string `xml:"y0,attr" json:"y0,omitempty"`
}

// Templates

type CSVTemplateData struct {
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
