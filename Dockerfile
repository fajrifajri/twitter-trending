FROM golang:latest

WORKDIR /app

COPY trending.go /app
RUN go mod init trending
RUN go mod tidy

RUN go build -o /trending

EXPOSE 2021

CMD [ "/trending" ]
