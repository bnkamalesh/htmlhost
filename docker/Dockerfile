FROM golang:1.18-alpine3.15 AS builder

COPY ${PWD} /app-src
WORKDIR /app-src

RUN CGO_ENABLED=0 go build -mod=vendor -ldflags '-s -w -extldflags "-static"' -o /app/htmlhost *.go
COPY internal/server/http/web /app/internal/server/http/web

FROM alpine:3.15
LABEL MAINTAINER Kamaleshwar <bnkamalesh@gmail.com>

# Run as a non-root/non-privileged user
RUN adduser -D appuser
USER appuser

COPY --from=builder /app /home/appuser/app

WORKDIR /home/appuser/app

EXPOSE 8000

CMD ["./htmlhost"]
