package event

import (
	"reflect"
)

type eventStore struct {
	fn    reflect.Value
	order float32
	once  bool
}

//eventsStore is array of  eventStore for sorting
type eventsStore []eventStore

func (slice eventsStore) Len() int {
	return len(slice)
}

func (slice eventsStore) Less(i, j int) bool {
	return slice[i].order < slice[j].order
}

func (slice eventsStore) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}
