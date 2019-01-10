FROM golang AS builder

WORKDIR /go/src/github.com/mikloslorinczi/snake-hub/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o snake-hub .

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /app/
COPY --from=builder /go/src/github.com/mikloslorinczi/snake-hub/snake-hub .
COPY ./snake-hub.yaml .

CMD ["./snake-hub"]