# Copyright 2017 Google Inc. All Rights Reserved.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
BASEDIR = $(shell pwd)

include ../../Makefile.properties

.DEFAULT_GOAL := app

app: build deploy

build:  env creds
	gcloud container builds submit "$(BASEDIR)/../containers/." \
	--tag=gcr.io/$(PROJECT)/hammer-strike

deploy:  creds deploy.config
	kubectl apply -f "$(BASEDIR)/hammer-strike-service.yaml"
	kubectl apply -f "$(BASEDIR)/hammer-strike-deployment.yaml"

deploy.config:
	-kubectl create configmap load-url \
	--from-literal=url=http://$(GAE_HOST)/

clean: clean.deployment clean.service clean.config

clean.deployment: env
	-kubectl delete -f "$(BASEDIR)/hammer-strike-deployment.yaml"

clean.service: env
	-kubectl delete -f "$(BASEDIR)/hammer-strike-service.yaml"

clean.config: env
	-kubectl delete configmap load-url

retry: clean build deploy
