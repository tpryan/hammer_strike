#	Copyright 2017, Google, Inc.
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#    http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
BASEDIR = $(shell pwd)

include ../../Makefile.properties
NODEDISKSIZE="100"
NODECOUNT="5"
NODETYPE="n1-standard-2"

.DEFAULT_GOAL := gke

gke: cluster rule


# Requests a GKE cluster
cluster: env
	-gcloud container clusters create $(CLUSTER) --num-nodes $(NODECOUNT) \
	--disk-size=$(NODEDISKSIZE) --machine-type=$(NODETYPE)

rule:
	gcloud compute --project $(PROJECT) firewall-rules create \
	"external-load-port"  --allow tcp:30000 --network "default" \
	--source-ranges "0.0.0.0/0"

# Runs all the clean commands.
clean: env clean.cluster clean.repo

# deletes the cluster.
clean.cluster:
	gcloud container clusters delete $(CLUSTER) -q

# Gets rid of all of the images in GKR
clean.repo:
	-gsutil -m rm -r gs://artifacts.$(PROJECT).appspot.com