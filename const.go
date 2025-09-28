package counters

const (
	DEFAULT_FONT_HEIGHT = 20.0

	DEFAULT_MARGINS_DISTANCE = 10

	DEFAULT_COUNTER_WIDTH  = 200.0
	DEFAULT_COUNTER_HEIGHT = 200.0

	DEFAULT_MODE = "tiles"

	TEMPLATE_MODE_TILES    = "tiles"
	TEMPLATE_MODE_TEMPLATE = "template"

	DEFAULT_FONT_COLOR       string = "black"
	DEFAULT_BACKGROUND_COLOR string = "white"

	CARD_AREA_FRAME_WIDTH = 4

	DEFAULT_CARD_MARGINS_DISTANCE  = 10
	DEFAULT_IMAGE_MARGINS_DISTANCE = 10
	DEFAULT_TEXT_BOX_MARGINS       = 10

	DEFAULT_IMAGE_WIDTH  = 800.0
	DEFAULT_IMAGE_HEIGHT = 1200.0

	DEFAULT_BORDER_WIDTH = 20
	DEFAULT_BORDER_COLOR = "white"

	IMAGE_SCALING_FIT_NONE   = "none"
	IMAGE_SCALING_FIT_WIDTH  = "fitWidth"
	IMAGE_SCALING_FIT_HEIGHT = "fitHeight"
	IMAGE_SCALING_FIT_WRAP   = "wrap"

	ALIGMENT_LEFT   = "left"
	ALIGMENT_RIGHT  = "right"
	ALIGMENT_CENTER = "center"

	BASE_FOLDER = "TemplateModule"

	STRIPE              = "assets/stripe.png"
	VassalInputXmlFile  = "template.xml"
	VassalOutputXmlFile = BASE_FOLDER + "/buildFile.xml"
)

type FileContent int

const (
	FileContent_CounterTemplate = iota
	FileContent_CardTemplate
	FileContent_Events
	FileContent_Quotes
	FileContent_CSV
	FileContent_JSON
)
