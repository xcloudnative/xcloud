// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	v1 "github.com/xcloudnative/xcloud/pkg/pkg/apis/jenkins.io/v1"
	scheme "github.com/xcloudnative/xcloud/pkg/pkg/client/clientset/versioned/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// GitServicesGetter has a method to return a GitServiceInterface.
// A group's client should implement this interface.
type GitServicesGetter interface {
	GitServices(namespace string) GitServiceInterface
}

// GitServiceInterface has methods to work with GitService resources.
type GitServiceInterface interface {
	Create(*v1.GitService) (*v1.GitService, error)
	Update(*v1.GitService) (*v1.GitService, error)
	Delete(name string, options *metav1.DeleteOptions) error
	DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error
	Get(name string, options metav1.GetOptions) (*v1.GitService, error)
	List(opts metav1.ListOptions) (*v1.GitServiceList, error)
	Watch(opts metav1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.GitService, err error)
	GitServiceExpansion
}

// gitServices implements GitServiceInterface
type gitServices struct {
	client rest.Interface
	ns     string
}

// newGitServices returns a GitServices
func newGitServices(c *JenkinsV1Client, namespace string) *gitServices {
	return &gitServices{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the gitService, and returns the corresponding gitService object, and an error if there is any.
func (c *gitServices) Get(name string, options metav1.GetOptions) (result *v1.GitService, err error) {
	result = &v1.GitService{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("gitservices").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of GitServices that match those selectors.
func (c *gitServices) List(opts metav1.ListOptions) (result *v1.GitServiceList, err error) {
	result = &v1.GitServiceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("gitservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested gitServices.
func (c *gitServices) Watch(opts metav1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("gitservices").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a gitService and creates it.  Returns the server's representation of the gitService, and an error, if there is any.
func (c *gitServices) Create(gitService *v1.GitService) (result *v1.GitService, err error) {
	result = &v1.GitService{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("gitservices").
		Body(gitService).
		Do().
		Into(result)
	return
}

// Update takes the representation of a gitService and updates it. Returns the server's representation of the gitService, and an error, if there is any.
func (c *gitServices) Update(gitService *v1.GitService) (result *v1.GitService, err error) {
	result = &v1.GitService{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("gitservices").
		Name(gitService.Name).
		Body(gitService).
		Do().
		Into(result)
	return
}

// Delete takes name of the gitService and deletes it. Returns an error if one occurs.
func (c *gitServices) Delete(name string, options *metav1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("gitservices").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *gitServices) DeleteCollection(options *metav1.DeleteOptions, listOptions metav1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("gitservices").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched gitService.
func (c *gitServices) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.GitService, err error) {
	result = &v1.GitService{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("gitservices").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}