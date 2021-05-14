package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"wizard/redis"
	"wizard/repository"

	"github.com/elliotchance/phpserialize"
)

type Handler struct{}

func (h Handler) report(w http.ResponseWriter, r *http.Request) {
	req := Request{}
	json.NewDecoder(r.Body).Decode(&req)
	if len(req.Key) == 0 {
		http.Error(w, "Empty key", http.StatusInternalServerError)
		return
	}
	phoneStorage, err := h.getPhoneStorage(req.Key)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	excludedList, err := h.getExcludedList(phoneStorage)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	report, err := h.getReport(phoneStorage, excludedList)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, _ := json.Marshal(report)
	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h Handler) getReport(phoneStorage PhoneStorage, excludedList ExcludedList) (report Report, err error) {
	mutex := sync.Mutex{}
	for key, source := range phoneStorage.Book {
		if key == "include" {
			for name := range source.(map[interface{}]interface{}) {
				redisKey := fmt.Sprintf("%s:book:%s", phoneStorage.Key, name)
				part := 0
				wg := sync.WaitGroup{}
				for {
					source, ok := redis.Cli.HGet(redisKey, strconv.Itoa(part))
					part++
					if !ok {
						break
					}
					wg.Add(1)
					go h.processPhones(source, &excludedList, &report, &mutex, &wg)
				}
				wg.Wait()
			}
		}
	}
	return
}

func (h Handler) processPhones(source []byte, list *ExcludedList, report *Report, mut *sync.Mutex, wg *sync.WaitGroup) {
	dto := SourceDto{}
	err := phpserialize.Unmarshal(source, &dto)
	if err != nil {
		wg.Done()
		fmt.Println(err)
		return
	}
	mut.Lock()
	for _, phones := range dto.Phones {
		arr := phones.([]interface{})
		phone := arr[0].(string)
		report.All++
		if len(phone) == 0 {
			report.Bad++
			continue
		}
		if _, ok := list.BlackList[phone]; ok {
			report.BlackList++
			continue
		}
		if _, ok := list.Unique[phone]; ok {
			report.Duplicate++
			continue
		}
		list.Unique[phone] = true
		report.Good++
	}
	wg.Done()
	mut.Unlock()
}

func (h Handler) getExcludedList(storage PhoneStorage) (list ExcludedList, err error) {
	phones, err := repository.BlackListRepository{}.FindBlackListPhones()
	if err != nil {
		return
	}
	list.BlackList = phones
	list.Unique = make(map[string]bool)
	return
}

func (h Handler) getPhoneStorage(key string) (storage PhoneStorage, err error) {
	res, ok := redis.Cli.Get(key)
	if !ok {
		err = errors.New("PhoneStorage not found")
		return
	}
	err = phpserialize.Unmarshal(res, &storage)
	return
}

type Request struct {
	Key              string `json:"dispatch_key"`
	RemoveDuplicates bool   `json:"remove_duplicates"`
	ResetStorage     bool   `json:"reset_storage"`
}

type Report struct {
	Duplicate int    `json:"duplicate"`
	Good      int    `json:"good"`
	Bad       int    `json:"bad"`
	All       int    `json:"all"`
	BlackList int    `json:"black_list"`
	Excluded  int    `json:"excluded"`
	Error     string `json:"error"`
}

type ExcludedList struct {
	Excluded  map[string]bool
	BlackList map[string]bool
	Unique    map[string]bool
}

type PhoneStorage struct {
	File   map[interface{}]interface{} `php:"file"`
	Manual map[interface{}]interface{} `php:"manual"`
	Book   map[interface{}]interface{} `php:"book"`
	Key    interface{}                 `php:"key"`
}

type SourceDto struct {
	Name   string        `php:"name"`
	Phones []interface{} `php:"phones"`
}
