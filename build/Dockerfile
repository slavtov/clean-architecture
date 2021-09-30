FROM golang:alpine as builder

WORKDIR /go/src/app
COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o app cmd/app/main.go

FROM alpine:latest  

RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/app/configs ./configs
COPY --from=builder /go/src/app/app ./
CMD ["./app"]
