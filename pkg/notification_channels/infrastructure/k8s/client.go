package k8s

import (
	"context"
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"github.com/go-logr/logr"
	client_go "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
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

func (c *Client) GetChannel(request reconcile.Request) (*v1alpha1.SlackNotificationChannel, error) {
	var instance v1alpha1.SlackNotificationChannel
	err := c.client.Get(context.TODO(), request.NamespacedName, &instance)
	if err != nil {
		return nil, err
	}

	return &instance, nil
}

func (c *Client) GetAllChannels() (v1alpha1.SlackNotificationChannelList, error) {
	var instance v1alpha1.SlackNotificationChannelList
	err := c.client.List(context.TODO(), &instance)
	if err != nil {
		return instance, err
	}

	return instance, nil
}

func (c *Client) GetPolicies(channel v1alpha1.SlackNotificationChannel) (v1alpha1.AlertPolicyList, error) {
	options := &client_go.ListOptions{
		LabelSelector: channel.Spec.PolicySelector.AsSelector(),
	}

	var result v1alpha1.AlertPolicyList
	err := c.client.List(context.TODO(), &result, options)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *Client) DeleteChannel(policy v1alpha1.SlackNotificationChannel) error {
	policy.ObjectMeta.Finalizers = []string{}
	err := c.client.Update(context.TODO(), &policy)
	if err != nil {
		c.logr.Error(err, "Error updating resource")
		return err
	}
	return nil
}

func (c *Client) UpdateChannel(channel v1alpha1.SlackNotificationChannel) error {
	err := c.client.Status().Patch(context.TODO(), &channel, client_go.MergeFrom(&v1alpha1.SlackNotificationChannel{}))
	if err != nil {
		c.logr.Error(err, "Error patching channel status")
		return err
	}

	channel.ObjectMeta.Finalizers = []string{"newrelic"}
	err = c.client.Update(context.TODO(), &channel)
	if err != nil {
		c.logr.Error(err, "Error updating channel")
		return err
	}

	return nil
}
