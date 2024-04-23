FROM golang:1.22-alpine

EXPOSE 8080

COPY . /app

WORKDIR /app

RUN go build .

CMD [ "./identity" ]