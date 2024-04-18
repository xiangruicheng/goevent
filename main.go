package main

import (
	"goevent/config"
	"goevent/event"
	"goevent/subscribe"
)

func main() {
	config.InitConfig()

	obj := new(event.EventObj)
	obj.Topic = "login_event"
	obj.Name = "登录事件"
	obj.Note = "this is note"

	/*	err := event.CreateEvent(obj)
		fmt.Println(err)

		for i := 0; i < 1200; i++ {
			event.Producter(obj.Topic, i)
		}*/

	subscribeObj := new(subscribe.SubscribeObj)
	subscribeObj.Topic = "order_event"
	subscribeObj.Key = "subscribe_3"
	subscribeObj.Callback = "http://127.0.0.1:9217/callback"
	subscribeObj.Status = 1
	subscribeObj.BeginOffset = 0
	subscribeObj.CurrOffset = 1199

	//err := subscribe.CreateSubscribe(subscribeObj)
	//fmt.Println(err)

	event.Consumer(obj.Topic)

}
