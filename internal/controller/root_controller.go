package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type RootController struct{}

func NewRootController() *RootController {
	return &RootController{}
}

// Root retorna informações sobre a API
//
//	@Summary		Informações da API
//	@Description	Retorna informações básicas sobre a API de Loterias
//	@Tags			Root
//	@Produce		json
//	@Success		200	{object}	map[string]interface{}
//	@Router			/ [get]
func (c *RootController) Root(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"name":        "Loterias API Golang",
		"version":     "1.0.0",
		"description": "API para consulta de resultados de loterias da Caixa Econômica Federal",
		"swagger":     "/swagger/index.html",
		"endpoints": gin.H{
			"lotteries":  "/api",
			"by_lottery": "/api/{loteria}",
			"by_contest": "/api/{loteria}/{concurso}",
			"latest":     "/api/{loteria}/latest",
		},
	})
}
