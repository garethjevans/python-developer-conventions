/*
Copyright 2020 VMware Inc.

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

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"k8s.io/apimachinery/pkg/util/json"

	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"go.uber.org/zap"
	corev1 "k8s.io/api/core/v1"

	"github.com/vmware-tanzu/cartographer-conventions/webhook"
)

func conventionHandler(template *corev1.PodTemplateSpec, images []webhook.ImageConfig) ([]string, error) {
	readinessProbe := getAnnotation(template, "garethjevans.org/readinessProbe")
	livenessProbe := getAnnotation(template, "garethjevans.org/livenessProbe")
	startupProbe := getAnnotation(template, "garethjevans.org/startupProbe")

	var applied []string

	for i := range template.Spec.Containers {
		log.Printf("Adding for container %s", template.Spec.Containers[i].Name)
		c := &template.Spec.Containers[i]

		if readinessProbe != "" {
			// readiness probe
			if c.ReadinessProbe == nil {
				c.ReadinessProbe = &corev1.Probe{}
			}

			if c.ReadinessProbe.ProbeHandler == (corev1.ProbeHandler{}) {
				probe := corev1.ProbeHandler{}
				err := json.Unmarshal([]byte(readinessProbe), &probe)
				if err != nil {
					return nil, err
				}
				log.Printf("Adding ReadinessProbe %+v", probe)
				c.ReadinessProbe.ProbeHandler = probe

				applied = append(applied, "garethjevans-readiness-probe")
			}
		}

		if livenessProbe != "" {
			// liveness probe
			if c.LivenessProbe == nil {
				c.LivenessProbe = &corev1.Probe{}
			}
			if c.LivenessProbe.ProbeHandler == (corev1.ProbeHandler{}) {
				probe := corev1.ProbeHandler{}
				err := json.Unmarshal([]byte(livenessProbe), &probe)
				if err != nil {
					return nil, err
				}
				log.Printf("Adding LivenessProbe %+v", probe)
				c.LivenessProbe.ProbeHandler = probe

				applied = append(applied, "garethjevans-liveness-probe")
			}
		}

		if startupProbe != "" {
			// startup probe
			if c.StartupProbe == nil {
				c.StartupProbe = &corev1.Probe{}
			}
			if c.StartupProbe.ProbeHandler == (corev1.ProbeHandler{}) {
				probe := corev1.ProbeHandler{}
				err := json.Unmarshal([]byte(startupProbe), &probe)
				if err != nil {
					return nil, err
				}
				log.Printf("Adding StartupProbe %+v", probe)
				c.StartupProbe.ProbeHandler = probe

				applied = append(applied, "garethjevans-startup-probe")
			}
		}
	}

	log.Printf("PodTemplateSpec: %+v", template)

	return applied, nil
}

// getLabel gets the label on PodTemplateSpec
func getLabel(pts *corev1.PodTemplateSpec, key string) string {
	if pts.Labels == nil || len(pts.Labels[key]) == 0 {
		return ""
	}
	return pts.Labels[key]
}

// getAnnotation gets the annotation on PodTemplateSpec
func getAnnotation(pts *corev1.PodTemplateSpec, key string) string {
	if pts.Annotations == nil || len(pts.Annotations[key]) == 0 {
		return ""
	}
	return pts.Annotations[key]
}

func main() {
	ctx := context.Background()
	port := os.Getenv("PORT")
	if port == "" {
		port = "9000"
	}

	zapLog, err := zap.NewProductionConfig().Build()
	if err != nil {
		log.Fatalf("failed to create logger: %v", err)
	}
	logger := zapr.NewLogger(zapLog)
	ctx = logr.NewContext(ctx, logger)

	http.HandleFunc("/", webhook.ConventionHandler(ctx, conventionHandler))
	log.Fatal(webhook.NewConventionServer(ctx, fmt.Sprintf(":%s", port)))
}
