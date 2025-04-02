# Use a imagem oficial do Go 1.24.2
FROM golang:1.24.2-alpine

# Instalar dependências necessárias
RUN apk add --no-cache gcc musl-dev

# Definir o diretório de trabalho
WORKDIR /app

# Copiar os arquivos do backend
COPY backend/ ./backend/

# Copiar os arquivos do frontend
COPY frontend/ ./frontend/

# Instalar dependências do Go
WORKDIR /app/backend
RUN go mod download

# Expor a porta 8080
EXPOSE 8080

# Comando para executar a aplicação
CMD ["go", "run", "main.go"] 