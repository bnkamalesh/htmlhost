FROM golang AS builder

COPY ${PWD} /app-src
WORKDIR /app-src

RUN CGO_ENABLED=0 go build -mod=vendor -ldflags '-s -w -extldflags "-static"' -o /app/htmlhost *.go

COPY internal/server/http/web internal/server/http/web

FROM alpine:3.12
LABEL MAINTAINER Kamaleshwar <bnkamalesh@gmail.com>

# Run as a non-root/non-privileged user
RUN adduser -D appuser
USER appuser

COPY --from=builder /app /home/appuser/app

WORKDIR /home/appuser/app

EXPOSE 8000

CMD ["htmlhost"]
