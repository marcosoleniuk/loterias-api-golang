package service

import (
	"log"
	"sync"

	"loterias-api-golang/internal/model"
)

const (
	batchSize      = 50 // Processar 50 concursos por vez
	maxConcurrency = 3  // Máximo 3 requisições simultâneas
)

type LoteriasUpdate struct {
	consumer         *Consumer
	resultadoService *ResultadoService
}

func NewLoteriasUpdate(consumer *Consumer, resultadoService *ResultadoService) *LoteriasUpdate {
	return &LoteriasUpdate{
		consumer:         consumer,
		resultadoService: resultadoService,
	}
}

func (l *LoteriasUpdate) UpdateAll() {
	log.Println("Starting lottery update...")

	var wg sync.WaitGroup
	loterias := model.AllLoterias()

	for _, loteria := range loterias {
		wg.Add(1)
		go func(lot string) {
			defer wg.Done()
			if err := l.updateLoteria(lot); err != nil {
				log.Printf("Error updating %s: %v", lot, err)
			}
		}(loteria)
	}

	wg.Wait()
	log.Println("Lottery update completed")
}

func (l *LoteriasUpdate) updateLoteria(loteria string) error {
	log.Printf("Updating %s...", loteria)

	latest, err := l.resultadoService.FindLatest(loteria)
	if err != nil {
		log.Printf("%s: Error finding latest in DB: %v", loteria, err)
		return err
	}

	latestAPI, err := l.consumer.GetLatestResultado(loteria)
	if err != nil {
		log.Printf("%s: Error fetching latest from API: %v", loteria, err)
		return err
	}

	log.Printf("%s: Latest in DB: %v, Latest in API: %d", loteria, latest, latestAPI.Concurso)

	startConcurso := 1
	if latest != nil && latest.Concurso > 0 {
		if latestAPI.Concurso <= latest.Concurso {
			log.Printf("%s: Already up to date (contest %d)", loteria, latest.Concurso)
			return nil
		}
		startConcurso = latest.Concurso + 1
	}

	totalConcursos := latestAPI.Concurso - startConcurso + 1
	log.Printf("%s: Fetching contests from %d to %d (%d total)", loteria, startConcurso, latestAPI.Concurso, totalConcursos)

	for batchStart := startConcurso; batchStart <= latestAPI.Concurso; batchStart += batchSize {
		batchEnd := batchStart + batchSize - 1
		if batchEnd > latestAPI.Concurso {
			batchEnd = latestAPI.Concurso
		}

		log.Printf("%s: Fetching batch %d-%d...", loteria, batchStart, batchEnd)
		resultados := l.fetchBatch(loteria, batchStart, batchEnd)
		log.Printf("%s: Fetched %d results from batch %d-%d", loteria, len(resultados), batchStart, batchEnd)

		if len(resultados) > 0 {
			if err := l.resultadoService.SaveAll(resultados); err != nil {
				log.Printf("%s: Error saving batch %d-%d: %v", loteria, batchStart, batchEnd, err)
				continue
			}
			log.Printf("%s: Saved batch %d-%d (%d contests)", loteria, batchStart, batchEnd, len(resultados))
		} else {
			log.Printf("%s: No results fetched for batch %d-%d", loteria, batchStart, batchEnd)
		}
	}

	log.Printf("%s: Update completed", loteria)
	return nil
}

func (l *LoteriasUpdate) fetchBatch(loteria string, start, end int) []model.Resultado {
	jobs := make(chan int, end-start+1)
	results := make(chan *model.Resultado, end-start+1)

	var wg sync.WaitGroup
	for w := 0; w < maxConcurrency; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for concurso := range jobs {
				resultado, err := l.consumer.GetResultado(loteria, concurso)
				if err != nil {
					log.Printf("Error fetching %s contest %d: %v", loteria, concurso, err)
					continue
				}
				results <- resultado
			}
		}()
	}

	for concurso := start; concurso <= end; concurso++ {
		jobs <- concurso
	}
	close(jobs)

	wg.Wait()
	close(results)

	var resultados []model.Resultado
	for resultado := range results {
		resultados = append(resultados, *resultado)
	}

	return resultados
}

func (l *LoteriasUpdate) UpdateOne(loteria string) error {
	return l.updateLoteria(loteria)
}
