package progressutil

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
)

// ProgressBar represents a progress bar that can be rendered to an output stream.
type ProgressBar struct {
	total   int
	current int
	width   int
	fill    string
	empty   string
	output  io.Writer
	mu      sync.Mutex
}

// NewProgressBar creates a new progress bar with the specified total value, width, fill and empty characters, and output writer.
// If fill is empty, it defaults to "=".
// If empty is empty, it defaults to " ".
// If output is nil, it defaults to os.Stdout.
func NewProgressBar(total int, width int, fill, empty string, output io.Writer) *ProgressBar {
	if fill == "" {
		fill = "="
	}
	if empty == "" {
		empty = " "
	}
	if output == nil {
		output = os.Stdout
	}
	return &ProgressBar{
		total:  total,
		width:  width,
		fill:   fill,
		empty:  empty,
		output: output,
	}
}

// SetProgress sets the current progress to the specified value.
// Returns an error if current is negative.
// If current exceeds total, it will be clamped to total.
func (p *ProgressBar) SetProgress(current int) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if current < 0 {
		return fmt.Errorf("current progress cannot be negative")
	}
	if current > p.total {
		current = p.total
	}
	p.current = current
	return nil
}

// Increment increases the current progress by 1.
// Returns an error if the progress is already complete.
func (p *ProgressBar) Increment() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.current >= p.total {
		return fmt.Errorf("progress already complete")
	}
	p.current++
	return nil
}

// Render writes the progress bar to the output stream.
// The progress bar is rendered as a single line, overwriting the current line.
// When progress is complete (current == total), a newline is added.
func (p *ProgressBar) Render() error {
	p.mu.Lock()
	defer p.mu.Unlock()

	percent := float64(p.current) / float64(p.total) * 100
	filled := int(percent / 100 * float64(p.width))
	bar := strings.Repeat(p.fill, filled) + strings.Repeat(p.empty, p.width-filled)

	_, err := fmt.Fprintf(p.output, "\r[%s] %.2f%%", bar, percent)
	if err != nil {
		return err
	}

	if p.current == p.total {
		_, err = fmt.Fprintln(p.output, " done!")
	}
	return err
}

// Show sets the current progress and immediately renders the progress bar to the output stream.
// It combines the functionality of SetProgress and Render in a single method call.
// Returns any error encountered while setting progress or rendering.
func (p *ProgressBar) Show(current int) error {
	if err := p.SetProgress(current); err != nil {
		return fmt.Errorf("failed to set progress: %w", err)
	}
	if err := p.Render(); err != nil {
		return fmt.Errorf("failed to render progress bar: %w", err)
	}
	return nil
}
