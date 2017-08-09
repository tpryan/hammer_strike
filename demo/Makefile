BASEDIR = $(shell pwd)
include ../Makefile.properties

INSTANCELIST = $(shell gcloud app instances list --format='value[terminator=" "](id)')
CURRENTVERSION = $(shell gcloud app versions list --hide-no-traffic --service=hammer-strike-load --format='value[terminator=" "](id)')
 
clean:
	for i in $(INSTANCELIST); do \
        gcloud app instances delete $$i -q --service=hammer-strike-load --version=$(CURRENTVERSION); \
    done
	

list:
	gcloud app instances list --service=hammer-strike-load

count: 
	gcloud app instances list --service=hammer-strike-load | sed 1d | wc -l	

report:
	gcloud app instances list --service=hammer-strike-load --format='value(id,instance.requests)'