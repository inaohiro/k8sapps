package models

import (
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeploymentCreate struct {
	Name  string
	Image string
}

func FromRequest(r DeploymentCreateRequest) DeploymentCreate {
	return DeploymentCreate{
		Name:  r.Name,
		Image: r.Image,
	}
}

func (d DeploymentCreate) ToKubeDeployment() *appsv1.Deployment {
	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: d.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{"app": d.Name},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"app": d.Name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Name:  d.Name,
						Image: d.Image,
					}},
				},
			},
		},
	}
}
