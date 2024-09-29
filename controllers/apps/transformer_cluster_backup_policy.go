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

package apps

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-logr/logr"
	"golang.org/x/exp/slices"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/gengo/examples/set-gen/sets"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	dpv1alpha1 "github.com/apecloud/kubeblocks/apis/dataprotection/v1alpha1"
	workloads "github.com/apecloud/kubeblocks/apis/workloads/v1"
	"github.com/apecloud/kubeblocks/pkg/common"
	"github.com/apecloud/kubeblocks/pkg/constant"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
	"github.com/apecloud/kubeblocks/pkg/controller/model"
	intctrlutil "github.com/apecloud/kubeblocks/pkg/controllerutil"
	dptypes "github.com/apecloud/kubeblocks/pkg/dataprotection/types"
	dputils "github.com/apecloud/kubeblocks/pkg/dataprotection/utils"
	"github.com/apecloud/kubeblocks/pkg/dataprotection/utils/boolptr"
	viper "github.com/apecloud/kubeblocks/pkg/viperx"
)

const (
	defaultCronExpression = "0 18 * * *"
)

// clusterBackupPolicyTransformer transforms the backup policy template to the data protection backup policy and backup schedule.
type clusterBackupPolicyTransformer struct{}

type backupPolicyCtx struct {
	ctx     context.Context
	cli     client.Reader
	logger  logr.Logger
	event   record.EventRecorder
	cluster *appsv1.Cluster

	tplCount          int
	tplIdentifier     string
	isDefaultTemplate string

	backupPolicyTpl *appsv1alpha1.BackupPolicyTemplate
	backupPolicy    *appsv1alpha1.BackupPolicy
}

type componentItem struct {
	compSpec *appsv1.ClusterComponentSpec
	// shardingSpec.Name or componentSpec.Name
	componentName string
	isSharding    bool
	// componentSpec.Name or component name label which creates by shardingSpec.
	fullComponentName string
}

var _ graph.Transformer = &clusterBackupPolicyTransformer{}

// Transform transforms the backup policy template to the backup policy and backup schedule.
func (r *clusterBackupPolicyTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*clusterTransformContext)
	if model.IsObjectDeleting(transCtx.OrigCluster) {
		return nil
	}

	if common.IsCompactMode(transCtx.OrigCluster.Annotations) {
		transCtx.V(1).Info("Cluster is in compact mode, no need to create backup related objects",
			"cluster", client.ObjectKeyFromObject(transCtx.OrigCluster))
		return nil
	}

	graphCli, _ := transCtx.Client.(model.GraphClient)
	backupPolicyTPLs, err := r.getBackupPolicyTemplates(transCtx)
	if err != nil {
		return err
	}

	bpCtx := &backupPolicyCtx{
		ctx:      transCtx.Context,
		cli:      transCtx.Client,
		logger:   transCtx.Logger,
		event:    transCtx.EventRecorder,
		cluster:  transCtx.OrigCluster,
		tplCount: len(backupPolicyTPLs.Items),
	}
	if err := r.reconcileBackupPolicyTemplates(dag, graphCli, bpCtx, backupPolicyTPLs); err != nil {
		return err
	}
	return nil
}

// getBackupPolicyTemplates gets the backupPolicyTemplate for the cluster.
func (r *clusterBackupPolicyTransformer) getBackupPolicyTemplates(transCtx *clusterTransformContext) (*appsv1alpha1.BackupPolicyTemplateList, error) {
	backupPolicyTPLs := &appsv1alpha1.BackupPolicyTemplateList{}
	tplMap := map[string]sets.Empty{}
	for _, v := range transCtx.ComponentDefs {
		tmpTPLs := &appsv1alpha1.BackupPolicyTemplateList{}
		// TODO: prefix match for componentDef name?
		if err := transCtx.Client.List(transCtx.Context, tmpTPLs, client.MatchingLabels{v.Name: v.Name}); err != nil {
			return nil, err
		}
		for i := range tmpTPLs.Items {
			if _, ok := tplMap[tmpTPLs.Items[i].Name]; !ok {
				backupPolicyTPLs.Items = append(backupPolicyTPLs.Items, tmpTPLs.Items[i])
				tplMap[tmpTPLs.Items[i].Name] = sets.Empty{}
			}
		}
	}
	return backupPolicyTPLs, nil
}

func (r *clusterBackupPolicyTransformer) reconcileBackupPolicyTemplates(dag *graph.DAG, graphCli model.GraphClient, backupPolicyCtx *backupPolicyCtx,
	bptList *appsv1alpha1.BackupPolicyTemplateList) error {
	backupPolicyMap := map[string]struct{}{}
	backupScheduleMap := map[string]struct{}{}
	for _, tpl := range bptList.Items {
		backupPolicyCtx.isDefaultTemplate = tpl.Annotations[dptypes.DefaultBackupPolicyTemplateAnnotationKey]
		backupPolicyCtx.tplIdentifier = tpl.Spec.Identifier
		backupPolicyCtx.backupPolicyTpl = &tpl
		for i := range tpl.Spec.BackupPolicies {
			backupPolicyCtx.backupPolicy = &tpl.Spec.BackupPolicies[i]
			compItems := r.getClusterComponentItems(backupPolicyCtx)
			if err := reconcileBackupPolicyTemplate(dag, graphCli, backupPolicyCtx, compItems, backupPolicyMap, backupScheduleMap); err != nil {
				return err
			}
		}
	}
	return nil
}

func (r *clusterBackupPolicyTransformer) getClusterComponentItems(backupPolicyCtx *backupPolicyCtx) []componentItem {
	matchedCompDef := func(compSpec appsv1.ClusterComponentSpec) bool {
		// TODO: support to create bp when using cluster topology and componentDef is empty
		if len(compSpec.ComponentDef) == 0 {
			return false
		}
		for _, compDef := range backupPolicyCtx.backupPolicy.ComponentDefs {
			if strings.HasPrefix(compSpec.ComponentDef, compDef) || strings.HasPrefix(compDef, compSpec.ComponentDef) {
				return true
			}
		}
		return false
	}
	var compSpecItems []componentItem
	for i, v := range backupPolicyCtx.cluster.Spec.ComponentSpecs {
		if matchedCompDef(v) {
			compSpecItems = append(compSpecItems, componentItem{
				compSpec:          &backupPolicyCtx.cluster.Spec.ComponentSpecs[i],
				componentName:     v.Name,
				fullComponentName: v.Name,
			})
		}
	}
	return compSpecItems
}

func reconcileBackupPolicyTemplate(dag *graph.DAG, graphCli model.GraphClient, backupPolicyCtx *backupPolicyCtx,
	compItems []componentItem, backupPolicyMap, backupScheduleMap map[string]struct{}) error {
	transformBackupPolicyWrap := func(comp componentItem) *dpv1alpha1.BackupPolicy {
		// build the data protection backup policy from the template.
		oldBackupPolicy, newBackupPolicy := transformBackupPolicy(backupPolicyCtx, comp)
		if newBackupPolicy == nil {
			return nil
		}

		// if exist multiple backup policy templates and duplicate spec.identifier,
		// the generated backupPolicy may have duplicate names, so it is
		// necessary to check if it already exists.
		if _, ok := backupPolicyMap[newBackupPolicy.Name]; ok {
			return nil
		}

		if oldBackupPolicy == nil {
			graphCli.Create(dag, newBackupPolicy)
		} else {
			graphCli.Patch(dag, oldBackupPolicy, newBackupPolicy)
		}
		backupPolicyMap[newBackupPolicy.Name] = struct{}{}
		return newBackupPolicy
	}

	transformBackupScheduleWrap := func(comp componentItem, backupPolicy *dpv1alpha1.BackupPolicy, needMergeClusterBackup bool) {
		// if backup policy is nil, it means that the backup policy template
		// is invalid, backup schedule depends on backup policy, so we do
		// not need to transform backup schedule.
		if backupPolicy == nil {
			return
		}

		// only create backup schedule for the default backup policy template
		// if there are more than one backup policy templates.
		if backupPolicyCtx.isDefaultTemplate != trueVal && backupPolicyCtx.tplCount > 1 {
			backupPolicyCtx.logger.V(1).Info("Skip creating backup schedule for non-default backup policy template", "template", backupPolicyCtx.backupPolicyTpl.Name)
			return
		}

		// build the data protection backup schedule from the template.
		oldBackupSchedule, newBackupSchedule := transformBackupSchedule(backupPolicyCtx, comp, backupPolicy)

		// merge cluster backup configuration into the backup schedule.
		// If the backup schedule is nil, create a new backup schedule
		// based on the cluster backup configuration.
		// For a cluster, the default backup schedule is created by backup
		// policy template, user can also configure cluster backup in the
		// cluster custom object, such as enable cluster backup, set backup
		// schedule, etc.
		// We always prioritize the cluster backup configuration in the
		// cluster object, so we need to merge the cluster backup configuration
		// into the default backup schedule created by backup policy template
		// if it exists.
		if needMergeClusterBackup {
			newBackupSchedule = mergeClusterBackup(backupPolicyCtx, comp, backupPolicy, newBackupSchedule)
		}
		if newBackupSchedule == nil {
			return
		}
		// if exist multiple backup policy templates and duplicate spec.identifier,
		// the backupSchedule that may be generated may have duplicate names,
		// and it is necessary to check if it already exists.
		if _, ok := backupScheduleMap[newBackupSchedule.Name]; ok {
			return
		}

		if oldBackupSchedule == nil {
			graphCli.Create(dag, newBackupSchedule)
		} else {
			graphCli.Patch(dag, oldBackupSchedule, newBackupSchedule)
		}
		graphCli.DependOn(dag, backupPolicy, newBackupSchedule)
		comps := graphCli.FindAll(dag, &appsv1.Component{})
		graphCli.DependOn(dag, backupPolicy, comps...)
		backupScheduleMap[newBackupSchedule.Name] = struct{}{}
	}

	// transform backup policy template to data protection backupPolicy and backupSchedule
	for j, v := range compItems {
		policy := transformBackupPolicyWrap(v)
		// only merge the first backupSchedule for the cluster backup.
		transformBackupScheduleWrap(v, policy, j == 0)
	}

	return nil
}

// transformBackupPolicy transforms backup policy template to backup policy.
func transformBackupPolicy(backupPolicyCtx *backupPolicyCtx, comp componentItem) (*dpv1alpha1.BackupPolicy, *dpv1alpha1.BackupPolicy) {
	backupPolicyName := generateBackupPolicyName(backupPolicyCtx.cluster.Name, comp.componentName, backupPolicyCtx.tplIdentifier)
	backupPolicy := &dpv1alpha1.BackupPolicy{}
	if err := backupPolicyCtx.cli.Get(backupPolicyCtx.ctx, client.ObjectKey{
		Namespace: backupPolicyCtx.cluster.Namespace,
		Name:      backupPolicyName,
	}, backupPolicy); client.IgnoreNotFound(err) != nil {
		backupPolicyCtx.logger.Error(err, "failed to get backup policy", "backupPolicy", backupPolicyName)
		return nil, nil
	}

	if len(backupPolicy.Name) == 0 {
		// build a new backup policy by the backup policy template.
		return nil, buildBackupPolicy(backupPolicyCtx, comp, backupPolicyName)
	}

	// sync the existing backup policy with the cluster changes
	old := backupPolicy.DeepCopy()
	syncBackupPolicy(backupPolicyCtx, comp, backupPolicy)
	return old, backupPolicy
}

func transformBackupSchedule(backupPolicyCtx *backupPolicyCtx, comp componentItem, backupPolicy *dpv1alpha1.BackupPolicy) (*dpv1alpha1.BackupSchedule, *dpv1alpha1.BackupSchedule) {
	scheduleName := generateBackupScheduleName(backupPolicyCtx.cluster.Name, comp.componentName, backupPolicyCtx.tplIdentifier)
	backupSchedule := &dpv1alpha1.BackupSchedule{}
	if err := backupPolicyCtx.cli.Get(backupPolicyCtx.ctx, client.ObjectKey{
		Namespace: backupPolicyCtx.cluster.Namespace,
		Name:      scheduleName,
	}, backupSchedule); client.IgnoreNotFound(err) != nil {
		backupPolicyCtx.logger.Error(err, "failed to get backup schedule", "backupSchedule", scheduleName)
		return nil, nil
	}

	// build a new backup schedule from the backup policy template.
	if len(backupSchedule.Name) == 0 {
		return nil, buildBackupSchedule(backupPolicyCtx, comp, scheduleName, backupPolicy)
	}

	old := backupSchedule.DeepCopy()
	syncBackupSchedule(backupPolicyCtx, backupSchedule)
	return old, backupSchedule
}

// syncBackupPolicy syncs labels and annotations of the backup policy with the cluster changes.
func syncBackupPolicy(backupPolicyCtx *backupPolicyCtx, comp componentItem, backupPolicy *dpv1alpha1.BackupPolicy) {
	// update labels and annotations of the backup policy.
	if backupPolicy.Annotations == nil {
		backupPolicy.Annotations = map[string]string{}
	}
	if backupPolicy.Labels == nil {
		backupPolicy.Labels = map[string]string{}
	}
	mergeMap(backupPolicy.Annotations, buildAnnotations(backupPolicyCtx))
	mergeMap(backupPolicy.Labels, buildLabels(backupPolicyCtx, comp, nil))

	// update backup repo of the backup policy.
	if backupPolicyCtx.cluster.Spec.Backup != nil && backupPolicyCtx.cluster.Spec.Backup.RepoName != "" {
		backupPolicy.Spec.BackupRepoName = &backupPolicyCtx.cluster.Spec.Backup.RepoName
	}
	backupPolicy.Spec.BackoffLimit = backupPolicyCtx.backupPolicy.BackoffLimit
	syncBackupMethods(backupPolicyCtx, backupPolicy, comp)
	syncBackupPolicyTargetSpec(backupPolicyCtx, backupPolicy, comp)
}

func setDefaultEncryptionConfig(backupPolicyCtx *backupPolicyCtx, backupPolicy *dpv1alpha1.BackupPolicy) {
	secretKeyRefJSON := viper.GetString(constant.CfgKeyDPBackupEncryptionSecretKeyRef)
	if secretKeyRefJSON == "" {
		return
	}
	secretKeyRef := &corev1.SecretKeySelector{}
	err := json.Unmarshal([]byte(secretKeyRefJSON), secretKeyRef)
	if err != nil {
		backupPolicyCtx.logger.Error(err, "failed to unmarshal secretKeyRef", "json", secretKeyRefJSON)
		return
	}
	if secretKeyRef.Name == "" || secretKeyRef.Key == "" {
		return
	}
	algorithm := viper.GetString(constant.CfgKeyDPBackupEncryptionAlgorithm)
	if algorithm == "" {
		algorithm = dpv1alpha1.DefaultEncryptionAlgorithm
	}
	backupPolicy.Spec.EncryptionConfig = &dpv1alpha1.EncryptionConfig{
		Algorithm:              algorithm,
		PassPhraseSecretKeyRef: secretKeyRef,
	}
}

func buildBackupSchedule(backupPolicyCtx *backupPolicyCtx, comp componentItem, name string, backupPolicy *dpv1alpha1.BackupPolicy) *dpv1alpha1.BackupSchedule {
	if len(backupPolicyCtx.backupPolicy.Schedules) == 0 {
		return nil
	}
	backupSchedule := &dpv1alpha1.BackupSchedule{
		ObjectMeta: metav1.ObjectMeta{
			Name:        name,
			Namespace:   backupPolicyCtx.cluster.Namespace,
			Labels:      buildLabels(backupPolicyCtx, comp, backupPolicy),
			Annotations: buildAnnotations(backupPolicyCtx),
		},
		Spec: dpv1alpha1.BackupScheduleSpec{
			BackupPolicyName: backupPolicy.Name,
		},
	}

	var schedules []dpv1alpha1.SchedulePolicy
	for _, s := range backupPolicyCtx.backupPolicy.Schedules {
		schedules = append(schedules, dpv1alpha1.SchedulePolicy{
			BackupMethod:    s.BackupMethod,
			CronExpression:  s.CronExpression,
			Enabled:         s.Enabled,
			RetentionPeriod: s.RetentionPeriod,
		})
	}
	backupSchedule.Spec.Schedules = schedules
	return backupSchedule
}

func syncBackupSchedule(backupPolicyCtx *backupPolicyCtx, backupSchedule *dpv1alpha1.BackupSchedule) {
	scheduleMethodMap := map[string]struct{}{}
	for _, s := range backupSchedule.Spec.Schedules {
		scheduleMethodMap[s.BackupMethod] = struct{}{}
	}
	mergeMap(backupSchedule.Annotations, buildAnnotations(backupPolicyCtx))
	// update backupSchedule annotation to reconcile it.
	backupSchedule.Annotations[constant.ReconcileAnnotationKey] = backupPolicyCtx.cluster.ResourceVersion
	// sync the newly added schedule policies.
	for _, s := range backupPolicyCtx.backupPolicy.Schedules {
		if _, ok := scheduleMethodMap[s.BackupMethod]; ok {
			continue
		}
		backupSchedule.Spec.Schedules = append(backupSchedule.Spec.Schedules, dpv1alpha1.SchedulePolicy{
			BackupMethod:    s.BackupMethod,
			CronExpression:  s.CronExpression,
			Enabled:         s.Enabled,
			RetentionPeriod: s.RetentionPeriod,
		})
	}
}

func syncRoleLabelSelector(backupPolicyCtx *backupPolicyCtx, comp componentItem, target *dpv1alpha1.BackupTarget, role, alternateRole string) {
	if len(role) == 0 || target == nil {
		return
	}
	podSelector := target.PodSelector
	if podSelector.LabelSelector == nil || podSelector.LabelSelector.MatchLabels == nil {
		podSelector.LabelSelector = &metav1.LabelSelector{MatchLabels: map[string]string{}}
	}
	if getCompReplicas(backupPolicyCtx, comp) == 1 {
		delete(podSelector.LabelSelector.MatchLabels, constant.RoleLabelKey)
		if podSelector.FallbackLabelSelector != nil && podSelector.FallbackLabelSelector.MatchLabels != nil {
			delete(podSelector.FallbackLabelSelector.MatchLabels, constant.RoleLabelKey)
		}
	} else {
		podSelector.LabelSelector.MatchLabels[constant.RoleLabelKey] = role
		if len(alternateRole) > 0 {
			if podSelector.FallbackLabelSelector == nil || podSelector.FallbackLabelSelector.MatchLabels == nil {
				podSelector.FallbackLabelSelector = &metav1.LabelSelector{MatchLabels: map[string]string{}}
			}
			podSelector.FallbackLabelSelector.MatchLabels[constant.RoleLabelKey] = alternateRole
		}
	}
}

func getCompReplicas(backupPolicyCtx *backupPolicyCtx, comp componentItem) int32 {
	its := &workloads.InstanceSet{}
	name := fmt.Sprintf("%s-%s", backupPolicyCtx.cluster.Name, comp.fullComponentName)
	if err := backupPolicyCtx.cli.Get(backupPolicyCtx.ctx, client.ObjectKey{Name: name, Namespace: backupPolicyCtx.cluster.Namespace}, its); err != nil {
		return comp.compSpec.Replicas
	}
	return *its.Spec.Replicas
}

// buildBackupPolicy builds a new backup policy by the backup policy template.
func buildBackupPolicy(backupPolicyCtx *backupPolicyCtx, comp componentItem, backupPolicyName string) *dpv1alpha1.BackupPolicy {
	backupPolicy := &dpv1alpha1.BackupPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:        backupPolicyName,
			Namespace:   backupPolicyCtx.cluster.Namespace,
			Labels:      buildLabels(backupPolicyCtx, comp, nil),
			Annotations: buildAnnotations(backupPolicyCtx),
		},
	}
	syncBackupMethods(backupPolicyCtx, backupPolicy, comp)
	bpSpec := backupPolicy.Spec
	// if cluster have backup repo, set backup repo name to backup policy.
	if backupPolicyCtx.cluster.Spec.Backup != nil && backupPolicyCtx.cluster.Spec.Backup.RepoName != "" {
		bpSpec.BackupRepoName = &backupPolicyCtx.cluster.Spec.Backup.RepoName
	}
	bpSpec.PathPrefix = buildBackupPathPrefix(backupPolicyCtx.cluster, comp.componentName)
	bpSpec.BackoffLimit = backupPolicyCtx.backupPolicy.BackoffLimit
	backupPolicy.Spec = bpSpec
	setDefaultEncryptionConfig(backupPolicyCtx, backupPolicy)
	syncBackupPolicyTargetSpec(backupPolicyCtx, backupPolicy, comp)
	return backupPolicy
}

// syncBackupMethods syncs the backupMethod of tpl to backupPolicy.
func syncBackupMethods(backupPolicyCtx *backupPolicyCtx, backupPolicy *dpv1alpha1.BackupPolicy, comp componentItem) {
	var backupMethods []dpv1alpha1.BackupMethod
	oldBackupMethodMap := map[string]dpv1alpha1.BackupMethod{}
	for _, v := range backupPolicy.Spec.BackupMethods {
		oldBackupMethodMap[v.Name] = v
	}
	for _, v := range backupPolicyCtx.backupPolicy.BackupMethods {
		backupMethod := v.BackupMethod
		if m, ok := oldBackupMethodMap[backupMethod.Name]; ok {
			backupMethod = m
			delete(oldBackupMethodMap, backupMethod.Name)
		} else if v.Target != nil {
			if comp.isSharding {
				backupMethod.Targets = buildBackupTargets(backupPolicyCtx, backupMethod.Targets, comp)
			} else {
				backupMethod.Target = buildBackupTarget(backupPolicyCtx, backupMethod.Target, *v.Target, comp)
			}
		}
		mappingEnv := doEnvMapping(comp.compSpec, v.EnvMapping)
		backupMethod.Env = dputils.MergeEnv(backupMethod.Env, mappingEnv)
		backupMethods = append(backupMethods, backupMethod)
	}
	for _, v := range oldBackupMethodMap {
		backupMethods = append(backupMethods, v)
	}
	backupPolicy.Spec.BackupMethods = backupMethods
}

func doEnvMapping(comp *appsv1.ClusterComponentSpec, envMapping []appsv1alpha1.EnvMappingVar) []corev1.EnvVar {
	var env []corev1.EnvVar
	for _, v := range envMapping {
		for _, cm := range v.ValueFrom.ComponentDef {
			if !slices.Contains(cm.Names, comp.ComponentDef) {
				continue
			}
			env = append(env, corev1.EnvVar{
				Name:  v.Key,
				Value: cm.MappingValue,
			})
		}
	}
	return env
}

func syncBackupPolicyTargetSpec(backupPolicyCtx *backupPolicyCtx, backupPolicy *dpv1alpha1.BackupPolicy, comp componentItem) {
	if comp.isSharding {
		backupPolicy.Spec.Targets = buildBackupTargets(backupPolicyCtx, backupPolicy.Spec.Targets, comp)
	} else {
		backupPolicy.Spec.Target = buildBackupTarget(backupPolicyCtx, backupPolicy.Spec.Target, backupPolicyCtx.backupPolicy.Target, comp)
	}
}

func buildBackupTargets(backupPolicyCtx *backupPolicyCtx, targets []dpv1alpha1.BackupTarget, comp componentItem) []dpv1alpha1.BackupTarget {
	shardComponents, _ := intctrlutil.ListShardingComponents(backupPolicyCtx.ctx, backupPolicyCtx.cli, backupPolicyCtx.cluster, comp.componentName)
	sourceTargetMap := map[string]*dpv1alpha1.BackupTarget{}
	for i := range targets {
		sourceTargetMap[targets[i].Name] = &targets[i]
	}
	var backupTargets []dpv1alpha1.BackupTarget
	for _, v := range shardComponents {
		// set ClusterComponentSpec name to component name
		comp.fullComponentName = v.Labels[constant.KBAppComponentLabelKey]
		target := buildBackupTarget(backupPolicyCtx, sourceTargetMap[comp.fullComponentName], backupPolicyCtx.backupPolicy.Target, comp)
		if target != nil {
			backupTargets = append(backupTargets, *target)
		}
	}
	return backupTargets
}

func buildBackupTarget(backupPolicyCtx *backupPolicyCtx, oldTarget *dpv1alpha1.BackupTarget,
	targetTpl appsv1alpha1.TargetInstance, comp componentItem) *dpv1alpha1.BackupTarget {
	if oldTarget != nil {
		// if the target already exists, only sync the role by component replicas automatically.
		syncRoleLabelSelector(backupPolicyCtx, comp, oldTarget, targetTpl.Role, targetTpl.FallbackRole)
		return oldTarget
	}
	clusterName := backupPolicyCtx.cluster.Name
	if targetTpl.Strategy == "" {
		targetTpl.Strategy = dpv1alpha1.PodSelectionStrategyAny
	}
	target := &dpv1alpha1.BackupTarget{
		PodSelector: &dpv1alpha1.PodSelector{
			Strategy: targetTpl.Strategy,
			LabelSelector: &metav1.LabelSelector{
				MatchLabels: buildTargetPodLabels(backupPolicyCtx, targetTpl.Role, comp),
			},
		},
		// dataProtection will use its dedicated service account if this field is empty.
		ServiceAccountName: "",
		ContainerPort:      targetTpl.ContainerPort,
	}
	if len(targetTpl.Role) != 0 && len(targetTpl.FallbackRole) != 0 {
		target.PodSelector.FallbackLabelSelector = &metav1.LabelSelector{
			MatchLabels: buildTargetPodLabels(backupPolicyCtx, targetTpl.FallbackRole, comp),
		}
	}
	if comp.isSharding {
		target.Name = comp.fullComponentName
	}
	// build the target connection credential
	if targetTpl.Account != "" {
		target.ConnectionCredential = &dpv1alpha1.ConnectionCredential{
			SecretName:  constant.GenerateAccountSecretName(clusterName, comp.fullComponentName, targetTpl.Account),
			PasswordKey: constant.AccountPasswdForSecret,
			UsernameKey: constant.AccountNameForSecret,
		}
	}
	return target
}

func mergeClusterBackup(backupPolicyCtx *backupPolicyCtx, comp componentItem, backupPolicy *dpv1alpha1.BackupPolicy, backupSchedule *dpv1alpha1.BackupSchedule) *dpv1alpha1.BackupSchedule {
	backupEnabled := func() bool {
		return backupPolicyCtx.cluster.Spec.Backup != nil && boolValue(backupPolicyCtx.cluster.Spec.Backup.Enabled)
	}

	if backupPolicy == nil || backupPolicyCtx.cluster.Spec.Backup == nil {
		// backup policy is nil, can not enable cluster backup, so record event and return.
		if backupEnabled() {
			backupPolicyCtx.event.Event(backupPolicyCtx.cluster, corev1.EventTypeWarning,
				"BackupPolicyNotFound", "backup policy is nil, can not enable cluster backup")
		}
		return backupSchedule
	}

	backup := backupPolicyCtx.cluster.Spec.Backup
	method := dputils.GetBackupMethodByName(backup.Method, backupPolicy)
	// the specified backup method should be in the backup policy, if not, record event and return.
	if method == nil {
		backupPolicyCtx.event.Event(backupPolicyCtx.cluster, corev1.EventTypeWarning,
			"BackupMethodNotFound", fmt.Sprintf("backup method %s is not found in backup policy", backup.Method))
		return backupSchedule
	}

	// there is no backup schedule created by backup policy template, so we need to
	// create a new backup schedule for cluster backup.
	if backupSchedule == nil {
		backupSchedule = &dpv1alpha1.BackupSchedule{
			ObjectMeta: metav1.ObjectMeta{
				Name:        generateBackupScheduleName(backupPolicyCtx.cluster.Name, comp.componentName, backupPolicyCtx.tplIdentifier),
				Namespace:   backupPolicyCtx.cluster.Namespace,
				Labels:      buildLabels(backupPolicyCtx, comp, backupPolicy),
				Annotations: buildAnnotations(backupPolicyCtx),
			},
			Spec: dpv1alpha1.BackupScheduleSpec{
				BackupPolicyName:        backupPolicy.Name,
				StartingDeadlineMinutes: backup.StartingDeadlineMinutes,
				Schedules:               []dpv1alpha1.SchedulePolicy{},
			},
		}
	}

	// build backup schedule policy by cluster backup spec
	sp := &dpv1alpha1.SchedulePolicy{
		Enabled:         backup.Enabled,
		RetentionPeriod: backup.RetentionPeriod,
		BackupMethod:    backup.Method,
		CronExpression:  backup.CronExpression,
	}

	// merge cluster backup schedule policy into backup schedule, if the backup
	// schedule with specified method already exists, we need to update it
	// using the cluster backup schedule policy. Otherwise, we need to append
	// it to the backup schedule.
	// If cluster backup method is changed, we need to disable previous backup
	// method, for instance, the method is changed from A to B, we need to
	// disable A and enable B.
	exist := false
	hasSyncPITRMethod := false
	enableAutoBackup := boolptr.IsSetToTrue(backup.Enabled)
	for i, s := range backupSchedule.Spec.Schedules {
		if s.BackupMethod == backup.Method {
			mergeSchedulePolicy(sp, &backupSchedule.Spec.Schedules[i])
			exist = true
			continue
		}

		m := dputils.GetBackupMethodByName(s.BackupMethod, backupPolicy)
		if m == nil {
			continue
		}
		if m.ActionSetName == "" {
			if boolptr.IsSetToTrue(m.SnapshotVolumes) && enableAutoBackup {
				// disable the automatic backup when the specified method is not a volume snapshot for volume-snapshot method
				backupSchedule.Spec.Schedules[i].Enabled = boolptr.False()
			}
			continue
		}

		as := &dpv1alpha1.ActionSet{}
		if err := backupPolicyCtx.cli.Get(backupPolicyCtx.ctx, client.ObjectKey{Name: m.ActionSetName}, as); err != nil {
			backupPolicyCtx.logger.Error(err, "failed to get ActionSet for backup.", "ActionSet", as.Name)
			continue
		}
		if as.Spec.BackupType == dpv1alpha1.BackupTypeContinuous && backup.PITREnabled != nil && !hasSyncPITRMethod {
			// auto-sync the first continuous backup for the 'pirtEnable' option.
			backupSchedule.Spec.Schedules[i].Enabled = backup.PITREnabled
			hasSyncPITRMethod = true
		}
		if as.Spec.BackupType == dpv1alpha1.BackupTypeFull && enableAutoBackup {
			// disable the automatic backup for other full backup method
			backupSchedule.Spec.Schedules[i].Enabled = boolptr.False()
		}
	}
	if !exist {
		if sp.CronExpression == "" {
			sp.CronExpression = defaultCronExpression
		}
		backupSchedule.Spec.Schedules = append(backupSchedule.Spec.Schedules, *sp)
	}
	return backupSchedule
}

func defaultPolicyAnnotationValue(backupPolicyCtx *backupPolicyCtx) string {
	if backupPolicyCtx.tplCount > 1 && backupPolicyCtx.isDefaultTemplate != trueVal {
		return "false"
	}
	return trueVal
}

func buildAnnotations(backupPolicyCtx *backupPolicyCtx) map[string]string {
	annotations := map[string]string{
		dptypes.DefaultBackupPolicyAnnotationKey:   defaultPolicyAnnotationValue(backupPolicyCtx),
		constant.BackupPolicyTemplateAnnotationKey: backupPolicyCtx.backupPolicyTpl.Name,
	}
	if backupPolicyCtx.backupPolicyTpl.Annotations[dptypes.ReconfigureRefAnnotationKey] != "" {
		annotations[dptypes.ReconfigureRefAnnotationKey] = backupPolicyCtx.backupPolicyTpl.Annotations[dptypes.ReconfigureRefAnnotationKey]
	}
	return annotations
}

func buildLabels(backupPolicyCtx *backupPolicyCtx, compItem componentItem, policy *dpv1alpha1.BackupPolicy) map[string]string {
	labels := map[string]string{
		constant.AppManagedByLabelKey:        constant.AppName,
		constant.AppInstanceLabelKey:         backupPolicyCtx.cluster.Name,
		constant.ComponentDefinitionLabelKey: compDefName(compItem.compSpec, policy),
	}
	if compItem.isSharding {
		labels[constant.KBAppShardingNameLabelKey] = compItem.componentName
	} else {
		labels[constant.KBAppComponentLabelKey] = compItem.componentName
	}
	return labels
}

func compDefName(comp *appsv1.ClusterComponentSpec, policy *dpv1alpha1.BackupPolicy) string {
	switch {
	case comp != nil:
		return compDefNameFromSpec(comp)
	case policy != nil:
		return compDefNameFromPolicy(policy)
	default:
		panic("runtime error - unexpected way to get component definition name")
	}
}

func compDefNameFromSpec(comp *appsv1.ClusterComponentSpec) string {
	return comp.ComponentDef
}

func compDefNameFromPolicy(policy *dpv1alpha1.BackupPolicy) string {
	compDefName := ""
	if policy.Labels != nil {
		compDefName = policy.Labels[constant.ComponentDefinitionLabelKey]
	}
	return compDefName
}

// buildTargetPodLabels builds the target labels for the backup policy that will be
// used to select the target pod.
func buildTargetPodLabels(backupPolicyCtx *backupPolicyCtx, role string, comp componentItem) map[string]string {
	labels := map[string]string{
		constant.AppInstanceLabelKey:    backupPolicyCtx.cluster.Name,
		constant.AppManagedByLabelKey:   constant.AppName,
		constant.KBAppComponentLabelKey: comp.fullComponentName,
	}
	// append label to filter specific role of the component.
	if len(role) > 0 && getCompReplicas(backupPolicyCtx, comp) > 1 {
		// the role only works when the component has multiple replicas.
		labels[constant.RoleLabelKey] = role
	}
	if comp.isSharding {
		labels[constant.KBAppShardingNameLabelKey] = comp.componentName
	}
	return labels
}

// generateBackupPolicyName generates the backup policy name which is created from backup policy template.
func generateBackupPolicyName(clusterName, componentName, identifier string) string {
	if len(identifier) == 0 {
		return fmt.Sprintf("%s-%s-backup-policy", clusterName, componentName)
	}
	return fmt.Sprintf("%s-%s-backup-policy-%s", clusterName, componentName, identifier)
}

// generateBackupScheduleName generates the backup schedule name which is created from backup policy template.
func generateBackupScheduleName(clusterName, componentName, identifier string) string {
	if len(identifier) == 0 {
		return fmt.Sprintf("%s-%s-backup-schedule", clusterName, componentName)
	}
	return fmt.Sprintf("%s-%s-backup-schedule-%s", clusterName, componentName, identifier)
}

func buildBackupPathPrefix(cluster *appsv1.Cluster, compName string) string {
	return fmt.Sprintf("/%s-%s/%s", cluster.Name, cluster.UID, compName)
}

func mergeSchedulePolicy(src *dpv1alpha1.SchedulePolicy, dst *dpv1alpha1.SchedulePolicy) {
	if src.Enabled != nil {
		dst.Enabled = src.Enabled
	}
	if src.RetentionPeriod.String() != "" {
		dst.RetentionPeriod = src.RetentionPeriod
	}
	if src.BackupMethod != "" {
		dst.BackupMethod = src.BackupMethod
	}
	if src.CronExpression != "" {
		dst.CronExpression = src.CronExpression
	}
}
