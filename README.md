# oncall-duty
Serviço responsável pela geração automatizada de escalas de plantão e sobreaviso

## Participantes File

O arquivo que contém os participantes deve seguir esse formato:

```
[
  {
    "name": "Nome do Participante",
    "unavailable": ["01-01-2025", "02-01-2025"] // dias indisponíveis
  },
  {...}
]
```
