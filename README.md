# Sistema de Cabeleireiro

Um sistema simples para gerenciamento de salão de beleza, desenvolvido com Go no backend e JavaScript/HTML/CSS no frontend.

## Funcionalidades

- Gerenciamento de serviços (preços e comissões)
- Cadastro de clientes
- Controle de comandas
- Fluxo de caixa (entradas e saídas)
- Controle de estoque
- Histórico de operações

## Requisitos

- Go 1.21 ou superior
- SQLite3
- Navegador web moderno

## Configuração

1. Clone o repositório:
```bash
git clone [URL_DO_REPOSITÓRIO]
cd cabeleleiro
```

2. Instale as dependências do backend:
```bash
cd backend
go mod download
```

3. Inicie o servidor backend:
```bash
go run main.go
```

4. Abra o arquivo `frontend/index.html` em seu navegador

## Estrutura do Projeto

```
cabeleleiro/
├── backend/
│   ├── main.go
│   └── salon.db
└── frontend/
    ├── index.html
    ├── styles.css
    └── script.js
```

## Uso

1. **Serviços**: Cadastre os serviços oferecidos pelo salão, definindo preço e comissão
2. **Clientes**: Cadastre os clientes do salão
3. **Comandas**: Crie comandas para os clientes, selecionando o serviço desejado
4. **Fluxo de Caixa**: Registre entradas e saídas de dinheiro
5. **Estoque**: Controle o estoque de produtos

## API Endpoints

- `GET /services`: Lista todos os serviços
- `POST /services`: Cria um novo serviço
- `GET /clients`: Lista todos os clientes
- `POST /clients`: Cria um novo cliente
- `GET /orders`: Lista todas as comandas
- `POST /orders`: Cria uma nova comanda
- `GET /cashflow`: Lista todas as movimentações financeiras
- `POST /cashflow`: Registra uma nova movimentação
- `GET /inventory`: Lista todos os itens do estoque
- `POST /inventory`: Atualiza o estoque

## Contribuição

Sinta-se à vontade para contribuir com o projeto através de pull requests ou reportando issues. 