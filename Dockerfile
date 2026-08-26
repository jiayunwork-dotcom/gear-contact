FROM golang:1.21-alpine
ENV GOTOOLCHAIN=local
ENV CGO_ENABLED=0
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /gear-contact .
EXPOSE 8080
CMD ["/gear-contact", "-http", ":8080"]
