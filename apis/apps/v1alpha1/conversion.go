/*
Copyright (C) 2022-2024 ApeCloud Co., Ltd

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import corev1 "k8s.io/api/core/v1"

func VolumeResourceConvert(resource corev1.ResourceRequirements) corev1.VolumeResourceRequirements {
	return corev1.VolumeResourceRequirements{
		Limits:   resource.Limits,
		Requests: resource.Requests,
	}
}

func RequireResourceConvert(resource corev1.VolumeResourceRequirements) corev1.ResourceRequirements {
	return corev1.ResourceRequirements{
		Limits:   resource.Limits,
		Requests: resource.Requests,
	}
}
