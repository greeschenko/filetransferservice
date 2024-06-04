FROM golang:1.18-alpine
WORKDIR /app
#COPY . .
#COPY /home/olex/prodev/php/yii2/polonex.ga/web/img /app/tmp
#ADD –-from=pubweb . /web

#RUN go mod download
#RUN go mod verify
RUN ls -lahS

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

#CMD ["./main"]
