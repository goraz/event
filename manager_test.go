package event

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

type MyEvent1 struct {
	val int
}

func TestEvent(t *testing.T) {
	Convey("event test - test order", t, func() {
		eventManager := NewManager()
		// event := My{}

		eventManager.Subscribe(func(event *MyEvent1) {
			event.val++
		}, 1)
		eventManager.Subscribe(func(event *MyEvent1) {
			event.val += 3
			//1+2+3 . run at End
			So(event.val, ShouldEqual, 6)
		}, 3)
		eventManager.Subscribe(func(event *MyEvent1) {
			event.val += 2
			// 1+2. event.val += 3 run after this event
			So(event.val, ShouldEqual, 3)
		}, 2)

		event := MyEvent1{}
		eventManager.Publish(&event)

		//1+2+3
		So(event.val, ShouldEqual, 6)

	})

	Convey("event test - test once", t, func() {
		eventManager := NewManager()
		// event := My{}

		eventManager.Subscribe(func(event *MyEvent1) {
			event.val++
		}, 1)
		eventManager.Subscribe(func(event *MyEvent1) {
			event.val += 3
		}, 3)
		eventManager.SubscribeOnce(func(event *MyEvent1) {
			event.val += 2
			So(event.val, ShouldEqual, 3)
		}, 2)

		event := MyEvent1{}
		eventManager.Publish(&event)
		eventManager.Publish(&event)

		//1+2+3 and again 1+3
		So(event.val, ShouldEqual, 10)

	})

}

func TestEventError(t *testing.T) {
	Convey("event test - panic not function kind", t, func() {
		eventManager := NewManager()
		// event := My{}
		So(func() {
			eventManager.Subscribe(1, 1)
		}, ShouldPanic)
	})

	Convey("event test - panic on not valid args", t, func() {
		eventManager := NewManager()
		// event := My{}
		So(func() {
			eventManager.Subscribe(func(event *MyEvent1, other int) {
				event.val++
			}, 1)
		}, ShouldPanic)

		So(func() {
			eventManager.SubscribeOnce(func(event *MyEvent1, other int) {
				event.val++
			}, 1)
		}, ShouldPanic)

	})

}
