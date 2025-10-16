package scheduler

import (
	"fmt"
	"math/rand"
	"os"
	"sort"
	"time"

	"oncall-duty/config"
	"oncall-duty/internal/model"
	"oncall-duty/internal/util"
)

func GenerateSchedule(participants []*model.Participant, startDate, endDate time.Time, cfg config.Config, debug bool) []model.Duty {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(participants), func(i, j int) {
		participants[i], participants[j] = participants[j], participants[i]
	})

	var duties []model.Duty
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		hours := 15
		if d.Weekday() == time.Saturday || d.Weekday() == time.Sunday {
			hours = 24
		}
		duties = append(duties, model.Duty{Date: d, Hours: hours})
	}

	for i := range duties {
		available := filterAvailable(participants, duties[i], cfg.DescansoDias)
		if len(available) == 0 {
			panic(fmt.Sprintf("Nenhum participante disponível para %s", duties[i].Date.Format(model.Layout)))
		}

		sort.Slice(available, func(i, j int) bool {
			if available[i].TotalHours == available[j].TotalHours {
				return available[i].LastDutyDate.Before(available[j].LastDutyDate)
			}
			return available[i].TotalHours < available[j].TotalHours
		})

		var chosen *model.Participant
		for _, candidate := range available {
			minHours, maxHours := util.GetMinMaxHours(participants, candidate, duties[i].Hours)
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
			fmt.Printf("[DEBUG] Escolhido: %s para %s (%dh), Total: %dh\n", chosen.Name, duties[i].Date.Format(model.Layout), duties[i].Hours, chosen.TotalHours)
		}
	}

	return duties
}

func filterAvailable(participants []*model.Participant, duty model.Duty, descansoDias int) []*model.Participant {
	var result []*model.Participant
	for _, p := range participants {
		if p.Unavailable[duty.Date.Format(model.Layout)] {
			continue
		}
		if !p.LastDutyDate.IsZero() && p.LastDutyDate.AddDate(0, 0, descansoDias).Equal(duty.Date) {
			continue
		}
		result = append(result, p)
	}
	return result
}

func WriteScheduleFile(duties []model.Duty, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, duty := range duties {
		line := fmt.Sprintf("%s - %s - %s (%dh) - %s|%s\n",
			duty.Date.Format(model.Layout),
			duty.Date.AddDate(0, 0, 1).Format(model.Layout),
			duty.Who,
			duty.Hours,
			util.GetWeekdayPt(duty.Date.Weekday()),
			util.GetWeekdayPt(duty.Date.AddDate(0, 0, 1).Weekday()))
		_, err := file.WriteString(line)
		if err != nil {
			return err
		}
	}
	return nil
}

// Help with tests
func TestFilterAvailable(participants []*model.Participant, duty model.Duty, descansoDias int) []*model.Participant {
	return filterAvailable(participants, duty, descansoDias)
}

// Help with tests
func TestGetMinMaxHours(participants []*model.Participant, candidate *model.Participant, addHours int) (int, int) {
	return util.GetMinMaxHours(participants, candidate, addHours)
}
