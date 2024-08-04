# First stage: build the Go application
FROM --platform=amd64 golang:alpine as builder
ENV VERSION v1.0.0
WORKDIR /app
COPY go.mod go.sum /app/
RUN go mod download
COPY . .
RUN go build -buildvcs=false -v -ldflags="-X 'main.Version=$VERSION'" -o telega_bot ./cmd/main.go

# Second stage: create the final lightweight image
FROM --platform=amd64 alpine as prod
WORKDIR /app
COPY --from=builder /app/telega_bot /app
CMD [ "/app/telega_bot" ]
