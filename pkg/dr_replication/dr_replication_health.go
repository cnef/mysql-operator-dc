package drrepl

import (
	"fmt"
	"os"
	"sync"

	"github.com/golang/glog"
	"github.com/heptiolabs/healthcheck"
)

type ReplicationStatus string

// replication status
const (
	ReplStatusNone     ReplicationStatus = ""
	ReplStatusEmpty                      = "empty slave status"
	ReplStatusON                         = "ON"
	ReplStatusOFF                        = "OFF"
	ReplStatusIOError                    = "replication IO thread Error"
	ReplStatusSQLError                   = "replication SQL thread Error"
	ReplStatusUnknow                     = "replication unknow status"
)

// IsON indicate that replication to DC or DR is ON.
func (s ReplicationStatus) IsON() bool {
	return s == ReplStatusON
}

// IsOFF indicate that replication to DC or DR is OFF.
func (s ReplicationStatus) IsOFF() bool {
	return s == ReplStatusOFF || s == ReplStatusEmpty
}

// IsError indicate that there is some errors in replication to DR or DC.
func (s ReplicationStatus) IsError() bool {
	return s == ReplStatusIOError || s == ReplStatusSQLError
}

// DRRepcliationStatus is the status of DR-replication
type DRRepcliationStatus struct {
	NeedDRRepl bool
	Status     bool
	Reason     ReplicationStatus
}

// DeepCopy takes a deep copy of an DRRepcliationStatus object.
func (s *DRRepcliationStatus) DeepCopy() *DRRepcliationStatus {
	return &DRRepcliationStatus{
		NeedDRRepl: s.NeedDRRepl,
		Status:     s.Status,
		Reason:     s.Reason,
	}
}

var (
	drStatus      *DRRepcliationStatus
	drStatusMutex sync.Mutex
)

// SetDRStatus sets the DR replication status of the local mysql cluster. The cluster manager
// controller is responsible for updating.
func SetDRStatus(new *DRRepcliationStatus) {
	drStatusMutex.Lock()
	defer drStatusMutex.Unlock()
	drStatus = new.DeepCopy()
}

// GetDRStatus fetches a copy of the latest DR replication status.
func GetDRStatus() *DRRepcliationStatus {
	drStatusMutex.Lock()
	defer drStatusMutex.Unlock()
	if drStatus == nil {
		return nil
	}
	return drStatus.DeepCopy()
}

// NewDRHealthCheck constructs a DR-healthcheck for the local instance which checks
// cluster replication status between clusters using mysqlsh.
func NewDRHealthCheck() (healthcheck.Check, error) {
	drHost := os.Getenv("MYSQL_CLUSTER_DR_HOST")
	if drHost == "" {
		return nil, fmt.Errorf("No DR-cluster, No need DR-health-check")
	}

	return func() error {
		s := GetDRStatus()
		if s == nil {
			return nil
		}
		if s.NeedDRRepl && !s.Status {
			glog.Errorf("DR-replication fail:%v", s.Reason)
			return fmt.Errorf("DR-replication fail:%v", s.Reason)
		}
		return nil
	}, nil
}
