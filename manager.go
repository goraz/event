package event

import (
	"sort"
	"sync"
)

//Manager is simple event manager
type Manager struct {
	lock   sync.Mutex
	events map[string]eventsStore
}

func (manager *Manager) listen(name string, fn func(interface{}) bool, order float32, once bool) {
	manager.lock.Lock()
	defer manager.lock.Unlock()

	manager.events[name] = append(manager.events[name], eventStore{
		fn:    fn,
		order: order,
		once:  once,
	})
	sort.Sort(manager.events[name])
	return
}

//Subscribe  listen to  event
func (manager *Manager) Subscribe(name string, fn func(interface{}) bool, order float32) {
	manager.listen(name, fn, order, false)
}

//SubscribeOnce listen to  event for once
func (manager *Manager) SubscribeOnce(name string, fn func(interface{}) bool, order float32) {
	manager.listen(name, fn, order, true)
}

// Publish executes callback defined for event. if func return false stop execute and return false
func (manager *Manager) Publish(name string, arg interface{}) bool {
	manager.lock.Lock()
	defer manager.lock.Unlock()
	if events, ok := manager.events[name]; ok {
		for index, event := range events {
			if event.once {
				manager.remove(name, index)
			}
			if event.fn(arg) == false {
				return false
			}
		}
	}
	return true
}

func (manager *Manager) remove(name string, index int) {
	if index >= 0 {
		manager.events[name] = append(manager.events[name][:index],
			manager.events[name][index+1:]...)
	}
}

//NewManager retrun event manager
func NewManager() *Manager {
	return &Manager{
		sync.Mutex{},
		make(map[string]eventsStore),
	}
}
