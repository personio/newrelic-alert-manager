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
	logr    logr.Logger
	client  client_go.Client
	factory v1alpha1.ChannelFactory
}

func NewClient(logr logr.Logger, client client_go.Client, factory v1alpha1.ChannelFactory) *Client {
	return &Client{
		logr:    logr,
		client:  client,
		factory: factory,
	}
}

func (c *Client) GetChannel(name types.NamespacedName) (v1alpha1.NotificationChannel, error) {
	instance := c.factory.NewChannel()
	err := c.client.Get(context.TODO(), name, instance)
	if err != nil {
		return nil, err
	}

	return instance, nil
}

func (c *Client) GetChannels() (v1alpha1.NotificationChannelList, error) {
	instance := c.factory.NewList()
	err := c.client.List(context.TODO(), instance)
	if err != nil {
		return instance, err
	}

	return instance, nil
}

func (c *Client) GetPolicies(channel v1alpha1.NotificationChannel) (v1alpha1.AlertPolicyList, error) {
	options := &client_go.ListOptions{
		LabelSelector: channel.GetPolicySelector(),
	}

	var result v1alpha1.AlertPolicyList
	err := c.client.List(context.TODO(), &result, options)
	if err != nil {
		return result, err
	}

	return result, nil
}

func (c *Client) DeleteChannel(channel v1alpha1.NotificationChannel) error {
	channel.SetFinalizers([]string{})
	err := c.client.Update(context.TODO(), channel)
	if err != nil {
		c.logr.Error(err, "Error deleting channel")
		return err
	}
	return nil
}

func (c *Client) UpdateChannelStatus(channel v1alpha1.NotificationChannel) error {
	return c.updateWithRetries(channel.GetNamespacedName(), channel)
}

func (c *Client) updateWithRetries(key types.NamespacedName, channel v1alpha1.NotificationChannel) error {
	err := c.client.Status().Update(context.TODO(), channel)

	if err != nil && errors.IsConflict(err) {
		c.logr.Info("Conflict updating channel status, retrying")
		serverChannel, err := c.GetChannel(key)
		if err != nil {
			c.logr.Error(err, "Error updating channel status")
			return err
		}

		serverChannel.SetStatus(channel.GetStatus())
		return c.updateWithRetries(key, serverChannel)
	}

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SetFinalizer(channel v1alpha1.NotificationChannel) error {
	channel.SetFinalizers([]string{"newrelic"})
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
