FROM golang:1.25 AS builder
WORKDIR /app
COPY frames.go main.go go.mod /app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/animation

FROM scratch
COPY --from=builder /app/animation /app/animation
ENTRYPOINT ["/app/animation"]