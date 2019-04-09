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
_, src, dest, date, logo_path, encrypted, pwd = argv

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
        can.drawString(x, y - 15, "Signed by Company")
        can.drawString(x, y - 28, "Reason: Confidential")
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
        can.drawString(0, 0, "Signed by company.nl")
        can.drawString(0, -14, "Reason: Confidential")
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

# add meta data to pdf, change this to meet your needs
data = {
    '/Title': 'Shared PDF',
    '/Author': 'Author Name',
    '/Subject': 'Shared PDF',
}
pdf_writer.addMetadata(data)

# add encryption to pdf if necessary
owner_pwd = 'secret'
if encrypted == 'yes':
    pdf_writer.encrypt(user_pwd=pwd, owner_pwd=owner_pwd, use_128bit=True)

# write to destination
with file(dest, "wb") as outputStream:
    pdf_writer.write(outputStream)

    # rewind to the beginning of file
    outputStream.seek(0)
