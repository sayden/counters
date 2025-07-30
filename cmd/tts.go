package main

import (
	"os"
	"path"
	"text/template"

	"github.com/alecthomas/kong"
	"github.com/pkg/errors"
	"github.com/sayden/counters/fsops"
)

var templ = `
function onLoad()
    tiles = {
		{{ range .Files }}
		"{{ . }}",{{ end }}
	}

    markers_new = {
        {
            path = "file:///home/mcastro/go/src/github.com/sayden/counters/newgame/markers/017.png",
            bpath = "file:///home/mcastro/go/src/github.com/sayden/counters/newgame/markers/018.png",
            multiplier = 10
        }
    }

    x = 0
    z = -10
    for i, path in ipairs(tiles) do
        if x - math.floor(x / 5) * 5 == 0 then
            z = z + 1
            x = 0
        end
        createTile(path, tiles[i], x, z)
        x = x + 1
    end
end

function createTileOnlyFront(path, x, z)
    params = {
        type = 'Custom_Tile',
        position = { x = x, y = 1, z = z },
        scale = { x = 0.3, y = 0.3, z = 0.3 },
        callback_function = function(obj) myCustomTile(obj, path, back) end
    }
    spawnObject(params)
end

function createTile(path, back, x, z)
    params = {
        type = 'Custom_Tile',
        position = { x = x, y = 1, z = z },
        scale = { x = 0.3, y = 0.3, z = 0.3 },
        callback_function = function(obj) myCustomTile(obj, path, back) end
    }
    spawnObject(params)
end

function myCustomTile(obj, path, back)
    params = {
        image = path,
        image_bottom = back,
        type = 0,
        thickness = 0.1,
        stackable = true
    }

    obj.setCustomObject(params)
    obj.setLock(false)
    obj.reload()
end
`

type TTS struct {
	InputFolder string `arg:"Input folder to read" required:"true"`
}

type ttsData struct {
	Files []string
}

func (t *TTS) Run(ctx *kong.Context) error {
	templ, err := template.New("tts").Parse(templ)
	if err != nil {
		return err
	}

	files, err := fsops.ListFiles(t.InputFolder)
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return errors.New("no files found in the specified path")
	}

	for i, file := range files {
		files[i] = "file://" + path.Join(t.InputFolder, file)
	}

	ttsData := ttsData{Files: files}

	if err = templ.Execute(os.Stdout, ttsData); err != nil {
		return err
	}

	return nil
}
