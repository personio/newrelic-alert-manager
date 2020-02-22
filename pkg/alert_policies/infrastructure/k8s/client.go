package k8s

import (
	"context"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/io/v1alpha1"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	client_go "sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	logr   logr.Logger
	client client_go.Client
}

func NewClient(logr logr.Logger, client client_go.Client) *Client {
	return &Client{
		logr:   logr,
		client: client,
	}
}

func (c *Client) GetPolicy(name types.NamespacedName) (*v1alpha1.AlertPolicy, error) {
	var instance v1alpha1.AlertPolicy
	err := c.client.Get(context.TODO(), name, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

func (c *Client) DeletePolicy(policy v1alpha1.AlertPolicy) error {
	policy.ObjectMeta.Finalizers = []string{}
	err := c.client.Update(context.TODO(), &policy)
	if err != nil {
		c.logr.Error(err, "Error updating resource")
		return err
	}
	return nil
}

func (c *Client) UpdatePolicyStatus(policy *v1alpha1.AlertPolicy) error {
	key := types.NamespacedName{
		Namespace: policy.Namespace,
		Name:      policy.Name,
	}

	err := c.updateWithRetries(key, policy)
	if err != nil {
		c.logr.Error(err, "Error updating status status")
		return err
	}

	return nil
}

func (c *Client) updateWithRetries(key types.NamespacedName, policy *v1alpha1.AlertPolicy) error {
	err := c.client.Status().Update(context.TODO(), policy)

	if err != nil && errors.IsConflict(err) {
		c.logr.Info("Conflict updating policy status, retrying")
		serverPolicy, err := c.GetPolicy(key)
		if err != nil {
			c.logr.Error(err, "Error updating policy status")
			return err
		}

		serverPolicy.Status = policy.Status
		return c.updateWithRetries(key, serverPolicy)
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetFinalizer(policy v1alpha1.AlertPolicy) error {
	policy.ObjectMeta.Finalizers = []string{"newrelic"}
	err := c.client.Update(context.TODO(), &policy)
	if err != nil {
		c.logr.Error(err, "Error updating status")
		return err

	}

	return nil
}