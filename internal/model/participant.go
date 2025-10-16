package model

import (
	"encoding/json"
	"os"
	"time"
)

var Layout = "02-01-2006"

// Estrutura para definir períodos indisponíveis
type UnavailableRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

// Agrupamento de todas as indisponibilidades
type Unavailability struct {
	SpecificDays []string           `json:"specific_days"`
	Ranges       []UnavailableRange `json:"ranges"`
	WeekDays     []time.Weekday     `json:"week_days"`
}

type Participant struct {
	Name           string         `json:"name"`
	Unavailability Unavailability `json:"unavailability"`
	TotalHours     int            `json:"total_hours"`
	LastDutyDate   time.Time      `json:"last_duty_date"`
}

func LoadParticipants(filename string) ([]*Participant, error) {
	var participants []*Participant

	// Lê todo o conteúdo do arquivo
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Deserializa direto para a nova struct
	err = json.Unmarshal(data, &participants)
	if err != nil {
		return nil, err
	}

	return participants, nil
}
