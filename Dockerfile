FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /gitinspect ./cmd/gitinspect

FROM alpine:3.20
RUN apk add --no-cache git ca-certificates
COPY --from=builder /gitinspect /usr/local/bin/gitinspect
ENTRYPOINT ["gitinspect"]
