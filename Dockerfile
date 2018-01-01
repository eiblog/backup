FROM alpine

RUN apk --no-cache add mongodb-tools
COPY backup /app/

WORKDIR /app
ENTRYPOINT ["./backup"]
