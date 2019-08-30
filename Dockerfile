FROM golang:nanoserver
EXPOSE 8081
RUN mkdir /app
ADD . /app/
WORKDIR /app
RUN go build -o main .
CMD ["/app/main"]
#CMD ["/bin/bash"]

