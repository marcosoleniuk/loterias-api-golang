package service

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	_ "strings"
	"time"

	"loterias-api-golang/internal/model"
)

type Consumer struct {
	client       *http.Client
	requestDelay time.Duration
	maxRetries   int
}

func NewConsumer() *Consumer {
	tr := &http.Transport{
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   30 * time.Second,
	}

	return &Consumer{
		client:       client,
		requestDelay: 2000 * time.Millisecond, // 2 segundos entre requisições
		maxRetries:   5,                       // Máximo 5 tentativas
	}
}

type CaixaResponse struct {
	Numero                         int                      `json:"numero"`
	DataApuracao                   string                   `json:"dataApuracao"`
	LocalSorteio                   string                   `json:"localSorteio"`
	NomeMunicipioUFSorteio         string                   `json:"nomeMunicipioUFSorteio"`
	DezenasSorteadasOrdemSorteio   []string                 `json:"dezenasSorteadasOrdemSorteio"`
	ListaDezenas                   []string                 `json:"listaDezenas"`
	ListaDezenasSegundoSorteio     []string                 `json:"listaDezenasSegundoSorteio"`
	TrevosSorteados                []string                 `json:"trevosSorteados"`
	NomeTimeCoracaoMesSorte        string                   `json:"nomeTimeCoracaoMesSorte"`
	ListaRateioPremio              []CaixaPremiacao         `json:"listaRateioPremio"`
	ListaMunicipioUFGanhadores     []CaixaMunicipioGanhador `json:"listaMunicipioUFGanhadores"`
	Observacao                     string                   `json:"observacao"`
	Acumulado                      bool                     `json:"acumulado"`
	DataProximoConcurso            string                   `json:"dataProximoConcurso"`
	ValorArrecadado                float64                  `json:"valorArrecadado"`
	ValorAcumuladoConcurso_0_5     float64                  `json:"valorAcumuladoConcurso_0_5"`
	ValorAcumuladoConcursoEspecial float64                  `json:"valorAcumuladoConcursoEspecial"`
	ValorAcumuladoProximoConcurso  float64                  `json:"valorAcumuladoProximoConcurso"`
	ValorEstimadoProximoConcurso   float64                  `json:"valorEstimadoProximoConcurso"`
	NumeroProximoConcurso          int                      `json:"numeroConcursoProximo"`
}

type CaixaPremiacao struct {
	DescricaoFaixa     string  `json:"descricaoFaixa"`
	Faixa              int     `json:"faixa"`
	NumeroDeGanhadores int     `json:"numeroDeGanhadores"`
	ValorPremio        float64 `json:"valorPremio"`
}

type CaixaMunicipioGanhador struct {
	Ganhadores     int    `json:"ganhadores"`
	Municipio      string `json:"municipio"`
	NomeFatansiaUL string `json:"nomeFatansiaUL"`
	Posicao        int    `json:"posicao"`
	Serie          string `json:"serie"`
	UF             string `json:"uf"`
}

func (c *Consumer) GetResultado(loteria string, concurso int) (*model.Resultado, error) {
	return c.getResultadoFromServiceBus(loteria, strconv.Itoa(concurso))
}

func (c *Consumer) GetLatestResultado(loteria string) (*model.Resultado, error) {
	return c.getResultadoFromServiceBus(loteria, "")
}

func (c *Consumer) getResultadoFromServiceBus(loteria, concurso string) (*model.Resultado, error) {
	baseURL := "https://servicebus2.caixa.gov.br/portaldeloterias/api/"
	url := fmt.Sprintf("%s%s/%s", baseURL, loteria, concurso)

	var lastErr error
	for attempt := 1; attempt <= c.maxRetries; attempt++ {
		if attempt > 1 {
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			log.Printf("Retry attempt %d for %s (backoff: %v)", attempt, url, backoff)
			time.Sleep(backoff)
		} else {
			time.Sleep(c.requestDelay)
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			lastErr = fmt.Errorf("failed to create request: %w", err)
			continue
		}

		// Headers mais completos para evitar bloqueio
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36")
		req.Header.Set("Accept", "application/json, text/plain, */*")
		req.Header.Set("Accept-Language", "pt-BR,pt;q=0.9,en-US;q=0.8,en;q=0.7")
		req.Header.Set("Accept-Encoding", "gzip, deflate, br")
		req.Header.Set("Connection", "keep-alive")
		req.Header.Set("Referer", "https://loterias.caixa.gov.br/")
		req.Header.Set("Origin", "https://loterias.caixa.gov.br")
		req.Header.Set("Sec-Fetch-Dest", "empty")
		req.Header.Set("Sec-Fetch-Mode", "cors")
		req.Header.Set("Sec-Fetch-Site", "same-site")
		req.Header.Set("Cache-Control", "no-cache")
		req.Header.Set("Pragma", "no-cache")

		resp, err := c.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("failed to fetch data: %w", err)
			log.Printf("HTTP error for %s: %v", url, err)
			continue
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()

		if err != nil {
			lastErr = fmt.Errorf("failed to read response body: %w", err)
			log.Printf("Read body error for %s: %v", url, err)
			continue
		}

		if resp.StatusCode == http.StatusTooManyRequests {
			waitTime := time.Duration(5+attempt*2) * time.Second
			lastErr = fmt.Errorf("rate limited (429), waiting %v before retry", waitTime)
			log.Printf("⚠ Rate limited (429) for %s, waiting %v", url, waitTime)
			time.Sleep(waitTime)
			continue
		}

		if resp.StatusCode == http.StatusForbidden {
			waitTime := time.Duration(5+attempt*3) * time.Second
			lastErr = fmt.Errorf("forbidden (403), waiting %v before retry", waitTime)
			log.Printf("⚠ Forbidden (403) for %s, waiting %v before retry", url, waitTime)
			time.Sleep(waitTime)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("unexpected status code: %d", resp.StatusCode)
			log.Printf("❌ Unexpected status %d for %s", resp.StatusCode, url)
			return nil, lastErr
		}

		var caixaResp CaixaResponse
		if err := json.Unmarshal(body, &caixaResp); err != nil {
			lastErr = fmt.Errorf("failed to unmarshal response: %w", err)
			log.Printf("Unmarshal error for %s: %v. Body: %s", url, err, string(body[:minimalV(200, len(body))]))
			return nil, lastErr
		}

		return c.convertToResultado(loteria, &caixaResp), nil
	}

	return nil, fmt.Errorf("max retries exceeded: %w", lastErr)
}

func minimalV(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (c *Consumer) convertToResultado(loteria string, resp *CaixaResponse) *model.Resultado {
	resultado := &model.Resultado{
		ID: model.ResultadoID{
			Loteria:  loteria,
			Concurso: resp.Numero,
		},
		Data:                           resp.DataApuracao,
		Local:                          resp.LocalSorteio + " em " + resp.NomeMunicipioUFSorteio,
		DezenasOrdemSorteio:            resp.DezenasSorteadasOrdemSorteio,
		Dezenas:                        c.processDezenas(loteria, resp),
		Trevos:                         resp.TrevosSorteados,
		Observacao:                     resp.Observacao,
		Acumulou:                       resp.Acumulado,
		ProximoConcurso:                resp.NumeroProximoConcurso,
		DataProximoConcurso:            resp.DataProximoConcurso,
		ValorArrecadado:                resp.ValorArrecadado,
		ValorAcumuladoConcurso_0_5:     resp.ValorAcumuladoConcurso_0_5,
		ValorAcumuladoConcursoEspecial: resp.ValorAcumuladoConcursoEspecial,
		ValorAcumuladoProximoConcurso:  resp.ValorAcumuladoProximoConcurso,
		ValorEstimadoProximoConcurso:   resp.ValorEstimadoProximoConcurso,
	}

	if resp.NomeTimeCoracaoMesSorte != "" {
		if loteria == string(model.DiaDeSorte) {
			resultado.MesSorte = c.convertMonthNumber(resp.NomeTimeCoracaoMesSorte)
		} else if loteria == string(model.Timemania) {
			resultado.TimeCoracao = resp.NomeTimeCoracaoMesSorte
		}
	}

	for _, p := range resp.ListaRateioPremio {
		resultado.Premiacoes = append(resultado.Premiacoes, model.Premiacao{
			Descricao:          p.DescricaoFaixa,
			Faixa:              p.Faixa,
			NumeroDeGanhadores: p.NumeroDeGanhadores,
			Valor:              p.ValorPremio,
		})
	}

	for _, mg := range resp.ListaMunicipioUFGanhadores {
		resultado.LocalGanhadores = append(resultado.LocalGanhadores, model.MunicipioUFGanhadores{
			Ganhadores: mg.Ganhadores,
			Municipio:  mg.Municipio,
			Posicao:    mg.Posicao,
			UF:         mg.UF,
			Serie:      mg.Serie,
		})
	}

	resultado.AfterFind()
	return resultado
}

func (c *Consumer) processDezenas(loteria string, resp *CaixaResponse) []string {
	dezenas := make([]string, len(resp.ListaDezenas))
	copy(dezenas, resp.ListaDezenas)

	if len(resp.ListaDezenasSegundoSorteio) > 0 {
		dezenas = append(dezenas, resp.ListaDezenasSegundoSorteio...)
	}

	if loteria == string(model.DuplaSena) && len(dezenas) == 12 {
		primeiro := make([]string, 6)
		segundo := make([]string, 6)
		copy(primeiro, dezenas[:6])
		copy(segundo, dezenas[6:])

		sort.Strings(primeiro)
		sort.Strings(segundo)

		dezenas = append(primeiro, segundo...)
	} else if loteria != string(model.SuperSete) && loteria != string(model.Federal) {
		sort.Strings(dezenas)
	}

	return dezenas
}

func (c *Consumer) convertMonthNumber(monthStr string) string {
	meses := []string{
		"Janeiro", "Fevereiro", "Março", "Abril", "Maio", "Junho",
		"Julho", "Agosto", "Setembro", "Outubro", "Novembro", "Dezembro",
	}

	monthNum, err := strconv.Atoi(monthStr)
	if err != nil || monthNum < 1 || monthNum > 12 {
		return monthStr
	}

	return meses[monthNum-1]
}
