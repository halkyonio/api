/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1beta1 "halkyon.io/api/capability-info/v1beta1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCapabilityInfos implements CapabilityInfoInterface
type FakeCapabilityInfos struct {
	Fake *FakeHalkyonV1beta1
	ns   string
}

var capabilityinfosResource = schema.GroupVersionResource{Group: "halkyon.io", Version: "v1beta1", Resource: "capabilityinfos"}

var capabilityinfosKind = schema.GroupVersionKind{Group: "halkyon.io", Version: "v1beta1", Kind: "CapabilityInfo"}

// Get takes name of the capabilityInfo, and returns the corresponding capabilityInfo object, and an error if there is any.
func (c *FakeCapabilityInfos) Get(name string, options v1.GetOptions) (result *v1beta1.CapabilityInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(capabilityinfosResource, c.ns, name), &v1beta1.CapabilityInfo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.CapabilityInfo), err
}

// List takes label and field selectors, and returns the list of CapabilityInfos that match those selectors.
func (c *FakeCapabilityInfos) List(opts v1.ListOptions) (result *v1beta1.CapabilityInfoList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(capabilityinfosResource, capabilityinfosKind, c.ns, opts), &v1beta1.CapabilityInfoList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1beta1.CapabilityInfoList{ListMeta: obj.(*v1beta1.CapabilityInfoList).ListMeta}
	for _, item := range obj.(*v1beta1.CapabilityInfoList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested capabilityInfos.
func (c *FakeCapabilityInfos) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(capabilityinfosResource, c.ns, opts))

}

// Create takes the representation of a capabilityInfo and creates it.  Returns the server's representation of the capabilityInfo, and an error, if there is any.
func (c *FakeCapabilityInfos) Create(capabilityInfo *v1beta1.CapabilityInfo) (result *v1beta1.CapabilityInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(capabilityinfosResource, c.ns, capabilityInfo), &v1beta1.CapabilityInfo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.CapabilityInfo), err
}

// Update takes the representation of a capabilityInfo and updates it. Returns the server's representation of the capabilityInfo, and an error, if there is any.
func (c *FakeCapabilityInfos) Update(capabilityInfo *v1beta1.CapabilityInfo) (result *v1beta1.CapabilityInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(capabilityinfosResource, c.ns, capabilityInfo), &v1beta1.CapabilityInfo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.CapabilityInfo), err
}

// Delete takes name of the capabilityInfo and deletes it. Returns an error if one occurs.
func (c *FakeCapabilityInfos) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(capabilityinfosResource, c.ns, name), &v1beta1.CapabilityInfo{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCapabilityInfos) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(capabilityinfosResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1beta1.CapabilityInfoList{})
	return err
}

// Patch applies the patch and returns the patched capabilityInfo.
func (c *FakeCapabilityInfos) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1beta1.CapabilityInfo, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(capabilityinfosResource, c.ns, name, pt, data, subresources...), &v1beta1.CapabilityInfo{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1beta1.CapabilityInfo), err
}
