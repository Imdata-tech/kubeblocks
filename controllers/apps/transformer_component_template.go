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
	"encoding/json"
	"fmt"
	"reflect"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"golang.org/x/exp/maps"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	workloads "github.com/apecloud/kubeblocks/apis/workloads/v1"
	"github.com/apecloud/kubeblocks/pkg/constant"
	"github.com/apecloud/kubeblocks/pkg/controller/builder"
	"github.com/apecloud/kubeblocks/pkg/controller/component"
	"github.com/apecloud/kubeblocks/pkg/controller/component/lifecycle"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
	"github.com/apecloud/kubeblocks/pkg/controller/model"
)

const (
	kubeBlockFileTemplate = "kubeblock-file-templates"
)

type componentFileTemplateTransformer struct{}

var _ graph.Transformer = &componentFileTemplateTransformer{}

func (t *componentFileTemplateTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*componentTransformContext)
	if model.IsObjectDeleting(transCtx.ComponentOrig) {
		return nil
	}

	runningObjs, err := t.getTemplateObjects(transCtx)
	if err != nil {
		return err
	}

	protoObjs, err := t.buildTemplateObjects(transCtx)
	if err != nil {
		return err
	}

	if err = t.transform(transCtx, runningObjs, protoObjs); err != nil {
		return err
	}

	toCreate, toDelete, toUpdate := mapDiff(runningObjs, protoObjs)
	graphCli, _ := transCtx.Client.(model.GraphClient)
	for name := range toCreate {
		graphCli.Create(dag, protoObjs[name], inDataContext4G())
	}
	for name := range toDelete {
		graphCli.Delete(dag, runningObjs[name], inDataContext4G())
	}
	for name := range toUpdate {
		runningObj, protoObj := runningObjs[name], protoObjs[name]
		if !reflect.DeepEqual(runningObj.Data, protoObj.Data) ||
			!reflect.DeepEqual(runningObj.Labels, protoObj.Labels) ||
			!reflect.DeepEqual(runningObj.Annotations, protoObj.Annotations) {
			graphCli.Update(dag, runningObj, protoObj, inDataContext4G())
		}
	}
	return nil
}

func (t *componentFileTemplateTransformer) templateObjectName(transCtx *componentTransformContext, tplName string) string {
	synthesizedComp := transCtx.SynthesizeComponent
	return fmt.Sprintf("%s-%s", synthesizedComp.FullCompName, tplName)
}

func (t *componentFileTemplateTransformer) getTemplateObjects(transCtx *componentTransformContext) (map[string]*corev1.ConfigMap, error) {
	synthesizedComp := transCtx.SynthesizeComponent

	opts := []client.ListOption{
		client.MatchingLabels(
			constant.GetCompLabels(synthesizedComp.ClusterName, synthesizedComp.Name, map[string]string{
				constant.KBAppFileTemplateLabelKey: "true",
			})),
		client.InNamespace(synthesizedComp.Namespace),
	}

	cmList := &corev1.ConfigMapList{}
	if err := transCtx.Client.List(transCtx.Context, cmList, opts...); err != nil {
		return nil, err
	}

	objs := make(map[string]*corev1.ConfigMap)
	for i, obj := range cmList.Items {
		objs[obj.Name] = &cmList.Items[i]
	}
	return objs, nil
}

func (t *componentFileTemplateTransformer) buildTemplateObjects(transCtx *componentTransformContext) (map[string]*corev1.ConfigMap, error) {
	objs := make(map[string]*corev1.ConfigMap)
	for _, tpl := range transCtx.CompDef.Spec.FileTemplates {
		obj, err := t.buildTemplateObject(transCtx, tpl)
		if err != nil {
			return nil, err
		}
		objs[obj.Name] = obj
	}
	return objs, nil
}

func (t *componentFileTemplateTransformer) buildTemplateObject(
	transCtx *componentTransformContext, tpl appsv1.ComponentFileTemplate) (*corev1.ConfigMap, error) {
	var (
		compDef         = transCtx.CompDef
		synthesizedComp = transCtx.SynthesizeComponent
	)
	data, err := t.buildTemplateData(transCtx, tpl)
	if err != nil {
		return nil, err
	}

	obj := builder.NewConfigMapBuilder(synthesizedComp.Name, t.templateObjectName(transCtx, tpl.Name)).
		AddLabelsInMap(synthesizedComp.StaticLabels).
		AddLabelsInMap(constant.GetCompLabelsWithDef(synthesizedComp.ClusterName, synthesizedComp.Name, compDef.Name)).
		AddLabels(constant.KBAppFileTemplateLabelKey, "true").
		AddAnnotationsInMap(synthesizedComp.StaticAnnotations).
		SetData(data).
		GetObject()

	if err = t.templateToObject(tpl, obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func (t *componentFileTemplateTransformer) namedKey(name, key string) string {
	return fmt.Sprintf("%s-%s", name, key)
}

func (t *componentFileTemplateTransformer) nameNKey(namedKey string) (string, string) {
	i := strings.Index(namedKey, "-") // TODO: error
	return namedKey[:i], namedKey[i+1:]
}

func (t *componentFileTemplateTransformer) buildTemplateData(
	transCtx *componentTransformContext, tpl appsv1.ComponentFileTemplate) (map[string]string, error) {
	// TODO: cluster API
	if len(tpl.Template) == 0 {
		return nil, nil
	}

	cmObj, err := func() (*corev1.ConfigMap, error) {
		cm := &corev1.ConfigMap{}
		cmKey := types.NamespacedName{
			Namespace: func() string {
				if len(tpl.Namespace) > 0 {
					return tpl.Namespace
				}
				return "default"
			}(),
			Name: tpl.Template,
		}
		if err := transCtx.Client.Get(transCtx.Context, cmKey, cm); err != nil {
			return nil, err
		}
		return cm, nil
	}()
	if err != nil {
		return nil, err
	}

	return t.renderTemplateData(transCtx, tpl, cmObj.Data)
}

func (t *componentFileTemplateTransformer) renderTemplateData(transCtx *componentTransformContext,
	fileTemplate appsv1.ComponentFileTemplate, data map[string]string) (map[string]string, error) {
	var (
		synthesizedComp = transCtx.SynthesizeComponent
		rendered        = make(map[string]string)
	)
	tpl := template.New(fileTemplate.Name).Option("missingkey=error").Funcs(sprig.TxtFuncMap())
	for key, val := range data {
		ptpl, err := tpl.Parse(val)
		if err != nil {
			return nil, err
		}
		var buf strings.Builder
		if err = ptpl.Execute(&buf, synthesizedComp.TemplateVars); err != nil {
			return nil, err
		}
		rendered[key] = buf.String()
	}
	return rendered, nil
}

func (t *componentFileTemplateTransformer) templateToObject(tpl appsv1.ComponentFileTemplate, obj *corev1.ConfigMap) error {
	data, err := json.Marshal(tpl)
	if err != nil {
		return err
	}
	if obj.Annotations == nil {
		obj.Annotations = make(map[string]string)
	}
	obj.Annotations[kubeBlockFileTemplate] = string(data)
	return nil
}

func (t *componentFileTemplateTransformer) templatesFromObject(obj *corev1.ConfigMap) (*appsv1.ComponentFileTemplate, error) {
	data := obj.Annotations[kubeBlockFileTemplate]
	if len(data) == 0 {
		return nil, nil
	}
	var tpl appsv1.ComponentFileTemplate
	if err := json.Unmarshal([]byte(data), &tpl); err != nil {
		return nil, err
	}
	return &tpl, nil
}

func (t *componentFileTemplateTransformer) transform(
	transCtx *componentTransformContext, runningObjs, protoObjs map[string]*corev1.ConfigMap) error {
	var (
		compDef         = transCtx.CompDef
		synthesizedComp = transCtx.SynthesizeComponent
	)

	if err := t.transformPodVolumes(transCtx); err != nil {
		return err
	}

	if compDef.Spec.LifecycleActions == nil || compDef.Spec.LifecycleActions.Reconfigure == nil {
		return nil
	}

	// if !reflect.DeepEqual(podSpecCopy.Volumes, synthesizedComp.PodSpec.Volumes) {
	//	// volumes changed, and the workload will be restarted
	//	return t.cancelQueuedReconfigure(transCtx)
	// }

	runningItems := sets.New(maps.Keys(runningObjs)...)
	protoItems := sets.New(maps.Keys(protoObjs)...)
	for _, item := range sets.List(runningItems.Intersection(protoItems)) {
		if !reflect.DeepEqual(runningObj.Data[item], protoObj.Data[item]) {
			return t.queueReconfigure(transCtx)
		}
	}

	// TODO: dynamic render task

	return t.reconfigure(transCtx)
}

func (t *componentFileTemplateTransformer) transformPodVolumes(transCtx *componentTransformContext) error {
	var (
		compDef         = transCtx.CompDef
		synthesizedComp = transCtx.SynthesizeComponent
	)
	if synthesizedComp.PodSpec.Volumes == nil {
		synthesizedComp.PodSpec.Volumes = []corev1.Volume{}
	}
	for _, tpl := range compDef.Spec.FileTemplates {
		objName := t.templateObjectName(transCtx, tpl.Name)
		synthesizedComp.PodSpec.Volumes = append(synthesizedComp.PodSpec.Volumes, t.newVolume(tpl, objName))
	}
	return nil
}

func (t *componentFileTemplateTransformer) newVolume(tpl appsv1.ComponentFileTemplate, objName string) corev1.Volume {
	vol := corev1.Volume{
		Name: tpl.VolumeName,
		VolumeSource: corev1.VolumeSource{
			ConfigMap: &corev1.ConfigMapVolumeSource{
				LocalObjectReference: corev1.LocalObjectReference{
					Name: objName,
				},
				DefaultMode: tpl.DefaultMode,
			},
		},
	}
	if vol.VolumeSource.ConfigMap.DefaultMode == nil {
		vol.VolumeSource.ConfigMap.DefaultMode = ptr.To[int32](0444)
	}
	return vol
}

func (t *componentFileTemplateTransformer) queueReconfigure(transCtx *componentTransformContext) error {
	return t.updateReconfigureStatus(transCtx, func(s *component.ReplicaStatus) {
		s.Reconfigured = ptr.To(false)
	})
}

// func (t *componentFileTemplateTransformer) cancelQueuedReconfigure(transCtx *componentTransformContext) error {
//	return t.updateReconfigureStatus(transCtx, func(s *component.ReplicaStatus) {
//		s.Reconfigured = nil
//	})
// }

func (t *componentFileTemplateTransformer) reconfigured(transCtx *componentTransformContext, replicas []string) error {
	replicasSet := sets.New(replicas...)
	return t.updateReconfigureStatus(transCtx, func(s *component.ReplicaStatus) {
		if replicasSet.Has(s.Name) {
			s.Reconfigured = ptr.To(true)
		}
	})
}

func (t *componentFileTemplateTransformer) updateReconfigureStatus(
	transCtx *componentTransformContext, f func(*component.ReplicaStatus)) error {
	var (
		synthesizedComp = transCtx.SynthesizeComponent
	)
	its := &workloads.InstanceSet{}
	itsKey := types.NamespacedName{
		Namespace: synthesizedComp.Namespace,
		Name:      constant.GenerateWorkloadNamePattern(synthesizedComp.ClusterName, synthesizedComp.Name),
	}
	if err := transCtx.Client.Get(transCtx.Context, itsKey, its); err != nil {
		return err
	}
	return component.UpdateReplicasStatusFunc(its, func(r *component.ReplicasStatus) error {
		for i := range r.Status {
			f(&r.Status[i])
		}
		return nil
	})
}

func (t *componentFileTemplateTransformer) reconfigure(transCtx *componentTransformContext) error {
	replicas, err := component.GetReplicasStatusFunc(nil, func(r component.ReplicaStatus) bool {
		return r.Reconfigured != nil && !*r.Reconfigured
	})
	if err != nil {
		return err
	}

	if len(replicas) > 0 {
		succeed, err1 := t.doReconfigure(transCtx, replicas)
		if len(succeed) > 0 {
			if err2 := t.reconfigured(transCtx, succeed); err2 != nil {
				return err2
			}
		}
		if err1 != nil {
			return err1
		}
		if len(succeed) < len(replicas) {
			succeedSet := sets.New(succeed...)
			remain := slices.DeleteFunc(replicas, func(r string) bool {
				return succeedSet.Has(r)
			})
			return fmt.Errorf("has %d replicas waiting for reconfigure: %s",
				len(replicas)-len(succeed), strings.Join(remain, ","))
		}
	}
	return nil
}

func (t *componentFileTemplateTransformer) doReconfigure(
	transCtx *componentTransformContext, replicas []string) ([]string, error) {
	var (
		synthesizedComp = transCtx.SynthesizeComponent
	)

	pods, err := func() (map[string]*corev1.Pod, error) {
		pods, err := component.ListOwnedPods(transCtx.Context, transCtx.Client,
			synthesizedComp.Namespace, synthesizedComp.ClusterName, synthesizedComp.Name)
		if err != nil {
			return nil, err
		}
		m := make(map[string]*corev1.Pod)
		for _, pod := range pods {
			m[pod.Name] = pod
		}
		return m, nil
	}()
	if err != nil {
		return nil, err
	}

	succeed := make([]string, 0)
	for _, replica := range replicas {
		pod := pods[replica]
		if pod != nil {
			lfa, err1 := lifecycle.New(synthesizedComp, pod)
			if err1 != nil {
				return succeed, err1
			}
			// TODO: changes to call Reconfigure
			if err1 = lfa.Reconfigure(transCtx.Context, transCtx.Client, nil); err1 != nil {
				return succeed, err1
			}
		}
		// pod is not found, take it as succeed
		succeed = append(succeed, replica)
	}
	return succeed, nil
}
