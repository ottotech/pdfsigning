FROM python:2.7.13

ENV PYTHONUNBUFFERED=1

# setup the Django app
RUN mkdir /pdfpy
COPY . /pdfpy
WORKDIR /pdfpy
RUN pip install --upgrade pip \
    && pip install reportlab==3.5.13 \
    && pip install PyPDF2==1.26.0