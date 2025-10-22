## Build
FROM golang:1.23-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# binario estático, sin CGO
RUN CGO_ENABLED=0 go build -o /out/api ./cmd/api

## Run (mínima)
FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /out/api /app/api
# 👇 COPIAR migrate.sql DESDE LA ETAPA DE BUILD (no desde el contexto)
COPY --from=build /src/internal/db/migrate.sql /app/internal/db/migrate.sql
ENV PORT=8081
EXPOSE 8081
USER nonroot:nonroot
ENTRYPOINT ["/app/api"]
