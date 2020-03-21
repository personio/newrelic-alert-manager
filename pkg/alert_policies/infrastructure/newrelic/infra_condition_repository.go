package newrelic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fpetkovski/newrelic-alert-manager/internal"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/alert_policies/domain"
	"github.com/go-logr/logr"
	"io/ioutil"
)

type infraConditionRepository struct {
	client internal.NewrelicClient
	log    logr.Logger
}

func newInfraConditionRepository(log logr.Logger, client internal.NewrelicClient) *infraConditionRepository {
	return &infraConditionRepository{
		client: client,
		log:    log,
	}
}

func (repository infraConditionRepository) getConditions(policyId int64) (*domain.InfraConditionList, error) {
	repository.log.Info("Getting infra conditions for policy", "PolicyId", policyId)

	endpoint := fmt.Sprintf("alerts/conditions?policy_id=%d", policyId)
	response, err := repository.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var conditionList domain.InfraConditionList
	err = json.NewDecoder(response.Body).Decode(&conditionList)
	if err != nil {
		return nil, err
	}

	return &conditionList, nil
}

func (repository infraConditionRepository) saveConditions(policy *domain.AlertPolicy) error {
	existingConditions, err := repository.getConditions(*policy.Policy.Id)
	if err != nil {
		return err
	}

	newConditionsSet := domain.NewInfraConditionSetFromSlice(policy.InfraConditions)
	for _, condition := range existingConditions.Condition {
		if newConditionsSet.Contains(condition) {
			continue
		}
		err := repository.deleteConditions(*condition.Id)
		if err != nil {
			return err
		}
	}

	existingConditionSet := domain.NewInfraConditionSet(*existingConditions)
	for _, newCondition := range policy.InfraConditions {
		if existingConditionSet.Contains(newCondition.Condition) {
			continue
		}

		err := repository.saveCondition(*policy.Policy.Id, newCondition)
		if err != nil {
			return err
		}
	}

	return nil
}

func (repository infraConditionRepository) deleteConditions(conditionId int64) error {
	repository.log.Info("Deleting alerts condition", "ConditionId", conditionId)

	endpoint := fmt.Sprintf("alerts/conditions/%d", conditionId)
	_, err := repository.client.Delete(endpoint)

	return err
}

func (repository infraConditionRepository) saveCondition(policyId int64, condition *domain.InfraCondition) error {
	repository.log.Info("Saving alerts condition", "Policy Id", policyId, "InfraConditionBody", condition)
	condition.Condition.PolicyId = policyId
	payload, err := json.Marshal(&condition)
	if err != nil {
		return err
	}

	response, err := repository.client.PostJson("alerts/conditions", payload)
	if response != nil && response.StatusCode >= 300 {
		responseContent, _ := ioutil.ReadAll(response.Body)
		return errors.New(string(responseContent))
	}

	if err != nil {
		return err
	}

	return nil
}
