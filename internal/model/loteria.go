package model

type Loteria string

const (
	MaisMilionaria Loteria = "maismilionaria"
	MegaSena       Loteria = "megasena"
	Lotofacil      Loteria = "lotofacil"
	Quina          Loteria = "quina"
	Lotomania      Loteria = "lotomania"
	Timemania      Loteria = "timemania"
	DuplaSena      Loteria = "duplasena"
	Federal        Loteria = "federal"
	DiaDeSorte     Loteria = "diadesorte"
	SuperSete      Loteria = "supersete"
)

func AllLoterias() []string {
	return []string{
		string(MaisMilionaria),
		string(MegaSena),
		string(Lotofacil),
		string(Quina),
		string(Lotomania),
		string(Timemania),
		string(DuplaSena),
		string(Federal),
		string(DiaDeSorte),
		string(SuperSete),
	}
}

func IsValid(loteria string) bool {
	for _, l := range AllLoterias() {
		if l == loteria {
			return true
		}
	}
	return false
}
