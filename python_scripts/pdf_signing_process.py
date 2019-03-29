# python imports
import StringIO
from sys import argv

# third party imports
from PyPDF2 import PdfFileWriter, PdfFileReader
from reportlab.pdfgen import canvas
from reportlab.lib.pagesizes import A4
from reportlab.lib.utils import ImageReader
from reportlab.lib.units import inch

# unpack variables
_, src, dest, date, logo_path = argv

# read existing PDF
existing_pdf = PdfFileReader(file(src, "rb"))

# instantiate Writer
pdf_writer = PdfFileWriter()

# get image
im = ImageReader(logo_path)

# add "watermark" (which is the new pdf) on the pages of the existing pdf
for i in range(0, existing_pdf.getNumPages()):
    page = existing_pdf.getPage(i)
    y = int(page.mediaBox.getHeight())  # page height
    buff = StringIO.StringIO()  # create buffer
    orientation = 'Portrait' if page.get('/Rotate') is None else 'Landscape'

    if orientation == 'Portrait':
        x = 7
        # create new PDF with Reportlab
        can = canvas.Canvas(buff, pagesize=A4)
        can.setFont("Helvetica-Bold", 10)
        can.drawString(x, y - 15, "Signed by lequest.nl")
        can.drawString(x, y - 28, "Reason: Confidential LG")
        can.drawString(x, y - 41, "Location: Rotterdam")
        can.drawString(x, y - 54, "Date: %s" % date)
        can.drawImage(
            image=im,
            x=130,
            y=y - 40,
            width=1.20 * inch,
            height=0.32 * inch,
            mask='auto'
        )
        can.save()
    elif orientation == 'Landscape':
        x = int(page.mediaBox.getUpperRight_x() - 15)
        y = y - 15
        # create new PDF with Reportlab
        can = canvas.Canvas(buff, pagesize=A4)
        can.translate(dx=x, dy=y)  # new origin (0, 0)
        can.rotate(-90)
        can.setFont("Helvetica-Bold", 10)
        can.drawString(0, 0, "Signed by lequest.nl")
        can.drawString(0, -14, "Reason: Confidential LG")
        can.drawString(0, -29, "Location: Rotterdam")
        can.drawString(0, -43, "Date: %s" % date)
        can.drawImage(
            image=im,
            x=130,
            y=-30,
            width=1.20 * inch,
            height=0.32 * inch,
            mask='auto'
        )
        can.save()

    # move to the beginning of the StringIO buffer
    buff.seek(0)
    new_pdf = PdfFileReader(buff)

    # merge new pdf into the page of existing pdf
    page.mergePage(new_pdf.getPage(0))
    pdf_writer.addPage(page)

# finally, write to destination
outputStream = file(dest, "wb")
pdf_writer.write(outputStream)
outputStream.seek(0)
outputStream.close()
