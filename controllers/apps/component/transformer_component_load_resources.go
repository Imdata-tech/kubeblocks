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
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"

	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	"github.com/apecloud/kubeblocks/pkg/controller/component"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
)

// componentLoadResourcesTransformer handles referenced resources validation and load them into context
type componentLoadResourcesTransformer struct{}

var _ graph.Transformer = &componentLoadResourcesTransformer{}

func (t *componentLoadResourcesTransformer) Transform(ctx graph.TransformContext, dag *graph.DAG) error {
	transCtx, _ := ctx.(*componentTransformContext)
	comp := transCtx.Component

	var err error
	defer func() {
		setProvisioningStartedCondition(&comp.Status.Conditions, comp.Name, comp.Generation, err)
	}()

	clusterName, err := component.GetClusterName(comp)
	if err != nil {
		return newRequeueError(requeueDuration, err.Error())
	}

	cluster := &appsv1.Cluster{}
	err = transCtx.Client.Get(transCtx.Context, types.NamespacedName{Name: clusterName, Namespace: comp.Namespace}, cluster)
	if err != nil {
		return newRequeueError(requeueDuration, err.Error())
	}
	transCtx.Cluster = cluster

	return t.transformForNativeComponent(transCtx)
}

func (t *componentLoadResourcesTransformer) transformForNativeComponent(transCtx *componentTransformContext) error {
	var (
		ctx  = transCtx.Context
		cli  = transCtx.Client
		comp = transCtx.Component
	)
	compDef, err := getNCheckCompDefinition(ctx, cli, comp.Spec.CompDef)
	if err != nil {
		return newRequeueError(requeueDuration, err.Error())
	}
	if err = component.UpdateCompDefinitionImages4ServiceVersion(ctx, cli, compDef, comp.Spec.ServiceVersion); err != nil {
		return newRequeueError(requeueDuration, err.Error())
	}
	transCtx.CompDef = compDef

	synthesizedComp, err := component.BuildSynthesizedComponent(ctx, transCtx.Client, compDef, comp, transCtx.Cluster)
	if err != nil {
		message := fmt.Sprintf("build synthesized component for %s failed: %s", comp.Name, err.Error())
		return newRequeueError(requeueDuration, message)
	}
	transCtx.SynthesizeComponent = synthesizedComp

	return nil
}

func getNCheckCompDefinition(ctx context.Context, cli client.Reader, name string) (*appsv1.ComponentDefinition, error) {
	compKey := types.NamespacedName{
		Name: name,
	}
	compDef := &appsv1.ComponentDefinition{}
	if err := cli.Get(ctx, compKey, compDef); err != nil {
		return nil, err
	}
	if compDef.Generation != compDef.Status.ObservedGeneration {
		return nil, fmt.Errorf("the referenced ComponentDefinition is not up to date: %s", compDef.Name)
	}
	if compDef.Status.Phase != appsv1.AvailablePhase {
		return nil, fmt.Errorf("the referenced ComponentDefinition is unavailable: %s", compDef.Name)
	}
	return compDef, nil
}
