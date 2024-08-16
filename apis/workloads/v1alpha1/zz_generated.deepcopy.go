//go:build !ignore_autogenerated

/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

This file is part of KubeBlocks project

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/intstr"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Action) DeepCopyInto(out *Action) {
	*out = *in
	if in.Command != nil {
		in, out := &in.Command, &out.Command
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Action.
func (in *Action) DeepCopy() *Action {
	if in == nil {
		return nil
	}
	out := new(Action)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Credential) DeepCopyInto(out *Credential) {
	*out = *in
	in.Username.DeepCopyInto(&out.Username)
	in.Password.DeepCopyInto(&out.Password)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Credential.
func (in *Credential) DeepCopy() *Credential {
	if in == nil {
		return nil
	}
	out := new(Credential)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CredentialVar) DeepCopyInto(out *CredentialVar) {
	*out = *in
	if in.ValueFrom != nil {
		in, out := &in.ValueFrom, &out.ValueFrom
		*out = new(v1.EnvVarSource)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CredentialVar.
func (in *CredentialVar) DeepCopy() *CredentialVar {
	if in == nil {
		return nil
	}
	out := new(CredentialVar)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceSet) DeepCopyInto(out *InstanceSet) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceSet.
func (in *InstanceSet) DeepCopy() *InstanceSet {
	if in == nil {
		return nil
	}
	out := new(InstanceSet)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstanceSet) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceSetList) DeepCopyInto(out *InstanceSetList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]InstanceSet, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceSetList.
func (in *InstanceSetList) DeepCopy() *InstanceSetList {
	if in == nil {
		return nil
	}
	out := new(InstanceSetList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *InstanceSetList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceSetSpec) DeepCopyInto(out *InstanceSetSpec) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	in.DefaultTemplateOrdinals.DeepCopyInto(&out.DefaultTemplateOrdinals)
	if in.Selector != nil {
		in, out := &in.Selector, &out.Selector
		*out = new(metav1.LabelSelector)
		(*in).DeepCopyInto(*out)
	}
	if in.Service != nil {
		in, out := &in.Service, &out.Service
		*out = new(v1.Service)
		(*in).DeepCopyInto(*out)
	}
	in.Template.DeepCopyInto(&out.Template)
	if in.Instances != nil {
		in, out := &in.Instances, &out.Instances
		*out = make([]InstanceTemplate, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.OfflineInstances != nil {
		in, out := &in.OfflineInstances, &out.OfflineInstances
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.VolumeClaimTemplates != nil {
		in, out := &in.VolumeClaimTemplates, &out.VolumeClaimTemplates
		*out = make([]v1.PersistentVolumeClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ParallelPodManagementConcurrency != nil {
		in, out := &in.ParallelPodManagementConcurrency, &out.ParallelPodManagementConcurrency
		*out = new(intstr.IntOrString)
		**out = **in
	}
	if in.UpdateStrategy != nil {
		in, out := &in.UpdateStrategy, &out.UpdateStrategy
		*out = new(InstanceUpdateStrategy)
		(*in).DeepCopyInto(*out)
	}
	if in.Roles != nil {
		in, out := &in.Roles, &out.Roles
		*out = make([]ReplicaRole, len(*in))
		copy(*out, *in)
	}
	if in.RoleProbe != nil {
		in, out := &in.RoleProbe, &out.RoleProbe
		*out = new(RoleProbe)
		(*in).DeepCopyInto(*out)
	}
	if in.MembershipReconfiguration != nil {
		in, out := &in.MembershipReconfiguration, &out.MembershipReconfiguration
		*out = new(MembershipReconfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.MemberUpdateStrategy != nil {
		in, out := &in.MemberUpdateStrategy, &out.MemberUpdateStrategy
		*out = new(MemberUpdateStrategy)
		**out = **in
	}
	if in.Credential != nil {
		in, out := &in.Credential, &out.Credential
		*out = new(Credential)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceSetSpec.
func (in *InstanceSetSpec) DeepCopy() *InstanceSetSpec {
	if in == nil {
		return nil
	}
	out := new(InstanceSetSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceSetStatus) DeepCopyInto(out *InstanceSetStatus) {
	*out = *in
	if in.Conditions != nil {
		in, out := &in.Conditions, &out.Conditions
		*out = make([]metav1.Condition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.MembersStatus != nil {
		in, out := &in.MembersStatus, &out.MembersStatus
		*out = make([]MemberStatus, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.CurrentRevisions != nil {
		in, out := &in.CurrentRevisions, &out.CurrentRevisions
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.UpdateRevisions != nil {
		in, out := &in.UpdateRevisions, &out.UpdateRevisions
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.TemplatesStatus != nil {
		in, out := &in.TemplatesStatus, &out.TemplatesStatus
		*out = make([]InstanceTemplateStatus, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceSetStatus.
func (in *InstanceSetStatus) DeepCopy() *InstanceSetStatus {
	if in == nil {
		return nil
	}
	out := new(InstanceSetStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceTemplate) DeepCopyInto(out *InstanceTemplate) {
	*out = *in
	if in.Replicas != nil {
		in, out := &in.Replicas, &out.Replicas
		*out = new(int32)
		**out = **in
	}
	in.Ordinals.DeepCopyInto(&out.Ordinals)
	if in.Annotations != nil {
		in, out := &in.Annotations, &out.Annotations
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Image != nil {
		in, out := &in.Image, &out.Image
		*out = new(string)
		**out = **in
	}
	if in.SchedulingPolicy != nil {
		in, out := &in.SchedulingPolicy, &out.SchedulingPolicy
		*out = new(SchedulingPolicy)
		(*in).DeepCopyInto(*out)
	}
	if in.Resources != nil {
		in, out := &in.Resources, &out.Resources
		*out = new(v1.ResourceRequirements)
		(*in).DeepCopyInto(*out)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]v1.EnvVar, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]v1.Volume, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeMounts != nil {
		in, out := &in.VolumeMounts, &out.VolumeMounts
		*out = make([]v1.VolumeMount, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.VolumeClaimTemplates != nil {
		in, out := &in.VolumeClaimTemplates, &out.VolumeClaimTemplates
		*out = make([]v1.PersistentVolumeClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceTemplate.
func (in *InstanceTemplate) DeepCopy() *InstanceTemplate {
	if in == nil {
		return nil
	}
	out := new(InstanceTemplate)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceTemplateStatus) DeepCopyInto(out *InstanceTemplateStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceTemplateStatus.
func (in *InstanceTemplateStatus) DeepCopy() *InstanceTemplateStatus {
	if in == nil {
		return nil
	}
	out := new(InstanceTemplateStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InstanceUpdateStrategy) DeepCopyInto(out *InstanceUpdateStrategy) {
	*out = *in
	if in.Partition != nil {
		in, out := &in.Partition, &out.Partition
		*out = new(int32)
		**out = **in
	}
	if in.MaxUnavailable != nil {
		in, out := &in.MaxUnavailable, &out.MaxUnavailable
		*out = new(intstr.IntOrString)
		**out = **in
	}
	if in.MemberUpdateStrategy != nil {
		in, out := &in.MemberUpdateStrategy, &out.MemberUpdateStrategy
		*out = new(MemberUpdateStrategy)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InstanceUpdateStrategy.
func (in *InstanceUpdateStrategy) DeepCopy() *InstanceUpdateStrategy {
	if in == nil {
		return nil
	}
	out := new(InstanceUpdateStrategy)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MemberStatus) DeepCopyInto(out *MemberStatus) {
	*out = *in
	if in.ReplicaRole != nil {
		in, out := &in.ReplicaRole, &out.ReplicaRole
		*out = new(ReplicaRole)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MemberStatus.
func (in *MemberStatus) DeepCopy() *MemberStatus {
	if in == nil {
		return nil
	}
	out := new(MemberStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MembershipReconfiguration) DeepCopyInto(out *MembershipReconfiguration) {
	*out = *in
	if in.SwitchoverAction != nil {
		in, out := &in.SwitchoverAction, &out.SwitchoverAction
		*out = new(Action)
		(*in).DeepCopyInto(*out)
	}
	if in.MemberJoinAction != nil {
		in, out := &in.MemberJoinAction, &out.MemberJoinAction
		*out = new(Action)
		(*in).DeepCopyInto(*out)
	}
	if in.MemberLeaveAction != nil {
		in, out := &in.MemberLeaveAction, &out.MemberLeaveAction
		*out = new(Action)
		(*in).DeepCopyInto(*out)
	}
	if in.LogSyncAction != nil {
		in, out := &in.LogSyncAction, &out.LogSyncAction
		*out = new(Action)
		(*in).DeepCopyInto(*out)
	}
	if in.PromoteAction != nil {
		in, out := &in.PromoteAction, &out.PromoteAction
		*out = new(Action)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MembershipReconfiguration.
func (in *MembershipReconfiguration) DeepCopy() *MembershipReconfiguration {
	if in == nil {
		return nil
	}
	out := new(MembershipReconfiguration)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ordinals) DeepCopyInto(out *Ordinals) {
	*out = *in
	if in.Ranges != nil {
		in, out := &in.Ranges, &out.Ranges
		*out = make([]Range, len(*in))
		copy(*out, *in)
	}
	if in.Discrete != nil {
		in, out := &in.Discrete, &out.Discrete
		*out = make([]int32, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ordinals.
func (in *Ordinals) DeepCopy() *Ordinals {
	if in == nil {
		return nil
	}
	out := new(Ordinals)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Range) DeepCopyInto(out *Range) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Range.
func (in *Range) DeepCopy() *Range {
	if in == nil {
		return nil
	}
	out := new(Range)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReplicaRole) DeepCopyInto(out *ReplicaRole) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReplicaRole.
func (in *ReplicaRole) DeepCopy() *ReplicaRole {
	if in == nil {
		return nil
	}
	out := new(ReplicaRole)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *RoleProbe) DeepCopyInto(out *RoleProbe) {
	*out = *in
	if in.BuiltinHandler != nil {
		in, out := &in.BuiltinHandler, &out.BuiltinHandler
		*out = new(string)
		**out = **in
	}
	if in.CustomHandler != nil {
		in, out := &in.CustomHandler, &out.CustomHandler
		*out = make([]Action, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new RoleProbe.
func (in *RoleProbe) DeepCopy() *RoleProbe {
	if in == nil {
		return nil
	}
	out := new(RoleProbe)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SchedulingPolicy) DeepCopyInto(out *SchedulingPolicy) {
	*out = *in
	if in.NodeSelector != nil {
		in, out := &in.NodeSelector, &out.NodeSelector
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Affinity != nil {
		in, out := &in.Affinity, &out.Affinity
		*out = new(v1.Affinity)
		(*in).DeepCopyInto(*out)
	}
	if in.Tolerations != nil {
		in, out := &in.Tolerations, &out.Tolerations
		*out = make([]v1.Toleration, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.TopologySpreadConstraints != nil {
		in, out := &in.TopologySpreadConstraints, &out.TopologySpreadConstraints
		*out = make([]v1.TopologySpreadConstraint, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SchedulingPolicy.
func (in *SchedulingPolicy) DeepCopy() *SchedulingPolicy {
	if in == nil {
		return nil
	}
	out := new(SchedulingPolicy)
	in.DeepCopyInto(out)
	return out
}
