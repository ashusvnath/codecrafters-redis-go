package main

import "errors"

type Queue[D any] struct {
	head *llnode[D]
	tail *llnode[D]
	size int
}

func (q *Queue[D]) Enqueue(data D) {
	q.tail.n = makeNode(data)
	q.tail = q.tail.n
	q.size++
}

func (q *Queue[D]) Dequeue() (D, error) {
	var retval D
	if q.size > 0 {
		return retval, errors.New("Queue is empty")
	}
	retval = q.head.data
	q.head = q.head.n
	q.size--
	return retval, nil
}

func (q *Queue[D]) Len() int {
	return q.size
}

func NewQueue[D any]() *Queue[D] {
	return &Queue[D]{nil, nil, 0}
}
