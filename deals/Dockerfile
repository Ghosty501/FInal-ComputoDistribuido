# Imagen base (Basada en Debian, usa apt-get)
FROM golang:1.25.3

WORKDIR /app

# 1. ACTUALIZAR LISTAS Y LUEGO INSTALAR 'git'
# Usamos apt-get en lugar de apk
RUN apt-get update && apt-get install -y --no-install-recommends git

# Copiamos los archivos
COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Compilamos
# ANTES (ERROR):
# RUN go build -o deals-service ./main.go 

# DESPUÉS (CORRECCIÓN):
# Indicamos a Go que compile el binario en el directorio deals
RUN go build -o deals-service ./deals
# Exponemos el puerto de deals
EXPOSE 50052

# Comando de ejecución
CMD ["./deals-service"]
