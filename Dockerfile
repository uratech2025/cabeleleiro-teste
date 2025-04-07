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

# Compilar a aplicação
RUN go build -o main .

# Expor a porta 8080
EXPOSE 8080

# Definir variável de ambiente para o host
ENV HOST=0.0.0.0

# Comando para executar a aplicação
CMD ["./main"] 