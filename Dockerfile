ARG GO_VERSION=1.12

FROM golang:${GO_VERSION}-alpine AS builder

RUN apk add --no-cache ca-certificates git

WORKDIR /src

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY ./ ./

RUN CGO_ENABLED=0 go build -ldflags '-extldflags "-static"' -o /app ./cmd/app/

FROM scratch AS final

COPY --from=builder /app /app

EXPOSE 9000

ENTRYPOINT ["/app"]
