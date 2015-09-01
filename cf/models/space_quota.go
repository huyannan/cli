package models

func NewSpaceQuota(name string, memory int64, routes int, services int, nonbasicservices bool, orgGuid string, bandwidth int64) (q SpaceQuota) {
	q.Name = name
	q.MemoryLimit = memory
	q.RoutesLimit = routes
	q.ServicesLimit = services
	q.NonBasicServicesAllowed = nonbasicservices
	q.OrgGuid = orgGuid
	q.BandwidthLimit = bandwidth
	return
}

type SpaceQuota struct {
	Guid                    string `json:"guid,omitempty"`
	Name                    string `json:"name"`
	MemoryLimit             int64  `json:"memory_limit"`          // in Megabytes
	InstanceMemoryLimit     int64  `json:"instance_memory_limit"` // in Megabytes
	RoutesLimit             int    `json:"total_routes"`
	ServicesLimit           int    `json:"total_services"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	OrgGuid                 string `json:"organization_guid"`
	BandwidthLimit          int64  `json:"bandwidth_in_kb_limit"`          // in Kilobits
	InstanceBandwidthLimit  int64  `json:"instance_bandwidth_in_kb_limit"` // in Kilobits
}
