FROM golang:1.14-alpine AS builder
WORKDIR /app
COPY go.mod  go.sum ./
RUN go mod download
COPY main.go main.go
COPY internal internal
RUN ls -lrt
RUN pwd
RUN go build -o main main.go
FROM alpine
WORKDIR /app/
COPY --from=builder /app/main .
EXPOSE 8090
CMD [ "./main" ]