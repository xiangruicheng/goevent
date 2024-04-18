package event

import (
	"encoding/json"
	"errors"
	"fmt"
	"goevent/config"
	"goevent/store"
	"goevent/subscribe"
	"goevent/util"
	"runtime"
	"sync"
	"time"
)

// EventObj event object
type EventObj struct {
	Topic  string `json:"topic"`
	Name   string `json:"name"`
	Note   string `json:"note"`
	Offset int32  `json:"offset"`
}

// CreateEvent create event
func CreateEvent(eventObj *EventObj) error {
	eventList, _ := GetAllEvent()
	for _, obj := range eventList {
		if obj.Topic == eventObj.Topic {
			return errors.New("event already existed")
		}
	}
	err := store.GetInstance().Topic(eventObj.Topic).ObjType(store.ObjTypeEvent).Add(eventObj)
	return err
}

func GetEvent(topic string) (*EventObj, error) {
	eventList, err := GetAllEvent()
	if err != nil {
		return nil, err
	}
	for _, obj := range eventList {
		if obj.Topic == topic {
			return obj, nil
		}
	}
	return nil, errors.New("event not exist")
}

// Producter producer data
func Producter(topic string, data any) error {
	eventObj, err := GetEvent(topic)
	if err != nil {
		return err
	}
	// write data
	err = store.GetInstance().Topic(topic).ObjType(store.ObjTypeData).Offset(eventObj.Offset).Add(data)
	if err != nil {
		return err
	}
	eventObj.Offset += 1

	// update
	err = Update(eventObj)
	return err
}

// Update update subscribe
func Update(obj *EventObj) error {
	list, _ := GetAllEvent()
	for key, tmpObj := range list {
		if tmpObj.Topic == obj.Topic {
			err := store.GetInstance().
				Topic(obj.Topic).
				ObjType(store.ObjTypeEvent).
				Offset(int32(key + 1)).Update(obj)
			return err
		}
	}
	return nil
}

// StartOneSubscribe Start one subscribe
func StartOneSubscribe(subscribeObj *subscribe.SubscribeObj, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		currFileOffset := subscribeObj.CurrOffset % config.Data.Event.DataFileMaxLines
		dataList, err := store.GetInstance().
			Topic(subscribeObj.Topic).
			ObjType(store.ObjTypeData).
			Offset(subscribeObj.CurrOffset).
			GetList(currFileOffset, 10)

		if len(dataList) <= 0 {
			time.Sleep(5 * time.Second)
			continue
		}
		for _, data := range dataList {
			logStr := subscribeObj.Topic + subscribeObj.Key + ":" + data
			fmt.Println(logStr)
			subscribeObj.CurrOffset += 1
		}
		err = subscribe.Update(subscribeObj)
		if err != nil {
			fmt.Println(err)
		}
	}
}

// GetAllEvent Get all event
func GetAllEvent() ([]*EventObj, error) {
	strList, err := store.GetInstance().ObjType(store.ObjTypeEvent).All()
	if err != nil {
		return nil, err
	}

	objs := []*EventObj{}
	for _, str := range strList {
		obj := new(EventObj)
		json.Unmarshal([]byte(str), &obj)
		objs = append(objs, obj)
	}
	return objs, nil
}

// Consumer consumer
func Consumer(topic string) {
	wg := new(sync.WaitGroup)
	tags := []string{}
	for {
		eventObjList, _ := GetAllEvent()
		for _, eventObj := range eventObjList {
			subscribeObjList, _ := subscribe.GetAllSubscribe(eventObj.Topic)
			for _, subscribeObj := range subscribeObjList {
				tag := subscribeObj.Topic + subscribeObj.Key
				isStart, _ := util.InArray(tag, tags)
				if !isStart {
					fmt.Printf("Start consumer: %s \n", tag)
					wg.Add(1)
					go StartOneSubscribe(subscribeObj, wg)
					tags = append(tags, tag)
				}
			}
		}

		// current goroutines
		numGoroutines := runtime.NumGoroutine()
		fmt.Printf("Number of goroutines: %d\n", numGoroutines)
		time.Sleep(time.Second * 10)
	}
	wg.Wait()

}
