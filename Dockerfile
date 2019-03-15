FROM golang:alpine AS builder
ADD main.go .
RUN go build -o poddle

FROM alpine
RUN apk add --no-cache ffmpeg ca-certificates
WORKDIR /srv
COPY --from=builder /go/poddle ./
COPY app ./app
EXPOSE 8080
CMD ./poddle
