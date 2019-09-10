FROM golang:latest

WORKDIR /go/src/github.com/dpgolang/PetBook

COPY . /go/src/github.com/dpgolang/PetBook

RUN go build -o ./server ./cmd/main.go

CMD "./server"

EXPOSE 8080


