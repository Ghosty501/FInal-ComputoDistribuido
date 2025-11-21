# ETAPA 1: BUILDER
# Le damos el alias 'builder' a esta primera etapa.
FROM golang:1.25.3 AS builder 

WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Compilamos el binario estático
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /gateway ./gateway/cmd

# ----------------------------------------------------------------------
# ETAPA 2: FINAL
# ----------------------------------------------------------------------
FROM scratch

# CAMBIO CLAVE: Referenciamos la etapa por su alias: 'builder'
COPY --from=builder /gateway /gateway

# Exponemos el puerto
EXPOSE 8080

# Definimos el punto de entrada
ENTRYPOINT ["/gateway"]