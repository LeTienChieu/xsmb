FROM golang:1.24
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -v -o /home/chieu/share-folder/back-end/go-lang/xsmb ./...
CMD ["app"]