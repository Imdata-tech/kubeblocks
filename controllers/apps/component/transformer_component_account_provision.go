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

package component

import (
	"slices"
	"strings"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	"github.com/apecloud/kubeblocks/pkg/common"
	"github.com/apecloud/kubeblocks/pkg/constant"
	"github.com/apecloud/kubeblocks/pkg/controller/component"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
	"github.com/apecloud/kubeblocks/pkg/controller/lifecycle"
	"github.com/apecloud/kubeblocks/pkg/controller/model"
)

const (
	accountProvisionConditionType             = "SystemAccountProvision"
	accountProvisionConditionReasonInProgress = "InProgress"
	accountProvisionConditionReasonDone       = "AllProvisioned"
)

// componentAccountProvisionTransformer provisions component system accounts.
type componentAccountProvisionTransformer struct{}

var _ graph.Transformer = &componentAccountProvisionTransformer{}

func (t *componentAccountProvisionTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*componentTransformContext)
	if model.IsObjectDeleting(transCtx.ComponentOrig) {
		return nil
	}
	if common.IsCompactMode(transCtx.ComponentOrig.Annotations) {
		transCtx.V(1).Info("Component is in compact mode, no need to create component account related objects",
			"component", client.ObjectKeyFromObject(transCtx.ComponentOrig))
		return nil
	}

	if len(transCtx.SynthesizeComponent.SystemAccounts) == 0 {
		return nil
	}
	if transCtx.Component.Status.Phase != appsv1.RunningComponentPhase {
		return nil
	}
	// TODO: (good-first-issue) if the component's account is deleted by user, we should re-provision it
	cond, provisioned := t.isProvisioned(transCtx)
	if provisioned {
		return nil
	}

	lifecycleActions := transCtx.CompDef.Spec.LifecycleActions
	if lifecycleActions == nil || lifecycleActions.AccountProvision == nil {
		return nil
	}

	lfa, err := t.lifecycleAction(transCtx)
	if err != nil {
		return err
	}
	for _, account := range transCtx.SynthesizeComponent.SystemAccounts {
		// The secret of initAccount should be rendered into the config file,
		// or injected into the container through specific account&password environment variables name supported by the engine.
		// When the engine starts up, it will automatically load and create this account.
		if account.InitAccount {
			continue
		}
		if t.isAccountProvisioned(cond, account) {
			continue
		}
		if transCtx.SynthesizeComponent.Annotations[constant.RestoreFromBackupAnnotationKey] == "" {
			// TODO: restore account secret from backup.
			// provision account when the component is not recovered from backup
			if err = t.provisionAccount(transCtx, cond, lfa, account); err != nil {
				t.markProvisionAsFailed(transCtx, &cond, err)
				return err
			}
		}
		t.markAccountProvisioned(&cond, account)
	}
	t.markProvisioned(transCtx, cond)

	return nil
}

func (t *componentAccountProvisionTransformer) isProvisioned(transCtx *componentTransformContext) (metav1.Condition, bool) {
	for _, cond := range transCtx.Component.Status.Conditions {
		if cond.Type == accountProvisionConditionType {
			if cond.Status == metav1.ConditionTrue {
				return cond, true
			}
			return cond, false
		}
	}
	return metav1.Condition{
		Type:               accountProvisionConditionType,
		Status:             metav1.ConditionFalse,
		ObservedGeneration: transCtx.Component.Generation,
		LastTransitionTime: metav1.Now(),
		Reason:             accountProvisionConditionReasonInProgress,
		Message:            "",
	}, false
}

func (t *componentAccountProvisionTransformer) markProvisionAsFailed(transCtx *componentTransformContext, cond *metav1.Condition, err error) {
	cond.Status = metav1.ConditionFalse
	cond.ObservedGeneration = transCtx.Component.Generation
	cond.LastTransitionTime = metav1.Now()
	// cond.Reason = err.Error() // TODO: error
}

func (t *componentAccountProvisionTransformer) markProvisioned(transCtx *componentTransformContext, cond metav1.Condition) {
	cond.Status = metav1.ConditionTrue
	cond.ObservedGeneration = transCtx.Component.Generation
	cond.LastTransitionTime = metav1.Now()
	cond.Reason = accountProvisionConditionReasonDone

	conditions := transCtx.Component.Status.Conditions
	if conditions == nil {
		conditions = make([]metav1.Condition, 0)
	}
	existed := false
	for i, c := range conditions {
		if c.Type == cond.Type {
			existed = true
			conditions[i] = cond
		}
	}
	if !existed {
		conditions = append(conditions, cond)
	}
	transCtx.Component.Status.Conditions = conditions
}

func (t *componentAccountProvisionTransformer) isAccountProvisioned(cond metav1.Condition, account appsv1.SystemAccount) bool {
	if len(cond.Message) == 0 {
		return false
	}
	accounts := strings.Split(cond.Message, ",")
	return slices.Contains(accounts, account.Name)
}

func (t *componentAccountProvisionTransformer) markAccountProvisioned(cond *metav1.Condition, account appsv1.SystemAccount) {
	if len(cond.Message) == 0 {
		cond.Message = account.Name
		return
	}
	accounts := strings.Split(cond.Message, ",")
	if slices.Contains(accounts, account.Name) {
		return
	}
	accounts = append(accounts, account.Name)
	cond.Message = strings.Join(accounts, ",")
}

func (t *componentAccountProvisionTransformer) lifecycleAction(transCtx *componentTransformContext) (lifecycle.Lifecycle, error) {
	synthesizedComp := transCtx.SynthesizeComponent
	pods, err := component.ListOwnedPods(transCtx.Context, transCtx.Client,
		synthesizedComp.Namespace, synthesizedComp.ClusterName, synthesizedComp.Name)
	if err != nil {
		return nil, err
	}
	lfa, err := lifecycle.New(synthesizedComp.Namespace, synthesizedComp.ClusterName, synthesizedComp.Name,
		synthesizedComp.LifecycleActions, synthesizedComp.TemplateVars, nil, pods...)
	if err != nil {
		return nil, err
	}
	return lfa, nil
}

func (t *componentAccountProvisionTransformer) provisionAccount(transCtx *componentTransformContext,
	_ metav1.Condition, lfa lifecycle.Lifecycle, account appsv1.SystemAccount) error {

	synthesizedComp := transCtx.SynthesizeComponent
	secret, err := t.getAccountSecret(transCtx, synthesizedComp, account)
	if err != nil {
		return err
	}

	username, password := secret.Data[constant.AccountNameForSecret], secret.Data[constant.AccountPasswdForSecret]
	if len(username) == 0 || len(password) == 0 {
		return nil
	}

	err = lfa.AccountProvision(transCtx.Context, transCtx.Client, nil, account.Statement, string(username), string(password))
	return lifecycle.IgnoreNotDefined(err)
}

func (t *componentAccountProvisionTransformer) getAccountSecret(ctx graph.TransformContext,
	synthesizeComp *component.SynthesizedComponent, account appsv1.SystemAccount) (*corev1.Secret, error) {
	secretKey := types.NamespacedName{
		Namespace: synthesizeComp.Namespace,
		Name:      constant.GenerateAccountSecretName(synthesizeComp.ClusterName, synthesizeComp.Name, account.Name),
	}
	secret := &corev1.Secret{}
	if err := ctx.GetClient().Get(ctx.GetContext(), secretKey, secret); err != nil {
		return nil, err
	}
	return secret, nil
}
