BASEDIR = $(shell pwd)

include ../../Makefile.properties
NODEDISKSIZE="100"
NODECOUNT="10"
NODETYPE="n1-standard-2"

.DEFAULT_GOAL := infrastructure

infrastructure: cluster rule


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