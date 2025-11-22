FROM golang:1.23

WORKDIR /app

# Copia o código
COPY . .

# Baixa dependências
RUN go mod download

# Comando padrão: rodar testes
CMD ["go", "test", "./...", "-v"]