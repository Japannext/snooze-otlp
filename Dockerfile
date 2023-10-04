FROM golang:1.20 as build
ARG VERSION=1.0.0
ARG COMMIT
ARG PROJECT=github.com/japannext/snooze-otlp
WORKDIR /app

# Local CA use
COPY .ca-bundle /usr/local/share/ca-certificates/
RUN chmod -R 644 /usr/local/share/ca-certificates/ && update-ca-certificates

COPY go.mod go.sum ./
RUN go mod download

ADD server/ ./server
ADD cmd/ ./cmd
ADD main.go .
RUN CGO_ENABLED=0 GOOS=linux go build -o ./build/snooze-otlp \
    -ldflags "-X ${PROJECT}/server.Version=${VERSION} -X ${PROJECT}/server.Commit=${COMMIT} -w -s"

FROM scratch
USER 1000
COPY --from=build --chown=1000 --chmod=755 /app/build/snooze-otlp /
CMD ["/snooze-otlp", "run"]
