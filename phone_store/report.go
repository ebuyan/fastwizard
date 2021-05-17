package phonestore

import (
	"strconv"
	"sync"
)

type PhoneReport struct {
	*sync.Mutex
	*Phones
	storage *PhoneStorage
}

type Phones struct {
	BlackList map[string]bool
	Unique    map[string]bool
}

func NewPhoneReport(storage *PhoneStorage, blackListPhones map[string]bool) *PhoneReport {
	phones := Phones{blackListPhones, make(map[string]bool)}
	return &PhoneReport{&sync.Mutex{}, &phones, storage}
}

func (p *PhoneReport) GetReport() (report Report) {
	for sourceType, sources := range p.storage.GetAll() {
		for included, source := range sources.(map[interface{}]interface{}) {
			source, ok := source.(map[interface{}]interface{})
			if included == "include" && ok {
				for name := range source {
					sourceKey := p.storage.GetSourceKey(sourceType, name.(string))
					part := 0
					wg := sync.WaitGroup{}
					for {
						source, ok := p.storage.GetSource(sourceKey, strconv.Itoa(part))
						part++
						if !ok {
							break
						}
						wg.Add(1)
						go p.processPhones(source, &report, &wg)
					}
					wg.Wait()
				}
			}
		}
	}
	return
}

func (p *PhoneReport) processPhones(source Source, report *Report, wg *sync.WaitGroup) {
	p.Lock()
	defer wg.Done()
	defer p.Unlock()
	for _, phones := range source.Phones {
		arr := phones.([]interface{})
		phone := arr[0].(string)
		report.All++
		if len(phone) == 0 {
			report.Bad++
			continue
		}
		if _, ok := p.BlackList[phone]; ok {
			report.BlackList++
			continue
		}
		if _, ok := p.Unique[phone]; ok {
			report.Duplicate++
			continue
		}
		p.Unique[phone] = true
		report.Good++
	}
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
