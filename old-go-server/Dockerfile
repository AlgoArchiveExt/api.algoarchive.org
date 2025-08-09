FROM golang:1.23

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN make build

EXPOSE 3000

CMD ["make", "run"]