# Utiliza uma imagem base oficial do Go como uma etapa de build
FROM golang:1.20 AS builder

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos go.mod e go.sum para o diretório de trabalho
COPY go.mod go.sum ./

# Baixa as dependências necessárias
RUN go mod download

# Copia o código-fonte da aplicação para o diretório de trabalho
COPY . .

# Compila o binário executável
RUN go build -o main .

# Utiliza uma imagem base oficial do Debian como a etapa final
FROM debian:bullseye-slim

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Instala as ferramentas de banco de dados MongoDB
RUN apt-get update && apt-get install -y wget && \
    wget https://fastdl.mongodb.org/tools/db/mongodb-database-tools-ubuntu2004-x86_64-100.7.4.deb && \
    apt-get install -y ./mongodb-database-tools-ubuntu2004-x86_64-100.7.4.deb && \
    rm ./mongodb-database-tools-ubuntu2004-x86_64-100.7.4.deb && \
    apt-get clean && rm -rf /var/lib/apt/lists/*

# Copia o binário executável da etapa de build para o diretório de trabalho
COPY --from=builder /app/main .

# Define o comando de entrada padrão para o contêiner
CMD ["./main"]
