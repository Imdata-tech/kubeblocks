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

package v1

import (
	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
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
		*out = new(corev1.EnvVarSource)
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
	in.Template.DeepCopyInto(&out.Template)
	if in.Instances != nil {
		in, out := &in.Instances, &out.Instances
		*out = make([]appsv1.InstanceTemplate, len(*in))
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
		*out = make([]corev1.PersistentVolumeClaim, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.ParallelPodManagementConcurrency != nil {
		in, out := &in.ParallelPodManagementConcurrency, &out.ParallelPodManagementConcurrency
		*out = new(intstr.IntOrString)
		**out = **in
	}
	in.UpdateStrategy.DeepCopyInto(&out.UpdateStrategy)
	if in.Roles != nil {
		in, out := &in.Roles, &out.Roles
		*out = make([]appsv1.ReplicaRole, len(*in))
		copy(*out, *in)
	}
	if in.MembershipReconfiguration != nil {
		in, out := &in.MembershipReconfiguration, &out.MembershipReconfiguration
		*out = new(MembershipReconfiguration)
		(*in).DeepCopyInto(*out)
	}
	if in.TemplateVars != nil {
		in, out := &in.TemplateVars, &out.TemplateVars
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
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
func (in *MemberStatus) DeepCopyInto(out *MemberStatus) {
	*out = *in
	if in.ReplicaRole != nil {
		in, out := &in.ReplicaRole, &out.ReplicaRole
		*out = new(appsv1.ReplicaRole)
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
	if in.Switchover != nil {
		in, out := &in.Switchover, &out.Switchover
		*out = new(appsv1.Action)
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
