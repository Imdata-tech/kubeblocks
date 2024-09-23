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
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CELExpression) DeepCopyInto(out *CELExpression) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CELExpression.
func (in *CELExpression) DeepCopy() *CELExpression {
	if in == nil {
		return nil
	}
	out := new(CELExpression)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EventAttributes) DeepCopyInto(out *EventAttributes) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EventAttributes.
func (in *EventAttributes) DeepCopy() *EventAttributes {
	if in == nil {
		return nil
	}
	out := new(EventAttributes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *FieldPath) DeepCopyInto(out *FieldPath) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new FieldPath.
func (in *FieldPath) DeepCopy() *FieldPath {
	if in == nil {
		return nil
	}
	out := new(FieldPath)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectChange) DeepCopyInto(out *ObjectChange) {
	*out = *in
	out.ObjectReference = in.ObjectReference
	if in.EventAttributes != nil {
		in, out := &in.EventAttributes, &out.EventAttributes
		*out = new(EventAttributes)
		**out = **in
	}
	if in.Timestamp != nil {
		in, out := &in.Timestamp, &out.Timestamp
		*out = (*in).DeepCopy()
	}
	if in.LocalDescription != nil {
		in, out := &in.LocalDescription, &out.LocalDescription
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectChange.
func (in *ObjectChange) DeepCopy() *ObjectChange {
	if in == nil {
		return nil
	}
	out := new(ObjectChange)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectChangeSummary) DeepCopyInto(out *ObjectChangeSummary) {
	*out = *in
	if in.Added != nil {
		in, out := &in.Added, &out.Added
		*out = new(int32)
		**out = **in
	}
	if in.Updated != nil {
		in, out := &in.Updated, &out.Updated
		*out = new(int32)
		**out = **in
	}
	if in.Deleted != nil {
		in, out := &in.Deleted, &out.Deleted
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectChangeSummary.
func (in *ObjectChangeSummary) DeepCopy() *ObjectChangeSummary {
	if in == nil {
		return nil
	}
	out := new(ObjectChangeSummary)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectReference) DeepCopyInto(out *ObjectReference) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectReference.
func (in *ObjectReference) DeepCopy() *ObjectReference {
	if in == nil {
		return nil
	}
	out := new(ObjectReference)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectSummary) DeepCopyInto(out *ObjectSummary) {
	*out = *in
	out.Type = in.Type
	if in.ChangeSummary != nil {
		in, out := &in.ChangeSummary, &out.ChangeSummary
		*out = new(ObjectChangeSummary)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectSummary.
func (in *ObjectSummary) DeepCopy() *ObjectSummary {
	if in == nil {
		return nil
	}
	out := new(ObjectSummary)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectTreeNode) DeepCopyInto(out *ObjectTreeNode) {
	*out = *in
	out.Primary = in.Primary
	if in.Secondaries != nil {
		in, out := &in.Secondaries, &out.Secondaries
		*out = make([]*ObjectTreeNode, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(ObjectTreeNode)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectTreeNode.
func (in *ObjectTreeNode) DeepCopy() *ObjectTreeNode {
	if in == nil {
		return nil
	}
	out := new(ObjectTreeNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ObjectType) DeepCopyInto(out *ObjectType) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ObjectType.
func (in *ObjectType) DeepCopy() *ObjectType {
	if in == nil {
		return nil
	}
	out := new(ObjectType)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OwnedResource) DeepCopyInto(out *OwnedResource) {
	*out = *in
	out.Secondary = in.Secondary
	in.Criteria.DeepCopyInto(&out.Criteria)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OwnedResource.
func (in *OwnedResource) DeepCopy() *OwnedResource {
	if in == nil {
		return nil
	}
	out := new(OwnedResource)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OwnershipCriteria) DeepCopyInto(out *OwnershipCriteria) {
	*out = *in
	if in.SelectorCriteria != nil {
		in, out := &in.SelectorCriteria, &out.SelectorCriteria
		*out = new(FieldPath)
		**out = **in
	}
	if in.LabelCriteria != nil {
		in, out := &in.LabelCriteria, &out.LabelCriteria
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.BuiltinRelationshipCriteria != nil {
		in, out := &in.BuiltinRelationshipCriteria, &out.BuiltinRelationshipCriteria
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OwnershipCriteria.
func (in *OwnershipCriteria) DeepCopy() *OwnershipCriteria {
	if in == nil {
		return nil
	}
	out := new(OwnershipCriteria)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OwnershipRule) DeepCopyInto(out *OwnershipRule) {
	*out = *in
	out.Primary = in.Primary
	if in.OwnedResources != nil {
		in, out := &in.OwnedResources, &out.OwnedResources
		*out = make([]OwnedResource, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OwnershipRule.
func (in *OwnershipRule) DeepCopy() *OwnershipRule {
	if in == nil {
		return nil
	}
	out := new(OwnershipRule)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *PlanSummary) DeepCopyInto(out *PlanSummary) {
	*out = *in
	if in.ObjectSummaries != nil {
		in, out := &in.ObjectSummaries, &out.ObjectSummaries
		*out = make([]ObjectSummary, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new PlanSummary.
func (in *PlanSummary) DeepCopy() *PlanSummary {
	if in == nil {
		return nil
	}
	out := new(PlanSummary)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationPlan) DeepCopyInto(out *ReconciliationPlan) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationPlan.
func (in *ReconciliationPlan) DeepCopy() *ReconciliationPlan {
	if in == nil {
		return nil
	}
	out := new(ReconciliationPlan)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ReconciliationPlan) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationPlanList) DeepCopyInto(out *ReconciliationPlanList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ReconciliationPlan, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationPlanList.
func (in *ReconciliationPlanList) DeepCopy() *ReconciliationPlanList {
	if in == nil {
		return nil
	}
	out := new(ReconciliationPlanList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ReconciliationPlanList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationPlanSpec) DeepCopyInto(out *ReconciliationPlanSpec) {
	*out = *in
	if in.TargetObject != nil {
		in, out := &in.TargetObject, &out.TargetObject
		*out = new(ObjectReference)
		**out = **in
	}
	if in.Depth != nil {
		in, out := &in.Depth, &out.Depth
		*out = new(int32)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationPlanSpec.
func (in *ReconciliationPlanSpec) DeepCopy() *ReconciliationPlanSpec {
	if in == nil {
		return nil
	}
	out := new(ReconciliationPlanSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationPlanStatus) DeepCopyInto(out *ReconciliationPlanStatus) {
	*out = *in
	in.CurrentObjectTree.DeepCopyInto(&out.CurrentObjectTree)
	in.DesiredObjectTree.DeepCopyInto(&out.DesiredObjectTree)
	in.Summary.DeepCopyInto(&out.Summary)
	if in.Plan != nil {
		in, out := &in.Plan, &out.Plan
		*out = make([]ObjectChange, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationPlanStatus.
func (in *ReconciliationPlanStatus) DeepCopy() *ReconciliationPlanStatus {
	if in == nil {
		return nil
	}
	out := new(ReconciliationPlanStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationView) DeepCopyInto(out *ReconciliationView) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	in.Status.DeepCopyInto(&out.Status)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationView.
func (in *ReconciliationView) DeepCopy() *ReconciliationView {
	if in == nil {
		return nil
	}
	out := new(ReconciliationView)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ReconciliationView) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationViewDefinition) DeepCopyInto(out *ReconciliationViewDefinition) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationViewDefinition.
func (in *ReconciliationViewDefinition) DeepCopy() *ReconciliationViewDefinition {
	if in == nil {
		return nil
	}
	out := new(ReconciliationViewDefinition)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ReconciliationViewDefinition) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationViewDefinitionList) DeepCopyInto(out *ReconciliationViewDefinitionList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ReconciliationViewDefinition, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationViewDefinitionList.
func (in *ReconciliationViewDefinitionList) DeepCopy() *ReconciliationViewDefinitionList {
	if in == nil {
		return nil
	}
	out := new(ReconciliationViewDefinitionList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ReconciliationViewDefinitionList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationViewDefinitionSpec) DeepCopyInto(out *ReconciliationViewDefinitionSpec) {
	*out = *in
	if in.OwnershipRules != nil {
		in, out := &in.OwnershipRules, &out.OwnershipRules
		*out = make([]OwnershipRule, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	in.StateEvaluationExpression.DeepCopyInto(&out.StateEvaluationExpression)
	if in.I18nResourceRef != nil {
		in, out := &in.I18nResourceRef, &out.I18nResourceRef
		*out = new(ObjectReference)
		**out = **in
	}
	if in.Locale != nil {
		in, out := &in.Locale, &out.Locale
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationViewDefinitionSpec.
func (in *ReconciliationViewDefinitionSpec) DeepCopy() *ReconciliationViewDefinitionSpec {
	if in == nil {
		return nil
	}
	out := new(ReconciliationViewDefinitionSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationViewDefinitionStatus) DeepCopyInto(out *ReconciliationViewDefinitionStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationViewDefinitionStatus.
func (in *ReconciliationViewDefinitionStatus) DeepCopy() *ReconciliationViewDefinitionStatus {
	if in == nil {
		return nil
	}
	out := new(ReconciliationViewDefinitionStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationViewList) DeepCopyInto(out *ReconciliationViewList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ReconciliationView, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationViewList.
func (in *ReconciliationViewList) DeepCopy() *ReconciliationViewList {
	if in == nil {
		return nil
	}
	out := new(ReconciliationViewList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ReconciliationViewList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationViewSpec) DeepCopyInto(out *ReconciliationViewSpec) {
	*out = *in
	if in.TargetObject != nil {
		in, out := &in.TargetObject, &out.TargetObject
		*out = new(ObjectReference)
		**out = **in
	}
	if in.Depth != nil {
		in, out := &in.Depth, &out.Depth
		*out = new(int32)
		**out = **in
	}
	if in.StateEvaluationExpression != nil {
		in, out := &in.StateEvaluationExpression, &out.StateEvaluationExpression
		*out = new(StateEvaluationExpression)
		(*in).DeepCopyInto(*out)
	}
	if in.Locale != nil {
		in, out := &in.Locale, &out.Locale
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationViewSpec.
func (in *ReconciliationViewSpec) DeepCopy() *ReconciliationViewSpec {
	if in == nil {
		return nil
	}
	out := new(ReconciliationViewSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ReconciliationViewStatus) DeepCopyInto(out *ReconciliationViewStatus) {
	*out = *in
	if in.InitialObjectTree != nil {
		in, out := &in.InitialObjectTree, &out.InitialObjectTree
		*out = new(ObjectTreeNode)
		(*in).DeepCopyInto(*out)
	}
	if in.DesiredObjectTree != nil {
		in, out := &in.DesiredObjectTree, &out.DesiredObjectTree
		*out = new(ObjectTreeNode)
		(*in).DeepCopyInto(*out)
	}
	if in.CurrentObjectTree != nil {
		in, out := &in.CurrentObjectTree, &out.CurrentObjectTree
		*out = new(ObjectTreeNode)
		(*in).DeepCopyInto(*out)
	}
	in.PlanSummary.DeepCopyInto(&out.PlanSummary)
	in.ViewSummary.DeepCopyInto(&out.ViewSummary)
	if in.Plan != nil {
		in, out := &in.Plan, &out.Plan
		*out = make([]ObjectChange, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.View != nil {
		in, out := &in.View, &out.View
		*out = make([]ObjectChange, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ReconciliationViewStatus.
func (in *ReconciliationViewStatus) DeepCopy() *ReconciliationViewStatus {
	if in == nil {
		return nil
	}
	out := new(ReconciliationViewStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StateEvaluationExpression) DeepCopyInto(out *StateEvaluationExpression) {
	*out = *in
	if in.CELExpression != nil {
		in, out := &in.CELExpression, &out.CELExpression
		*out = new(CELExpression)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StateEvaluationExpression.
func (in *StateEvaluationExpression) DeepCopy() *StateEvaluationExpression {
	if in == nil {
		return nil
	}
	out := new(StateEvaluationExpression)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ViewSummary) DeepCopyInto(out *ViewSummary) {
	*out = *in
	if in.ObjectSummaries != nil {
		in, out := &in.ObjectSummaries, &out.ObjectSummaries
		*out = make([]ObjectSummary, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ViewSummary.
func (in *ViewSummary) DeepCopy() *ViewSummary {
	if in == nil {
		return nil
	}
	out := new(ViewSummary)
	in.DeepCopyInto(out)
	return out
}
