FROM golang
ADD . /go/src/github.com/pjvds/publichost/server
RUN go install github.com/pjvds/publichost/server
ENTRYPOINT /go/bin/server
EXPOSE 5000
EXPOSE 8080
