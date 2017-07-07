package distributor

import (
	"encoding/json"
	"os"
	"strings"

	"golang.org/x/net/context"

	"golang.org/x/oauth2/google"

	"google.golang.org/api/compute/v1"
	"google.golang.org/appengine"
	"google.golang.org/appengine/log"
	"google.golang.org/appengine/memcache"
)

const cachekey = "LoadNodeList"

// LoadNode is a struct that contains the only fields we need from LoadNode's.
type LoadNode struct {
	Name string `json:"Name"`
	IP   string `json:"ip"`
}

func List(c context.Context) ([]LoadNode, error) {
	lns, err := getFromCache(c)

	if err == nil {
		return lns, nil
	}

	lns, err = getFromProject(c)

	if err = cache(c, lns); err != nil {
		log.Warningf(c, "Could not cache the VM list: "+err.Error())
	}

	return lns, nil

}

func getFromCache(c context.Context) ([]LoadNode, error) {
	var lns []LoadNode

	item, err := memcache.Get(c, cachekey)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(item.Value, &lns)
	return lns, err
}

func getFromProject(c context.Context) ([]LoadNode, error) {

	client, err := google.DefaultClient(c, compute.ComputeScope)
	if err != nil {
		return nil, err
	}
	srv, err := compute.New(client)
	if err != nil {
		return nil, err
	}

	ins, err := srv.Instances.List(appengine.AppID(c), os.Getenv("ZONE")).Do()

	if err != nil {
		return nil, err
	}

	log.Debugf(c, "Response from service: %v", ins.Items)
	var lns []LoadNode
	for _, i := range ins.Items {
		if strings.Index(i.Name, "gke-"+os.Getenv("LOAD_CLUSTER")+"-") != 0 {
			continue
		}
		lns = append(lns, LoadNode{i.Name, publicIP(i)})
	}
	return lns, nil
}

func publicIP(instance *compute.Instance) string {
	for _, network := range instance.NetworkInterfaces {
		for _, cfg := range network.AccessConfigs {
			if cfg.Name == "external-nat" {
				return cfg.NatIP
			}
		}
	}
	return ""
}

func cache(c context.Context, lns []LoadNode) error {
	b, err := json.Marshal(lns)

	if err != nil {
		return err
	}

	item := &memcache.Item{
		Key:   cachekey,
		Value: b,
	}
	return memcache.Set(c, item)
}
