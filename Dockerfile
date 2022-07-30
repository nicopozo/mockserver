FROM golang:latest AS buildContainer
WORKDIR /go/src/app

COPY . .

WORKDIR ./cmd/mocks

RUN CGO_ENABLED=0 GOOS=linux go build -v -mod mod -ldflags "-s -w" .

FROM alpine:latest
WORKDIR /app
COPY --from=buildContainer /go/src/app/cmd/mocks/mocks .
COPY --from=buildContainer /go/src/app/web/dist web/dist

ENV GIN_MODE release

ENV HOST 0.0.0.0
ENV PORT 8080
EXPOSE 8080

CMD ["./mocks"]