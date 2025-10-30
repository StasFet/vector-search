FROM golang:1.25-alpine

WORKDIR /app

COPY . .

EXPOSE 3000

CMD [ "go", "run", "." ]