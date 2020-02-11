package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"testing"
)

func readyStatus(conditions ...*DependentCondition) Status {
	status := status(conditions...)
	status.Reason = ReasonReady
	return status
}

func status(conditions ...*DependentCondition) Status {
	status := Status{}
	dc := make([]DependentCondition, 0, len(conditions))
	for _, condition := range conditions {
		dc = append(dc, *condition)
	}
	status.Conditions = dc
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
		{"setting nil shouldn't do anything", Status{}, "", nil, 0, false},
		{"properly add condition", Status{}, ReasonReady, condition(DependentReady), 1, true},
		{"updating a condition to ready when it was pending before should result in ready status", status(condition(DependentPending)), ReasonReady, condition(DependentReady), 1, true},
		{"failed condition should result in failed status", status(condition(DependentReady, "ready"), condition(DependentPending, "pending"), condition(DependentReady, "failed")), ReasonFailed, condition(DependentFailed, "failed"), 3, true},
		{"not ready condition should result in not ready status", status(condition(DependentReady, "ready"), condition(DependentPending, "pending"), condition(DependentFailed, "failed")), ReasonPending, condition(DependentReady, "failed"), 3, true},
		{"status should be ready if all dependents become ready", status(condition(DependentPending, "pending"), condition(DependentReady, "ready"), condition(DependentReady, "ready2")), ReasonReady, condition(DependentReady, "pending"), 3, true},
		{"status should be ready only if all dependents become ready", status(condition(DependentPending, "pending"), condition(DependentReady, "ready"), condition(DependentPending, "pending2")), ReasonPending, condition(DependentReady, "pending"), 3, true},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			changed := tt.status.SetCondition(tt.condition)
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
