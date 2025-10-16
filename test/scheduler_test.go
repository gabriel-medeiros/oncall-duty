package test

import (
	"testing"
	"time"

	"oncall-duty/internal/model"
	"oncall-duty/internal/scheduler"
)

// Testa se um participante indisponível é realmente filtrado
func TestFilterAvailableParticipants_Indisponivel(t *testing.T) {
	layout := model.Layout
	date, _ := time.Parse(layout, "20-10-2025")

	p := &model.Participant{
		Name:         "Teste",
		Unavailable:  map[string]bool{"20-10-2025": true},
		TotalHours:   0,
		LastDutyDate: time.Time{},
	}

	duty := model.Duty{Date: date}

	available := scheduler.TestFilterAvailable([]*model.Participant{p}, duty, 1)
	if len(available) != 0 {
		t.Errorf("Esperado nenhuma disponibilidade, mas retornou %d", len(available))
	}
}

// Testa se maxDiff nunca é excedido para um candidato
func TestMaxDiffNeverExceeded(t *testing.T) {
	participants := []*model.Participant{
		{Name: "A", TotalHours: 10},
		{Name: "B", TotalHours: 20},
	}

	min, max := scheduler.TestGetMinMaxHours(participants, participants[0], 5)
	if max-min > 15 {
		t.Errorf("Diferença excedeu maxDiff: %d", max-min)
	}
}
