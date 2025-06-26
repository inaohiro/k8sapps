package models

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
)

type ServiceCreate struct {
	Name  string
	Ports []ServicePort
	Type  string
}

func FromRequest(r ServiceCreateRequest) ServiceCreate {
	return ServiceCreate{
		Name:  r.Name,
		Ports: r.Ports,
		Type:  r.Type,
	}
}

func (s ServiceCreate) ToKubeService() *corev1.Service {
	ports := make([]corev1.ServicePort, 0, len(s.Ports))
	for _, p := range s.Ports {
		ports = append(ports, corev1.ServicePort{
			// TODO: Port が１つの場合 Name は省略可能
			Name:       *p.Name,
			Port:       int32(p.Port),
			TargetPort: intstr.FromInt(p.TargetPort),
			Protocol:   corev1.Protocol(p.Protocol),
		})
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:   s.Name,
			Labels: map[string]string{"app": s.Name},
		},
		Spec: corev1.ServiceSpec{
			Type:  corev1.ServiceType(s.Type),
			Ports: ports,
		},
	}
}

func FromKubeServicePort(port corev1.ServicePort) ServicePort {
	svcPort := ServicePort{
		Port:       int(port.Port),
		TargetPort: int(port.TargetPort.IntVal),
		Protocol:   string(port.Protocol),
	}
	if port.Name != "" {
		svcPort.Name = &port.Name
	}

	return svcPort
}

func FromKubeService(svc *corev1.Service) Service {
	ports := make([]ServicePort, 0, len(svc.Spec.Ports))
	for _, p := range svc.Spec.Ports {
		ports = append(ports, FromKubeServicePort(p))
	}
	return Service{
		Id:        string(svc.UID),
		Name:      svc.Name,
		Type:      string(svc.Spec.Type),
		ClusterIP: svc.Spec.ClusterIP,
		Ports:     ports,
		CreatedAt: svc.CreationTimestamp.Time,
	}
}

func FromKubeServiceList(list []corev1.Service) []Service {
	result := make([]Service, 0, len(list))
	for _, s := range list {
		copy := s
		result = append(result, FromKubeService(&copy))
	}
	return result
}
