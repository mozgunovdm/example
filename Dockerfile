FROM golang:latest

ADD . /github.com/mozgunovdm/example/cmd
ADD . /github.com/mozgunovdm/example/internal
ADD . /github.com/mozgunovdm/example/pkg
ADD . /github.com/mozgunovdm/example/vendor

#ADD . /go/src/github.com/mozgunovdm/example

COPY .env .

RUN go install github.com/mozgunovdm/example/cmd/employe@latest
# RUN cd /github.com/mozgunovdm/example
# WORKDIR /github.com/mozgunovdm/example
# RUN go mod init
# RUN go install ./cmd/employe

ENTRYPOINT ["/go/bin/employe"]
#ENTRYPOINT ["/github.com/mozgunovdm/example/main"]
