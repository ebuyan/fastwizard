package phonestore

import (
	"strconv"
	"sync"
	"wizard/pkg/memory"
)

type PhoneReport struct {
	sync.Mutex
	Phones
	storage PhoneStorage
}

type Phones struct {
	BlackList map[int]bool
	Unique    map[int]bool
}

func NewPhoneReport(storage PhoneStorage, blackListPhones map[int]bool) *PhoneReport {
	return &PhoneReport{sync.Mutex{}, Phones{blackListPhones, make(map[int]bool)}, storage}
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
						memory.Usage()
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
	p.Phones.Unique = nil
	p.Phones.BlackList = nil
	return
}

func (p *PhoneReport) processPhones(source Source, report *Report, wg *sync.WaitGroup) {
	p.Lock()
	defer wg.Done()
	defer p.Unlock()
	for _, phones := range source.Phones {
		ph := phones.([]interface{})[0].(string)
		phone, err := strconv.Atoi(ph)

		report.All++
		if err != nil || phone == 0 {
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
