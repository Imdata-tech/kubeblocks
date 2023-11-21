/*
Copyright (C) 2022-2023 ApeCloud Co., Ltd

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
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	ctrl "sigs.k8s.io/controller-runtime"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/pkg/constant"
	"github.com/apecloud/kubeblocks/pkg/controller/graph"
	"github.com/apecloud/kubeblocks/pkg/controller/model"
	intctrlutil "github.com/apecloud/kubeblocks/pkg/controllerutil"
	testapps "github.com/apecloud/kubeblocks/pkg/testutil/apps"
)

var tlog = ctrl.Log.WithName("component_testing")

var _ = Describe("Component PostStart Test", func() {
	Context("has the BuildComponent function", func() {
		const (
			clusterDefName     = "test-clusterdef"
			clusterVersionName = "test-clusterversion"
			clusterName        = "test-cluster"
			mysqlCompDefName   = "replicasets"
			mysqlCompName      = "mysql"
			mockScriptSelector = "mock-script-selecotr"
		)

		var (
			clusterDef     *appsv1alpha1.ClusterDefinition
			clusterVersion *appsv1alpha1.ClusterVersion
			cluster        *appsv1alpha1.Cluster
		)

		BeforeEach(func() {
			clusterDef = testapps.NewClusterDefFactory(clusterDefName).
				AddComponentDef(testapps.StatefulMySQLComponent, mysqlCompDefName).
				GetObject()
			clusterVersion = testapps.NewClusterVersionFactory(clusterVersionName, clusterDefName).
				AddComponentVersion(mysqlCompDefName).
				AddContainerShort("mysql", testapps.ApeCloudMySQLImage).
				GetObject()
			pvcSpec := testapps.NewPVCSpec("1Gi")
			cluster = testapps.NewClusterFactory(testCtx.DefaultNamespace, clusterName,
				clusterDef.Name, clusterVersion.Name).
				SetUID(clusterName).
				AddComponent(mysqlCompName, mysqlCompDefName).
				AddVolumeClaimTemplate(testapps.DataVolumeName, pvcSpec).
				GetObject()
		})

		mockPodsForTest := func(cluster *appsv1alpha1.Cluster, number int) []corev1.Pod {
			clusterDefName := cluster.Spec.ClusterDefRef
			componentName := cluster.Spec.ComponentSpecs[0].Name
			clusterName := cluster.Name
			stsName := cluster.Name + "-" + componentName
			pods := make([]corev1.Pod, 0)
			for i := 0; i < number; i++ {
				pod := &corev1.Pod{
					ObjectMeta: metav1.ObjectMeta{
						Name:      stsName + "-" + strconv.Itoa(i),
						Namespace: testCtx.DefaultNamespace,
						Labels: map[string]string{
							constant.AppManagedByLabelKey:         constant.AppName,
							constant.AppNameLabelKey:              clusterDefName,
							constant.AppInstanceLabelKey:          clusterName,
							constant.KBAppComponentLabelKey:       componentName,
							appsv1.ControllerRevisionHashLabelKey: "mock-version",
						},
					},
					Spec: corev1.PodSpec{
						Containers: []corev1.Container{{
							Name:  "mock-container",
							Image: "mock-container",
						}},
					},
				}
				pods = append(pods, *pod)
			}
			return pods
		}

		It("should work as expected with various inputs", func() {
			By("test component definition without poststart")
			reqCtx := intctrlutil.RequestCtx{
				Ctx: ctx,
				Log: tlog,
			}
			synthesizeComp, err := BuildSynthesizedComponentWrapper4Test(
				reqCtx,
				testCtx.Cli,
				clusterDef,
				clusterVersion,
				cluster,
				&cluster.Spec.ComponentSpecs[0])
			Expect(err).Should(Succeed())
			Expect(synthesizeComp).ShouldNot(BeNil())
			Expect(synthesizeComp.LifecycleActions).Should(BeNil())
			Expect(synthesizeComp.LifecycleActions.PostStart).ShouldNot(BeNil())
			Expect(synthesizeComp.LifecycleActions.PostStart.CustomHandler).ShouldNot(BeNil())

			comp, err := BuildComponent(cluster, &cluster.Spec.ComponentSpecs[0])
			Expect(err).Should(Succeed())
			Expect(comp).ShouldNot(BeNil())

			dag := graph.NewDAG()
			dag.AddVertex(&model.ObjectVertex{Obj: cluster, Action: model.ActionUpdatePtr()})
			err = ReconcileCompPostStart(testCtx.Ctx, testCtx.Cli, cluster, comp, synthesizeComp, dag)
			Expect(err).Should(Succeed())

			By("build component with poststartSpec without PodList, do not need to do postStartAction")
			synthesizeComp.LifecycleActions = &appsv1alpha1.ComponentLifecycleActions{}
			postStart := appsv1alpha1.LifecycleActionSpec{
				LifecycleActionHandler: appsv1alpha1.LifecycleActionHandler{
					CustomHandler: &appsv1alpha1.Action{
						Image: constant.KBToolsImage,
						Exec: &appsv1alpha1.ExecAction{
							Command: []string{"echo", "mock"},
							Args:    []string{},
						},
					},
				},
				ScriptSpecSelectors: []appsv1alpha1.ScriptSpecSelector{
					{
						Name: mockScriptSelector,
					},
				},
			}
			synthesizeComp.LifecycleActions.PostStart = &postStart
			err = ReconcileCompPostStart(testCtx.Ctx, testCtx.Cli, cluster, comp, synthesizeComp, dag)
			Expect(err).ShouldNot(Succeed())

			By("build component with poststartSpec with PodList, do postStartAction and requeue waiting job")
			pods := mockPodsForTest(cluster, 1)
			for _, pod := range pods {
				Expect(testCtx.CheckedCreateObj(testCtx.Ctx, &pod)).Should(Succeed())
				// mock the status to pass the isReady(pod) check in consensus_set
				pod.Status.Conditions = []corev1.PodCondition{{
					Type:   corev1.PodReady,
					Status: corev1.ConditionTrue,
				}}
				Expect(k8sClient.Status().Update(ctx, &pod)).Should(Succeed())
			}
			err = ReconcileCompPostStart(testCtx.Ctx, testCtx.Cli, cluster, comp, synthesizeComp, dag)
			Expect(err).Should(Succeed())

			jobName := genPostStartJobName(cluster.Name, synthesizeComp.Name)
			err = CheckJobSucceed(testCtx.Ctx, testCtx.Cli, cluster, jobName)
			Expect(err).ShouldNot(Succeed())
			Expect(err.Error()).Should(ContainSubstring("requeue to waiting for job"))

			By("test setPostStartDoneAnnotation without error")
			err = setPostStartDoneAnnotation(testCtx.Cli, *comp, dag)
			Expect(err).Should(Succeed())
		})
	})
})
