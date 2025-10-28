package model

import _ "time"

type ResultadoID struct {
	Loteria  string `bson:"loteria" json:"loteria"`
	Concurso int    `bson:"concurso" json:"concurso"`
}

type Resultado struct {
	ID                             ResultadoID             `bson:"_id" json:"-"`
	Loteria                        string                  `bson:"-" json:"loteria"`
	Concurso                       int                     `bson:"-" json:"concurso"`
	Data                           string                  `bson:"data" json:"data"`
	Local                          string                  `bson:"local" json:"local,omitempty"`
	DezenasOrdemSorteio            []string                `bson:"dezenasOrdemSorteio,omitempty" json:"dezenasOrdemSorteio,omitempty"`
	Dezenas                        []string                `bson:"dezenas" json:"dezenas"`
	Trevos                         []string                `bson:"trevos,omitempty" json:"trevos,omitempty"`
	TimeCoracao                    string                  `bson:"timeCoracao,omitempty" json:"timeCoracao,omitempty"`
	MesSorte                       string                  `bson:"mesSorte,omitempty" json:"mesSorte,omitempty"`
	Premiacoes                     []Premiacao             `bson:"premiacoes" json:"premiacoes"`
	LocalGanhadores                []MunicipioUFGanhadores `bson:"localGanhadores,omitempty" json:"municipiosUFGanhadores,omitempty"`
	EstadosPremiados               []Estado                `bson:"estadosPremiados,omitempty" json:"estadosPremiados,omitempty"`
	Observacao                     string                  `bson:"observacao,omitempty" json:"observacao,omitempty"`
	Acumulou                       bool                    `bson:"acumulou" json:"acumulou"`
	ProximoConcurso                int                     `bson:"proximoConcurso,omitempty" json:"proximoConcurso,omitempty"`
	DataProximoConcurso            string                  `bson:"dataProximoConcurso,omitempty" json:"dataProximoConcurso,omitempty"`
	ValorArrecadado                float64                 `bson:"valorArrecadado,omitempty" json:"valorArrecadado,omitempty"`
	ValorAcumuladoConcurso_0_5     float64                 `bson:"valorAcumuladoConcurso_0_5,omitempty" json:"valorAcumuladoConcurso_0_5,omitempty"`
	ValorAcumuladoConcursoEspecial float64                 `bson:"valorAcumuladoConcursoEspecial,omitempty" json:"valorAcumuladoConcursoEspecial,omitempty"`
	ValorAcumuladoProximoConcurso  float64                 `bson:"valorAcumuladoProximoConcurso,omitempty" json:"valorAcumuladoProximoConcurso,omitempty"`
	ValorEstimadoProximoConcurso   float64                 `bson:"valorEstimadoProximoConcurso,omitempty" json:"valorEstimadoProximoConcurso,omitempty"`
}

type Premiacao struct {
	Descricao          string  `bson:"descricao" json:"descricao"`
	Faixa              int     `bson:"faixa" json:"faixa"`
	NumeroDeGanhadores int     `bson:"numeroDeGanhadores" json:"numeroDeGanhadores"`
	Valor              float64 `bson:"valor" json:"valor"`
}

type MunicipioUFGanhadores struct {
	Ganhadores    int    `bson:"ganhadores" json:"ganhadores"`
	Municipio     string `bson:"municipio" json:"municipio"`
	Posicao       int    `bson:"posicao" json:"posicao"`
	UF            string `bson:"uf" json:"uf"`
	Serie         string `bson:"serie,omitempty" json:"serie,omitempty"`
	NumeroBilhete string `bson:"numeroBilhete,omitempty" json:"numeroBilhete,omitempty"`
}

type Estado struct {
	Nome       string `bson:"nome" json:"nome"`
	UF         string `bson:"uf" json:"uf"`
	Ganhadores int    `bson:"ganhadores" json:"ganhadores"`
}

func (r *Resultado) BeforeSave() {
	r.Loteria = r.ID.Loteria
	r.Concurso = r.ID.Concurso
}

func (r *Resultado) AfterFind() {
	r.Loteria = r.ID.Loteria
	r.Concurso = r.ID.Concurso
}
