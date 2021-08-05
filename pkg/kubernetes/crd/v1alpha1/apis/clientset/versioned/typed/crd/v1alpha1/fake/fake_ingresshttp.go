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

// FakeIngressHTTPs implements IngressHTTPInterface
type FakeIngressHTTPs struct {
	Fake *FakeGotwayV1alpha1
	ns   string
}

var ingresshttpsResource = schema.GroupVersionResource{Group: "gotway.io", Version: "v1alpha1", Resource: "ingresshttps"}

var ingresshttpsKind = schema.GroupVersionKind{Group: "gotway.io", Version: "v1alpha1", Kind: "IngressHTTP"}

// Get takes name of the ingressHTTP, and returns the corresponding ingressHTTP object, and an error if there is any.
func (c *FakeIngressHTTPs) Get(ctx context.Context, name string, options v1.GetOptions) (result *v1alpha1.IngressHTTP, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(ingresshttpsResource, c.ns, name), &v1alpha1.IngressHTTP{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IngressHTTP), err
}

// List takes label and field selectors, and returns the list of IngressHTTPs that match those selectors.
func (c *FakeIngressHTTPs) List(ctx context.Context, opts v1.ListOptions) (result *v1alpha1.IngressHTTPList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(ingresshttpsResource, ingresshttpsKind, c.ns, opts), &v1alpha1.IngressHTTPList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.IngressHTTPList{ListMeta: obj.(*v1alpha1.IngressHTTPList).ListMeta}
	for _, item := range obj.(*v1alpha1.IngressHTTPList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested ingressHTTPs.
func (c *FakeIngressHTTPs) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(ingresshttpsResource, c.ns, opts))

}

// Create takes the representation of a ingressHTTP and creates it.  Returns the server's representation of the ingressHTTP, and an error, if there is any.
func (c *FakeIngressHTTPs) Create(ctx context.Context, ingressHTTP *v1alpha1.IngressHTTP, opts v1.CreateOptions) (result *v1alpha1.IngressHTTP, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(ingresshttpsResource, c.ns, ingressHTTP), &v1alpha1.IngressHTTP{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IngressHTTP), err
}

// Update takes the representation of a ingressHTTP and updates it. Returns the server's representation of the ingressHTTP, and an error, if there is any.
func (c *FakeIngressHTTPs) Update(ctx context.Context, ingressHTTP *v1alpha1.IngressHTTP, opts v1.UpdateOptions) (result *v1alpha1.IngressHTTP, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(ingresshttpsResource, c.ns, ingressHTTP), &v1alpha1.IngressHTTP{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IngressHTTP), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeIngressHTTPs) UpdateStatus(ctx context.Context, ingressHTTP *v1alpha1.IngressHTTP, opts v1.UpdateOptions) (*v1alpha1.IngressHTTP, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(ingresshttpsResource, "status", c.ns, ingressHTTP), &v1alpha1.IngressHTTP{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IngressHTTP), err
}

// Delete takes name of the ingressHTTP and deletes it. Returns an error if one occurs.
func (c *FakeIngressHTTPs) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(ingresshttpsResource, c.ns, name), &v1alpha1.IngressHTTP{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeIngressHTTPs) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(ingresshttpsResource, c.ns, listOpts)

	_, err := c.Fake.Invokes(action, &v1alpha1.IngressHTTPList{})
	return err
}

// Patch applies the patch and returns the patched ingressHTTP.
func (c *FakeIngressHTTPs) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v1alpha1.IngressHTTP, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(ingresshttpsResource, c.ns, name, pt, data, subresources...), &v1alpha1.IngressHTTP{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.IngressHTTP), err
}
