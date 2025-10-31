package service

import (
	"log"
	"sync"
	"time"

	"loterias-api-golang/internal/model"
)

const (
	batchSize      = 10 // Processar 10 concursos por vez
	maxConcurrency = 1  // Máximo 1 requisição simultânea por loteria
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

	loterias := model.AllLoterias()

	// Processar sequencialmente com delay para evitar bloqueio da API
	for i, loteria := range loterias {
		if i > 0 {
			// Aguardar 3 segundos entre cada loteria
			log.Printf("Waiting 3 seconds before updating next lottery...")
			time.Sleep(3 * time.Second)
		}

		if err := l.updateLoteria(loteria); err != nil {
			log.Printf("Error updating %s: %v", loteria, err)
		}
	}

	log.Println("Lottery update completed")
}

func (l *LoteriasUpdate) updateLoteria(loteria string) error {
	log.Printf("========== Updating %s ==========", loteria)

	// Buscar último concurso no banco de dados
	latest, err := l.resultadoService.FindLatest(loteria)
	if err != nil {
		log.Printf("%s: Error finding latest in DB: %v", loteria, err)
		return err
	}

	// Buscar último concurso disponível na API
	latestAPI, err := l.consumer.GetLatestResultado(loteria)
	if err != nil {
		log.Printf("%s: Error fetching latest from API: %v", loteria, err)
		return err
	}

	var latestDBConcurso int
	if latest != nil && latest.Concurso > 0 {
		latestDBConcurso = latest.Concurso
	}

	log.Printf("%s: Latest in DB: %d, Latest in API: %d", loteria, latestDBConcurso, latestAPI.Concurso)

	// Se já está atualizado, não fazer nada
	if latestDBConcurso >= latestAPI.Concurso {
		log.Printf("%s: ✓ Already up to date (contest %d)", loteria, latestDBConcurso)
		return nil
	}

	// Determinar de qual concurso começar
	startConcurso := latestDBConcurso + 1
	if latestDBConcurso == 0 {
		startConcurso = 1
	}

	totalConcursos := latestAPI.Concurso - startConcurso + 1
	log.Printf("%s: 📥 Fetching contests from %d to %d (%d new contests)", loteria, startConcurso, latestAPI.Concurso, totalConcursos)

	// Processar em lotes
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
				log.Printf("%s: ❌ Error saving batch %d-%d: %v", loteria, batchStart, batchEnd, err)
				continue
			}
			log.Printf("%s: ✓ Saved batch %d-%d (%d contests)", loteria, batchStart, batchEnd, len(resultados))
		} else {
			log.Printf("%s: ⚠ No results fetched for batch %d-%d", loteria, batchStart, batchEnd)
		}
	}

	log.Printf("%s: ========== Update completed ==========", loteria)
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
