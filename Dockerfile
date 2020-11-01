FROM golang:alpine AS builder

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Move to working directory /build
WORKDIR /build

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY cmd ./cmd
COPY configs ./configs
COPY internal ./internal
COPY docs ./docs
COPY web/dist ./dist

# Build the application
WORKDIR /build/cmd/mocks
RUN go build

# Move to /dist directory as the place for resulting binary folder
WORKDIR /dist

# Copy binary from build to main folder
RUN cp /build/cmd/mocks/mocks .
RUN cp -R /build/docs .
RUN cp -R /build/dist .

# Build a small image
FROM scratch

COPY --from=builder /dist/mocks /
COPY --from=builder /dist/docs /docs
COPY --from=builder /dist/dist /dist


EXPOSE 8081

# Command to run
ENTRYPOINT ["/mocks"]