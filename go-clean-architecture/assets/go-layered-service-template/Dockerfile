FROM golang:1.22 AS build

WORKDIR /src

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/app ./cmd/app

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=build /out/app /app
EXPOSE 8080

ENTRYPOINT ["/app"]

