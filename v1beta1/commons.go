package v1beta1

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
	"time"
)

const (
	// See https://kubernetes.io/docs/concepts/overview/working-with-objects/common-labels/
	NameLabelKey           = "app.kubernetes.io/name"
	VersionLabelKey        = "app.kubernetes.io/version"
	InstanceLabelKey       = "app.kubernetes.io/instance"
	PartOfLabelKey         = "app.kubernetes.io/part-of"
	ComponentLabelKey      = "app.kubernetes.io/component"
	ManagedByLabelKey      = "app.kubernetes.io/managed-by"
	RuntimeLabelKey        = "app.openshift.io/runtime"
	RuntimeVersionLabelKey = "app.openshift.io/version"
)

type NameValuePair struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// DependentConditionType is a valid value for DependentCondition.Type
type DependentConditionType string

// These are valid dependent conditions.
const (
	// DependentReady means the dependent is able to service requests.
	DependentReady DependentConditionType = "Ready"
	// DependentPending means that the dependent is still processing.
	DependentPending DependentConditionType = "Pending"
	// DependentFailed means that the dependent is in error and probably requires user intervention to get back to working state.
	DependentFailed DependentConditionType = "Failed"
)

// DependentCondition contains details for the current condition of the associated DependentResource.
type DependentCondition struct {
	// Type of the condition.
	Type DependentConditionType `json:"type"`
	// Type of the dependent associated with the condition.
	DependentType schema.GroupVersionKind `json:"dependentType"`
	// Name of the dependent associated with the condition.
	DependentName string `json:"dependentName"`
	// Records the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime v1.Time `json:"lastTransitionTime,omitempty"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty"`
	// Human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty"`
	// Additional information that the condition wishes to convey/record as name-value pairs.
	// +optional
	Attributes []NameValuePair `json:"attributes,omitempty"`
}

const (
	// ReasonPending means the entity has been accepted by the system, but it is still being processed. This includes time
	// being instantiated.
	ReasonPending = "Pending"
	// ReasonReady means the entity has been instantiated to a node and all of its dependencies are available. The
	// entity is able to process requests.
	ReasonReady = "Ready"
	// ReasonFailed means that the entity or some of its dependencies have failed and cannot be salvaged without user
	// intervention.
	ReasonFailed = "Failed"
)

type Status struct {
	LastUpdate v1.Time              `json:"lastUpdate,omitempty"`
	Reason     string               `json:"reason,omitempty"`
	Message    string               `json:"message,omitempty"`
	Conditions []DependentCondition `json:"conditions,omitempty"`
}

func (in *Status) GetConditionsWith(gvk schema.GroupVersionKind) []DependentCondition {
	if gvk.Empty() {
		return []DependentCondition{}
	}

	conditions := make([]DependentCondition, 0, len(in.Conditions))
	for _, condition := range in.Conditions {
		if condition.DependentType == gvk {
			conditions = append(conditions, condition)
		}
	}
	return conditions
}

func (in *Status) GetConditionFor(name string, gvk schema.GroupVersionKind) (existingOrNew *DependentCondition) {
	_, condition := in.indexAndConditionWith(name, gvk)
	return condition
}

func (in *Status) indexAndConditionWith(name string, gvk schema.GroupVersionKind) (index int, existingOrNew *DependentCondition) {
	if len(name) == 0 || gvk.Empty() {
		panic(fmt.Errorf("a condition needs a name and a gvk"))
	}
	if in.Conditions == nil {
		in.Conditions = make([]DependentCondition, 0, 15)
	}
	for i, condition := range in.Conditions {
		if condition.DependentName == name && condition.DependentType == gvk {
			return i, &condition
		}
	}
	existingOrNew = &DependentCondition{
		DependentType: gvk,
		DependentName: name,
	}
	index = len(in.Conditions)
	in.Conditions = append(in.Conditions, *existingOrNew)
	return
}

// SetCondition sets the given condition on the Status, returning true if the status has been modified as a result
func (in *Status) SetCondition(condition *DependentCondition) (updated bool) {
	if condition == nil {
		return false
	}

	// if, for some reason, the index is not set on the condition, retrieve the condition again from the array as this sets the index
	index, previous := in.indexAndConditionWith(condition.DependentName, condition.DependentType)
	if previous == nil || condition.Type != previous.Type || condition.Message != previous.Message {
		condition.Reason = string(condition.Type)
		now := v1.NewTime(time.Now())
		condition.LastTransitionTime = now
		in.LastUpdate = now
		in.Conditions[index] = *condition // update the array with the new condition
		updated = true
	}

	// re-compute overall status only if the set condition has changed or if we don't already have an overall status
	if updated || len(in.Reason) == 0 {
		overall := ReasonReady
		conditionMessages := make([]string, 0, len(in.Conditions))
		for _, c := range in.Conditions {
			// if the condition isn't ready, then the overall status should be pending
			if !c.IsReady() {
				if len(c.Message) > 0 {
					conditionMessages = append(conditionMessages, c.Message)
				}
				overall = ReasonPending
			}
			// if the condition is failed, then the overall status should be failed
			if c.IsFailed() {
				overall = ReasonFailed
			}
		}
		if in.Reason != overall {
			in.Reason = overall
			updated = true
		}
		msg := strings.Join(conditionMessages, ", ")
		if in.Message != msg {
			in.Message = msg
			updated = true
		}
	}

	return
}

func (in *DependentCondition) GetAttribute(name string) string {
	for _, attribute := range in.Attributes {
		if attribute.Name == name {
			return attribute.Value
		}
	}
	return ""
}

func (in *DependentCondition) SetAttribute(name, value string) string {
	for i, attribute := range in.Attributes {
		if attribute.Name == name {
			in.Attributes[i] = NameValuePair{
				Name:  attribute.Name,
				Value: value,
			}
			return attribute.Value
		}
	}
	in.Attributes = append(in.Attributes, NameValuePair{Name: name, Value: value})
	return ""
}

func (in *DependentCondition) IsReady() bool {
	return in.Type == DependentReady
}

func (in *DependentCondition) IsFailed() bool {
	return in.Type == DependentFailed
}

type StatusAware interface {
	GetStatus() Status
	SetStatus(status Status)
	Handle(err error) (bool, Status)
}

func (in *DependentCondition) DeepCopyInto(out *DependentCondition) {
	*out = *in
	out.DependentType = in.DependentType
	in.LastTransitionTime.DeepCopyInto(&out.LastTransitionTime)
	if in.Attributes != nil {
		in, out := &in.Attributes, &out.Attributes
		*out = make([]NameValuePair, len(*in))
		copy(*out, *in)
	}
	return
}

func (in *DependentCondition) DeepCopy() *DependentCondition {
	if in == nil {
		return nil
	}
	out := new(DependentCondition)
	in.DeepCopyInto(out)
	return out
}

func (in *Status) DeepCopyInto(out *Status) {
	*out = *in
	in.LastUpdate.DeepCopyInto(&out.LastUpdate)
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]DependentCondition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	return
}

func (in *Status) DeepCopy() *Status {
	if in == nil {
		return nil
	}
	out := new(Status)
	in.DeepCopyInto(out)
	return out
}
