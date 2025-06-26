package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodCreate struct {
	Name  string
	Image string
}

func FromRequest(r PodCreateRequest) PodCreate {
	return PodCreate{
		Name:  r.Name,
		Image: r.Image,
	}
}

func (p PodCreate) ToKubePod() *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:   p.Name,
			Labels: map[string]string{"app": p.Name},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Name:  p.Name,
				Image: p.Image,
			}},
		},
	}
}

// client-go の Pod から API 用 Pod DTO へ変換
func FromKubePod(pod *corev1.Pod) Pod {
	status := string(pod.Status.Phase)
	image := ""
	if len(pod.Spec.Containers) > 0 {
		image = pod.Spec.Containers[0].Image
	}
	return Pod{
		Id:        string(pod.UID),
		Name:      pod.Name,
		Status:    status,
		Image:     image,
		CreatedAt: pod.CreationTimestamp.Time,
	}
}

// 複数変換
func FromKubePodList(list []corev1.Pod) []Pod {
	result := make([]Pod, 0, len(list))
	for _, p := range list {
		copy := p // avoid pointer issue
		result = append(result, FromKubePod(&copy))
	}
	return result
}
