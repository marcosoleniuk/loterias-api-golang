package service

import (
	"loterias-api-golang/internal/model"
	"loterias-api-golang/internal/repository"
)

type ResultadoService struct {
	repository *repository.ResultadoRepository
}

func NewResultadoService(repository *repository.ResultadoRepository) *ResultadoService {
	return &ResultadoService{
		repository: repository,
	}
}

func (s *ResultadoService) FindByLoteria(loteria string) ([]model.Resultado, error) {
	return s.repository.FindByLoteria(loteria)
}

func (s *ResultadoService) FindByLoteriaAndConcurso(loteria string, concurso int) (*model.Resultado, error) {
	return s.repository.FindByID(loteria, concurso)
}

func (s *ResultadoService) FindLatest(loteria string) (*model.Resultado, error) {
	return s.repository.FindLatest(loteria)
}

func (s *ResultadoService) Save(resultado *model.Resultado) error {
	return s.repository.Save(resultado)
}

func (s *ResultadoService) SaveAll(resultados []model.Resultado) error {
	return s.repository.SaveAll(resultados)
}
