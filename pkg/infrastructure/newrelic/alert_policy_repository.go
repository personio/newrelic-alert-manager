package newrelic

import (
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"github.com/go-logr/logr"
)

type AlertPolicyRepository struct {
	client *Client
	log    logr.Logger
}

func NewAlertPolicyRepository(logr logr.Logger, client *Client) *AlertPolicyRepository {
	return &AlertPolicyRepository{
		client: client,
		log:    logr,
	}
}

func (repository AlertPolicyRepository) Save(policy *domain.NewrelicPolicy) error {
	repository.log.Info("Saving policy", "Policy", policy)
	if policy.Policy.Id == nil {
		return repository.create(policy)
	} else {
		return repository.update(policy)
	}
}

func (repository AlertPolicyRepository) Delete(policy *domain.NewrelicPolicy) error {
	repository.log.Info("Deleting policy", "Policy", policy)

	endpoint := fmt.Sprintf("%s/%d.json", "alerts_policies", *policy.Policy.Id)
	_, err := repository.client.Delete(endpoint)

	return err
}

func (repository AlertPolicyRepository) create(policy *domain.NewrelicPolicy) error {
	repository.log.Info("Creating policy", "Policy", policy)
	payload, err := json.Marshal(&policy)
	if err != nil {
		return err
	}

	response, err := repository.client.Post("alerts_policies.json", payload)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(policy)
	if err != nil {
		return err
	}

	return nil
}

func (repository AlertPolicyRepository) update(policy *domain.NewrelicPolicy) error {
	repository.log.Info("Updating policy", "Policy", policy)

	endpoint := fmt.Sprintf("%s/%d.json", "alerts_policies", *policy.Policy.Id)
	payload, err := repository.marshal(policy)
	if err != nil {
		return err
	}

	response, err := repository.client.Put(endpoint, payload)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(policy)
	if err != nil {
		return err
	}

	return nil
}

func (repository AlertPolicyRepository) marshal(policy *domain.NewrelicPolicy) ([]byte, error) {
	result := *policy
	result.Policy.Id = nil

	payload, err := json.Marshal(&result)
	return payload, err
}
