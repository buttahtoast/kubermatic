// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"
	"time"

	scheme "k8c.io/kubermatic/v2/pkg/crd/client/clientset/versioned/scheme"
	v1 "k8c.io/kubermatic/v2/pkg/crd/kubermatic/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// KubermaticSettingsGetter has a method to return a KubermaticSettingInterface.
// A group's client should implement this interface.
type KubermaticSettingsGetter interface {
	KubermaticSettings() KubermaticSettingInterface
}

// KubermaticSettingInterface has methods to work with KubermaticSetting resources.
type KubermaticSettingInterface interface {
	Create(ctx context.Context, kubermaticSetting *v1.KubermaticSetting, opts metav1.CreateOptions) (*v1.KubermaticSetting, error)
	Update(ctx context.Context, kubermaticSetting *v1.KubermaticSetting, opts metav1.UpdateOptions) (*v1.KubermaticSetting, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.KubermaticSetting, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.KubermaticSettingList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.KubermaticSetting, err error)
	KubermaticSettingExpansion
}

// kubermaticSettings implements KubermaticSettingInterface
type kubermaticSettings struct {
	client rest.Interface
}

// newKubermaticSettings returns a KubermaticSettings
func newKubermaticSettings(c *KubermaticV1Client) *kubermaticSettings {
	return &kubermaticSettings{
		client: c.RESTClient(),
	}
}

// Get takes name of the kubermaticSetting, and returns the corresponding kubermaticSetting object, and an error if there is any.
func (c *kubermaticSettings) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.KubermaticSetting, err error) {
	result = &v1.KubermaticSetting{}
	err = c.client.Get().
		Resource("kubermaticsettings").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of KubermaticSettings that match those selectors.
func (c *kubermaticSettings) List(ctx context.Context, opts metav1.ListOptions) (result *v1.KubermaticSettingList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v1.KubermaticSettingList{}
	err = c.client.Get().
		Resource("kubermaticsettings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested kubermaticSettings.
func (c *kubermaticSettings) Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Resource("kubermaticsettings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a kubermaticSetting and creates it.  Returns the server's representation of the kubermaticSetting, and an error, if there is any.
func (c *kubermaticSettings) Create(ctx context.Context, kubermaticSetting *v1.KubermaticSetting, opts metav1.CreateOptions) (result *v1.KubermaticSetting, err error) {
	result = &v1.KubermaticSetting{}
	err = c.client.Post().
		Resource("kubermaticsettings").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kubermaticSetting).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a kubermaticSetting and updates it. Returns the server's representation of the kubermaticSetting, and an error, if there is any.
func (c *kubermaticSettings) Update(ctx context.Context, kubermaticSetting *v1.KubermaticSetting, opts metav1.UpdateOptions) (result *v1.KubermaticSetting, err error) {
	result = &v1.KubermaticSetting{}
	err = c.client.Put().
		Resource("kubermaticsettings").
		Name(kubermaticSetting.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(kubermaticSetting).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the kubermaticSetting and deletes it. Returns an error if one occurs.
func (c *kubermaticSettings) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("kubermaticsettings").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *kubermaticSettings) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Resource("kubermaticsettings").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched kubermaticSetting.
func (c *kubermaticSettings) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.KubermaticSetting, err error) {
	result = &v1.KubermaticSetting{}
	err = c.client.Patch(pt).
		Resource("kubermaticsettings").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
