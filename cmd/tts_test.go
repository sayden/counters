package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTTS(t *testing.T) {
	tts := TTS{}
	err := tts.Run(nil)
	require.NoError(t, err)
}
