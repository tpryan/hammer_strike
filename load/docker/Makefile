BASEDIR = $(shell pwd)

.DEFAULT_GOAL := app

app: clean build serve

build:
	docker build -t hammer-strike-load "$(BASEDIR)/../containers/."

serve:
	docker run --name=hammer-strike-load -d -P -p 80:80 \
	-e LOAD_URL="http://localhost/healthz" hammer-strike-load

clean:
	-docker stop hammer-strike-load
	-docker rm hammer-strike-load
	-docker rmi hammer-strike-load