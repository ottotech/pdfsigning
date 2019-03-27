FROM golang

# install pip and python libs
RUN curl https://bootstrap.pypa.io/get-pip.py -o get-pip.py \
    && python get-pip.py \
    && pip install reportlab==3.5.13 \
    && pip install PyPDF2==1.26.0

ADD . /go/src/pdfsigning
WORKDIR /go/src/pdfsigning

CMD go run main.go