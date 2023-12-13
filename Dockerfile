FROM golang:1.21-alpine as build
RUN addgroup appgroup
RUN adduser -G appgroup -D appuser
WORKDIR /tmp/build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY  *.go ./
RUN CGO_ENABLED=0 go build -o ./acceptor

FROM scratch as deploy
WORKDIR /app
USER appuser
COPY --from=build /etc/passwd /etc/passwd
COPY --from=build /tmp/build/acceptor ./

EXPOSE 80/tcp

ENTRYPOINT [ "/app/acceptor" ]