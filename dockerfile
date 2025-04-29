# Etapa 1: Build da aplicação
FROM golang:1.24-alpine AS builder

# Variáveis de ambiente para build
ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Diretório de trabalho
WORKDIR /app

# Copia os arquivos do projeto
COPY . .

# Baixa as dependências e compila o binário
RUN go mod tidy && \
    go build -o app

# Etapa 2: Imagem final, só com o binário
FROM alpine:latest

# Define diretório de trabalho
WORKDIR /app

# Copia o binário da etapa de build
COPY --from=builder /app/app .
# Copia o arquivo .env para o container
COPY --from=builder /app/.env .

# Expõe a porta usada pela aplicação (se aplicável)
EXPOSE 7888

# Comando de inicialização
CMD ["./app"]
