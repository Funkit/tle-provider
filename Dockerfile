# syntax=docker/dockerfile:1
FROM golang:1.18-bullseye AS build
LABEL stage=builder
WORKDIR /app
COPY . .
RUN go build -o TLEProvider ./main.go

FROM debian:bullseye-slim
RUN apt-get update \
 && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates
RUN useradd -ms /bin/bash appuser
WORKDIR /home/appuser
COPY --from=build --chown=appuser:appuser /app/TLEProvider .
USER appuser
CMD ["sh", "-c", "./TLEProvider serve --config /home/appuser/configuration.yml"]

EXPOSE 5000
