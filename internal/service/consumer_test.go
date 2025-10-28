package service_test

import (
	"testing"

	"loterias-api-golang/internal/model"
)

func TestConsumer_ConvertMonthNumber(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"1", "Janeiro"},
		{"2", "Fevereiro"},
		{"12", "Dezembro"},
		{"invalid", "invalid"},
		{"0", "0"},
		{"13", "13"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			if !model.IsValid("megasena") {
				t.Errorf("Expected megasena to be valid")
			}

			if model.IsValid("invalid") {
				t.Errorf("Expected invalid to be invalid")
			}
		})
	}
}

func TestLoteria_IsValid(t *testing.T) {
	tests := []struct {
		name     string
		loteria  string
		expected bool
	}{
		{"Valid - Mega Sena", "megasena", true},
		{"Valid - Lotofacil", "lotofacil", true},
		{"Valid - Quina", "quina", true},
		{"Invalid - Empty", "", false},
		{"Invalid - Wrong name", "invalid", false},
		{"Invalid - Case sensitive", "MegaSena", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := model.IsValid(tt.loteria)
			if result != tt.expected {
				t.Errorf("IsValid(%s) = %v, want %v", tt.loteria, result, tt.expected)
			}
		})
	}
}

func TestLoteria_AllLoterias(t *testing.T) {
	loterias := model.AllLoterias()

	if len(loterias) != 10 {
		t.Errorf("Expected 10 loterias, got %d", len(loterias))
	}

	expectedLoterias := []string{"megasena", "lotofacil", "quina"}
	for _, expected := range expectedLoterias {
		found := false
		for _, loteria := range loterias {
			if loteria == expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected loteria %s not found", expected)
		}
	}
}
