BASEDIR = $(shell pwd)
include ../Makefile.properties

create: env
	gcloud app create --project $(PROJECT) --region $(REGION)

deploy: deploy.ui deploy.distributor deploy.load deploy.dispatch deploy.queue

deploy.distributor: env
	gcloud app deploy distributor/app.yaml --project $(PROJECT) -q

deploy.load: env
	gcloud app deploy load/app.yaml --project $(PROJECT) -q

deploy.ui: env
	gcloud app deploy ui/app.yaml --project $(PROJECT) -q

deploy.queue: env
	gcloud app deploy queue.yaml --project $(PROJECT) -q

deploy.dispatch: env
	gcloud app deploy dispatch.yaml --project $(PROJECT) -q