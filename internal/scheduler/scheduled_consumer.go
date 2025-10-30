package scheduler

import (
	"log"
	"os"

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
	// Obter schedule do ambiente ou usar padrão (a cada hora)
	schedule := os.Getenv("CRON_SCHEDULE")
	if schedule == "" {
		schedule = "0 * * * *" // Padrão: a cada hora no minuto 0
	}

	_, err := s.cron.AddFunc(schedule, func() {
		log.Println("========================================")
		log.Println("Running scheduled lottery update...")
		log.Println("========================================")
		s.loteriasUpdate.UpdateAll()
		log.Println("========================================")
		log.Println("Scheduled lottery update completed")
		log.Println("========================================")
	})

	if err != nil {
		log.Printf("Error scheduling task: %v", err)
		return
	}

	s.cron.Start()
	log.Printf("Scheduler started - lottery updates will run with schedule: %s", schedule)
	log.Println("Running initial lottery update...")

	// Executar atualização inicial em background
	go s.loteriasUpdate.UpdateAll()
}

func (s *ScheduledConsumer) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}
