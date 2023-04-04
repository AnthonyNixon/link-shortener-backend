FROM golang:1.20-alpine AS build

WORKDIR /go/src/github.com/AnthonyNixon/link-shortener-backend

COPY go.mod ./
COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go get ./...
RUN go build -o /bin/link-shortener-backend ./cmd/link-shortener-backend

FROM alpine:3.17 AS deploy
RUN apk --no-cache add ca-certificates
RUN update-ca-certificates
COPY --from=build /bin/link-shortener-backend /bin/link-shortener-backend
ENTRYPOINT ["/bin/link-shortener-backend"]
EXPOSE 8080
