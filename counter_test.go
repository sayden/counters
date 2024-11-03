package counters

import (
	"bytes"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCounterFilename(t *testing.T) {
	tests := []struct {
		name           string
		counter        Counter
		position       int
		suffix         string
		filenumber     int
		filenamesInUse map[string]bool
		expected       string
	}{
		{
			name: "Basic case",
			counter: Counter{
				Texts: []Text{
					{Settings: Settings{Position: 0}, String: "Test"}},
			},
			position:       0,
			suffix:         "suffix",
			filenumber:     -1,
			filenamesInUse: map[string]bool{},
			expected:       "Test suffix.png",
		},
		{
			name: "With Extra Title",
			counter: Counter{
				Texts: []Text{
					{Settings: Settings{Position: 0}, String: "Test"}},
				Extra: &Extra{Title: stringP("ExtraTitle")},
			},
			position:       0,
			suffix:         "suffix",
			filenumber:     -1,
			filenamesInUse: map[string]bool{},
			expected:       "Test ExtraTitle suffix.png",
		},
		{
			name: "With Title Position",
			counter: Counter{
				Texts: []Text{
					{Settings: Settings{Position: 0}, String: "Test"},
					{Settings: Settings{Position: 1}, String: "TitlePosition"},
				},
				Extra: &Extra{TitlePosition: intP(1)},
			},
			position:       0,
			suffix:         "suffix",
			filenumber:     -1,
			filenamesInUse: map[string]bool{},
			expected:       "TitlePosition suffix.png",
		},
		{
			name: "Filename in use",
			counter: Counter{
				Texts: []Text{
					{Settings: Settings{Position: 0}, String: "Test"},
				},
			},
			position:       0,
			suffix:         "suffix",
			filenumber:     1,
			filenamesInUse: map[string]bool{"Test suffix": true},
			expected:       "Test suffix 001.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.counter.GetCounterFilename(tt.position, tt.suffix, tt.filenumber, tt.filenamesInUse)
			if got != tt.expected {
				t.Errorf("GetCounterFilename() = '%v', want '%v'", got, tt.expected)
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
