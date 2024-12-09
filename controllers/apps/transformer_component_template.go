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
	"golang.org/x/exp/maps"
	"reflect"
	"slices"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
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
	kubeblockFileTemplate = "kubeblock-file-templates"
)

type componentFileTemplateTransformer struct{}

var _ graph.Transformer = &componentFileTemplateTransformer{}

func (t *componentFileTemplateTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*componentTransformContext)
	if model.IsObjectDeleting(transCtx.ComponentOrig) {
		return nil
	}
	return t.transform(transCtx, dag)
}

func (t *componentFileTemplateTransformer) transform(transCtx *componentTransformContext, dag *graph.DAG) error {
	runningObj, err := t.getFileTemplateObject(transCtx)
	if err != nil {
		return err
	}

	protoObj, err := t.buildFileTemplateObject(transCtx)
	if err != nil {
		return err
	}

	isCreate, isDelete, isUpdate, err1 := t.transformData(transCtx, runningObj, protoObj)
	if err1 != nil {
		return err1
	}

	graphCli, _ := transCtx.Client.(model.GraphClient)
	switch {
	case isCreate:
		graphCli.Create(dag, protoObj, inDataContext4G())
	case isDelete:
		graphCli.Delete(dag, runningObj, inDataContext4G())
	case isUpdate:
		graphCli.Update(dag, runningObj, protoObj, inDataContext4G())
	}
	return nil
}

func (t *componentFileTemplateTransformer) fileTemplateObjectName(transCtx *componentTransformContext) string {
	synthesizedComp := transCtx.SynthesizeComponent
	return fmt.Sprintf("%s-%s", synthesizedComp.FullCompName, "file-template")
}

func (t *componentFileTemplateTransformer) getFileTemplateObject(transCtx *componentTransformContext) (*corev1.ConfigMap, error) {
	cm := &corev1.ConfigMap{}
	cmKey := types.NamespacedName{
		Namespace: transCtx.SynthesizeComponent.Namespace,
		Name:      t.fileTemplateObjectName(transCtx),
	}
	if err := transCtx.Client.Get(transCtx.Context, cmKey, cm); err != nil {
		return nil, client.IgnoreNotFound(err)
	}
	return cm, nil
}

func (t *componentFileTemplateTransformer) buildFileTemplateObject(transCtx *componentTransformContext) (*corev1.ConfigMap, error) {
	var (
		compDef         = transCtx.CompDef
		synthesizedComp = transCtx.SynthesizeComponent
	)

	mergedData := make(map[string]string)
	for _, tpl := range compDef.Spec.FileTemplates {
		data, err := t.buildFileTemplateData(transCtx, tpl)
		if err != nil {
			return nil, err
		}
		for key, val := range data {
			namedKey := t.namedKey(tpl.Name, key)
			if _, ok := mergedData[namedKey]; ok {
				return nil, fmt.Errorf("duplicated template and file: %s", namedKey)
			}
			mergedData[namedKey] = val
		}
	}
	obj := builder.NewConfigMapBuilder(synthesizedComp.Name, t.fileTemplateObjectName(transCtx)).
		AddLabelsInMap(synthesizedComp.StaticLabels).
		AddLabelsInMap(constant.GetCompLabelsWithDef(synthesizedComp.ClusterName, synthesizedComp.Name, compDef.Name)).
		AddAnnotationsInMap(synthesizedComp.StaticAnnotations).
		SetData(mergedData).
		GetObject()

	if err := t.templatesToObject(compDef.Spec.FileTemplates, obj); err != nil {
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

func (t *componentFileTemplateTransformer) buildFileTemplateData(
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

	return t.renderFileTemplateData(transCtx, tpl, cmObj.Data)
}

func (t *componentFileTemplateTransformer) renderFileTemplateData(transCtx *componentTransformContext,
	fileTemplate appsv1.ComponentFileTemplate, data map[string]string) (map[string]string, error) {
	var (
		synthesizedComp = transCtx.SynthesizeComponent
		renderedData    = make(map[string]string)
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
		renderedData[key] = buf.String()
	}
	return renderedData, nil
}

func (t *componentFileTemplateTransformer) templatesToObject(templates []appsv1.ComponentFileTemplate, obj *corev1.ConfigMap) error {
	tpl, err := json.Marshal(templates)
	if err != nil {
		return err
	}
	obj.Data[kubeblockFileTemplate] = string(tpl)
	return nil
}

func (t *componentFileTemplateTransformer) templatesFromObject(obj *corev1.ConfigMap) ([]appsv1.ComponentFileTemplate, error) {
	tplData := obj.Data[kubeblockFileTemplate]
	if len(tplData) == 0 {
		return nil, nil
	}
	var templates []appsv1.ComponentFileTemplate
	if err := json.Unmarshal([]byte(tplData), &templates); err != nil {
		return nil, err
	}
	return templates, nil
}

func (t *componentFileTemplateTransformer) transformData(
	transCtx *componentTransformContext, runningObj, protoObj *corev1.ConfigMap) (bool, bool, bool, error) {
	var (
		isCreate, isDelete, isUpdate, isDataUpdate bool
		err                                        error
	)

	if runningObj == nil && protoObj == nil {
		return false, false, false, nil
	}
	if runningObj == nil {
		t.createPodVolumes(transCtx, protoObj)
		isCreate = true
	}
	if protoObj == nil {
		if err = t.deletePodVolumes(transCtx, runningObj); err != nil {
			return false, false, false, err
		}
		isDelete = true
	}
	if runningObj != nil && protoObj != nil {
		isUpdate, isDataUpdate, err = t.updatePodVolumes(transCtx, runningObj, protoObj)
		if err != nil {
			return false, false, false, err
		}
	}

	// TODO: dynamic render task

	if isDelete || isDataUpdate {
		if err = t.pendingReconfigure(transCtx); err != nil {
			return false, false, false, err
		}
	}
	return isCreate, isDelete, isUpdate, nil
}

func (t *componentFileTemplateTransformer) createPodVolumes(transCtx *componentTransformContext, protoObj *corev1.ConfigMap) {
	var (
		compDef         = transCtx.CompDef
		synthesizedComp = transCtx.SynthesizeComponent
	)
	volumes := t.volumes(compDef.Spec.FileTemplates, protoObj)
	if synthesizedComp.PodSpec.Volumes == nil {
		synthesizedComp.PodSpec.Volumes = []corev1.Volume{}
	}
	synthesizedComp.PodSpec.Volumes = append(synthesizedComp.PodSpec.Volumes, volumes...)
}

func (t *componentFileTemplateTransformer) deletePodVolumes(transCtx *componentTransformContext, runningObj *corev1.ConfigMap) error {
	var (
		synthesizedComp = transCtx.SynthesizeComponent
	)
	templates, err := t.templatesFromObject(runningObj)
	if err != nil {
		return err
	}
	for _, tpl := range templates {
		synthesizedComp.PodSpec.Volumes = slices.DeleteFunc(synthesizedComp.PodSpec.Volumes, func(vol corev1.Volume) bool {
			return vol.Name == tpl.VolumeName
		})
	}
	return nil
}

func (t *componentFileTemplateTransformer) updatePodVolumes(
	transCtx *componentTransformContext, runningObj, protoObj *corev1.ConfigMap) (bool, bool, error) {
	volumes := func(obj *corev1.ConfigMap) (map[string]corev1.Volume, sets.Set[string], error) {
		templates, err := t.templatesFromObject(obj)
		if err != nil {
			return nil, nil, err
		}
		vols := t.volumes(templates, obj)

		volsMap := make(map[string]corev1.Volume)
		volsSet := sets.New[string]()
		for i, vol := range vols {
			volsMap[vol.Name] = vols[i]
			volsSet.Insert(vol.Name)
		}
		return volsMap, volsSet, nil
	}

	runningVolumes, runningVolumesSet, err1 := volumes(runningObj)
	if err1 != nil {
		return false, false, err1
	}
	protoVolumes, protoVolumesSet, err2 := volumes(protoObj)
	if err2 != nil {
		return false, false, err2
	}

	createSet, deleteSet, updateSet := setDiff(runningVolumesSet, protoVolumesSet)

	synthesizedComp := transCtx.SynthesizeComponent
	for _, name := range sets.List(deleteSet) {
		synthesizedComp.PodSpec.Volumes = slices.DeleteFunc(synthesizedComp.PodSpec.Volumes, func(vol corev1.Volume) bool {
			return vol.Name == name
		})
	}
	for _, name := range sets.List(updateSet) {
		for i, vol := range synthesizedComp.PodSpec.Volumes {
			if vol.Name == name {
				synthesizedComp.PodSpec.Volumes[i] = protoVolumes[name]
				break
			}
		}
	}
	for _, name := range sets.List(createSet) {
		synthesizedComp.PodSpec.Volumes = append(synthesizedComp.PodSpec.Volumes, protoVolumes[name])
	}

	isUpdate := !reflect.DeepEqual(runningObj.Data, protoObj.Data) ||
		!reflect.DeepEqual(runningObj.Labels, protoObj.Labels) ||
		!reflect.DeepEqual(runningObj.Annotations, protoObj.Annotations)

	isDataUpdate := false

	return isUpdate, isDataUpdate, nil
}

func (t *componentFileTemplateTransformer) volumes(templates []appsv1.ComponentFileTemplate, obj *corev1.ConfigMap) []corev1.Volume {
	keys := maps.Keys(obj.Data)
	slices.Sort(keys)

	vols := make([]corev1.Volume, 0)
	for _, tpl := range templates {
		vol := corev1.Volume{
			Name: tpl.VolumeName,
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: obj.Name,
					},
					Items:       []corev1.KeyToPath{},
					DefaultMode: tpl.DefaultMode,
				},
			},
		}
		for _, key := range keys {
			name, path := t.nameNKey(key)
			if tpl.Name == name {
				vol.VolumeSource.ConfigMap.Items = append(vol.VolumeSource.ConfigMap.Items,
					corev1.KeyToPath{
						Key:  key,
						Path: path,
					},
				)
			}
		}
		if vol.VolumeSource.ConfigMap.DefaultMode == nil {
			vol.VolumeSource.ConfigMap.DefaultMode = ptr.To[int32](0444)
		}
		vols = append(vols, vol)
	}
	return vols
}

func (t *componentFileTemplateTransformer) pendingReconfigure(transCtx *componentTransformContext) error {
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
		for _, s := range r.Status {
			s.Reconfigured = ptr.To(false)
		}
		return nil
	})
}

func (t *componentFileTemplateTransformer) reconfigured(transCtx *componentTransformContext, replicas []string) error {
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

	replicasSet := sets.New(replicas...)
	return component.UpdateReplicasStatusFunc(its, func(r *component.ReplicasStatus) error {
		for _, s := range r.Status {
			if replicasSet.Has(s.Name) {
				s.Reconfigured = ptr.To(true)
			}
		}
		return nil
	})
}

func (t *componentFileTemplateTransformer) reconfigure(transCtx *componentTransformContext, protoObj *corev1.ConfigMap) error {
	replicas, err := component.GetReplicasStatusFunc(nil, func(r component.ReplicaStatus) bool {
		return r.Reconfigured != nil && !*r.Reconfigured
	})
	if err != nil {
		return err
	}

	if len(replicas) > 0 {
		succeed, err1 := t.doReconfigure(transCtx, protoObj, replicas)
		if len(succeed) > 0 {
			if err2 := t.reconfigured(transCtx, succeed); err2 != nil {
				return err2
			}
		}
		if err1 != nil {
			return err1
		}
	}
	return nil
}

func (t *componentFileTemplateTransformer) doReconfigure(transCtx *componentTransformContext,
	protoObj *corev1.ConfigMap, replicas []string) ([]string, error) {
	var (
		compDef         = transCtx.CompDef
		synthesizedComp = transCtx.SynthesizeComponent
	)

	if compDef.Spec.LifecycleActions == nil || compDef.Spec.LifecycleActions.Reconfigure == nil {
		return replicas, nil
	}

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
			if err1 = lfa.Reconfigure(transCtx.Context, transCtx.Client, nil); err1 != nil {
				return succeed, err1
			}
		}
		// pod is not found, take it as succeed
		succeed = append(succeed, replica)
	}
	return succeed, nil
}
