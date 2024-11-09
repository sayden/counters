package counters

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCounterPrototypeToCounter(t *testing.T) {
	proto := CounterPrototype{
		Counter: Counter{
			Texts: []Text{{String: "text"}},
		},
		TextPrototypes: []TextPrototype{
			{StringList: []string{"text1", "text2"}},
		},
		ImagePrototypes: []ImagePrototype{
			{PathList: []string{"../assets/binoculars.png", "../assets/stripe.png"}},
		},
	}

	filenamesInUse := &sync.Map{}
	counter, err := proto.ToCounters(filenamesInUse, "side", "testing", 0)
	if assert.NoError(t, err) {
		assert.Equal(t, 2, len(counter))
		assert.Equal(t, "text", counter[0].Texts[0].String)
		assert.Equal(t, "text1", counter[0].Texts[1].String)
		assert.Equal(t, "../assets/binoculars.png", counter[0].Images[0].Path)

		assert.Equal(t, "text", counter[1].Texts[0].String)
		assert.Equal(t, "text2", counter[1].Texts[1].String)
		assert.Equal(t, "../assets/stripe.png", counter[1].Images[0].Path)
	}
}

func TestJSONPrototypes(t *testing.T) {
	// read the json file with testing data
	filepath := "./testdata/prototype.json"
	byt, err := os.ReadFile(filepath)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	filenamesInUse := &sync.Map{}
	newTempl, err := ParseCounterTemplate(byt, filenamesInUse)
	if !assert.NoError(t, err, "error parsing template") {
		t.FailNow()
	}

	// Extract the counters from the prototypes into a new template
	newTempl, err = newTempl.ParsePrototype()
	assert.NoError(t, err)

	// check the new template
	assert.Equal(t, 8, len(newTempl.Counters))
	assert.Equal(t, 2, len(newTempl.Counters[0].Texts))

	// check the marshalling of the template to an expected byte slice
	actualBytes, err := json.MarshalIndent(newTempl, "", "  ")
	assert.NoError(t, err)
	if !assert.Equal(t, 18356, len(actualBytes)) {
		t.FailNow()
	}

	expectedFile, err := os.Open("./testdata/parse_template_01.json")
	assert.NoError(t, err)
	defer expectedFile.Close()

	expectedBytes, err := io.ReadAll(expectedFile)
	assert.NoError(t, err)

	// // ensure we are using the expected file and that it has not been altered by mistake
	if !assert.Equal(t, 18356, len(expectedBytes), "expected file has been altered, aborting test") {
		t.FailNow()
	}

	// compare the bytes of the expected file data and the actual data
	assert.Equal(t, len(expectedBytes), len(actualBytes), "Expected and actual byte slices have different lengths")
	assert.Equal(t, string(expectedBytes), string(actualBytes), "Expected and actual byte data are different")
}
