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
	snapshotv1 "github.com/kubernetes-csi/external-snapshotter/client/v6/apis/volumesnapshot/v1"
	batchv1 "k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"

	appsv1 "github.com/apecloud/kubeblocks/apis/apps/v1"
	appsv1alpha1 "github.com/apecloud/kubeblocks/apis/apps/v1alpha1"
	appsv1beta1 "github.com/apecloud/kubeblocks/apis/apps/v1beta1"
	dpv1alpha1 "github.com/apecloud/kubeblocks/apis/dataprotection/v1alpha1"
	extensionsv1alpha1 "github.com/apecloud/kubeblocks/apis/extensions/v1alpha1"
	workloads "github.com/apecloud/kubeblocks/apis/workloads/v1"
)

var (
	rscheme = runtime.NewScheme()
)

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(rscheme))
	utilruntime.Must(appsv1alpha1.AddToScheme(rscheme))
	utilruntime.Must(appsv1beta1.AddToScheme(rscheme))
	utilruntime.Must(appsv1.AddToScheme(rscheme))
	utilruntime.Must(dpv1alpha1.AddToScheme(rscheme))
	utilruntime.Must(snapshotv1.AddToScheme(rscheme))
	utilruntime.Must(extensionsv1alpha1.AddToScheme(rscheme))
	utilruntime.Must(batchv1.AddToScheme(rscheme))
	utilruntime.Must(workloads.AddToScheme(rscheme))
}
