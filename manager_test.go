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

		eventManager.Subscribe("MyEvent1", func(event interface{}) {
			myEvent := event.(*MyEvent1)
			myEvent.val++
		}, 1)
		eventManager.Subscribe("MyEvent1", func(event interface{}) {
			myEvent := event.(*MyEvent1)
			myEvent.val += 3
			//1+2+3 . run at End
			So(myEvent.val, ShouldEqual, 6)
		}, 3)
		eventManager.Subscribe("MyEvent1", func(event interface{}) {
			myEvent := event.(*MyEvent1)
			myEvent.val += 2
			// 1+2. event.val += 3 run after this event
			So(myEvent.val, ShouldEqual, 3)
		}, 2)

		event := MyEvent1{}
		eventManager.Publish("MyEvent1", &event)

		//1+2+3
		So(event.val, ShouldEqual, 6)

	})

	Convey("event test - test once", t, func() {
		eventManager := NewManager()

		eventManager.Subscribe("MyEvent1", func(event interface{}) {
			myEvent := event.(*MyEvent1)
			myEvent.val++
		}, 1)
		eventManager.Subscribe("MyEvent1", func(event interface{}) {
			myEvent := event.(*MyEvent1)
			myEvent.val += 3
		}, 3)
		eventManager.SubscribeOnce("MyEvent1", func(event interface{}) {
			myEvent := event.(*MyEvent1)
			myEvent.val += 2
			So(myEvent.val, ShouldEqual, 3)
		}, 2)

		event := MyEvent1{}
		eventManager.Publish("MyEvent1", &event)
		eventManager.Publish("MyEvent1", &event)

		//1+2+3 and again 1+3
		So(event.val, ShouldEqual, 10)

	})

}

func BenchmarkEvent(b *testing.B) {
	for n := 0; n < b.N; n++ {
		eventManager := NewManager()

		eventManager.Subscribe("MyEvent1", func(event interface{}) {
			myEvent := event.(*MyEvent1)
			myEvent.val++
		}, 1)
		eventManager.Subscribe("MyEvent1", func(event interface{}) {
			myEvent := event.(*MyEvent1)
			myEvent.val += 3
		}, 3)
		eventManager.Subscribe("MyEvent1", func(event interface{}) {
			myEvent := event.(*MyEvent1)
			myEvent.val += 2
		}, 2)

		event := MyEvent1{}
		eventManager.Publish("MyEvent1", &event)

	}
}
