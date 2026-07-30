# Build

FROM golang:1.22-alpine AS build

WORKDIR /app

COPY go.mod ./
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/server


# Imagem final - run

FROM alpine:3.19

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=build /app/server .

USER appuser
EXPOSE 8080

HEALTHCHECK --interval=5s --timeout=3s --retries=3 \ 
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

CMD ["./server"]
