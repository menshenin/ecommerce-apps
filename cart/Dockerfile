FROM golang:1.22 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0  go build -mod=mod -o /bin/cart-server ./cmd/cart-server

FROM scratch

COPY --from=builder /bin/cart-server /bin/cart-server

CMD [ "/bin/cart-server" ]
