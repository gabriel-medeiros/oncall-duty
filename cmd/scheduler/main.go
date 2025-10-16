package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"oncall-duty/config"
	"oncall-duty/internal/model"
	"oncall-duty/internal/scheduler"
)

const layout = "02-01-2006"

func main() {
	startStr := flag.String("inicio", "", "Data inicial no formato DD-MM-YYYY")
	endStr := flag.String("fim", "", "Data final no formato DD-MM-YYYY")
	jsonFile := flag.String("participants", "participants.json", "Arquivo JSON com participantes")
	configFile := flag.String("config", "config.json", "Arquivo JSON com configuração")
	debug := flag.Bool("debug", false, "Exibir logs detalhados")
	flag.Parse()

	if *startStr == "" || *endStr == "" {
		log.Fatal("Uso: go run main.go --inicio <data-inicial> --fim <data-final>")
	}

	startDate, err := time.Parse(layout, *startStr)
	if err != nil {
		log.Fatalf("Data inicial inválida: %v", err)
	}
	endDate, err := time.Parse(layout, *endStr)
	if err != nil {
		log.Fatalf("Data final inválida: %v", err)
	}

	cfg, err := config.Load(*configFile)
	if err != nil {
		log.Fatalf("Erro ao carregar configuração: %v", err)
	}

	participants, err := model.LoadParticipants(*jsonFile)
	if err != nil {
		log.Fatalf("Erro ao carregar participantes: %v", err)
	}

	duties := scheduler.GenerateSchedule(participants, startDate, endDate, cfg, *debug)

	err = scheduler.WriteScheduleFile(duties, "plantao.txt")
	if err != nil {
		log.Fatalf("Erro ao escrever escala: %v", err)
	}

	fmt.Println("\n📊 Total de Horas por Participante:")
	for _, p := range participants {
		fmt.Printf("%s: %dh\n", p.Name, p.TotalHours)
	}
}
