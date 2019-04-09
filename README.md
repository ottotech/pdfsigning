## pdfsigning

*pdfsigning* is a small tool to sign PDFs files with a simple generic stamp on each page at the 
upper-left-corner, regardless the orientation of the page. It also draws a logo (274x82) next to the stamp. And it 
also encrypts the pdf if the user requires it so.

## Install

To run the app you need Docker and Docker Compose. 
Simply run: 
```
docker-compose up -d
```
at root level to start the process and go to http://localhost/ to see the app in action.

## Common Application Directories

### `/python_scripts`

*/python_scripts* contains a python script that is used under the hood by the Golang app to sign the pages of the pdf
file with the reportlab and PyPDF2 libraries. This directory also has a Dockerfile and a docker-compose.yml file you can 
run if you want to test or change the pdf signing functionality with python.

## Docker

service name: **gopdfapp**

container name: **gopdfapp**

image name: **gopdfimage**

exposed port: **0.0.0.0:80->8080/tcp**

#### Volumes

* `- .:/go/src/pdfsigning` (bind mounting)

## Built With

* go version go1.11.5 linux/amd64
* Python 2.7.13
* reportlab==3.5.13
* PyPDF2==1.26.0

## Contributing

## Authors 
* Otto Schuldt - *Initial work*

## TODO

* tests

## License

This project is licensed under the MIT License.
