FROM golang:1.10
RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN go get github.com/gin-gonic/gin
RUN go get github.com/jinzhu/gorm
RUN go get github.com/mattn/go-sqlite3
RUN go build -o main .
CMD ["./main"]