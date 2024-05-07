FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -mod=mod -o /bin/consumer ./cmd/consumer

FROM scratch

COPY --from=builder /bin/consumer /bin/consumer

CMD [ "/bin/consumer" ]
