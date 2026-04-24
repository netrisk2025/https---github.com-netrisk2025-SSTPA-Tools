FROM golang:1.24 AS build

WORKDIR /src

COPY go.mod ./
COPY cmd ./cmd
COPY internal ./internal

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/sstpa-api ./cmd/api

FROM gcr.io/distroless/static-debian12

COPY --from=build /out/sstpa-api /sstpa-api

EXPOSE 8080

ENTRYPOINT ["/sstpa-api"]
