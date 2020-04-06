package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

const initialStatusMessage = "Initial"

func status(conditions ...*DependentCondition) Status {
	return statusWithInitialReason("", conditions...)
}

func statusWithInitialReason(reason string, conditions ...*DependentCondition) Status {
	status := Status{}
	dc := make([]DependentCondition, 0, len(conditions))
	for _, condition := range conditions {
		dc = append(dc, *condition)
	}
	status.Conditions = dc
	status.Message = initialStatusMessage
	if len(reason) > 0 {
		status.Reason = reason
	}
	return status
}

func condition(conditionType DependentConditionType, name ...string) *DependentCondition {
	d := &DependentCondition{
		Type:          conditionType,
		DependentType: schema.GroupVersionKind{Group: "foo", Version: "bar"},
		DependentName: "default-condition",
	}
	if len(name) == 1 {
		d.DependentName = name[0]
	}
	return d
}

func TestSetCondition(t *testing.T) {
	tests := []struct {
		testName       string
		status         Status
		expectedReason string
		condition      *DependentCondition
		conditionNb    int
		changed        bool
	}{
		{
			testName: "setting nil shouldn't do anything",
			status:   status(),
		},
		{
			testName:       "properly add condition",
			status:         status(),
			expectedReason: ReasonReady,
			condition:      condition(DependentReady),
			conditionNb:    1,
			changed:        true,
		},
		{
			testName:       "updating a condition to ready when it was pending before should result in ready status",
			status:         status(condition(DependentPending)),
			expectedReason: ReasonReady,
			condition:      condition(DependentReady),
			conditionNb:    1,
			changed:        true,
		},
		{
			testName:       "failed condition should result in failed status",
			status:         status(condition(DependentReady, "ready"), condition(DependentPending, "pending"), condition(DependentReady, "failed")),
			expectedReason: ReasonFailed,
			condition:      condition(DependentFailed, "failed"),
			conditionNb:    3,
			changed:        true,
		},
		{
			testName:       "not ready condition should result in not ready status",
			status:         status(condition(DependentReady, "ready"), condition(DependentPending, "pending"), condition(DependentFailed, "failed")),
			expectedReason: ReasonPending,
			condition:      condition(DependentReady, "failed"),
			conditionNb:    3,
			changed:        true,
		},
		{
			testName:       "status should be ready if all dependents become ready",
			status:         status(condition(DependentPending, "pending"), condition(DependentReady, "ready"), condition(DependentReady, "ready2")),
			expectedReason: ReasonReady,
			condition:      condition(DependentReady, "pending"),
			conditionNb:    3,
			changed:        true,
		},
		{
			testName:       "status should be ready only if all dependents become ready",
			status:         status(condition(DependentPending, "pending"), condition(DependentReady, "ready"), condition(DependentPending, "pending2")),
			expectedReason: ReasonPending,
			condition:      condition(DependentReady, "pending"),
			conditionNb:    3,
			changed:        true,
		},
		{
			testName:       "status should not change when setting unchanged condition",
			status:         statusWithInitialReason(ReasonPending, condition(DependentPending, "pending"), condition(DependentReady, "ready")),
			expectedReason: ReasonPending,
			condition:      condition(DependentPending, "pending"),
			conditionNb:    2,
			changed:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			changed := tt.status.SetCondition(tt.condition)
			if changed && tt.status.Message == initialStatusMessage {
				t.Errorf("status has changed so its message should have been updated but hasn't")
			}
			if changed != tt.changed {
				t.Errorf("expected changed status to be %t, got %t", tt.changed, changed)
			}
			i := len(tt.status.Conditions)
			if i != tt.conditionNb {
				t.Errorf("expected to have %d conditions, got %d", tt.conditionNb, i)
			}
			if tt.expectedReason != tt.status.Reason {
				t.Errorf("expected to status to be have %s status, got %s", tt.expectedReason, tt.status.Reason)
			}
		})
	}
}
