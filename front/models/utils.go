package model

type SortableString []string

func (this SortableString) Len() int {
	return len(this)
}

func (this SortableString) Less(i, j int) bool {
	return (this)[i] < (this)[j]
}

func (this SortableString) Swap(i, j int) {
	temp := this[i]
	this[i] = this[j]
	this[j] = temp
}
