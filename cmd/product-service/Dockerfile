FROM golang:alpine3.20 as base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download


COPY . .
RUN apk add make

RUN make clean
RUN make build-product

FROM scratch
COPY --from=base /app/product-service-bin /app/

WORKDIR /app

CMD [ "./product-service-bin" ]
