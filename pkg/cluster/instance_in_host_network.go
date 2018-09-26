package cluster

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/oracle/mysql-operator/pkg/cluster/innodb"
	"github.com/pkg/errors"
)

// InstanceInHostNetwork represents the local MySQL instance.
type InstanceInHostNetwork struct {
	// Namespace is the Kubernetes Namespace in which the instance is running.
	namespace string
	// ClusterName is the name of the Cluster to which the instance
	// belongs.
	clusterName string
	// ordinalIsValid indicate the ordinal is valid or not
	ordinalIsValid bool
	// Ordinal is the StatefulSet ordinal of the instances Pod.
	ordinal int
	// Port is the port on which MySQLDB is listening.
	port int
	// MultiMaster specifies if all, or just a single, instance is configured to be read/write.
	multiMaster bool

	// IP is the IP address of the Kubernetes Pod.
	ip net.IP

	// podname
	podName string
	// hostname
	hostName string
}

// create local instance when use cluster network
func newLocalInstanceInHostNetwork() (*InstanceInHostNetwork, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	podName := os.Getenv("MY_POD_NAME")
	if podName == "" {
		return nil, errors.Errorf("use host-network, but pod name not set")
	}
	_, ordinal := GetParentNameAndOrdinal(podName)
	multiMaster, _ := strconv.ParseBool(os.Getenv("MYSQL_CLUSTER_MULTI_MASTER"))
	return &InstanceInHostNetwork{
		clusterName:    os.Getenv("MYSQL_CLUSTER_NAME"),
		namespace:      os.Getenv("POD_NAMESPACE"),
		ordinalIsValid: true,
		ordinal:        ordinal,
		port:           innodb.MySQLDBPort,
		multiMaster:    multiMaster,
		ip:             net.ParseIP(os.Getenv("MY_POD_IP")),
		hostName:       hostname,
		podName:        podName,
	}, nil
}

// seed should be hostname:mysql-port
func newInstanceFromGroupSeedInHostNetwork(seed string) (*InstanceInHostNetwork, error) {
	hostName, err := podNameFromSeed(seed)
	if err != nil {
		return nil, errors.Wrap(err, "getting pod name from group seed")
	}
	// We don't care about the returned port here as the Instance's port its
	// MySQLDB port not its group replication port.
	multiMaster, _ := strconv.ParseBool(os.Getenv("MYSQL_CLUSTER_MULTI_MASTER"))
	return &InstanceInHostNetwork{
		clusterName: os.Getenv("MYSQL_CLUSTER_NAME"),
		namespace:   os.Getenv("POD_NAMESPACE"),
		port:        innodb.MySQLDBPort,
		multiMaster: multiMaster,
		hostName:    hostName,
		podName:     os.Getenv("MY_POD_NAME"),
	}, nil
}

// GetUser returns the username of the MySQL operator's management
// user.
func (i *InstanceInHostNetwork) GetUser() string {
	return "root"
}

// GetPassword returns the password of the MySQL operator's
// management user.
func (i *InstanceInHostNetwork) GetPassword() string {
	return os.Getenv("MYSQL_ROOT_PASSWORD")
}

// GetShellURI returns the MySQL shell URI for the local MySQL instance.
func (i *InstanceInHostNetwork) GetShellURI() string {
	return fmt.Sprintf("%s:%s@%s:%d", i.GetUser(), i.GetPassword(), i.hostName, i.port)
}

// GetAddr returns the addr of the instance
func (i *InstanceInHostNetwork) GetAddr() string {
	return fmt.Sprintf("%s:%d", i.hostName, i.port)
}

// Namespace returns the namespace of the instance
func (i *InstanceInHostNetwork) Namespace() string {
	return i.namespace
}

// ClusterName returns the clustername of the instance
func (i *InstanceInHostNetwork) ClusterName() string {
	return i.clusterName
}

// Name returns the name of the instance.
func (i *InstanceInHostNetwork) Name() string {
	return i.hostName
}

// PodName returns the pod name of the instance.
func (i *InstanceInHostNetwork) PodName() string {
	return i.podName
}

// Ordinal returns the ordinal of the instance.
func (i *InstanceInHostNetwork) Ordinal() (int, error) {
	if i.ordinalIsValid {
		return i.ordinal, nil
	}
	return 0, errors.Errorf("invalid ordinal in host-network")
}

// Port returns the port of the instance.
func (i *InstanceInHostNetwork) Port() int {
	return i.port
}

// WhitelistCIDR returns the CIDR range to whitelist for GR based on the Pod's IP.
func (i *InstanceInHostNetwork) WhitelistCIDR() (string, error) {
	var privateRanges []*net.IPNet

	for _, addrRange := range []string{
		"10.0.0.0/8",
		"172.0.0.0/8",
		"192.168.0.0/16",
		"100.64.0.0/10", // IPv4 shared address space (RFC 6598), improperly used by kops
	} {
		_, block, _ := net.ParseCIDR(addrRange)
		privateRanges = append(privateRanges, block)
	}

	for _, block := range privateRanges {
		if block.Contains(i.ip) {
			return block.String(), nil
		}
	}

	return "", errors.Errorf("pod IP %q is not a private IPv4 address", i.ip.String())
}

// MultiMaster indicate the mysql-cluster is in multi-master mode or not
func (i *InstanceInHostNetwork) MultiMaster() bool {
	return i.multiMaster
}
