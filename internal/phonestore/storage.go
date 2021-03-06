package phonestore

import (
	"errors"
	"fmt"
	"log"
	"wizard/pkg/redis"

	"github.com/elliotchance/phpserialize"
)

type PhoneStorage struct {
	Store
	*redis.Client
}

func NewPhoneStorage(redis *redis.Client, key string) (storage PhoneStorage, err error) {
	res, ok := redis.Get(key)
	if !ok {
		err = errors.New("PhoneStorage not found")
		return
	}
	store := Store{}
	err = phpserialize.Unmarshal(res, &store)
	storage = PhoneStorage{}
	storage.Store = store
	storage.Client = redis
	return
}

func (p PhoneStorage) GetSource(sourceKey, sourcePart string) (source Source, ok bool) {
	res, ok := p.HGet(sourceKey, sourcePart)
	if !ok {
		return
	}
	err := phpserialize.Unmarshal(res, &source)
	if err != nil {
		log.Println(err)
		ok = false
		return
	}
	return
}

func (p PhoneStorage) GetSourceKey(sourceType, sourceName string) string {
	return fmt.Sprintf("%s:%s:%s", p.Key, sourceType, sourceName)
}

func (p PhoneStorage) GetAll() (s Sources) {
	s = Sources{}
	s["book"] = p.Book
	s["file"] = p.File
	s["manual"] = p.Manual
	return
}

type Sources map[string]interface{}

type Store struct {
	File   map[interface{}]interface{} `php:"file"`
	Manual map[interface{}]interface{} `php:"manual"`
	Book   map[interface{}]interface{} `php:"book"`
	Key    interface{}                 `php:"key"`
}

type Source struct {
	Name   string        `php:"name"`
	Phones []interface{} `php:"phones"`
}
