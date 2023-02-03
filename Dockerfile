FROM golang:1.19

EXPOSE 5050

WORKDIR /app

COPY go.mod ./
# COPY go.sum ./

RUN go mod download

COPY . .

RUN go build main.go 

CMD ["./main"]
