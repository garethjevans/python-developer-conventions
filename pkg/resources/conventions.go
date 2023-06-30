package resources

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"

	"github.com/garethjevans/simple-conventions/pkg/conventions"
)

const Prefix = "garethjevans.org"

var Conventions = []conventions.Convention{
	&conventions.BasicConvention{
		Id: fmt.Sprintf("%s-readiness", Prefix),
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata conventions.ImageMetadata) bool {
			return getAnnotation(target, fmt.Sprintf("%s/readinessProbe", Prefix)) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata conventions.ImageMetadata, imageName string) error {
			readinessProbe := getAnnotation(target, fmt.Sprintf("%s/readinessProbe", Prefix))

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				if c.ReadinessProbe == nil {
					p, err := getProbe(readinessProbe)
					if err != nil {
						return err
					}
					log.Printf("Adding ReadinessProbe %+v", p)
					c.ReadinessProbe = p
				}
			}
			return nil
		},
	},

	&conventions.BasicConvention{
		Id: fmt.Sprintf("%s-liveness", Prefix),
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata conventions.ImageMetadata) bool {
			return getAnnotation(target, fmt.Sprintf("%s/livenessProbe", Prefix)) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata conventions.ImageMetadata, imageName string) error {
			livenessProbe := getAnnotation(target, fmt.Sprintf("%s/livenessProbe", Prefix))

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				if c.LivenessProbe == nil {
					p, err := getProbe(livenessProbe)
					if err != nil {
						return err
					}
					log.Printf("Adding LivenessProbe %+v", p)
					c.LivenessProbe = p
				}
			}
			return nil
		},
	},

	&conventions.BasicConvention{
		Id: fmt.Sprintf("%s-startup", Prefix),
		Applicable: func(ctx context.Context, target *corev1.PodTemplateSpec, metadata conventions.ImageMetadata) bool {
			return getAnnotation(target, fmt.Sprintf("%s/startupProbe", Prefix)) != ""
		},
		Apply: func(ctx context.Context, target *corev1.PodTemplateSpec, containerIdx int, metadata conventions.ImageMetadata, imageName string) error {
			startupProbe := getAnnotation(target, fmt.Sprintf("%s/startupProbe", Prefix))

			for i := range target.Spec.Containers {
				c := &target.Spec.Containers[i]

				if c.StartupProbe == nil {
					p, err := getProbe(startupProbe)
					if err != nil {
						return err
					}
					log.Printf("Adding StartupProbe %+v", p)
					c.StartupProbe = p
				}
			}
			return nil
		},
	},
}

// getAnnotation gets the annotation on PodTemplateSpec
func getAnnotation(pts *corev1.PodTemplateSpec, key string) string {
	if pts.Annotations == nil || len(pts.Annotations[key]) == 0 {
		return ""
	}
	return pts.Annotations[key]
}

func getProbe(config string) (*corev1.Probe, error) {
	probe := corev1.Probe{}
	err := json.Unmarshal([]byte(config), &probe)
	return &probe, err
}
