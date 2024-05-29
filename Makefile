#Makefile

build:
	docker build . -t filetransferservice

run:
	docker run -d --name filetransferservicecontainer -p 8080:8080 -t filetransferservice

stop:
	docker stop filetransferservicecontainer
