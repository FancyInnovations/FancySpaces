package cluster

import (
	"log/slog"

	"github.com/OliverSchlueter/goutils/sloki"
)

type DB interface {
	GetCluster(name string) (*Cluster, error)
	GetAllClusters() ([]*Cluster, error)
	CreateCluster(cluster *Cluster) error
	UpdateCluster(cluster *Cluster) error
	DeleteCluster(name string) error
}

type Store struct {
	db DB
}

type Configuration struct {
	DB DB
}

func NewStore(cfg *Configuration) *Store {
	s := &Store{
		db: cfg.DB,
	}

	s.loadLocalClusters()

	return s
}

func (s *Store) loadLocalClusters() {
	err := s.db.CreateCluster(&Cluster{
		Name: "test-cluster",
		Services: []Service{
			{
				Name:     "storage",
				Image:    "oliverschlueter/fancyspaces-storage:latest",
				Replicas: 1,
				Port:     8090,
				ResourceLimits: ResourceLimits{
					CPU:    "500m",
					Memory: "512Mi",
				},
			},
		},
	})

	if err != nil {
		slog.Error("Failed to create local cluster", sloki.WrapError(err))
	}
}
