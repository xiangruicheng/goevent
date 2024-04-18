package store

import "sync"

type ObjType string

const (
	ObjTypeEvent     ObjType = "event"
	ObjTypeData      ObjType = "data"
	ObjTypeSubscribe ObjType = "subscribe"
)

var (
	SubscribeFileMutex *sync.Mutex
	EventFileMutex     *sync.Mutex
	DataFileMutex      *sync.Mutex
)

func init() {
	EventFileMutex = new(sync.Mutex)
	SubscribeFileMutex = new(sync.Mutex)
	DataFileMutex = new(sync.Mutex)
}

type Store interface {
	Topic(topic string) Store
	ObjType(objType ObjType) Store
	Offset(offset int32) Store
	Add(data any) error
	GetList(offset int32, length int32) ([]string, error)
	Get(offset int32) (string, error)
	Update(obj any) error
	All() ([]string, error)
}

func GetInstance() Store {
	return newFileStroe()
}

var fileStore *FileStore

func newFileStroe() Store {
	fileStore = new(FileStore)
	return fileStore
}

/*
var redisStore *RedisStore

func newRedisStroe() Store {
	redisStore = new(RedisStore)
	return redisStore
}*/
