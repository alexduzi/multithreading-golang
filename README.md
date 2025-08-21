# Comparador de APIs de CEP com Multithreading

Este projeto demonstra o uso de multithreading em Go para comparar a velocidade entre duas APIs de CEP diferentes. O programa faz requisições simultâneas e retorna o resultado da API que responder mais rapidamente.

## Como funciona

O programa lança duas goroutines simultaneamente, uma para cada API:
- ViaCEP: `http://viacep.com.br/ws/{cep}/json/`
- BrasilAPI: `https://brasilapi.com.br/api/cep/v1/{cep}`

A primeira resposta que chegar é processada e exibida. Se nenhuma API responder em 1 segundo, o programa exibe timeout.

## Como executar

```bash
go run main.go 01153000
```

## Exemplo de saída

```
2025/08/20 15:30:45 chamada https://brasilapi.com.br/api/cep/v1/01153000 finalizada com sucesso em 180ms!
2025/08/20 15:30:45 resposta mais rápida veio da: BRAZILAPI com tempo de 180ms
2025/08/20 15:30:45 url: https://brasilapi.com.br/api/cep/v1/01153000

2025/08/20 15:30:45 dados do endereço BRAZILAPI

2025/08/20 15:30:45 CEP: 01153-000
2025/08/20 15:30:45 Rua: Rua Visconde de Parnaíba
2025/08/20 15:30:45 Bairro: Barra Funda
2025/08/20 15:30:45 Cidade: São Paulo
2025/08/20 15:30:45 Estado: SP
```

## Exemplos de CEPs para testar

```bash
go run main.go 01153000    # São Paulo
go run main.go 20040020    # Rio de Janeiro  
go run main.go 30112000    # Belo Horizonte
```

## Estrutura do projeto

```
.
├── main.go           # Ponto de entrada da aplicação
├── model/
│   └── model.go      # Estruturas de dados
├── service/
│   └── service.go    # Lógica das requisições HTTP
└── go.mod           # Módulo Go
```

## Conceitos demonstrados

- Goroutines para execução concorrente
- Channels para comunicação entre goroutines  
- Select para capturar a primeira resposta
- Context com timeout
- Requisições HTTP com cancelamento