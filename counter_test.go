package counters

import (
	"bytes"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCounterFilename(t *testing.T) {
	filenamesInUse := new(sync.Map)

	tests := []struct {
		name           string
		counter        Counter
		position       int
		filenumber     int
		expected       string
		expectedPretty string
	}{
		{
			name: "Basic case",
			counter: Counter{
				Texts: []Text{
					{Settings: Settings{Position: 0}, String: "Test"}},
			},
			position:       0,
			filenumber:     -1,
			expected:       "side_Test.png",
			expectedPretty: "Test",
		},
		{
			name: "With Extra Title",
			counter: Counter{
				Texts: []Text{
					{Settings: Settings{Position: 0}, String: "Test"}},
				Metadata: &Metadata{Title: "ExtraTitle"},
			},
			position:       0,
			filenumber:     -1,
			expected:       "side_Test ExtraTitle.png",
			expectedPretty: "Test ExtraTitle",
		},
		{
			name: "With Title Position",
			counter: Counter{
				Texts: []Text{
					{Settings: Settings{Position: 0}, String: "Test"},
					{Settings: Settings{Position: 1}, String: "TitlePosition"},
				},
				Metadata: &Metadata{TitlePosition: intP(1)},
			},
			position:       0,
			filenumber:     -1,
			expected:       "side_TitlePosition.png",
			expectedPretty: "TitlePosition",
		},
		{
			name: "Filename in use",
			counter: Counter{
				Metadata: &Metadata{Title: "ExtraTitle"},
				Texts: []Text{
					{Settings: Settings{Position: 3}, String: "Test"},
				},
			},
			position:       3,
			filenumber:     1,
			expected:       "side_Test ExtraTitle_0000.png",
			expectedPretty: "Test ExtraTitle",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.counter.GenerateCounterFilename("side", tt.position, filenamesInUse)
			if tt.counter.Filename != tt.expected {
				t.Errorf("GetCounterFilename() = '%v', want '%v'", tt.counter.Filename, tt.expected)
			}
			if tt.counter.PrettyName != tt.expectedPretty {
				t.Errorf("PrettyName = '%v', want '%v'", tt.counter.PrettyName, tt.expectedPretty)
			}
		})
	}
}

func TestCounterEncode(t *testing.T) {
	counter := Counter{
		Settings: Settings{
			Width:           100,
			Height:          100,
			FontPath:        "assets/freesans.ttf",
			FontColorS:      "black",
			BackgroundColor: stringP("black"),
			StrokeWidth:     floatP(2),
			StrokeColorS:    "white",
			FontHeight:      15,
			BorderWidth:     floatP(2),
			BorderColorS:    "red",
		},
		Texts: []Text{
			{String: "Area text"},
		},
	}

	byt := make([]byte, 0, 10000)
	buf := bytes.NewBuffer(byt)

	err := counter.EncodeCounter(buf, false)
	if err != nil {
		t.Fatal(err)
	}
	byt = buf.Bytes()

	expectedByt, err := os.ReadFile("testdata/counter_01.png")
	if err != nil {
		t.FailNow()
	}

	if assert.Equal(t, len(expectedByt), len(byt)) {
		assert.Equal(t, expectedByt, byt)
	}
}
