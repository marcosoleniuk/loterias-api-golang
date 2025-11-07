package service

import (
	"fmt"
	"log"
	"time"

	"loterias-api-golang/internal/model"
)

//const (
//	batchSize      = 10 // Processar 10 concursos por vez
//	maxConcurrency = 1  // Máximo 1 requisição simultânea por loteria
//)

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
		log.Printf("%s: ❌ Error finding latest in DB: %v", loteria, err)
		return err
	}

	// Buscar último concurso disponível na API (com retry)
	var latestAPI *model.Resultado
	var apiErr error
	for i := 0; i < 3; i++ {
		latestAPI, apiErr = l.consumer.GetLatestResultado(loteria)
		if apiErr == nil {
			break
		}
		log.Printf("%s: ⚠ Attempt %d to fetch latest from API failed: %v", loteria, i+1, apiErr)
		if i < 2 {
			time.Sleep(2 * time.Second)
		}
	}

	if apiErr != nil {
		log.Printf("%s: ❌ Error fetching latest from API after 3 attempts: %v", loteria, apiErr)
		return apiErr
	}

	if latestAPI == nil {
		log.Printf("%s: ❌ API returned nil result", loteria)
		return fmt.Errorf("API returned nil result for %s", loteria)
	}

	var latestDBConcurso int
	if latest != nil && latest.Concurso > 0 {
		latestDBConcurso = latest.Concurso
	}

	log.Printf("%s: 🔍 Latest in DB: %d | Latest in API: %d | Difference: %d", loteria, latestDBConcurso, latestAPI.Concurso, latestAPI.Concurso-latestDBConcurso)

	// Se o concurso é igual - atualizar apenas os dados (como em Java)
	if latestDBConcurso == latestAPI.Concurso {
		log.Printf("%s: 🔄 Same contest (%d), updating prize data...", loteria, latestDBConcurso)

		// Atualizar dados do concurso existente
		latest.Data = latestAPI.Data
		latest.Local = latestAPI.Local
		latest.Premiacoes = latestAPI.Premiacoes
		latest.LocalGanhadores = latestAPI.LocalGanhadores
		latest.Acumulou = latestAPI.Acumulou
		latest.DataProximoConcurso = latestAPI.DataProximoConcurso
		latest.ValorAcumuladoProximoConcurso = latestAPI.ValorAcumuladoProximoConcurso
		latest.ValorEstimadoProximoConcurso = latestAPI.ValorEstimadoProximoConcurso

		if err := l.resultadoService.Save(latest); err != nil {
			log.Printf("%s: ❌ Error updating contest %d: %v", loteria, latestDBConcurso, err)
			return err
		}
		log.Printf("%s: ✓ Contest %d data updated", loteria, latestDBConcurso)
		return nil
	}

	// Se já está atualizado (já tem novos concursos)
	if latestDBConcurso > latestAPI.Concurso {
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

	// Processar com retry (como em Java)
	retriesMap := make(map[int]int)
	for concurso := startConcurso; concurso <= latestAPI.Concurso; {
		resultado, err := l.consumer.GetResultado(loteria, concurso)
		if err != nil {
			retries := retriesMap[concurso]
			if retries < 20 {
				retries++
				retriesMap[concurso] = retries
				log.Printf("%s: ⚠ Error fetching contest %d (attempt %d/20): %v", loteria, concurso, retries, err)
				time.Sleep(2 * time.Second) // Aguardar antes de retry
				continue
			} else {
				log.Printf("%s: ❌ Stopped fetching from contest %d (max retries reached)", loteria, concurso)
				break
			}
		}

		if err := l.resultadoService.Save(resultado); err != nil {
			log.Printf("%s: ❌ Error saving contest %d: %v", loteria, concurso, err)
			// Não para, continua tentando outros
		} else {
			log.Printf("%s: ✓ Saved contest %d", loteria, concurso)
		}

		concurso++
	}

	log.Printf("%s: ========== Update completed ==========", loteria)
	return nil
}

func (l *LoteriasUpdate) UpdateOne(loteria string) error {
	return l.updateLoteria(loteria)
}
