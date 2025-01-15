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
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/uuid"

	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	workloads "github.com/apecloud/kubeblocks/apis/workloads/v1"
	"github.com/apecloud/kubeblocks/pkg/constant"
	"github.com/apecloud/kubeblocks/pkg/controller/component"
	"github.com/apecloud/kubeblocks/pkg/controller/factory"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
	"github.com/apecloud/kubeblocks/pkg/controller/model"
	testapps "github.com/apecloud/kubeblocks/pkg/testutil/apps"
	viper "github.com/apecloud/kubeblocks/pkg/viperx"
)

var _ = Describe("object rbac transformer test.", func() {
	const (
		compDefName = "test-compdef"
		clusterName = "test-cluster"
		compName    = "test-comp"
	)

	var (
		serviceAccountName = constant.GenerateDefaultServiceAccountName(clusterName)
		transCtx           graph.TransformContext
		dag                *graph.DAG
		graphCli           model.GraphClient
		transformer        graph.Transformer
		compDefObj         *appsv1.ComponentDefinition
		clusterUID         = string(uuid.NewUUID())
		compObj            *appsv1.Component
		synthesizedComp    *component.SynthesizedComponent
		saKey              types.NamespacedName
		settings           map[string]interface{}
	)

	BeforeEach(func() {
		By("Create a component definition")
		compDefObj = testapps.NewComponentDefinitionFactory(compDefName).
			WithRandomName().
			SetDefaultSpec().
			Create(&testCtx).
			GetObject()

		By("Creating a component")
		fullCompName := constant.GenerateClusterComponentName(clusterName, compName)
		compObj = testapps.NewComponentFactory(testCtx.DefaultNamespace, fullCompName, compDefName).
			AddAnnotations(constant.KBAppClusterUIDKey, clusterUID).
			AddLabels(constant.AppInstanceLabelKey, clusterName).
			SetReplicas(1).
			SetServiceAccountName(serviceAccountName).
			GetObject()

		saKey = types.NamespacedName{
			Namespace: testCtx.DefaultNamespace,
			Name:      serviceAccountName,
		}

		graphCli = model.NewGraphClient(k8sClient)

		var err error
		synthesizedComp, err = component.BuildSynthesizedComponent(ctx, k8sClient, compDefObj, compObj)
		Expect(err).Should(Succeed())

		transCtx = &componentTransformContext{
			Context:             ctx,
			Client:              graphCli,
			EventRecorder:       nil,
			Logger:              logger,
			CompDef:             compDefObj,
			Component:           compObj,
			ComponentOrig:       compObj.DeepCopy(),
			SynthesizeComponent: synthesizedComp,
		}

		dag = mockDAG(graphCli, compObj)
		transformer = &componentRBACTransformer{}
		settings = viper.AllSettings()
		viper.SetDefault(constant.EnableRBACManager, true)
	})

	AfterEach(func() {
		viper.SetDefault(constant.EnableRBACManager, false)
		if settings != nil {
			Expect(viper.MergeConfigMap(settings)).ShouldNot(HaveOccurred())
			settings = nil
		}
	})

	Context("transformer rbac manager", func() {
		It("create serviceaccount, rolebinding if not exist", func() {
			Eventually(testapps.CheckObjExists(&testCtx, saKey,
				&corev1.ServiceAccount{}, false)).Should(Succeed())
			Expect(transformer.Transform(transCtx, dag)).Should(BeNil())

			serviceAccount := factory.BuildServiceAccount(synthesizedComp, serviceAccountName)
			roleBinding := factory.BuildRoleBinding(synthesizedComp, serviceAccount.Name)

			dagExpected := mockDAG(graphCli, compObj)
			graphCli.Create(dagExpected, serviceAccount)
			graphCli.Create(dagExpected, roleBinding)
			graphCli.DependOn(dagExpected, roleBinding, serviceAccount)
			itsList := graphCli.FindAll(dagExpected, &workloads.InstanceSet{})
			for i := range itsList {
				graphCli.DependOn(dagExpected, itsList[i], serviceAccount)
			}
			Expect(dag.Equals(dagExpected, model.DefaultLess)).Should(BeTrue())
		})
	})
})

func mockDAG(graphCli model.GraphClient, comp *appsv1.Component) *graph.DAG {
	d := graph.NewDAG()
	graphCli.Root(d, comp, comp, model.ActionStatusPtr())
	its := &workloads.InstanceSet{}
	graphCli.Create(d, its)
	return d
}
