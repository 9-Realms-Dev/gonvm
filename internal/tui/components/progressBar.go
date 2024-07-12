package tui

import (
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type progressModel struct {
	progress progress.Model
	percent  float64
	total    int64
	current  int64
	done     bool
}

func (m progressModel) Init() tea.Cmd {
	return nil
}

func (m progressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - 4
		if m.progress.Width > 80 {
			m.progress.Width = 80
		}
		return m, nil
	case int64:
		m.current = m.current + msg
		m.percent = float64(m.current) / float64(m.total)
		if m.current >= m.total {
			m.done = true
			return m, tea.Quit
		}
		return m, nil
	}
	return m, nil
}

func (m progressModel) View() string {
	if m.done {
		return "Download complete!\n"
	}
	pad := strings.Repeat(" ", 2)
	return "\n" +
		pad + fmt.Sprintf("%.2f%% of %.2f MB", m.percent*100, float64(m.total)/1024/1024) + "\n" +
		pad + m.progress.ViewAs(m.percent) + "\n\n"
}

func newProgressModel(total int64) progressModel {
	return progressModel{
		progress: progress.New(progress.WithDefaultGradient()),
		total:    total,
	}
}

func CopyWithProgress(dst io.Writer, src io.Reader, size int64, description string) error {
	model := newProgressModel(size)
	p := tea.NewProgram(model)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		buf := make([]byte, 32*1024)
		for {
			n, err := src.Read(buf)
			if n > 0 {
				_, writeErr := dst.Write(buf[:n])
				if writeErr != nil {
					p.Send(writeErr)
					return
				}
				p.Send(model.current + int64(n))
			}
			if err != nil {
				if err != io.EOF {
					p.Send(err)
				} else {
					p.Send(nil) // Signal successful completion
				}
				return
			}
		}
	}()

	// Run the program in a separate goroutine
	errChan := make(chan error, 1)
	go func() {
		_, err := p.Run()
		fmt.Printf("Program exited: %v\n", err)
		errChan <- err
	}()

	// Wait for either the copy operation to complete or the program to exit
	var err error
	select {
	case err = <-errChan:
		// Program exited, cancel the copy operation if it's still running
		// You might need to implement a cancellation mechanism in your model
	case <-waitGroupToChan(&wg):
		// Copy operation completed, quit the program
		p.Quit()
		err = <-errChan
	}

	return err
}

// Helper function to convert WaitGroup to a channel
func waitGroupToChan(wg *sync.WaitGroup) <-chan struct{} {
	ch := make(chan struct{})
	go func() {
		wg.Wait()
		close(ch)
	}()
	return ch
}
