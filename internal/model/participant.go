package model

import (
	"encoding/json"
	"os"
	"time"
)

const Layout = "02-01-2006"

type Participant struct {
	Name         string          `json:"name"`
	Unavailable  map[string]bool `json:"-"`
	UnavailableS []string        `json:"unavailable"`
	TotalHours   int             `json:"-"`
	LastDutyDate time.Time       `json:"-"`
}

func LoadParticipants(filename string) ([]*Participant, error) {
	var raw []*Participant
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return nil, err
	}

	for _, p := range raw {
		uMap := make(map[string]bool)
		for _, date := range p.UnavailableS {
			uMap[date] = true
		}
		p.Unavailable = uMap
	}
	return raw, nil
}
