package workers

import "sync/atomic"

type Progress struct {
	total int64
	done  atomic.Int64
}

func NewProgress(total int) *Progress {
	return &Progress{total: int64(total)}
}

func (p *Progress) Done() {
	p.done.Add(1)
}

func (p *Progress) Snapshot() (int64, int64) {
	return p.done.Load(), p.total
}
