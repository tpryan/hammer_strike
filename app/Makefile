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
include ../Makefile.properties

.DEFAULT_GOAL := app

app: create deploy


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