FROM golang:1.26
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY *.go ./
RUN go build -o /service-a
CMD ["/service-a"]

