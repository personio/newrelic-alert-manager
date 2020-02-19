package k8s_test

import (
	"encoding/json"
	"fmt"
	"github.com/fpetkovski/newrelic-operator/pkg/apis/io/v1alpha1"
	"gomodules.xyz/jsonpatch/v2"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"os"
	"testing"
)

const GroupName = "newrelic.fpetkovski.io"
const GroupVersion = "v1alpha1"

var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: GroupVersion}

var (
	SchemeBuilder = runtime.NewSchemeBuilder(addKnownTypes)
	AddToScheme   = SchemeBuilder.AddToScheme
)

func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(
		SchemeGroupVersion,
		&v1alpha1.SlackNotificationChannelList{},
		&v1alpha1.SlackNotificationChannel{},
	)

	v1.AddToGroupVersion(scheme, SchemeGroupVersion)
	return nil
}

func Test(t *testing.T) {
	var config *rest.Config
	var err error

	config, err = clientcmd.BuildConfigFromFlags("", os.Getenv("KUBECONFIG"))
	if err != nil {
		panic(err)
	}

	AddToScheme(scheme.Scheme)

	crdConfig := *config
	crdConfig.ContentConfig.GroupVersion = &SchemeGroupVersion
	crdConfig.APIPath = "/apis"
	crdConfig.NegotiatedSerializer = scheme.Codecs
	crdConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	exampleRestClient, err := rest.UnversionedRESTClientFor(&crdConfig)
	if err != nil {
		panic(err)
	}

	from := v1alpha1.SlackNotificationChannel{}
	from.Status.NewrelicChannelId = new(int64)
	fromJson, _ := json.Marshal(from)

	to := v1alpha1.SlackNotificationChannel{}
	to.Status.NewrelicChannelId = new(int64)
	*to.Status.NewrelicChannelId = 10
	to.Status.Reason = "because"
	toJson, _ := json.Marshal(to)

	patch, _ := jsonpatch.CreatePatch(fromJson, toJson)
	pb, _ := json.MarshalIndent(patch, "", "  ")

	statusCode := new(int)
	r := exampleRestClient.
		Patch(types.JSONPatchType).
		Name("fp-test-123").
		Namespace("default").
		Resource("slacknotificationchannels").
		Body(pb).
		Do().StatusCode(statusCode)

	fmt.Println(*statusCode)

	if r.Error() != nil {
		fmt.Println(r.Error().Error())
	}

	get, err := r.Get()
	if err != nil {
		fmt.Println("error: " + err.Error())
	}
	fmt.Println(get)
}
