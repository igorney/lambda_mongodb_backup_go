# Utiliza uma imagem base oficial do Go como uma etapa de build
FROM golang:1.21.0-bullseye AS builder

# Defina o diretório de trabalho dentro do contêiner
WORKDIR /app

# Copia os arquivos go.mod e go.sum para o diretório de trabalho
COPY go.mod go.sum ./

# Baixa as dependências necessárias
RUN go mod download

# Copia o código-fonte da aplicação para o diretório de trabalho
COPY . .

# Compila o binário executável para o Lambda
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o main .

# Utiliza uma imagem base oficial do Amazon Linux para Lambda
FROM public.ecr.aws/lambda/provided:al2

# Copia o binário executável da etapa de build para o diretório de trabalho
COPY --from=builder /app/main /var/task/main

# Define o comando de entrada padrão para o contêiner Lambda
CMD ["main"]