package controller

import (
	"net/http"
	"strconv"

	"loterias-api-golang/internal/model"
	"loterias-api-golang/internal/service"

	"github.com/gin-gonic/gin"
)

type ApiController struct {
	resultadoService *service.ResultadoService
}

func NewApiController(resultadoService *service.ResultadoService) *ApiController {
	return &ApiController{
		resultadoService: resultadoService,
	}
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// GetLotteries retorna todas as loterias disponíveis
//
//	@Summary		Lista todas as loterias
//	@Description	Retorna todas as loterias disponíveis para pesquisa
//	@Tags			Loterias
//	@Produce		json
//	@Success		200	{array}	string
//	@Router			/ [get]
func (c *ApiController) GetLotteries(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, model.AllLoterias())
}

// GetResultsByLottery retorna todos os resultados de uma loteria
//
//	@Summary		Busca resultados por loteria
//	@Description	Retorna todos os resultados já realizados da loteria especificada
//	@Tags			Loterias
//	@Produce		json
//	@Param			loteria	path		string	true	"ID da Loteria"	Enums(maismilionaria, megasena, lotofacil, quina, lotomania, timemania, duplasena, federal, diadesorte, supersete)
//	@Success		200		{array}		model.Resultado
//	@Failure		404		{object}	ErrorResponse
//	@Router			/{loteria} [get]
func (c *ApiController) GetResultsByLottery(ctx *gin.Context) {
	loteria := ctx.Param("loteria")

	if !model.IsValid(loteria) {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Resource Not Found",
			Message: c.getInvalidLotteryMessage(loteria),
		})
		return
	}

	resultados, err := c.resultadoService.FindByLoteria(loteria)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resultados)
}

// GetResultByID retorna um resultado específico
//
//	@Summary		Busca resultado por loteria e concurso
//	@Description	Retorna o resultado da loteria e concurso especificado
//	@Tags			Loterias
//	@Produce		json
//	@Param			loteria		path		string	true	"ID da Loteria"	Enums(maismilionaria, megasena, lotofacil, quina, lotomania, timemania, duplasena, federal, diadesorte, supersete)
//	@Param			concurso	path		int		true	"Número do Concurso"
//	@Success		200			{object}	model.Resultado
//	@Failure		404			{object}	ErrorResponse
//	@Router			/{loteria}/{concurso} [get]
func (c *ApiController) GetResultByID(ctx *gin.Context) {
	loteria := ctx.Param("loteria")
	concursoStr := ctx.Param("concurso")

	if !model.IsValid(loteria) {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Resource Not Found",
			Message: c.getInvalidLotteryMessage(loteria),
		})
		return
	}

	concurso, err := strconv.Atoi(concursoStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Bad Request",
			Message: "Invalid contest number",
		})
		return
	}

	resultado, err := c.resultadoService.FindByLoteriaAndConcurso(loteria, concurso)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	if resultado == nil {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Resource Not Found",
			Message: "Result not found",
		})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}

// GetLatestResult retorna o resultado mais recente de uma loteria
//
//	@Summary		Busca resultado mais recente
//	@Description	Retorna o resultado mais recente da loteria especificada
//	@Tags			Loterias
//	@Produce		json
//	@Param			loteria	path		string	true	"ID da Loteria"	Enums(maismilionaria, megasena, lotofacil, quina, lotomania, timemania, duplasena, federal, diadesorte, supersete)
//	@Success		200		{object}	model.Resultado
//	@Failure		404		{object}	ErrorResponse
//	@Router			/{loteria}/latest [get]
func (c *ApiController) GetLatestResult(ctx *gin.Context) {
	loteria := ctx.Param("loteria")

	if !model.IsValid(loteria) {
		ctx.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Resource Not Found",
			Message: c.getInvalidLotteryMessage(loteria),
		})
		return
	}

	resultado, err := c.resultadoService.FindLatest(loteria)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Internal Server Error",
			Message: err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, resultado)
}

func (c *ApiController) getInvalidLotteryMessage(loteria string) string {
	loterias := model.AllLoterias()
	return "'" + loteria + "' não é o id de nenhuma das loterias suportadas. Loterias suportadas: " +
		"[" + join(loterias, ", ") + "]"
}

func join(arr []string, sep string) string {
	if len(arr) == 0 {
		return ""
	}
	result := arr[0]
	for i := 1; i < len(arr); i++ {
		result += sep + arr[i]
	}
	return result
}
