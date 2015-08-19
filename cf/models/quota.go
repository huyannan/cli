package models

func NewQuotaFields(name string, memory int64, InstanceMemoryLimit int64, routes int, services int, nonbasicservices bool, bandwidth int64, InstanceBandwidthLimit int64) (q QuotaFields) {
	q.Name = name
	q.MemoryLimit = memory
	q.InstanceMemoryLimit = InstanceMemoryLimit
	q.RoutesLimit = routes
	q.ServicesLimit = services
	q.NonBasicServicesAllowed = nonbasicservices
	q.BandwidthLimit = bandwidth
	q.InstanceBandwidthLimit = InstanceBandwidthLimit
	return
}

type QuotaFields struct {
	Guid                    string `json:"guid,omitempty"`
	Name                    string `json:"name"`
	MemoryLimit             int64  `json:"memory_limit"`          // in Megabytes
	InstanceMemoryLimit     int64  `json:"instance_memory_limit"` // in Megabytes
	RoutesLimit             int    `json:"total_routes"`
	ServicesLimit           int    `json:"total_services"`
	NonBasicServicesAllowed bool   `json:"non_basic_services_allowed"`
	BandwidthLimit          int64  `json:"bandwith_in_kb_limit"`          // in Kilobits
	InstanceBandwidthLimit  int64  `json:"instance_bandwith_in_kb_limit"` // in Kilobits
}
