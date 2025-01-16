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
	"fmt"
	"time"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/pkg/constant"
	"github.com/apecloud/kubeblocks/pkg/controller/component"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
	"github.com/apecloud/kubeblocks/pkg/controller/model"
)

type componentNotifyDependentsTransformer struct{}

var _ graph.Transformer = &componentNotifyDependentsTransformer{}

func (t *componentNotifyDependentsTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*componentTransformContext)
	if model.IsObjectDeleting(transCtx.ComponentOrig) || model.IsObjectStatusUpdating(transCtx.ComponentOrig) {
		return nil
	}

	dependents, err := t.dependents(transCtx)
	if err != nil {
		return err
	}

	graphCli, _ := transCtx.Client.(model.GraphClient)
	for _, compName := range dependents {
		if err = t.notify(transCtx, graphCli, dag, compName); err != nil {
			return err
		}
	}
	return nil
}

func (t *componentNotifyDependentsTransformer) dependents(transCtx *componentTransformContext) ([]string, error) {
	synthesizedComp := transCtx.SynthesizeComponent
	dependents := make([]string, 0)
	for compName, compDefName := range synthesizedComp.Comp2CompDefs {
		if compName == synthesizedComp.Name {
			continue // skip self
		}
		depended, err := t.depended(transCtx, compDefName)
		if err != nil {
			return nil, err
		}
		if depended {
			dependents = append(dependents, compName)
		}
	}
	return dependents, nil
}

func (t *componentNotifyDependentsTransformer) depended(transCtx *componentTransformContext, compDefName string) (bool, error) {
	compDefReferenced := func(v appsv1.EnvVar) string {
		if v.ValueFrom != nil {
			if v.ValueFrom.HostNetworkVarRef != nil {
				return v.ValueFrom.HostNetworkVarRef.CompDef
			}
			if v.ValueFrom.ServiceVarRef != nil {
				return v.ValueFrom.ServiceVarRef.CompDef
			}
			if v.ValueFrom.CredentialVarRef != nil {
				return v.ValueFrom.CredentialVarRef.CompDef
			}
			if v.ValueFrom.TLSVarRef != nil {
				return v.ValueFrom.TLSVarRef.CompDef
			}
			if v.ValueFrom.ServiceRefVarRef != nil {
				return v.ValueFrom.ServiceRefVarRef.CompDef
			}
			if v.ValueFrom.ComponentVarRef != nil {
				return v.ValueFrom.ComponentVarRef.CompDef
			}
		}
		return ""
	}

	compDef, err := getNCheckCompDefinition(transCtx.Context, transCtx.Client, compDefName)
	if err != nil {
		return false, err
	}

	synthesizedComp := transCtx.SynthesizeComponent
	for _, v := range compDef.Spec.Vars {
		compDefPattern := compDefReferenced(v)
		if len(compDefPattern) > 0 {
			if component.PrefixOrRegexMatched(synthesizedComp.CompDefName, compDefPattern) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (t *componentNotifyDependentsTransformer) notify(transCtx *componentTransformContext,
	graphCli model.GraphClient, dag *graph.DAG, compName string) error {
	synthesizedComp := transCtx.SynthesizeComponent

	comp := &appsv1alpha1.Component{}
	compKey := types.NamespacedName{
		Namespace: synthesizedComp.Namespace,
		Name:      constant.GenerateClusterComponentName(synthesizedComp.ClusterName, compName),
	}
	if err := transCtx.Client.Get(transCtx.Context, compKey, comp); err != nil {
		return client.IgnoreNotFound(err)
	}

	if comp.Annotations == nil {
		comp.Annotations = make(map[string]string)
	}
	comp.Annotations[constant.ReconcileAnnotationKey] =
		fmt.Sprintf("%s@%s", synthesizedComp.Name, time.Now().Format(time.RFC3339))

	graphCli.Update(dag, nil, comp)

	return nil
}
