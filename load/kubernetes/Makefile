BASEDIR = $(shell pwd)

include ../../Makefile.properties

app: build deploy

build:  env creds
	gcloud container builds submit "$(BASEDIR)/../containers/." \
	--tag=gcr.io/$(PROJECT)/hammer-strike

deploy:  creds deploy.config
	kubectl apply -f "$(BASEDIR)/hammer-strike-service.yaml"
	kubectl apply -f "$(BASEDIR)/hammer-strike-deployment.yaml"

deploy.config:
	-kubectl create configmap load-url \
	--from-literal=url=https://$(GAE_HOST)/load/

clean: clean.deployment clean.service clean.config

clean.deployment: env
	-kubectl delete -f "$(BASEDIR)/hammer-strike-deployment.yaml"

clean.service: env
	-kubectl delete -f "$(BASEDIR)/hammer-strike-service.yaml"

clean.config: env
	-kubectl delete configmap load-url

retry: clean build deploy
