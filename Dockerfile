FROM scratch

COPY ./dist/linux-amd64/esevaluator /app

EXPOSE 8080

CMD ["-http", "0.0.0.0:8080"]

ENTRYPOINT ["/app"]
