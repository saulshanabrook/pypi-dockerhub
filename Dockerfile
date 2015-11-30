FROM golang:1.5

ADD . /go/src/github.com/saulshanabrook/pypi-dockerhub
WORKDIR /go/src/github.com/saulshanabrook/pypi-dockerhub
ENV GO15VENDOREXPERIMENT 1
RUN go install
CMD "go run *.go""
