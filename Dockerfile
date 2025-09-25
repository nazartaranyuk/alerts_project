FROM golang:1.24.6 AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/main.go

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /out/app /app/app
COPY configs ./configs
EXPOSE 8080
ENV PORT=:8080
CMD ["./app"]