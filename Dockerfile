FROM alpine:3.15

RUN addgroup -S ory; \
    adduser -S ory -G ory -D  -h /home/ory -s /bin/nologin; \
    chown -R ory:ory /home/ory

RUN apk add -U --no-cache ca-certificates

COPY ory /usr/bin/ory

USER ory

ENTRYPOINT ["ory"]
