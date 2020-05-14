package newrelic

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/personio/newrelic-alert-manager/internal"
	"github.com/personio/newrelic-alert-manager/pkg/alert_policies/domain"
	"github.com/go-logr/logr"
	"io/ioutil"
)

type apmConditionRepository struct {
	client internal.NewrelicClient
	log    logr.Logger
}

func newApmConditionRepository(log logr.Logger, client internal.NewrelicClient) *apmConditionRepository {
	return &apmConditionRepository{
		client: client,
		log:    log,
	}
}

func (repository apmConditionRepository) getConditions(policyId int64) (*domain.ApmConditionList, error) {
	endpoint := fmt.Sprintf("alerts_conditions.json?policy_id=%d", policyId)
	response, err := repository.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var conditionList domain.ApmConditionList
	err = json.NewDecoder(response.Body).Decode(&conditionList)
	if err != nil {
		return nil, err
	}

	return &conditionList, nil
}

func (repository apmConditionRepository) saveConditions(policy *domain.AlertPolicy) error {
	existingConditions, err := repository.getConditions(*policy.Policy.Id)
	if err != nil {
		return err
	}

	newConditionsSet := domain.NewApmConditionSetFromSlice(policy.ApmConditions)
	for _, condition := range existingConditions.Condition {
		if newConditionsSet.Contains(condition) {
			continue
		}
		err := repository.deleteConditions(*condition.Id)
		if err != nil {
			return err
		}
	}

	existingConditionSet := domain.NewApmConditionSet(*existingConditions)
	for _, newCondition := range policy.ApmConditions {
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

func (repository apmConditionRepository) deleteConditions(conditionId int64) error {
	repository.log.Info("Deleting alert condition", "ConditionId", conditionId)

	endpoint := fmt.Sprintf("alerts_conditions/%d.json", conditionId)
	_, err := repository.client.Delete(endpoint)

	return err
}

func (repository apmConditionRepository) saveCondition(policyId int64, condition *domain.ApmCondition) error {
	repository.log.Info("Saving alert condition", "Policy Id", policyId, "NrqlConditionBody", condition)
	payload, err := json.Marshal(&condition)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("alerts_conditions/policies/%d.json", policyId)
	response, err := repository.client.PostJson(endpoint, payload)
	if response != nil && response.StatusCode >= 300 {
		responseContent, _ := ioutil.ReadAll(response.Body)
		return errors.New(string(responseContent))
	}

	if err != nil {
		return err
	}

	return nil
}
