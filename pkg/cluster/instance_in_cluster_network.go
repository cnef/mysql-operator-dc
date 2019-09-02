package cluster

import (
	"fmt"
	"net"
	"os"
	"strconv"

	"github.com/oracle/mysql-operator/pkg/cluster/innodb"
	"github.com/pkg/errors"
)

// InstanceInClusterNetwork represents the local MySQL instance.
type InstanceInClusterNetwork struct {
	// Namespace is the Kubernetes Namespace in which the instance is running.
	namespace string
	// ClusterName is the name of the Cluster to which the instance
	// belongs.
	clusterName string
	// ParentName is the name of the StatefulSet to which the instance belongs.
	parentName string
	// Ordinal is the StatefulSet ordinal of the instances Pod.
	ordinal int
	// Port is the port on which MySQLDB is listening.
	port int
	// MultiMaster specifies if all, or just a single, instance is configured to be read/write.
	multiMaster bool

	// IP is the IP address of the Kubernetes Pod.
	ip net.IP
}

// create local instance when use cluster network
func newLocalInstanceInClusterNetwork() (*InstanceInClusterNetwork, error) {
	hostname, err := os.Hostname()
	if err != nil {
		return nil, err
	}
	name, ordinal := GetParentNameAndOrdinal(hostname)
	multiMaster, _ := strconv.ParseBool(os.Getenv("MYSQL_CLUSTER_MULTI_MASTER"))
	return &InstanceInClusterNetwork{
		namespace:   os.Getenv("POD_NAMESPACE"),
		clusterName: os.Getenv("MYSQL_CLUSTER_NAME"),
		parentName:  name,
		ordinal:     ordinal,
		port:        innodb.MySQLDBPort,
		multiMaster: multiMaster,
		ip:          net.ParseIP(os.Getenv("MY_POD_IP")),
	}, nil
}

func newInstanceFromGroupSeedInClusterNetwork(seed string) (*InstanceInClusterNetwork, error) {
	podName, err := podNameFromSeed(seed)
	if err != nil {
		return nil, errors.Wrap(err, "getting pod name from group seed")
	}
	// We don't care about the returned port here as the Instance's port its
	// MySQLDB port not its group replication port.
	parentName, ordinal := GetParentNameAndOrdinal(podName)
	multiMaster, _ := strconv.ParseBool(os.Getenv("MYSQL_CLUSTER_MULTI_MASTER"))
	return &InstanceInClusterNetwork{
		clusterName: os.Getenv("MYSQL_CLUSTER_NAME"),
		namespace:   os.Getenv("POD_NAMESPACE"),
		parentName:  parentName,
		ordinal:     ordinal,
		port:        innodb.MySQLDBPort,
		multiMaster: multiMaster,
	}, nil
}

// GetUser returns the username of the MySQL operator's management
// user.
func (i *InstanceInClusterNetwork) GetUser() string {
	return "root"
}

// GetPassword returns the password of the MySQL operator's
// management user.
func (i *InstanceInClusterNetwork) GetPassword() string {
	return os.Getenv("MYSQL_ROOT_PASSWORD")
}

// GetShellURI returns the MySQL shell URI for the local MySQL instance.
func (i *InstanceInClusterNetwork) GetShellURI() string {
	return fmt.Sprintf("%s:%s@%s:%d", i.GetUser(), i.GetPassword(), i.Name(), i.port)
}

// GetAddr returns the addr of the instance
func (i *InstanceInClusterNetwork) GetAddr() string {
	return fmt.Sprintf("%s:%d", i.Name(), i.port)
}

// Namespace returns the namespace of the instance
func (i *InstanceInClusterNetwork) Namespace() string {
	return i.namespace
}

// ClusterName returns the clustername of the instance
func (i *InstanceInClusterNetwork) ClusterName() string {
	return i.clusterName
}

// Name returns the name of the instance.
func (i *InstanceInClusterNetwork) Name() string {
	return fmt.Sprintf("%s.%s", i.PodName(), i.parentName)
}

// PodName returns the name of the instance's Pod.
func (i *InstanceInClusterNetwork) PodName() string {
	return fmt.Sprintf("%s-%d", i.parentName, i.ordinal)
}

// Ordinal returns the ordinal of the instance.
func (i *InstanceInClusterNetwork) Ordinal() int {
	return i.ordinal
}

// Port returns the pods of the instance.
func (i *InstanceInClusterNetwork) Port() int {
	return i.port
}

// WhitelistCIDR returns the CIDR range to whitelist for GR based on the Pod's IP.
func (i *InstanceInClusterNetwork) WhitelistCIDR() (string, error) {
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
func (i *InstanceInClusterNetwork) MultiMaster() bool {
	return i.multiMaster
}
