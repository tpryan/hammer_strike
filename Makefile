BASEDIR = $(shell pwd)

include Makefile.properties

deploy.dispatch: env 
	appcfg.py update_dispatch -A $(PROJECT) .

deploy.cron: env 
	appcfg.py update_cron -A $(PROJECT) .				