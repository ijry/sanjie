FROM golang:1.25-alpine AS build

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/sanjie .

FROM alpine:3.22

WORKDIR /app

RUN addgroup -S sanjie \
  && adduser -S -G sanjie sanjie \
  && mkdir -p /data \
  && chown -R sanjie:sanjie /app /data

COPY --from=build /out/sanjie /app/sanjie

USER sanjie
EXPOSE 8080

ENV SANJIE_SERVER_ADDR=:8080
ENV SANJIE_DATABASE_DSN=/data/sanjie.db
ENV SANJIE_DEMO_SEED=true

CMD ["/app/sanjie"]
