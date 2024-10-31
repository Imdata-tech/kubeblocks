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

package configuration

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"

	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	appsv1beta1 "github.com/apecloud/kubeblocks/apis/apps/v1beta1"
	parametersv1alpha1 "github.com/apecloud/kubeblocks/apis/parameters/v1alpha1"
	"github.com/apecloud/kubeblocks/pkg/configuration/core"
	intctrlutil "github.com/apecloud/kubeblocks/pkg/controllerutil"
)

func ApplyParameters(item parametersv1alpha1.ConfigTemplateItemDetail, orig *corev1.ConfigMap, configRender *parametersv1alpha1.ParameterDrivenConfigRender, paramsDefs []*parametersv1alpha1.ParametersDefinition, revision string) (*corev1.ConfigMap, error) {
	if configRender == nil || len(configRender.Spec.Configs) == 0 {
		return nil, fmt.Errorf("not support parameter reconfigure")
	}

	newData, err := DoMerge(orig.Data, item.ConfigFileParams, paramsDefs, configRender.Spec.Configs, *item.ConfigSpec)
	if err != nil {
		return nil, err
	}

	expected := orig.DeepCopy()
	expected.Data = newData
	if err := updateConfigMetaForCM(expected, item, revision); err != nil {
		return nil, err
	}
	return expected, nil
}

func DoMerge(baseData map[string]string, patch map[string]parametersv1alpha1.ParametersInFile, paramsDefs []*parametersv1alpha1.ParametersDefinition, configDescs []parametersv1alpha1.ComponentConfigDescription, configSpec appsv1.ComponentTemplateSpec) (map[string]string, error) {
	var (
		updatedFiles  = make(map[string]string, len(patch))
		updatedParams = make([]core.ParamPairs, 0, len(patch))
	)

	for key, params := range patch {
		if params.Content != nil {
			updatedFiles[key] = *params.Content
		}
		if len(params.Parameters) > 0 {
			updatedParams = append(updatedParams, core.ParamPairs{
				Key:           key,
				UpdatedParams: core.FromStringMap(params.Parameters),
			})
		}
	}
	return mergeUpdatedParams(baseData, updatedFiles, updatedParams, cc, configSpec)
}

func mergeUpdatedParams(base map[string]string,
	updatedFiles map[string]string,
	updatedParams []core.ParamPairs,
	cc *appsv1beta1.ConfigConstraint,
	tpl appsv1.ComponentConfigSpec) (map[string]string, error) {
	updatedConfig := base

	// merge updated files into configmap
	if len(updatedFiles) != 0 {
		updatedConfig = core.MergeUpdatedConfig(base, updatedFiles)
	}
	if cc == nil {
		return updatedConfig, nil
	}
	return intctrlutil.MergeAndValidateConfigs(cc.Spec, updatedConfig, tpl.Keys, updatedParams)
}
