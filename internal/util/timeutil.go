package util

import (
	"time"

	"oncall-duty/internal/model"
)

func GetMinMaxHours(participants []*model.Participant, candidate *model.Participant, addHours int) (int, int) {
	min, max := 999999, 0
	for _, p := range participants {
		h := p.TotalHours
		if p == candidate {
			h += addHours
		}
		if h < min {
			min = h
		}
		if h > max {
			max = h
		}
	}
	return min, max
}

func GetWeekdayPt(wd time.Weekday) string {
	switch wd {
	case time.Sunday:
		return "domingo"
	case time.Monday:
		return "segunda"
	case time.Tuesday:
		return "terça"
	case time.Wednesday:
		return "quarta"
	case time.Thursday:
		return "quinta"
	case time.Friday:
		return "sexta"
	case time.Saturday:
		return "sábado"
	}
	return ""
}
