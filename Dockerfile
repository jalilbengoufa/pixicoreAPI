FROM golang:1.10
RUN mkdir /app 
ADD . /app/
WORKDIR /app 
RUN go get github.com/gin-gonic/gin
RUN go get github.com/ghodss/yaml
RUN go get golang.org/x/crypto/ssh 
RUN go build -o main .
CMD ["./main"]