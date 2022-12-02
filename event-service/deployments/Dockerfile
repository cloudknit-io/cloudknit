# syntax=docker/dockerfile:1
FROM gcr.io/distroless/base:debug
USER nonroot:nonroot
COPY event-service /
COPY db/migrations /db/migrations
EXPOSE 8080
ENTRYPOINT [ "/event-service" ]
