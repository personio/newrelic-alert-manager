package newrelic

import (
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-alert-manager/internal"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/dashboards/domain"
	"github.com/go-logr/logr"
)

type Repository struct {
	logr   logr.Logger
	client internal.NewrelicClient
}

func NewRepository(logr logr.Logger, client internal.NewrelicClient) *Repository {
	return &Repository{
		logr:   logr,
		client: client,
	}
}

func (repository Repository) Save(dashboard *domain.Dashboard) error {
	if dashboard.DashboardBody.Id == nil {
		return repository.create(dashboard)
	} else {
		return repository.update(dashboard)
	}
}

func (repository Repository) create(dashboard *domain.Dashboard) error {
	repository.logr.Info("Creating dashboard", "DashboardBody", dashboard)

	payload, err := marshal(*dashboard)
	if err != nil {
		return err
	}

	response, err := repository.client.PostJson("/dashboards.json", payload)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(&dashboard)
	if err != nil {
		return err
	}

	return nil
}

func (repository Repository) update(dashboard *domain.Dashboard) error {
	repository.logr.Info("Updating dashboard", "DashboardBody", dashboard)

	existingDashboard, err := repository.get(*dashboard.DashboardBody.Id)
	if err != nil {
		return err
	}

	if existingDashboard == nil {
		return repository.create(dashboard)
	}

	if existingDashboard.Equals(*dashboard) {
		return nil
	}

	payload, err := marshal(*dashboard)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("/dashboards/%d.json", *dashboard.DashboardBody.Id)
	_, err = repository.client.PutJson(endpoint, payload)
	if err != nil {
		return err
	}

	return nil
}

func marshal(dashboard domain.Dashboard) ([]byte, error) {
	dashboard.DashboardBody.Id = nil
	return json.Marshal(dashboard)
}

func (repository Repository) get(channelId int64) (*domain.Dashboard, error) {
	endpoint := fmt.Sprintf("/dashboards/%d.json", channelId)
	response, err := repository.client.GetJson(endpoint)

	if response != nil && response.StatusCode == 404 {
		return nil, nil
	}

	if err != nil {
		return nil, err
	}

	var dashboard domain.Dashboard
	err = json.NewDecoder(response.Body).Decode(&dashboard)
	if err != nil {
		return nil, err
	}

	return &dashboard, nil
}

func (repository *Repository) Delete(dashboard domain.Dashboard) error {
	repository.logr.Info("Deleting dashboard", "DashboardBody", dashboard)
	if dashboard.DashboardBody.Id == nil {
		return nil
	}

	endpoint := fmt.Sprintf("/dashboards/%d.json", *dashboard.DashboardBody.Id)
	_, err := repository.client.Delete(endpoint)

	return err
}