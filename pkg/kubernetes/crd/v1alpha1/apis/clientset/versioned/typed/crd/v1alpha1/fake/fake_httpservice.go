/*
MIT License

Copyright (c) 2021 Gotway

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v1alpha1 "github.com/gotway/gotway/pkg/kubernetes/crd/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeHTTPServices implements HTTPServiceInterface
type FakeHTTPServices struct {
	Fake *FakeGotwayV1alpha1
	ns   string
}

var httpservicesResource = schema.GroupVersionResource{Group: "gotway.io", Version: "v1alpha1", Resource: "httpservices"}

var httpservicesKind = schema.GroupVersionKind{Group: "gotway.io", Version: "v1alpha1", Kind: "HTTPService"}

// Get takes name of the hTTPService, and returns the corresponding hTTPService object, and an error if there is any.
func (c *FakeHTTPServices) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.HTTPService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(httpservicesResource, c.ns, name), &v1alpha1.HTTPService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HTTPService), err
}

// List takes label and field selectors, and returns the list of HTTPServices that match those selectors.
func (c *FakeHTTPServices) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.HTTPServiceList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(httpservicesResource, httpservicesKind, c.ns, opts), &v1alpha1.HTTPServiceList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.HTTPServiceList{ListMeta: obj.(*v1alpha1.HTTPServiceList).ListMeta}
	for _, item := range obj.(*v1alpha1.HTTPServiceList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested hTTPServices.
func (c *FakeHTTPServices) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(httpservicesResource, c.ns, opts))

}

// Create takes the representation of a hTTPService and creates it.  Returns the server's representation of the hTTPService, and an error, if there is any.
func (c *FakeHTTPServices) Create(ctx context.Context, hTTPService *v1alpha1.HTTPService, opts v1.CreateOptions) (result *v1alpha1.HTTPService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(httpservicesResource, c.ns, hTTPService), &v1alpha1.HTTPService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HTTPService), err
}

// Update takes the representation of a hTTPService and updates it. Returns the server's representation of the hTTPService, and an error, if there is any.
func (c *FakeHTTPServices) Update(ctx context.Context, hTTPService *v1alpha1.HTTPService, opts v1.UpdateOptions) (result *v1alpha1.HTTPService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(httpservicesResource, c.ns, hTTPService), &v1alpha1.HTTPService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HTTPService), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeHTTPServices) UpdateStatus(ctx context.Context, hTTPService *v1alpha1.HTTPService, opts v1.UpdateOptions) (*v1alpha1.HTTPService, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(httpservicesResource, "status", c.ns, hTTPService), &v1alpha1.HTTPService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HTTPService), err
}

// Delete takes name of the hTTPService and deletes it. Returns an error if one occurs.
func (c *FakeHTTPServices) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(httpservicesResource, c.ns, name), &v1alpha1.HTTPService{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeHTTPServices) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(httpservicesResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.HTTPServiceList{})
	return err
}

// Patch applies the patch and returns the patched hTTPService.
func (c *FakeHTTPServices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.HTTPService, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(httpservicesResource, c.ns, name, pt, data, subresources...), &v1alpha1.HTTPService{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.HTTPService), err
}