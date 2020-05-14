package applications

import (
	"encoding/json"
	"fmt"
	"github.com/personio/newrelic-alert-manager/internal"
)

type Repository struct {
	client internal.NewrelicClient
}

func NewRepository(client internal.NewrelicClient) *Repository {
	return &Repository{
		client: client,
	}
}

func (repository Repository) GetApplicationByName(name string) (*Application, error) {
	endpoint := fmt.Sprintf("/applications.json?filter[name]=%s", name)
	response, err := repository.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var applications ApplicationList
	err = json.NewDecoder(response.Body).Decode(&applications)
	if err != nil {
		return nil, err
	}

	application := findApplicationByName(applications, name)
	if application == nil {
		return nil, fmt.Errorf("application with name %s does not exist", name)
	}
	return application, nil
}

func findApplicationByName(applications ApplicationList, name string) *Application {
	for _, application := range applications.Applications {
		if application.Name == name {
			return &application
		}
	}

	return nil
}
