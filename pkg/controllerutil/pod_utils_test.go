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

package controllerutil

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	"github.com/apecloud/kubeblocks/pkg/constant"
	testapps "github.com/apecloud/kubeblocks/pkg/testutil/apps"
	testk8s "github.com/apecloud/kubeblocks/pkg/testutil/k8s"
	viper "github.com/apecloud/kubeblocks/pkg/viperx"
)

type TestResourceUnit struct {
	pvc               corev1.PersistentVolumeClaimSpec
	container         corev1.Container
	expectMemorySize  int64
	expectCPU         int
	expectStorageSize int64
}

func TestPodIsReady(t *testing.T) {
	pod := testk8s.NewFakePod("foo", 1)
	pod.Status.Conditions = []corev1.PodCondition{
		{
			Type:   corev1.PodReady,
			Status: corev1.ConditionTrue,
		},
	}
	pod.Labels = map[string]string{constant.RoleLabelKey: "leader"}
	if !PodIsReadyWithLabel(*pod) {
		t.Errorf("isReady returned false negative")
	}

	pod.DeletionTimestamp = &metav1.Time{Time: time.Now()}
	if PodIsReadyWithLabel(*pod) {
		t.Errorf("isReady returned false positive")
	}

	pod.Labels = nil
	if PodIsReadyWithLabel(*pod) {
		t.Errorf("isReady returned false positive")
	}

	pod.Status.Conditions = nil
	if PodIsReadyWithLabel(*pod) {
		t.Errorf("isReady returned false positive")
	}

	pod.Status.Conditions = []corev1.PodCondition{}
	if PodIsReadyWithLabel(*pod) {
		t.Errorf("isReady returned false positive")
	}
}

func TestGetPodRevision(t *testing.T) {
	pod := testk8s.NewFakePod("foo", 1)
	if GetPodRevision(pod) != "" {
		t.Errorf("revision should be empty")
	}

	pod.Labels = make(map[string]string, 0)
	pod.Labels[appsv1.StatefulSetRevisionLabel] = "bar"

	if GetPodRevision(pod) != "bar" {
		t.Errorf("revision not matched")
	}
}

var _ = Describe("pod utils", func() {
	var (
		statefulSet     *appsv1.StatefulSet
		pod             *corev1.Pod
		configTemplates = []appsv1alpha1.ComponentConfigSpec{
			{
				ComponentTemplateSpec: appsv1alpha1.ComponentTemplateSpec{
					Name:       "xxxxx",
					VolumeName: "config1",
				},
			},
			{
				ComponentTemplateSpec: appsv1alpha1.ComponentTemplateSpec{
					Name:       "xxxxx2",
					VolumeName: "config2",
				},
			},
		}

		foundInitContainerConfigTemplates = []appsv1alpha1.ComponentConfigSpec{
			{
				ComponentTemplateSpec: appsv1alpha1.ComponentTemplateSpec{
					Name:       "xxxxx",
					VolumeName: "config1_init_container",
				},
			},
			{
				ComponentTemplateSpec: appsv1alpha1.ComponentTemplateSpec{
					Name:       "xxxxx2",
					VolumeName: "config2_init_container",
				},
			},
		}

		notFoundConfigTemplates = []appsv1alpha1.ComponentConfigSpec{
			{
				ComponentTemplateSpec: appsv1alpha1.ComponentTemplateSpec{
					Name:       "xxxxx",
					VolumeName: "config1_not_fount",
				},
			},
			{
				ComponentTemplateSpec: appsv1alpha1.ComponentTemplateSpec{
					Name:       "xxxxx2",
					VolumeName: "config2_not_fount",
				},
			},
		}
	)

	const (
		testContainers = `
{
  "name": "mysql",
  "imagePullPolicy": "IfNotPresent",
  "ports": [
    {
      "containerPort": 3306,
      "protocol": "TCP",
      "name": "mysql"
    },
    {
      "containerPort": 13306,
      "protocol": "TCP",
      "name": "paxos"
    }
  ],
  "volumeMounts": [
    {
      "mountPath": "/data/config",
      "name": "config1"
    },
    {
      "mountPath": "/data/config",
      "name": "config2"
    },
    {
      "mountPath": "/data",
      "name": "data"
    },
    {
      "mountPath": "/log",
      "name": "log"
    }
  ],
  "env": []
}
`
	)

	BeforeEach(func() {
		// Add any steup steps that needs to be executed before each test
		statefulSet = &appsv1.StatefulSet{}
		statefulSet.ObjectMeta.Name = "stateful_test"
		statefulSet.ObjectMeta.Namespace = "stateful_test_ns"

		container := corev1.Container{}
		if err := json.Unmarshal([]byte(testContainers), &container); err != nil {
			Fail("convert container failed!")
		}

		container2 := container.DeepCopy()
		container2.Name = "mysql2"
		container2.VolumeMounts[1].Name += "_not_found"
		container3 := container.DeepCopy()
		container3.Name = "mysql3"
		container3.VolumeMounts[0].Name += "_not_found"
		container3.EnvFrom = []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: "test-config-env"},
				},
			},
		}

		container4 := container.DeepCopy()
		container4.Name = "mysql4"
		container4.VolumeMounts = nil
		container4.EnvFrom = []corev1.EnvFromSource{
			{
				ConfigMapRef: &corev1.ConfigMapEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{Name: "test-config-env"},
				},
			},
		}

		statefulSet.Spec.Template.Spec.Containers = []corev1.Container{
			*container2, *container3, container}

		// init container
		initContainer := container.DeepCopy()
		initContainer.Name = "init_mysql"
		initContainer2 := container.DeepCopy()
		initContainer2.Name = "init_mysql_2"
		initContainer3 := container.DeepCopy()
		initContainer3.Name = "init_mysql_3"
		initContainer.VolumeMounts[0].Name += "_init_container"
		initContainer.VolumeMounts[1].Name += "_init_container"
		statefulSet.Spec.Template.Spec.InitContainers = []corev1.Container{
			*initContainer, *initContainer2, *initContainer3}

		// init pod
		pod = &corev1.Pod{}
		pod.ObjectMeta.Name = "pod_test"
		pod.ObjectMeta.Namespace = "pod_test_ns"
		pod.Spec.Containers = []corev1.Container{container, *container2, *container3, *container4}
		pod.Spec.Volumes = []corev1.Volume{
			{
				Name: "config1",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: "stateful_test-config1"},
					},
				},
			},
			{
				Name: "config2",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{Name: "stateful_test-config2"},
					},
				},
			},
		}

	})

	// for test GetContainerByConfigSpec
	Context("GetContainerByConfigSpec test", func() {
		// found name: mysql3
		It("Should succeed with no error", func() {
			podSpec := &statefulSet.Spec.Template.Spec
			Expect(GetContainerByConfigSpec(podSpec, configTemplates)).To(Equal(&podSpec.Containers[2]))
		})
		// found name: init_mysql
		It("Should succeed with no error", func() {
			podSpec := &statefulSet.Spec.Template.Spec
			Expect(GetContainerByConfigSpec(podSpec, foundInitContainerConfigTemplates)).To(Equal(&podSpec.InitContainers[0]))
		})
		// not found container
		It("Should fail", func() {
			podSpec := &statefulSet.Spec.Template.Spec
			Expect(GetContainerByConfigSpec(podSpec, notFoundConfigTemplates)).To(BeNil(), "get container is nil!")
		})
	})

	// for test GetVolumeMountName
	Context("GetPodContainerWithVolumeMount test", func() {
		It("Should succeed with no error", func() {
			mountedContainers := GetPodContainerWithVolumeMount(&pod.Spec, "config1")
			Expect(len(mountedContainers)).To(Equal(2))
			Expect(mountedContainers[0].Name).To(Equal("mysql"))
			Expect(mountedContainers[1].Name).To(Equal("mysql2"))

			//
			mountedContainers = GetPodContainerWithVolumeMount(&pod.Spec, "config2")
			Expect(len(mountedContainers)).To(Equal(2))
			Expect(mountedContainers[0].Name).To(Equal("mysql"))
			Expect(mountedContainers[1].Name).To(Equal("mysql3"))
		})
		It("Should fail", func() {
			Expect(len(GetPodContainerWithVolumeMount(&pod.Spec, "not_exist_cm"))).To(Equal(0))

			emptyPod := corev1.Pod{}
			emptyPod.ObjectMeta.Name = "empty_test"
			emptyPod.ObjectMeta.Namespace = "empty_test_ns"
			Expect(GetPodContainerWithVolumeMount(&emptyPod.Spec, "not_exist_cm")).To(BeNil())

		})
	})

	// for test GetContainerWithVolumeMount
	Context("GetVolumeMountName test", func() {
		It("Should succeed with no error", func() {
			volume := GetVolumeMountName(pod.Spec.Volumes, "stateful_test-config1")
			Expect(volume).NotTo(BeNil())
			Expect(volume.Name).To(Equal("config1"))

			Expect(GetVolumeMountName(pod.Spec.Volumes, "stateful_test-config1")).To(Equal(&pod.Spec.Volumes[0]))
		})
		It("Should fail", func() {
			Expect(GetVolumeMountName(pod.Spec.Volumes, "not_exist_resource")).To(BeNil())
		})
	})

	// for test MemorySize or CoreNum
	Context("Get Resource test", func() {
		It("Resource exists limit", func() {
			testResources := []TestResourceUnit{
				// memory unit: Gi
				{
					container: corev1.Container{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory:  resource.MustParse("10Gi"),
								corev1.ResourceCPU:     resource.MustParse("6"),
								corev1.ResourceStorage: resource.MustParse("100G"),
							},
						},
					},
					expectMemorySize: 10 * 1024 * 1024 * 1024,
					expectCPU:        6,
				},
				// memory unit: G
				{
					container: corev1.Container{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory:  resource.MustParse("10G"),
								corev1.ResourceCPU:     resource.MustParse("16"),
								corev1.ResourceStorage: resource.MustParse("100G"),
							},
						},
					},
					expectMemorySize: 10 * 1000 * 1000 * 1000,
					expectCPU:        16,
				},
				// memory unit: no
				{
					container: corev1.Container{
						Resources: corev1.ResourceRequirements{
							Limits: corev1.ResourceList{
								corev1.ResourceMemory:  resource.MustParse("1024000"),
								corev1.ResourceCPU:     resource.MustParse("26"),
								corev1.ResourceStorage: resource.MustParse("100G"),
							},
						},
					},
					expectMemorySize: 1024000,
					expectCPU:        26,
				},
			}

			for i := range testResources {
				Expect(GetMemorySize(testResources[i].container)).To(BeEquivalentTo(testResources[i].expectMemorySize))
				Expect(GetCoreNum(testResources[i].container)).To(BeEquivalentTo(testResources[i].expectCPU))
			}
		})
		It("Resource not limit", func() {
			container := corev1.Container{}
			Expect(GetMemorySize(container)).To(BeEquivalentTo(0))
			Expect(GetCoreNum(container)).To(BeEquivalentTo(0))
		})
	})

	// for test MemorySize or CoreNum
	Context("Get pvc test", func() {
		It("Resource exists request", func() {
			testResources := []TestResourceUnit{
				// memory unit: Gi
				{
					pvc: corev1.PersistentVolumeClaimSpec{
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("100Gi"),
							},
						},
					},
					expectStorageSize: 100 * 1024 * 1024 * 1024,
				},
				// memory unit: G
				{
					pvc: corev1.PersistentVolumeClaimSpec{
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("100G"),
							},
						},
					},
					expectStorageSize: 100 * 1000 * 1000 * 1000,
				},
				// memory unit: no
				{
					pvc: corev1.PersistentVolumeClaimSpec{
						Resources: corev1.VolumeResourceRequirements{
							Requests: corev1.ResourceList{
								corev1.ResourceStorage: resource.MustParse("10000"),
							},
						},
					},
					expectStorageSize: 10000,
				},
			}

			for i := range testResources {
				Expect(GetStorageSizeFromPersistentVolume(corev1.PersistentVolumeClaimTemplate{
					Spec: testResources[i].pvc,
				})).To(BeEquivalentTo(testResources[i].expectStorageSize))
			}
		})
		It("Resource not request", func() {
			pvcTpl := corev1.PersistentVolumeClaimTemplate{}
			Expect(GetStorageSizeFromPersistentVolume(pvcTpl)).To(BeEquivalentTo(-1))
		})
	})

	Context("testGetContainerID", func() {
		It("Should succeed with no error", func() {
			pods := []*corev1.Pod{{
				Status: corev1.PodStatus{
					ContainerStatuses: []corev1.ContainerStatus{
						{
							Name:        "a",
							ContainerID: "docker://27d1586d53ef9a6af5bd983831d13b6a38128119fadcdc22894d7b2397758eb5",
						},
						{
							Name:        "b",
							ContainerID: "docker://6f5ca0f22cd151943ba1b70f618591ad482cdbbc019ed58d7adf4c04f6d0ca7a",
						},
					},
				},
			}, {
				Status: corev1.PodStatus{
					ContainerStatuses: []corev1.ContainerStatus{},
				},
			}}

			type args struct {
				pod           *corev1.Pod
				containerName string
			}
			tests := []struct {
				name string
				args args
				want string
			}{{
				name: "test1",
				args: args{
					pod:           pods[0],
					containerName: "b",
				},
				want: "6f5ca0f22cd151943ba1b70f618591ad482cdbbc019ed58d7adf4c04f6d0ca7a",
			}, {
				name: "test2",
				args: args{
					pod:           pods[0],
					containerName: "f",
				},
				want: "",
			}, {
				name: "test3",
				args: args{
					pod:           pods[1],
					containerName: "a",
				},
				want: "",
			}}
			for _, tt := range tests {
				Expect(GetContainerID(tt.args.pod, tt.args.containerName)).Should(BeEquivalentTo(tt.want))
			}

		})
	})

	Context("common funcs test", func() {
		It("GetContainersByConfigmap Should succeed with no error", func() {
			type args struct {
				containers []corev1.Container
				volumeName string
				envFrom    string
				filters    []containerNameFilter
			}
			tests := []struct {
				name string
				args args
				want []string
			}{{
				name: "test1",
				args: args{
					containers: pod.Spec.Containers,
					volumeName: "config1",
				},
				want: []string{"mysql", "mysql2"},
			}, {
				name: "test1",
				args: args{
					containers: pod.Spec.Containers,
					volumeName: "config1",
					filters: []containerNameFilter{
						func(name string) bool {
							return name != "mysql"
						},
					},
				},
				want: []string{"mysql"},
			}, {
				name: "test1",
				args: args{
					containers: pod.Spec.Containers,
					volumeName: "config2",
					filters: []containerNameFilter{
						func(name string) bool {
							return strings.HasPrefix(name, "mysql")
						},
					},
				},
				want: []string{},
			}, {
				name: "test_env",
				args: args{
					containers: pod.Spec.Containers,
					volumeName: "not-config2",
					envFrom:    "test-config-env",
					filters: []containerNameFilter{
						func(name string) bool {
							return false
						},
					},
				},
				want: []string{"mysql3", "mysql4"},
			}}
			for _, tt := range tests {
				Expect(GetContainersByConfigmap(tt.args.containers, tt.args.volumeName, tt.args.envFrom, tt.args.filters...)).Should(BeEquivalentTo(tt.want))
			}

		})

		It("GetIntOrPercentValue Should succeed with no error", func() {
			fn := func(v intstr.IntOrString) *intstr.IntOrString { return &v }
			tests := []struct {
				name      string
				args      *intstr.IntOrString
				want      int
				isPercent bool
				wantErr   bool
			}{{
				name:      "test",
				args:      fn(intstr.FromString("10")),
				want:      0,
				isPercent: false,
				wantErr:   true,
			}, {
				name:      "test",
				args:      fn(intstr.FromString("10%")),
				want:      10,
				isPercent: true,
				wantErr:   false,
			}, {
				name:      "test",
				args:      fn(intstr.FromInt(60)),
				want:      60,
				isPercent: false,
				wantErr:   false,
			}}

			for _, tt := range tests {
				val, isPercent, err := GetIntOrPercentValue(tt.args)
				Expect(err != nil).Should(BeEquivalentTo(tt.wantErr))
				Expect(val).Should(BeEquivalentTo(tt.want))
				Expect(isPercent).Should(BeEquivalentTo(tt.isPercent))
			}
		})
	})

	Context("test sort by pod name", func() {
		It("Should succeed with no error", func() {
			pods := []corev1.Pod{{
				ObjectMeta: metav1.ObjectMeta{Name: "pod-2"},
			}, {
				ObjectMeta: metav1.ObjectMeta{Name: "pod-3"},
			}, {
				ObjectMeta: metav1.ObjectMeta{Name: "pod-0"},
			}, {
				ObjectMeta: metav1.ObjectMeta{Name: "pod-1"},
			}}
			sort.Sort(ByPodName(pods))
			Expect(pods[0].Name).Should(Equal("pod-0"))
			Expect(pods[3].Name).Should(Equal("pod-3"))
		})
	})

	Context("ResolvePodSpecDefaultFields", func() {
		It("test sync pod spec default values set by k8s", func() {
			var (
				clusterName = "cluster"
				compName    = "component"
				podName     = "pod"
				role        = "leader"
			)
			pod := testapps.MockInstanceSetPod(&testCtx, nil, clusterName, compName, podName, role)
			ppod := testapps.NewPodFactory(testCtx.DefaultNamespace, "pod").
				AddAppInstanceLabel(clusterName).
				AddAppComponentLabel(compName).
				AddAppManagedByLabel().
				AddRoleLabel(role).
				AddControllerRevisionHashLabel("").
				AddVolume(corev1.Volume{
					Name: testapps.DataVolumeName,
					VolumeSource: corev1.VolumeSource{
						PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
							ClaimName: fmt.Sprintf("%s-%s", testapps.DataVolumeName, podName),
						},
					},
				}).
				AddContainer(corev1.Container{
					Name:  testapps.DefaultMySQLContainerName,
					Image: testapps.ApeCloudMySQLImage,
					LivenessProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							HTTPGet: &corev1.HTTPGetAction{
								Path: "/hello",
								Port: intstr.FromInt(1024),
							},
						},
						TimeoutSeconds:   1,
						PeriodSeconds:    1,
						FailureThreshold: 1,
					},
					StartupProbe: &corev1.Probe{
						ProbeHandler: corev1.ProbeHandler{
							TCPSocket: &corev1.TCPSocketAction{
								Port: intstr.FromInt(1024),
							},
						},
					},
					VolumeMounts: []corev1.VolumeMount{
						{Name: testapps.DataVolumeName, MountPath: "/test"},
					},
				}).
				GetObject()
			ResolvePodSpecDefaultFields(pod.Spec, &ppod.Spec)
			Expect(reflect.DeepEqual(pod.Spec, ppod.Spec)).Should(BeTrue())
		})
	})
})

func TestBuildImagePullSecretsByEnv(t *testing.T) {
	tests := []struct {
		value    string
		expected []corev1.LocalObjectReference
	}{
		{
			value:    "",
			expected: nil,
		},
		{
			value: "[{\"name\":\"test\"}]",
			expected: []corev1.LocalObjectReference{
				{
					Name: "test",
				},
			},
		},
	}

	Context("test BuildImagePullSecrets", func() {
		It("Should succeed with no error", func() {
			for _, t := range tests {
				viper.Set(constant.KBImagePullSecrets, t.value)
				secrets := BuildImagePullSecrets()
				if t.value == "" {
					Expect(len(secrets)).To(Equal(0))
				} else {
					Expect(secrets).To(Equal(t.expected))
				}
			}
		})
	})
}
