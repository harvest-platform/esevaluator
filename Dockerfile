FROM scratch

COPY ./dist/linux-amd64/esevaluator /app

EXPOSE 8080

ENTRYPOINT ["/app"]
