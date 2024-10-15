FROM golang:alpine3.20 AS base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN apk add make

RUN make build-user

FROM scratch
COPY --from=base /app/user-service-bin /app/

WORKDIR /app

CMD [ "./user-service-bin" ]
