package scheduler

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"oncall-duty/config"
)

func GenerateSchedule(participants []*Participant, startDate, endDate time.Time, cfg config.Config, debug bool) []Duty {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})

	var duties []Duty
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		hours := 15
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			hours = 24
		}
		duties = append(duties, Duty{Date: d, Hours: hours})
	}

	for i := range duties {
		available := filterAvailable(participants, duties[i], cfg.DescansoDias)
		if len(available) == 0 {
			panic(fmt.Sprintf("Nenhum participante disponível para %s", duties[i].Date.Format(Layout)))
		}

		sort.Slice(available, func(i, j int) bool {
			if available[i].TotalHours == available[j].TotalHours {
				return available[i].LastDutyDate.Before(available[j].LastDutyDate)
			}
			return available[i].TotalHours < available[j].TotalHours
		})

		var chosen *Participant
		for _, candidate := range available {
			minHours, maxHours := GetMinMaxHours(participants, candidate, duties[i].Hours)
			if maxHours-minHours <= cfg.MaxDiff {
				chosen = candidate
				break
			}
		}

		if chosen == nil {
			chosen = available[0]
		}

		chosen.TotalHours += duties[i].Hours
		chosen.LastDutyDate = duties[i].Date
		duties[i].Who = chosen.Name

		if debug {
			fmt.Printf("[DEBUG] Escolhido: %s para %s (%dh), Total: %dh\n", chosen.Name, duties[i].Date.Format(Layout), duties[i].Hours, chosen.TotalHours)
		}
	}

	return duties
}

func filterAvailable(participants []*Participant, duty Duty, descansoDias int) []*Participant {
	var available []*Participant
	layout := Layout
	dutyDate := duty.Date

	for _, p := range participants {
		u := p.Unavailability

		// 🚫 Dias fixos da semana
		for _, day := range u.WeekDays {
			if day == dutyDate.Weekday() {
				goto nextParticipant
			}
		}

		// 🚫 Datas específicas
		for _, s := range u.SpecificDays {
			if s == dutyDate.Format(layout) {
				goto nextParticipant
			}
		}

		// 🚫 Períodos completos
		for _, r := range u.Ranges {
			start, _ := time.Parse(layout, r.Start)
			end, _ := time.Parse(layout, r.End)
			if !dutyDate.Before(start) && !dutyDate.After(end) {
				goto nextParticipant
			}
		}

		// ⏳ Regra de descanso
		if !p.LastDutyDate.IsZero() && dutyDate.Sub(p.LastDutyDate).Hours() < float64(descansoDias*24) {
			goto nextParticipant
		}

		available = append(available, p)

	nextParticipant:
	}
	return available
}

func WriteScheduleFile(duties []Duty, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, duty := range duties {
		line := fmt.Sprintf("%s - %s - %s (%dh) - %s|%s\n",
			duty.Date.Format(Layout),
			duty.Date.AddDate(0, 0, 1).Format(Layout),
			duty.Who,
			duty.Hours,
			GetWeekdayPt(duty.Date.Weekday()),
			GetWeekdayPt(duty.Date.AddDate(0, 0, 1).Weekday()))
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetMinMaxHours(participants []*Participant, candidate *Participant, addHours int) (int, int) {
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

// Help with tests
func TestFilterAvailable(participants []*Participant, duty Duty, descansoDias int) []*Participant {
	return filterAvailable(participants, duty, descansoDias)
}

// Help with tests
func TestGetMinMaxHours(participants []*Participant, candidate *Participant, addHours int) (int, int) {
	return GetMinMaxHours(participants, candidate, addHours)
}
