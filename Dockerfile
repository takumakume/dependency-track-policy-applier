FROM alpine:3.18.0
RUN apk update && apk add --upgrade libcrypto3 libssl3 curl jq

COPY dependency-track-policy-applier /usr/local/bin/dependency-track-policy-applier

ENTRYPOINT ["dependency-track-policy-applier"]
