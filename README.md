# Oncall Duty
Serviço responsável pela geração automatizada de escalas de plantão e sobreaviso

## How to run

- Configure the participants.json file with the below format
- Run $ `go mod init oncall-duty`
- Run $ `go mod tidy`
- Run $ `go run ./cmd/scheduler --inicio <data-inicial> --fim <data-final> --participantes participants.json --config config.json --debug`

Obs: the data format is DD-MM-YYYY

## Participants file

O arquivo que contém os participantes deve seguir esse formato:

```
[
  {
    "name": "Nome do Participante",
    "unavailability": {
      "specific_days": ["28-10-2025"],
      "ranges": [
        { "start": "30-10-2025", "end": "02-11-2025" }
      ],
      "week_days": [0] // Domingos
    }
  },
  {...}
]
```
