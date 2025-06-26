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

// client-go の Deployment から API 用 Deployment DTO へ変換
func FromKubeDeployment(dep *appsv1.Deployment) Deployment {
	status := "Unknown"
	if dep.Status.AvailableReplicas > 0 {
		status = "Running"
	} else if dep.Status.Replicas > 0 {
		status = "Pending"
	}
	image := ""
	if len(dep.Spec.Template.Spec.Containers) > 0 {
		image = dep.Spec.Template.Spec.Containers[0].Image
	}
	return Deployment{
		Id:        string(dep.UID),
		Name:      dep.Name,
		Status:    status,
		Image:     image,
		CreatedAt: dep.CreationTimestamp.Time,
	}
}

// 複数変換
func FromKubeDeploymentList(list []appsv1.Deployment) []Deployment {
	result := make([]Deployment, 0, len(list))
	for _, d := range list {
		copy := d // avoid pointer issue
		result = append(result, FromKubeDeployment(&copy))
	}
	return result
}
