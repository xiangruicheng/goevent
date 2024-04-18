package subscribe

import (
	"encoding/json"
	"errors"
	"goevent/store"
)

// SubscribeObj Subscribe object
type SubscribeObj struct {
	Topic       string `json:"topic"`        //主题
	Key         string `json:"key"`          //唯一标识
	Callback    string `json:"callback"`     //回调url
	OtherData   string `json:"other_data"`   //附加数据
	Status      int8   `json:"status"`       //状态 1启用2暂停
	BeginOffset int32  `json:"begin_offset"` //开始消费的游标
	CurrOffset  int32  `json:"curr_offset"`  //当前游标
}

// CreateSubscribe Create Subscribe
func CreateSubscribe(subscribeObj *SubscribeObj) error {
	subscribeObjList, err := GetAllSubscribe(subscribeObj.Topic)
	if err != nil {
		return store.GetInstance().Topic(subscribeObj.Topic).ObjType(store.ObjTypeSubscribe).Add(subscribeObj)
	}
	for _, obj := range subscribeObjList {
		if obj.Key == subscribeObj.Key {
			return errors.New("Key is exits")
		}
	}
	return store.GetInstance().Topic(subscribeObj.Topic).ObjType(store.ObjTypeSubscribe).Add(subscribeObj)
}

// GetAllSubscribe Get all subscribe
func GetAllSubscribe(topic string) ([]*SubscribeObj, error) {
	strList, err := store.GetInstance().Topic(topic).ObjType(store.ObjTypeSubscribe).All()
	if err != nil {
		return nil, err
	}

	objs := []*SubscribeObj{}
	for _, str := range strList {
		obj := new(SubscribeObj)
		json.Unmarshal([]byte(str), &obj)
		objs = append(objs, obj)
	}
	return objs, nil
}

// Update update subscribe
func Update(obj *SubscribeObj) error {
	list, _ := GetAllSubscribe(obj.Topic)
	for key, tmpObj := range list {
		if tmpObj.Key == obj.Key {
			err := store.GetInstance().
				Topic(obj.Topic).
				ObjType(store.ObjTypeSubscribe).
				Offset(int32(key + 1)).Update(obj)
			return err
		}
	}
	return nil
}
