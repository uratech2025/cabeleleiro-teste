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

# Expor a porta padrão do Render
EXPOSE 10000

# Definir variável de ambiente para a porta
ENV PORT=10000

# Definir o diretório de trabalho para a raiz do projeto
WORKDIR /app

# Comando para executar a aplicação
CMD ["backend/main"] 