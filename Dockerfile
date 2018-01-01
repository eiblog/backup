FROM alpine

RUN apk --no-cache add ca-certificates mongodb-tools
COPY backup /app/

WORKDIR /app
ENTRYPOINT ["./backup"]
