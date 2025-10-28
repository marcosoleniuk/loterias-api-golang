package scheduler

import (
	"log"

	"loterias-api-golang/internal/service"

	"github.com/robfig/cron/v3"
)

type ScheduledConsumer struct {
	cron           *cron.Cron
	loteriasUpdate *service.LoteriasUpdate
}

func NewScheduledConsumer(loteriasUpdate *service.LoteriasUpdate) *ScheduledConsumer {
	c := cron.New()
	return &ScheduledConsumer{
		cron:           c,
		loteriasUpdate: loteriasUpdate,
	}
}

func (s *ScheduledConsumer) Start() {
	_, err := s.cron.AddFunc("0 * * * *", func() {
		log.Println("Running scheduled lottery update...")
		s.loteriasUpdate.UpdateAll()
	})

	if err != nil {
		log.Printf("Error scheduling task: %v", err)
		return
	}

	s.cron.Start()
	log.Println("Scheduler started - lottery updates will run every hour")

	go s.loteriasUpdate.UpdateAll()
}

func (s *ScheduledConsumer) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}
