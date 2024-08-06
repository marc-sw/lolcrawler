package crawler

import "fmt"

type Progress struct {
	total     int
	completed int
	stepRange int
}

func NewProgress(total int) Progress {
	stepRange := total / 100
	if stepRange == 0 {
		stepRange = 1
	}
	return Progress{total: total, completed: 0, stepRange: stepRange}
}

func (p *Progress) Increase() {
	p.completed++
	if p.completed%p.stepRange == 0 {
		fmt.Printf("\r%d%c", p.completed*100/p.total, '%')
	}
}

func (p *Progress) Done() bool {
	return p.completed == p.total
}
