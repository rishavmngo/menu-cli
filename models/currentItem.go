package models

type CurrentItem struct {
	value int
}

func Get() *CurrentItem {

	return &CurrentItem{}

}

func (ci *CurrentItem) Get() int {
	return ci.value
}
func (ci *CurrentItem) Increment() {
	ci.value++
}

func (ci *CurrentItem) Decrement() {
	ci.value--
}
func (ci *CurrentItem) Set(value int) {
	ci.value = value
}
func (ci *CurrentItem) Reset() {
	ci.value = 0
}
