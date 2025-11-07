package scheduler

import (
	"log"
	"os"
	_ "time"

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
	// Múltiplos horários de verificação (como em Java)
	schedules := []string{
		"0 12 * * MON-SAT",  // 12h de segunda a sábado
		"0 21 * * MON-SAT",  // 21h
		"15 21 * * MON-SAT", // 21:15
		"0 22 * * MON-SAT",  // 22h (horário principal)
		"10 23 * * MON-SAT", // 23:10
		"20 0 * * MON-SAT",  // 00:20
		"0 1 * * MON-SAT",   // 01h
	}

	// Permitir override via variável de ambiente
	customSchedule := os.Getenv("CRON_SCHEDULE")
	if customSchedule != "" {
		schedules = []string{customSchedule}
		log.Printf("Using custom schedule: %s", customSchedule)
	}

	for _, schedule := range schedules {
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
			log.Printf("Error scheduling task %s: %v", schedule, err)
			continue
		}
		log.Printf("✓ Scheduled: %s", schedule)
	}

	s.cron.Start()
	log.Println("Scheduler started with multiple update times (like Java version)")
	log.Println("Running initial lottery update...")

	// Executar atualização inicial em background
	go s.loteriasUpdate.UpdateAll()
}

func (s *ScheduledConsumer) Stop() {
	s.cron.Stop()
	log.Println("Scheduler stopped")
}
