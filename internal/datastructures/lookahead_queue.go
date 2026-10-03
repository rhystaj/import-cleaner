package datastructures

import (
	"errors"
	"iter"
)

type Lookahead[E any] struct {
	queue     []*E
	lookahead int
}

func CreateNewLookahead[E any](lookahead int) (*Lookahead[E], error) {
	if lookahead < 1 {
		return nil, errors.New("lookahead must be at least 1")
	}

	return (&Lookahead[E]{
		queue:     make([]*E, lookahead),
		lookahead: lookahead,
	}), nil
}

func (q *Lookahead[E]) PushItem(value *E) {
	for i := 0; i < q.lookahead-1; i++ {
		q.queue[i] = q.queue[i+1]
	}
	q.queue[q.lookahead-1] = value
}

func (q *Lookahead[E]) PeekStart() *E {
	return q.queue[0]
}

func (q *Lookahead[E]) PeekEnd() *E {
	return q.queue[q.lookahead-1]
}

func (q *Lookahead[E]) ProcessSequence(sequence iter.Seq[E]) iter.Seq[*Lookahead[E]] {
	return func(yield func(*Lookahead[E]) bool) {
		prefillIndex := 1

		for item := range sequence {
			// Prefill the queue with lookahead-1 items so that the queue is always full
			// when the first item is yielded.
			if prefillIndex <= q.lookahead-1 {
				q.queue[prefillIndex] = &item
				prefillIndex++
				continue
			}

			q.PushItem(&item)
			if !yield(q) {
				return
			}
		}

		for i := 0; i < q.lookahead-1; i++ {
			q.PushItem(nil)
			if !yield(q) {
				return
			}
		}
	}
}
