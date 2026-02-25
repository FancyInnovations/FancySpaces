package idp

import (
	"encoding/json"
	"strings"
)

const ApiKeyPrefix = "fancyspaces_api_key"

func (s *Service) ValidateApiKey(apiKeyStr string) (*User, error) {
	if !strings.HasPrefix(apiKeyStr, ApiKeyPrefix) {
		return nil, ErrInvalidApiKeyFormat
	}

	if strings.Count(apiKeyStr, ".") != 2 {
		return nil, ErrInvalidApiKeyFormat
	}

	if len(apiKeyStr) > 300 {
		return nil, ErrInvalidApiKeyFormat
	}

	resp, err := s.broker.Request("idp.apikeys.validate", []byte(apiKeyStr))
	if err != nil {
		return nil, err
	}

	var u User
	if err := json.Unmarshal(resp.Data, &u); err != nil {
		return nil, err
	}

	s.usersCache.UpsertUser(&u)

	return &u, nil
}
