FROM golang:1.13.7-stretch
RUN mkdir /app 
ADD ./app /app/
WORKDIR /app 
RUN go mod download
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]


