package newrelic

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/internal"
	"github.com/fpetkovski/newrelic-operator/pkg/alert_policies/domain"
	"github.com/go-logr/logr"
	"github.com/opentracing/opentracing-go"
)

type AlertPolicyRepository struct {
	client                  internal.NewrelicClient
	log                     logr.Logger
	nrqlConditionRepository *nrqlConditionRepository
	apmConditionRepository  *apmConditionRepository
}

func NewAlertPolicyRepository(log logr.Logger, client internal.NewrelicClient) *AlertPolicyRepository {
	return &AlertPolicyRepository{
		client:                  client,
		log:                     log,
		nrqlConditionRepository: newNrqlConditionRepository(log, client),
		apmConditionRepository:  newApmConditionRepository(log, client),
	}
}

func (repository AlertPolicyRepository) Save(ctx context.Context, policy *domain.AlertPolicy) error {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "save-alert-policy-newrelic")
	defer sp.Finish()

	if policy.Policy.Id == nil {
		err := repository.createPolicy(ctx, policy)
		if err != nil {
			return err
		}
	} else {
		err := repository.updatePolicy(ctx, policy)
		if err != nil {
			return err
		}
	}

	err := repository.nrqlConditionRepository.saveConditions(ctx, policy)
	if err != nil {
		return err
	}

	err = repository.apmConditionRepository.saveConditions(ctx, policy)
	if err != nil {
		return err
	}

	return nil
}

func (repository AlertPolicyRepository) Delete(ctx context.Context, policy *domain.AlertPolicy) error {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "delete-alert-policy-newrelic")
	defer sp.Finish()

	endpoint := fmt.Sprintf("%s/%d.json", "alerts_policies", *policy.Policy.Id)
	response, err := repository.client.Delete(endpoint)
	if response != nil && response.StatusCode == 404 {
		fmt.Println(response.StatusCode)
		return nil
	}

	return err
}

func (repository AlertPolicyRepository) createPolicy(ctx context.Context, policy *domain.AlertPolicy) error {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "create-alert-policy-newrelic")
	defer sp.Finish()

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

func (repository AlertPolicyRepository) updatePolicy(ctx context.Context, policy *domain.AlertPolicy) error {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "update-alert-policy-newrelic")
	defer sp.Finish()

	repository.log.Info("Updating policy", "Policy", policy)

	existingPolicy, err := repository.getPolicy(ctx, *policy.Policy.Id)
	if err != nil {
		return err
	}

	if existingPolicy == nil {
		return repository.createPolicy(ctx, policy)
	}

	if existingPolicy.Equals(*policy) {
		return nil
	}

	endpoint := fmt.Sprintf("%s/%d.json", "alerts_policies", *policy.Policy.Id)
	payload, err := marshal(*policy)
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

func (repository AlertPolicyRepository) getPolicy(ctx context.Context, policyId int64) (*domain.AlertPolicy, error) {
	sp, ctx := opentracing.StartSpanFromContext(ctx, "get-alert-policy-newrelic")
	defer sp.Finish()

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

	payload, err := json.Marshal(result)
	return payload, err
}
