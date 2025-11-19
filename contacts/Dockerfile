# Imagen base
FROM golang:1.25.3

WORKDIR /app

# Copiamos los archivos
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Compilamos
RUN go build -o contacts-service ./contacts

# Exponemos el puerto
EXPOSE 50051

# Comando de ejecución
CMD ["./contacts-service"]
