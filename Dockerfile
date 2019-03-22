FROM golang


ADD . /go/src/pdfsigning
WORKDIR /go/src/pdfsigning

RUN go get github.com/jung-kurt/gofpdf

CMD go run main.go