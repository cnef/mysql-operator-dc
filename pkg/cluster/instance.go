// Copyright 2018 Oracle and/or its affiliates. All rights reserved.
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

package cluster

import (
	"net"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

type Instance interface {
	GetUser() string
	GetPassword() string
	GetShellURI() string
	GetAddr() string
	Namespace() string
	ClusterName() string
	Name() string
	PodName() string
	Ordinal() int
	Port() int
	MultiMaster() bool
	WhitelistCIDR() (string, error)
}

// NewLocalInstance creates a new instance of this structure, with it's name and index
// populated from os.Hostname().
func NewLocalInstance() (Instance, error) {
	useHostNetwork, err := strconv.ParseBool(os.Getenv("MYSQL_CLUSTER_USE_HOST_NETWORK"))
	if err != nil {
		return nil, err
	}
	glog.V(6).Infoln("NewLocalInstance useHostNetwork: ", useHostNetwork)
	if useHostNetwork {
		return newLocalInstanceInHostNetwork()
	}
	return newLocalInstanceInClusterNetwork()
}

// NewInstanceFromGroupSeed creates an Instance from a fully qualified group
// seed.
func NewInstanceFromGroupSeed(seed string) (Instance, error) {
	useHostNetwork, err := strconv.ParseBool(os.Getenv("MYSQL_CLUSTER_USE_HOST_NETWORK"))
	if err != nil {
		return nil, err
	}
	glog.V(6).Infof("NewInstanceFromGroupSeed useHostNetwork:%v seed:%v", useHostNetwork, seed)
	if useHostNetwork {
		return newInstanceFromGroupSeedInHostNetwork(seed)
	}
	return newInstanceFromGroupSeedInClusterNetwork(seed)
}

// statefulPodRegex is a regular expression that extracts the parent StatefulSet
// and ordinal from StatefulSet Pod's hostname.
var statefulPodRegex = regexp.MustCompile("(.*)-([0-9]+)$")

// GetParentNameAndOrdinal gets the name of a Pod's parent StatefulSet and Pod's
// ordinal from the Pods name (or hostname). If the Pod was not created by a
// StatefulSet, its parent is considered to be empty string, and its ordinal is
// considered to be -1.
func GetParentNameAndOrdinal(name string) (string, int) {
	parent := ""
	ordinal := -1
	subMatches := statefulPodRegex.FindStringSubmatch(name)
	if len(subMatches) < 3 {
		return parent, ordinal
	}
	parent = subMatches[1]
	if i, err := strconv.ParseInt(subMatches[2], 10, 32); err == nil {
		ordinal = int(i)
	}
	return parent, ordinal
}

func podNameFromSeed(seed string) (string, error) {
	host, _, err := net.SplitHostPort(seed)
	if err != nil {
		return "", errors.Wrap(err, "splitting host and port")
	}
	return strings.SplitN(host, ".", 2)[0], nil
}
