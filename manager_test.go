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

		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val++
			return true
		}, 1)
		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 3
			//1+2+3 . run at End
			So(myEvent.val, ShouldEqual, 6)
			return true
		}, 3)
		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 2
			// 1+2. event.val += 3 run after this event
			So(myEvent.val, ShouldEqual, 3)
			return true
		}, 2)

		event := MyEvent1{}
		eventManager.Publish("MyEvent1", &event)

		//1+2+3
		So(event.val, ShouldEqual, 6)

	})

	Convey("event test - test once", t, func() {
		eventManager := NewManager()

		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val++
			return true
		}, 1)
		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 3
			return true
		}, 3)
		eventManager.SubscribeOnce("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 2
			So(myEvent.val, ShouldEqual, 3)
			return true
		}, 2)

		event := MyEvent1{}
		eventManager.Publish("MyEvent1", &event)
		eventManager.Publish("MyEvent1", &event)

		//1+2+3 and again 1+3
		So(event.val, ShouldEqual, 10)

	})

	Convey("event test - test stop", t, func() {
		eventManager := NewManager()

		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val++
			return true
		}, 1)
		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 3
			return true
		}, 3)
		eventManager.SubscribeOnce("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 2
			So(myEvent.val, ShouldEqual, 3)
			return false
		}, 2)

		event := MyEvent1{}
		success := eventManager.Publish("MyEvent1", &event)

		//1+2  - 3 not start
		So(event.val, ShouldEqual, 3)
		So(success, ShouldEqual, false)

	})

	Convey("event test - test unsubscribe", t, func() {
		eventManager := NewManager()

		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val++
			return true
		}, 1)
		eventID := eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 3
			return true
		}, 3)
		eventManager.UnSubscribe("MyEvent1", eventID)

		event := MyEvent1{}

		success := eventManager.Publish("MyEvent1", &event)

		So(event.val, ShouldEqual, 1)
		So(success, ShouldEqual, true)

	})

	Convey("event test - test unsubscribeAll", t, func() {
		eventManager := NewManager()

		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val++
			return true
		}, 1)
		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 3
			return true
		}, 3)
		eventManager.UnSubscribeAll("MyEvent1")

		event := MyEvent1{}

		eventManager.Publish("MyEvent1", &event)

		So(event.val, ShouldEqual, 0)

	})

}

func BenchmarkEvent(b *testing.B) {
	for n := 0; n < b.N; n++ {
		eventManager := NewManager()

		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val++
			return true
		}, 1)
		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 3
			return true
		}, 3)
		eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
			myEvent := event.(*MyEvent1)
			myEvent.val += 2
			return true
		}, 2)

		event := MyEvent1{}
		eventManager.Publish("MyEvent1", &event)

	}
}

func BenchmarkPublish(b *testing.B) {
	eventManager := NewManager()

	eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
		myEvent := event.(*MyEvent1)
		myEvent.val++
		return true
	}, 1)
	eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
		myEvent := event.(*MyEvent1)
		myEvent.val += 3
		return true
	}, 3)
	eventManager.Subscribe("MyEvent1", func(event interface{}) bool {
		myEvent := event.(*MyEvent1)
		myEvent.val += 2
		return true
	}, 2)
	for n := 0; n < b.N; n++ {
		event := MyEvent1{}
		eventManager.Publish("MyEvent1", &event)

	}
}
