#Makefile

build:
	docker build . -t filetransferservice

run:
	@docker run -d\
        --name filetransferservicecontainer\
        --network docker_default\
        -p 8000:8000\
        -e MYSQLUSER=${MYSQLUSER}\
        -e MYSQLPASS=${MYSQLPASS}\
        -e MYSQLDOMEN=${MYSQLDOMEN}\
        -e MYSQLDBNAME=${MYSQLDBNAME}\
        -t filetransferservice

stop:
	docker stop filetransferservicecontainer
	docker rm filetransferservicecontainer
