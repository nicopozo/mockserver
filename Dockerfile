# Build Frontend
FROM node:20-alpine AS web-builder
WORKDIR /app
COPY web/package*.json ./
RUN npm install
COPY web/ .
RUN npm run build

# Build Backend
FROM golang:latest AS buildcontainer
WORKDIR /go/src/app
COPY . .
WORKDIR /go/src/app/cmd/mocks
RUN CGO_ENABLED=0 GOOS=linux go build -v -mod mod -ldflags "-s -w" .

# Final Stage
FROM alpine:latest
WORKDIR /app
COPY --from=buildcontainer /go/src/app/cmd/mocks/mocks .
COPY --from=web-builder /app/dist web/dist

ENV MOCKS_MODE=release
ENV HOST=0.0.0.0
ENV PORT=8080
EXPOSE 8080

CMD ["./mocks"]