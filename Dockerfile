FROM golang


ADD . /go/src/pdfsigning
WORKDIR /go/src/pdfsigning

CMD go run main.go