package k8s

import (
	"context"
	"github.com/fpetkovski/newrelic-alert-manager/pkg/apis/newrelic/v1alpha1"
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

func (c *Client) GetDashboard(name types.NamespacedName) (*v1alpha1.Dashboard, error) {
	var instance v1alpha1.Dashboard
	err := c.client.Get(context.TODO(), name, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

func (c *Client) DeleteDashboard(dashboard v1alpha1.Dashboard) error {
	dashboard.ObjectMeta.Finalizers = []string{}
	err := c.client.Update(context.TODO(), &dashboard)
	if err != nil {
		c.logr.Error(err, "Error deleting dashboard")
		return err
	}
	return nil
}

func (c *Client) UpdateDashboardStatus(dashboard *v1alpha1.Dashboard) error {
	key := types.NamespacedName{
		Namespace: dashboard.Namespace,
		Name:      dashboard.Name,
	}

	return c.updateWithRetries(key, dashboard)
}

func (c *Client) updateWithRetries(key types.NamespacedName, dashboard *v1alpha1.Dashboard) error {
	err := c.client.Status().Update(context.TODO(), dashboard)

	if err != nil && errors.IsConflict(err) {
		c.logr.Info("Conflict updating dashboard status, retrying")
		serverDashboard, err := c.GetDashboard(key)
		if err != nil {
			c.logr.Error(err, "Error updating dashboard status")
			return err
		}

		serverDashboard.Status = dashboard.Status
		return c.updateWithRetries(key, serverDashboard)
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetFinalizer(dashboard v1alpha1.Dashboard) error {
	dashboard.ObjectMeta.Finalizers = []string{"newrelic"}
	err := c.client.Update(context.TODO(), &dashboard)
	if err != nil {
		if errors.IsConflict(err) {
			c.logr.Info("Conflict adding dashboard finalizer, retrying")
		} else {
			c.logr.Error(err, "Error setting dashboard finalizer")
		}
		return err

	}

	return nil
}
