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
	"github.com/sayden/counters"
)

const (
	padding  = 2
	maxWidth = 80
)

type globalState struct {
	counterPos int
	fileNumber int
	// filenamesInUse map[string]bool
	filenamesInUse *sync.Map
	template       *counters.CounterTemplate
	sync.RWMutex
}

func (gs *globalState) setCounterPos(pos int) {
	gs.Lock()
	defer gs.Unlock()
	gs.counterPos = pos
}

func (gs *globalState) getCounterPos() int {
	gs.RLock()
	defer gs.RUnlock()
	return gs.counterPos
}

func (gs *globalState) incrFilenumber() {
	gs.Lock()
	defer gs.Unlock()
	gs.fileNumber++
}

func (gs *globalState) getFileNumber() int {
	gs.RLock()
	defer gs.RUnlock()
	return gs.fileNumber
}

func newGlobalState(template *counters.CounterTemplate) *globalState {
	return &globalState{
		// filenamesInUse: make(map[string]bool),
		filenamesInUse: new(sync.Map),
		fileNumber:     1,
		template:       template,
	}
}

// CountersToPNG generates PNG images based on the provided CounterTemplate.
func CountersToPNG(template *counters.CounterTemplate) {
	_ = os.MkdirAll(template.OutputFolder, 0750)

	gs := newGlobalState(template)

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

	ch := make(chan *counters.Counter, 50)
	for i := 0; i < 50; i++ {
		go generateCounterToFile(ch, template.DrawGuides, gs, pbar, template.Vassal.SideName)
	}

	for i := 0; i < total; i++ {
		counter := template.Counters[i]
		ch <- &counter
	}

	pbar.Send(100)
	pbar.Wait()
	close(ch)
}

func generateCounterToFile(ch <-chan *counters.Counter, drawGuides bool, gs *globalState, pbar *tea.Program, vassalSide string) {
	for counter := range ch {
		if counter.Skip {
			pbar.Send(1)
			continue
		}

		counterCanvas, err := counter.Canvas(drawGuides)
		if err != nil {
			log.Error("error trying to create counter canvas", err)
			pbar.Send(1)
			continue
		}

		iw := imageWriter{canvas: counterCanvas, template: gs.template}
		if err = iw.createFile(counter, gs); err != nil {
			log.Error("error trying to write counter to file", err)
			pbar.Send(1)
			continue
		}

		pbar.Send(1)
	}
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
		if counter.Skip {
			continue
		}

		filepath := path.Join(iw.template.OutputFolder, counter.Filename)

		if err := iw.canvas.SavePNG(filepath); err != nil {
			return fmt.Errorf("could not save PNG file: %w", err)
		}

		gs.incrFilenumber()
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
