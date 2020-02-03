package repositories

import (
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/client"
	"github.com/go-logr/logr"
)

type AlertPolicyRepository struct {
	client                  *client.Client
	logr                    logr.Logger
	nrqlConditionRepository *nrqlConditionRepository
}

func NewAlertPolicyRepository(logr logr.Logger, client *client.Client) *AlertPolicyRepository {
	return &AlertPolicyRepository{
		client:                  client,
		logr:                    logr,
		nrqlConditionRepository: newNrqlConditionRepository(logr, client),
	}
}

func (repository AlertPolicyRepository) Save(policy *domain.NewrelicPolicy) error {
	if policy.Policy.Id == nil {
		err := repository.createPolicy(policy)
		if err != nil {
			return err
		}
	} else {
		err := repository.updatePolicy(policy)
		if err != nil {
			return err
		}
	}

	existingConditions, err := repository.nrqlConditionRepository.getConditions(*policy.Policy.Id)
	if err != nil {
		return err
	}

	for _, condition := range existingConditions.Condition {
		err := repository.nrqlConditionRepository.deleteConditions(*condition.Id)
		if err != nil {
			return err
		}
	}

	for _, newCondition := range policy.NrqlConditions {
		err := repository.nrqlConditionRepository.saveCondition(*policy.Policy.Id, newCondition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository AlertPolicyRepository) Delete(policy *domain.NewrelicPolicy) error {
	repository.logr.Info("Deleting policy", "Policy", policy)

	endpoint := fmt.Sprintf("%s/%d.json", "alerts_policies", *policy.Policy.Id)
	_, err := repository.client.Delete(endpoint)

	return err
}

func (repository AlertPolicyRepository) createPolicy(policy *domain.NewrelicPolicy) error {
	repository.logr.Info("Creating policy", "Policy", policy)
	payload, err := json.Marshal(&policy)
	if err != nil {
		return err
	}

	response, err := repository.client.PostJson("alerts_policies.json", payload)
	if err != nil {
		return err
	}

	err = json.NewDecoder(response.Body).Decode(policy)
	if err != nil {
		return err
	}

	return nil
}

func (repository AlertPolicyRepository) updatePolicy(policy *domain.NewrelicPolicy) error {
	repository.logr.Info("Updating policy", "Policy", policy)

	endpoint := fmt.Sprintf("%s/%d.json", "alerts_policies", *policy.Policy.Id)
	payload, err := repository.marshal(policy)
	if err != nil {
		return err
	}

	response, err := repository.client.PutJson(endpoint, payload)
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
	result.NrqlConditions = nil

	payload, err := json.Marshal(&result)
	return payload, err
}
