package cluster

type Cluster struct {
	Name     string    `json:"name"`
	Services []Service `json:"services"`
}

type Service struct {
	Name           string         `json:"name"`
	Image          string         `json:"image"`
	Replicas       int            `json:"replicas"`
	Port           int            `json:"port"`
	ResourceLimits ResourceLimits `json:"resourceLimits"`
}

type ResourceLimits struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
