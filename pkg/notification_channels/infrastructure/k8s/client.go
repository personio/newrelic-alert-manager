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

func (c *Client) GetChannel(name types.NamespacedName) (*v1alpha1.SlackNotificationChannel, error) {
	var instance v1alpha1.SlackNotificationChannel
	err := c.client.Get(context.TODO(), name, &instance)
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
		c.logr.Error(err, "Error deleting channel")
		return err
	}
	return nil
}

func (c *Client) UpdateChannelStatus(channel *v1alpha1.SlackNotificationChannel) error {
	key := types.NamespacedName{
		Namespace: channel.Namespace,
		Name:      channel.Name,
	}

	return c.updateWithRetries(key, channel)
}

func (c *Client) updateWithRetries(key types.NamespacedName, channel *v1alpha1.SlackNotificationChannel) error {
	err := c.client.Status().Update(context.TODO(), channel)

	if err != nil && errors.IsConflict(err) {
		c.logr.Info("Conflict updating channel status, retrying")
		serverChannel, err := c.GetChannel(key)
		if err != nil {
			c.logr.Error(err, "Error updating channel status")
			return err
		}

		serverChannel.Status = channel.Status
		return c.updateWithRetries(key, serverChannel)
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetFinalizer(channel *v1alpha1.SlackNotificationChannel) error {
	channel.ObjectMeta.Finalizers = []string{"newrelic"}
	err := c.client.Update(context.TODO(), channel)
	if err != nil {
		if errors.IsConflict(err) {
			c.logr.Info("Conflict adding channel finalizer, retrying")
		} else {
			c.logr.Error(err, "Error channel policy finalizer")
		}
		return err
	}

	return nil
}
