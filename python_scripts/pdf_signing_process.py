# python imports
import StringIO
from sys import argv

# third party imports
from PyPDF2 import PdfFileWriter, PdfFileReader
from reportlab.pdfgen import canvas
from reportlab.lib.pagesizes import A4

# unpack variables
_, src, dest, date = argv

# create buffer
buff = StringIO.StringIO()

# create new PDF with Reportlab
can = canvas.Canvas(buff, pagesize=A4)
can.drawString(7, 800, "Signed by lequest.nl")
can.drawString(7, 780, "Reason: Confidential LG")
can.drawString(7, 760, "Location: Rotterdam")
can.drawString(7, 740, "Date: %s" % date)
can.save()

# move to the beginning of the StringIO buffer
buff.seek(0)
new_pdf = PdfFileReader(buff)

# read existing PDF
existing_pdf = PdfFileReader(file(src, "rb"))
output = PdfFileWriter()

# add "watermark" (which is the new pdf) on the pages of the existing pdf
for i in range(0, existing_pdf.getNumPages()):
    page = existing_pdf.getPage(i)
    page.mergePage(new_pdf.getPage(0))
    output.addPage(page)

# finally, write "output" to destination
outputStream = file(dest, "wb")
output.write(outputStream)
outputStream.close()