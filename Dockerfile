FROM mongo:3.2

COPY backup /app/

WORKDIR /app
ENTRYPOINT ["./backup"]
