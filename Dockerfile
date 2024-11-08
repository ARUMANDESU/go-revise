ARG GOLANG_VERSION="1.23"
ARG ALPINE_VERSION="3.20"

FROM golang:${GOLANG_VERSION}-alpine${ALPINE_VERSION} AS build

ENV GOPATH=/
ENV CGO_ENABLED=0
ENV GOOS=linux

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /bin/main /app/cmd/

FROM alpine:${ALPINE_VERSION}

COPY --from=build /bin/main /bin/main

# Default application environment variables
ENV ENV_MODE=dev
ENV TELEGRAM_URL=localhost:4000
ENV HTTP_PORT=5000

EXPOSE 4000
EXPOSE 5000

CMD ["/bin/main"]
