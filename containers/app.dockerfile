FROM golang:1.25-alpine

WORKDIR /app

COPY go.* ./

RUN go get

COPY . .

EXPOSE 3000

CMD [ "go", "run", "." ]