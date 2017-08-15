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

INSTANCELIST = $(shell gcloud app instances list --format='value[terminator=" "](id)')
CURRENTVERSION = $(shell gcloud app versions list --hide-no-traffic --service=hammer-strike-load --format='value[terminator=" "](id)')
 
clean: env
	for i in $(INSTANCELIST); do \
        gcloud app instances delete $$i -q --service=hammer-strike-load --version=$(CURRENTVERSION); \
    done
	

list: env
	gcloud app instances list --service=hammer-strike-load

count: env
	gcloud app instances list --service=hammer-strike-load | sed 1d | wc -l	

report: env
	gcloud app instances list --service=hammer-strike-load --format='value(id,instance.requests)'