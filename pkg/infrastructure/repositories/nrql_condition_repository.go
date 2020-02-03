package repositories

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/pkg/domain"
	"github.com/fpetkovski/newrelic-operator/pkg/infrastructure/client"
	"github.com/go-logr/logr"
	"io/ioutil"
)

type nrqlConditionRepository struct {
	client *client.Client
	logr   logr.Logger
}

func newNrqlConditionRepository(logr logr.Logger, client *client.Client) *nrqlConditionRepository {
	return &nrqlConditionRepository{
		client: client,
		logr:   logr,
	}
}

func (repository nrqlConditionRepository) getConditions(policyId int64) (*domain.NrqlConditionList, error) {
	repository.logr.Info("Getting conditions for policy", "PolicyId", policyId)

	endpoint := fmt.Sprintf("alerts_nrql_conditions.json?policy_id=%d", policyId)
	response, err := repository.client.Get(endpoint)
	if err != nil {
		return nil, err
	}

	var conditionList domain.NrqlConditionList
	err = json.NewDecoder(response.Body).Decode(&conditionList)
	if err != nil {
		return nil, err
	}

	fmt.Printf("%+v\n", conditionList)

	return &conditionList, nil
}

func (repository nrqlConditionRepository) deleteConditions(conditionId int64) error {
	repository.logr.Info("Deleting condition", "ConditionId", conditionId)

	fmt.Printf("Deleting condition %d\n", conditionId)
	endpoint := fmt.Sprintf("alerts_nrql_conditions/%d.json", conditionId)
	_, err := repository.client.Delete(endpoint)

	return err
}

func (repository nrqlConditionRepository) saveCondition(policyId int64, condition *domain.NrqlCondition) error {
	repository.logr.Info("Saving condition", "Policy Id", policyId, "Condition", condition)
	payload, err := json.Marshal(&condition)
	if err != nil {
		return err
	}

	endpoint := fmt.Sprintf("alerts_nrql_conditions/policies/%d.json", policyId)
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