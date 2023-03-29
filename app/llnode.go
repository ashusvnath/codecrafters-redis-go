package main

type llnode[D any] struct {
	data D
	n    *llnode[D]
}

func (llnd *llnode[D]) HasNext() bool {
	return llnd.n != nil
}

func (llnd *llnode[D]) Get() D {
	return llnd.data
}

func (llnd *llnode[D]) Next() llnode[D] {
	return *llnd.n
}

func (llnd *llnode[D]) Append(data D) *llnode[D]{
	if llnd == nil {
		return makeNode(data)
	}
	if llnd.n != nil {
		llnd.n.Append(data)
		return llnd
	}
	llnd.n = makeNode(data)
	return llnd
}

func makeNode[D any](data D) *llnode[D] {
	return &llnode[D]{data, nil}
}
