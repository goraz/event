package event

import (
	"fmt"
	"reflect"
	"sort"
	"sync"
)

//Manager is simple event manager
type Manager struct {
	lock   sync.Mutex
	events map[reflect.Type]eventsStore
}

func (manager *Manager) listen(fn interface{}, order float32, once bool) error {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	value := reflect.ValueOf(fn)
	if value.Kind() != reflect.Func {
		return fmt.Errorf("%s is not of type reflect.Func", value.Kind())
	}

	argsType := value.Type()
	if argsType.NumIn() != 1 {
		return fmt.Errorf("Handler must have a single argument")
	}
	structType := argsType.In(0)
	manager.events[structType] = append(manager.events[structType], eventStore{
		fn:    value,
		order: order,
		once:  once,
	})
	sort.Sort(manager.events[structType])
	return nil
}

//Subscribe  listen to  event
func (manager *Manager) Subscribe(fn interface{}, order float32) {
	err := manager.listen(fn, order, false)
	if err != nil {
		panic(err)
	}
}

//SubscribeOnce listen to  event for once
func (manager *Manager) SubscribeOnce(fn interface{}, order float32) {
	err := manager.listen(fn, order, true)
	if err != nil {
		panic(err)
	}
}

// Publish executes callback defined for a event
func (manager *Manager) Publish(event interface{}) {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	structType := reflect.TypeOf(event)
	if events, ok := manager.events[structType]; ok {
		args := [...]reflect.Value{reflect.ValueOf(event)}
		for index, event := range events {
			if event.once {
				manager.remove(structType, index)
			}
			event.fn.Call(args[:])

		}
	}
}

func (manager *Manager) remove(structType reflect.Type, index int) {
	if index >= 0 {
		manager.events[structType] = append(manager.events[structType][:index],
			manager.events[structType][index+1:]...)
	}
}

//NewManager retrun event manager
func NewManager() *Manager {
	return &Manager{
		sync.Mutex{},
		make(map[reflect.Type]eventsStore),
	}
}
