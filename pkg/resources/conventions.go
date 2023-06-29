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
					c.ReadinessProbe = &corev1.Probe{}
				}

				if c.ReadinessProbe.ProbeHandler == (corev1.ProbeHandler{}) {
					probe := corev1.ProbeHandler{}
					err := json.Unmarshal([]byte(readinessProbe), &probe)
					if err != nil {
						return err
					}
					log.Printf("Adding ReadinessProbe %+v", probe)
					c.ReadinessProbe.ProbeHandler = probe
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
					c.LivenessProbe = &corev1.Probe{}
				}

				if c.LivenessProbe.ProbeHandler == (corev1.ProbeHandler{}) {
					probe := corev1.ProbeHandler{}
					err := json.Unmarshal([]byte(livenessProbe), &probe)
					if err != nil {
						return err
					}
					log.Printf("Adding LivenessProbe %+v", probe)
					c.LivenessProbe.ProbeHandler = probe
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
					c.StartupProbe = &corev1.Probe{}
				}

				if c.StartupProbe.ProbeHandler == (corev1.ProbeHandler{}) {
					probe := corev1.ProbeHandler{}
					err := json.Unmarshal([]byte(startupProbe), &probe)
					if err != nil {
						return err
					}
					log.Printf("Adding StartupProbe %+v", probe)
					c.StartupProbe.ProbeHandler = probe
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
