package spaces

import (
	"encoding/json"

	"github.com/OliverSchlueter/goutils/broker"
)

type Service struct {
	broker      broker.Broker
	spacesCache *spacesCache
}

type Configuration struct {
	Broker broker.Broker
}

func NewService(cfg Configuration) *Service {
	return &Service{
		broker:      cfg.Broker,
		spacesCache: newSpacesCache(),
	}
}

func (s *Service) GetSpace(id string) (*InternalSpace, error) {
	spaceFromCache, err := s.spacesCache.GetByID(id)
	if err == nil {
		return spaceFromCache, nil
	}

	resp, err := s.broker.Request("fancyspaces.core.spaces.get", []byte(id))
	if err != nil {
		return nil, err
	}

	var space InternalSpace
	if err := json.Unmarshal(resp.Data, &space); err != nil {
		return nil, err
	}

	s.spacesCache.UpsertSpace(&space)

	return &space, nil
}

func (s *Service) GetDecryptedSecret(spaceID, key string) (string, error) {
	req := GetSecretReqDTO{
		SpaceID: spaceID,
		Key:     key,
	}
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return "", err
	}

	resp, err := s.broker.Request("fancyspaces.core.spaces.get_secret", reqBytes)
	if err != nil {
		return "", err
	}

	var secretResp GetSecretRespDTO
	if err := json.Unmarshal(resp.Data, &secretResp); err != nil {
		return "", err
	}

	return secretResp.Value, nil
}
