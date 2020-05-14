package newrelic

import (
	"encoding/json"
	"fmt"
	"github.com/personio/newrelic-alert-manager/internal"
	"github.com/personio/newrelic-alert-manager/pkg/alert_policies/domain"
	"github.com/go-logr/logr"
)

type AlertPolicyRepository struct {
	client                   internal.NewrelicClient
	infraClient              internal.NewrelicClient
	log                      logr.Logger
	nrqlConditionRepository  *nrqlConditionRepository
	apmConditionRepository   *apmConditionRepository
	infraConditionRepository *infraConditionRepository
}

func NewAlertPolicyRepository(log logr.Logger, client internal.NewrelicClient, infraClient internal.NewrelicClient) *AlertPolicyRepository {
	return &AlertPolicyRepository{
		client:                   client,
		infraClient:              infraClient,
		log:                      log,
		nrqlConditionRepository:  newNrqlConditionRepository(log, client),
		apmConditionRepository:   newApmConditionRepository(log, client),
		infraConditionRepository: newInfraConditionRepository(log, infraClient),
	}
}

func (repository AlertPolicyRepository) Save(policy *domain.AlertPolicy) error {
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

	err := repository.nrqlConditionRepository.saveConditions(policy)
	if err != nil {
		return err
	}

	err = repository.apmConditionRepository.saveConditions(policy)
	if err != nil {
		return err
	}

	err = repository.infraConditionRepository.saveConditions(policy)
	if err != nil {
		return err
	}

	return nil
}

func (repository AlertPolicyRepository) Delete(policy *domain.AlertPolicy) error {
	if policy.Policy.Id == nil {
		return nil
	}

	repository.log.Info("Deleting policy", "PolicyId", *policy.Policy.Id)
	endpoint := fmt.Sprintf("%s/%d.json", "alerts_policies", *policy.Policy.Id)
	response, err := repository.client.Delete(endpoint)
	if response != nil && response.StatusCode == 404 {
		return nil
	}

	return err
}

func (repository AlertPolicyRepository) createPolicy(policy *domain.AlertPolicy) error {
	repository.log.Info("Creating policy", "Policy", policy)
	payload, err := marshal(*policy)
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

func (repository AlertPolicyRepository) updatePolicy(policy *domain.AlertPolicy) error {
	existingPolicy, err := repository.getPolicy(*policy.Policy.Id)
	if err != nil {
		return err
	}

	if existingPolicy == nil {
		return repository.createPolicy(policy)
	}

	if existingPolicy.Equals(*policy) {
		return nil
	}

	endpoint := fmt.Sprintf("%s/%d.json", "alerts_policies", *policy.Policy.Id)
	payload, err := marshal(*policy)
	if err != nil {
		return err
	}

	repository.log.Info("Updating policy", "Policy", policy)
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

func (repository AlertPolicyRepository) getPolicy(policyId int64) (*domain.AlertPolicy, error) {
	var policyList domain.NewrelicPolicyList

	response, err := repository.client.GetJson("alerts_policies.json")
	if err != nil {
		return nil, err
	}

	err = json.NewDecoder(response.Body).Decode(&policyList)
	if err != nil {
		return nil, err
	}

	for _, policy := range policyList.Policies {
		if *policy.Id == policyId {
			return &domain.AlertPolicy{
				Policy: policy,
			}, nil
		}
	}

	return nil, nil
}

func marshal(policy domain.AlertPolicy) ([]byte, error) {
	result := policy
	result.Policy.Id = nil
	result.NrqlConditions = nil
	result.InfraConditions = nil

	payload, err := json.Marshal(result)
	return payload, err
}
