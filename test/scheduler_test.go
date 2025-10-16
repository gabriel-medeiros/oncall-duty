package test

import (
	"testing"
	"time"

	"oncall-duty/internal/model"
	"oncall-duty/internal/scheduler"
)

// Helper to create dates
func dateTimeParse(layout, value string) time.Time {
	t, err := time.Parse(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}

func TestFilterAvailable_SpecificDay(t *testing.T) {
	layout := model.Layout
	date := dateTimeParse(layout, "20-01-2025")

	p := &model.Participant{
		Name: "Teste",
		Unavailability: model.Unavailability{
			SpecificDays: []string{"20-01-2025"},
		},
	}

	duty := model.Duty{Date: date}
	available := scheduler.TestFilterAvailable([]*model.Participant{p}, duty, 1)

	if len(available) != 0 {
		t.Errorf("Esperado indisponível para data específica, mas retornou disponível")
	}
}

func TestFilterAvailable_Range(t *testing.T) {
	layout := model.Layout
	date := dateTimeParse(layout, "05-02-2025")

	p := &model.Participant{
		Name: "Teste",
		Unavailability: model.Unavailability{
			Ranges: []model.UnavailableRange{{Start: "01-02-2025", End: "10-02-2025"}},
		},
	}

	duty := model.Duty{Date: date}
	available := scheduler.TestFilterAvailable([]*model.Participant{p}, duty, 1)

	if len(available) != 0 {
		t.Errorf("Esperado indisponível no range de datas, mas retornou disponível")
	}
}

func TestFilterAvailable_WeekDay(t *testing.T) {
	// 3 == Wednesday
	date := dateTimeParse(model.Layout, "05-02-2025") // Essa data é quarta-feira

	p := &model.Participant{
		Name: "Teste",
		Unavailability: model.Unavailability{
			WeekDays: []time.Weekday{time.Wednesday},
		},
	}

	duty := model.Duty{Date: date}
	available := scheduler.TestFilterAvailable([]*model.Participant{p}, duty, 1)

	if len(available) != 0 {
		t.Errorf("Esperado indisponível no dia fixo da semana, mas retornou disponível")
	}
}

func TestFilterAvailable_Available(t *testing.T) {
	date := dateTimeParse(model.Layout, "10-03-2025")

	p := &model.Participant{
		Name: "Teste",
		Unavailability: model.Unavailability{
			SpecificDays: []string{},
			Ranges:       []model.UnavailableRange{},
			WeekDays:     []time.Weekday{},
		},
	}

	duty := model.Duty{Date: date}
	available := scheduler.TestFilterAvailable([]*model.Participant{p}, duty, 1)

	if len(available) != 1 {
		t.Errorf("Esperado disponibilidade, mas não retornou")
	}
}

// Testa se a regra de dias de descanso é aplicada (gerado pelo HubAI)
func TestFilterAvailable_DescansoDias(t *testing.T) {
	date := dateTimeParse(model.Layout, "15-03-2025")

	p := &model.Participant{
		Name:         "Teste",
		LastDutyDate: dateTimeParse(model.Layout, "14-03-2025"), // Menos de 2 dias de descanso
		Unavailability: model.Unavailability{
			SpecificDays: []string{},
			Ranges:       []model.UnavailableRange{},
			WeekDays:     []time.Weekday{},
		},
	}

	duty := model.Duty{Date: date}
	available := scheduler.TestFilterAvailable([]*model.Participant{p}, duty, 2)

	if len(available) != 0 {
		t.Errorf("Esperado indisponível pelo descanso, mas retornou disponível")
	}
}

// Testa se maxDiff nunca é excedido para um candidato (gerado pelo HubAI)
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
