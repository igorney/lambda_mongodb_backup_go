# Utilizando a imagem oficial do Golang para construir o aplicativo
FROM golang:1.18 AS builder

# Defina o diretório de trabalho
WORKDIR /app

# Copie os arquivos de módulo Go e baixe as dependências
COPY go.mod go.sum ./
RUN go mod download

# Copie o código-fonte para o contêiner
COPY . .

# Compile o binário Go
RUN GOOS=linux GOARCH=amd64 go build -o main main.go

# Utilizando uma imagem mínima do Ubuntu para a execução do binário
FROM ubuntu:20.04

# Instale dependências
RUN apt-get update && apt-get install -y ca-certificates mongodb-clients && rm -rf /var/lib/apt/lists/*

# Defina o diretório de trabalho
WORKDIR /app

# Copie o binário compilado do estágio anterior
COPY --from=builder /app/main .

# Crie a pasta de backups
RUN mkdir -p /app/backups

# Defina o ponto de entrada para o contêiner
ENTRYPOINT ["./main"]
