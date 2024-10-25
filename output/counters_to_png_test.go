package output

import (
	"os"
	"testing"

	"github.com/sayden/counters"
	"github.com/stretchr/testify/assert"
)

func TestCountersToPNG(t *testing.T) {
	byt, err := os.ReadFile("../testdata/counter_template.json")
	if err != nil {
		t.Fatal(err)
	}
	template, err := counters.ParseCounterTemplate(byt)
	if err != nil {
		t.Fatal(err)
	}

	if err != nil {
		t.Fatal(err)
	}
	t.Skip("Skipping test TestCountersToPNG because progress bar provoques a deadlock")

	err = CountersToPNG(template)
	if err != nil {
		t.Fatal(err)
	}

	// Two files should have been created in /tmp/generated folder
	f1, err := os.ReadFile(template.OutputFolder + "/counter_1.png")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(template.OutputFolder + "/counter_1.png")

	e1, err := os.ReadFile("../testdata/001.png")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, e1, f1)

	f2, err := os.ReadFile(template.OutputFolder + "/counter_2.png")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(template.OutputFolder + "/counter_2.png")

	e2, err := os.ReadFile("../testdata/002.png")
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, e2, f2)

}
