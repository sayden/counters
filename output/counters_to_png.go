package output

import (
	"fmt"
	"os"
	"path"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	"github.com/fogleman/gg"
	"github.com/pkg/errors"
	"github.com/sayden/counters"
)

const (
	padding  = 2
	maxWidth = 80
)

type globalState struct {
	counterPos     int
	fileNumber     int
	filenamesInUse map[string]bool
	template       *counters.CounterTemplate
	canvas         *gg.Context
}

func newGlobalState(template *counters.CounterTemplate, canvas *gg.Context) *globalState {
	return &globalState{
		filenamesInUse: make(map[string]bool),
		fileNumber:     1,
		canvas:         canvas,
		template:       template,
	}
}

// CountersToPNG generates PNG images based on the provided CounterTemplate.
func CountersToPNG(template *counters.CounterTemplate) error {
	var canvas *gg.Context
	_ = os.MkdirAll(template.OutputFolder, 0750)

	gs := newGlobalState(template, canvas)

	// Progress bar
	total := 0
	for _, c := range template.Counters {
		if !c.Skip {
			total += *c.Multiplier
		}
	}
	log.Info("Total number of counters: ", "total", total+1)

	prog := progress.New(progress.WithScaledGradient("#FF7CCB", "#FDFF8C"))
	pbar := tea.NewProgram(&progressBar{progress: prog, total: float64(total + 1)})
	go pbar.Run()
	defer pbar.Quit()

	var wg sync.WaitGroup
	wg.Add(total)
	for i := 0; i < total; i++ {
		counter := template.Counters[i]
		go func(counter *counters.Counter, pbar *tea.Program) {
			defer wg.Done()
			if counter.Skip {
				return
			}

			counterCanvas, err := counter.Canvas(template.DrawGuides)
			if err != nil {
				log.Error("error trying to create counter canvas", err)
			}

			if err = writeCounterToFile(counterCanvas, counter, gs); err != nil {
				log.Error("error trying to write counter to file", err)
			}

			pbar.Send(i + 1)
		}(&counter, pbar)
	}

	wg.Wait()
	pbar.Send(100)
	pbar.Wait()

	return nil
}

func writeCounterToFile(dc *gg.Context, counter *counters.Counter, gs *globalState) error {
	iw := imageWriter{
		canvas:   dc,
		template: gs.template,
	}

	if err := iw.createFile(counter, gs); err != nil {
		return errors.Wrap(err, "error trying to write counter file")
	}

	return nil
}

type imageWriter struct {
	canvas   *gg.Context
	template *counters.CounterTemplate
}

// createFile creates a file with the counter image. Filenumber is the filename, a pointer is passed to be able to use
// the multiplier to create more than one file with the same counter
func (iw *imageWriter) createFile(counter *counters.Counter, gs *globalState) error {
	// Use sequencing of numbers or a position in the counter texts to name files
	for i := 0; i < *counter.Multiplier; i++ {
		counterFilename := counter.GetCounterFilename(iw.template.PositionNumberForFilename, "",
			gs.fileNumber, gs.filenamesInUse)

		filepath := path.Join(iw.template.OutputFolder, counterFilename)
		if counter.Skip {
			continue
		}

		log.Debug("Saving file: ", filepath)
		if err := iw.canvas.SavePNG(filepath); err != nil {
			log.Error("file", filepath, "could not save PNG file")
			return err
		}
		gs.fileNumber++
	}

	return nil
}

type progressBar struct {
	total    float64
	percent  float64
	current  int32
	progress progress.Model
	sync.Mutex
}

func (m *progressBar) Init() tea.Cmd {
	return nil
}

func (m *progressBar) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case int:
		cur := atomic.AddInt32(&m.current, 1)
		m.Lock()
		defer m.Unlock()
		percent := float64(cur) / m.total
		m.percent = percent
		if percent >= 1.0 || cur >= int32(m.total) {
			m.percent = 1.0
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width >= maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	default:
		return m, tea.Quit
	}
}

func (m *progressBar) View() string {
	pad := strings.Repeat(" ", padding)
	return fmt.Sprintf("\n%s%s (%d/%d)\n", pad, m.progress.ViewAs(m.percent), m.current, int(m.total))
}
