// Code generated by applyconfiguration-gen. DO NOT EDIT.

package v1beta1

import (
	v1beta1 "github.com/openshift/api/machine/v1beta1"
)

// MachineSetStatusApplyConfiguration represents an declarative configuration of the MachineSetStatus type for use
// with apply.
type MachineSetStatusApplyConfiguration struct {
	Replicas               *int32                         `json:"replicas,omitempty"`
	FullyLabeledReplicas   *int32                         `json:"fullyLabeledReplicas,omitempty"`
	ReadyReplicas          *int32                         `json:"readyReplicas,omitempty"`
	AvailableReplicas      *int32                         `json:"availableReplicas,omitempty"`
	ObservedGeneration     *int64                         `json:"observedGeneration,omitempty"`
	ErrorReason            *v1beta1.MachineSetStatusError `json:"errorReason,omitempty"`
	ErrorMessage           *string                        `json:"errorMessage,omitempty"`
	AuthoritativeAPI       *v1beta1.MachineAuthority      `json:"authoritativeAPI,omitempty"`
	SynchronizedGeneration *int64                         `json:"synchronizedGeneration,omitempty"`
}

// MachineSetStatusApplyConfiguration constructs an declarative configuration of the MachineSetStatus type for use with
// apply.
func MachineSetStatus() *MachineSetStatusApplyConfiguration {
	return &MachineSetStatusApplyConfiguration{}
}

// WithReplicas sets the Replicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the Replicas field is set to the value of the last call.
func (b *MachineSetStatusApplyConfiguration) WithReplicas(value int32) *MachineSetStatusApplyConfiguration {
	b.Replicas = &value
	return b
}

// WithFullyLabeledReplicas sets the FullyLabeledReplicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the FullyLabeledReplicas field is set to the value of the last call.
func (b *MachineSetStatusApplyConfiguration) WithFullyLabeledReplicas(value int32) *MachineSetStatusApplyConfiguration {
	b.FullyLabeledReplicas = &value
	return b
}

// WithReadyReplicas sets the ReadyReplicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ReadyReplicas field is set to the value of the last call.
func (b *MachineSetStatusApplyConfiguration) WithReadyReplicas(value int32) *MachineSetStatusApplyConfiguration {
	b.ReadyReplicas = &value
	return b
}

// WithAvailableReplicas sets the AvailableReplicas field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AvailableReplicas field is set to the value of the last call.
func (b *MachineSetStatusApplyConfiguration) WithAvailableReplicas(value int32) *MachineSetStatusApplyConfiguration {
	b.AvailableReplicas = &value
	return b
}

// WithObservedGeneration sets the ObservedGeneration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ObservedGeneration field is set to the value of the last call.
func (b *MachineSetStatusApplyConfiguration) WithObservedGeneration(value int64) *MachineSetStatusApplyConfiguration {
	b.ObservedGeneration = &value
	return b
}

// WithErrorReason sets the ErrorReason field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ErrorReason field is set to the value of the last call.
func (b *MachineSetStatusApplyConfiguration) WithErrorReason(value v1beta1.MachineSetStatusError) *MachineSetStatusApplyConfiguration {
	b.ErrorReason = &value
	return b
}

// WithErrorMessage sets the ErrorMessage field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the ErrorMessage field is set to the value of the last call.
func (b *MachineSetStatusApplyConfiguration) WithErrorMessage(value string) *MachineSetStatusApplyConfiguration {
	b.ErrorMessage = &value
	return b
}

// WithAuthoritativeAPI sets the AuthoritativeAPI field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the AuthoritativeAPI field is set to the value of the last call.
func (b *MachineSetStatusApplyConfiguration) WithAuthoritativeAPI(value v1beta1.MachineAuthority) *MachineSetStatusApplyConfiguration {
	b.AuthoritativeAPI = &value
	return b
}

// WithSynchronizedGeneration sets the SynchronizedGeneration field in the declarative configuration to the given value
// and returns the receiver, so that objects can be built by chaining "With" function invocations.
// If called multiple times, the SynchronizedGeneration field is set to the value of the last call.
func (b *MachineSetStatusApplyConfiguration) WithSynchronizedGeneration(value int64) *MachineSetStatusApplyConfiguration {
	b.SynchronizedGeneration = &value
	return b
}