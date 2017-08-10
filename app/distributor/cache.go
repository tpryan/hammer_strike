// Copyright 2017 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package distributor

import (
	"errors"
	"strconv"
	"strings"

	"golang.org/x/net/context"

	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

type Instance struct {
	Name     string `json:"name"`
	Requests int    `json:"requests"`
}

type Report struct {
	RequestCount  int        `json:"request_count"`
	InstanceCount int        `json:"instance_count"`
	Instances     []Instance `json:"instances"`
	Start         string     `json:"start"`
	End           string     `json:"end"`
}

func GetInstances(c context.Context, token string) ([]Instance, error) {
	var result []Instance
	instanceList, err := getInstanceList(c, token)
	if err != nil {
		return result, errors.New("Getting Instance List: " + err.Error())
	}

	result2, err := memcache.GetMulti(c, instanceList)
	log.Infof(c, "Get Instances result: %+v", result2)

	for key, item := range result2 {
		count, err := strconv.Atoi(string(item.Value))
		if err != nil {
			return result, errors.New("Conversion Error: " + err.Error())
		}

		result = append(result, Instance{key, count})
	}

	return result, nil
}

func GetReport(c context.Context, token string) (Report, error) {
	var rep Report
	var err error

	if rep.RequestCount, err = count(c, token+"_total"); err != nil {
		return rep, errors.New("Total Count: " + err.Error())
	}

	if rep.InstanceCount, err = count(c, token+"_totalInstances"); err != nil {
		return rep, errors.New("Instance Count: " + err.Error())
	}

	if rep.Instances, err = GetInstances(c, token); err != nil {
		return rep, errors.New("Getting Instances: " + err.Error())
	}

	if rep.Start, err = value(c, token+"_start"); err != nil {
		return rep, errors.New("Getting Start: " + err.Error())
	}

	if rep.End, err = value(c, token+"_end"); err != nil {
		return rep, errors.New("Getting End: " + err.Error())
	}
	return rep, nil
}

func getInstanceList(c context.Context, token string) ([]string, error) {
	var result []string
	instancesList := ""
	var ins []string

	for i := 0; i <= 15; i++ {
		suffix := strconv.FormatInt(int64(i), 16)
		ins = append(ins, token+"_instances_"+suffix)
	}
	result2, err := memcache.GetMulti(c, ins)
	if err != nil {
		return result, errors.New("Memcache error: " + err.Error())
	}
	for _, item := range result2 {
		instancesList += string(item.Value) + "|"
	}
	instancesList = strings.TrimSuffix(instancesList, "|")
	result = strings.Split(instancesList, "|")

	return result, nil
}

func count(c context.Context, key string) (int, error) {
	item, err := memcache.Get(c, key)
	if err == memcache.ErrCacheMiss {
		return 0, nil
	}
	if err != nil {
		return 0, errors.New("Memcache Error: " + err.Error())
	}
	count, err := strconv.Atoi(string(item.Value))
	if err != nil {
		return 0, errors.New("Conversion Error: " + err.Error())
	}
	return count, nil
}

func value(c context.Context, key string) (string, error) {
	item, err := memcache.Get(c, key)
	if err == memcache.ErrCacheMiss {
		return "", nil
	}
	if err != nil {
		return "", errors.New("Memcache Error: " + err.Error())
	}
	return string(item.Value), nil
}
