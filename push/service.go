package push

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/rus-sharafiev/go-push/db"
)

var appError *error

type service struct {
	db *db.Postgres
	wg *sync.WaitGroup
}

func (s service) Start() error {
	fmt.Println("started")

	interval := 60
	if Config.Interval != nil {
		interval = *Config.Interval
	}

	scheduleTicker := time.NewTicker(time.Duration(interval) * time.Minute)
	eventTicker := time.NewTicker(time.Minute)
	var newWg sync.WaitGroup
	s.wg = &newWg
	s.wg.Add(1)

	go s.handleEvents()
	go s.handleShedule()

	go func() {
		for {
			select {
			case <-eventTicker.C:
				s.handleEvents()
			case <-scheduleTicker.C:
				s.handleShedule()
			}
		}
	}()

	s.wg.Wait()
	fmt.Println("stopped")

	scheduleTicker.Stop()
	eventTicker.Stop()

	return *appError
}

func (s service) handleEvents() {
	fmt.Println("run events")
	e := make(chan error)

	pushes, err := s.getEventPushes()
	if err != nil {
		appError = &err
		s.wg.Done()
	}

	for _, push := range pushes {
		if *push.Install {
			go s.getUserTokensAndSendByEvent(&push, "installed_at", e)
		}
		if *push.Reg {
			go s.getUserTokensAndSendByEvent(&push, "registered_at", e)
		}
		if *push.Dep {
			go s.getUserTokensAndSendByEvent(&push, "deposit_made_at", e)
		}
	}

	err = <-e
	if err != nil {
		fmt.Println(err)
	}
}

func (s service) handleShedule() {
	fmt.Println("run schedule")
	c := make(chan []Push)
	e := make(chan error)
	var pushes []Push

	go s.getOneTimePushes(c, e)
	go s.getSchedulePushes(c, e)

	for i := 0; i < 2; i++ {
		select {
		case p := <-c:
			pushes = append(pushes, p...)
		case e := <-e:
			appError = &e
			s.wg.Done()
			return
		}
	}

	for _, push := range pushes {
		go s.SendPush(push)
	}
}

func (s service) SendPush(push Push) {
	qty, err := s.getUsersQty(push.AdvertId)
	if err != nil {
		appError = &err
		s.wg.Done()
	}

	if qty == 0 {
		return
	}

	take := 5
	if Config.BatchSize != nil {
		take = *Config.BatchSize
	}
	batches := qty/take + 1

	e := make(chan error)
	for i := 0; i < batches; i++ {
		fmt.Println(take, take*i)
		go s.getUserTokensAndSend(push, take, take*i, e)
	}

	for i := 0; i < batches; i++ {
		err := <-e
		if err != nil {
			fmt.Println(err)
		}
	}
}

func (s service) SendBatch(data map[string]string, tokens []string) error {
	message := Message{
		Data:   data,
		Tokens: tokens,
	}

	b, _ := json.MarshalIndent(message, "", "  ")
	fmt.Println(string(b))
	return nil
}

var Service = &service{db: &db.Instance}
