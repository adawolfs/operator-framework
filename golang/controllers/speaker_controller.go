/*


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

package controllers

import (
	"bytes"
	"context"
	"text/template"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	horadek8sv1 "github.com/adawolfs/operator-framework/golang/api/v1"
)

// SpeakerReconciler reconciles a Speaker object
type SpeakerReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=hora.de.k8s.adawolfs.github.io,resources=speakers,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hora.de.k8s.adawolfs.github.io,resources=speakers/status,verbs=get;update;patch

func (r *SpeakerReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("speaker", req.NamespacedName)

	speaker := &horadek8sv1.Speaker{}
	err := r.Get(ctx, req.NamespacedName, speaker)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Speaker resource not found. Ignoring since object must be deleted")
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get Speaker")
		return ctrl.Result{}, err
	}

	cm_found := &corev1.ConfigMap{}
	err = r.Get(ctx, types.NamespacedName{Name: speaker.Name, Namespace: speaker.Namespace}, cm_found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Pod
		cm := r.deployConfigmap(speaker)
		log.Info("Creating a new Configmap", "Configmap.Namespace", cm.Namespace, "Configmap.Name", cm.Name)
		err = r.Create(ctx, cm)
		if err != nil {
			log.Error(err, "Failed to create new Pod", "Configmap.Namespace", cm.Namespace, "Configmap.Name", cm.Name)
			return ctrl.Result{}, err
		}
		// Pod created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Pod")
		return ctrl.Result{}, err
	}

	pod_found := &corev1.Pod{}
	err = r.Get(ctx, types.NamespacedName{Name: speaker.Name, Namespace: speaker.Namespace}, pod_found)
	if err != nil && errors.IsNotFound(err) {
		// Define a new Pod
		pod := r.deploySpeaker(speaker)
		log.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.Create(ctx, pod)
		if err != nil {
			log.Error(err, "Failed to create new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
			return ctrl.Result{}, err
		}
		// Pod created successfully - return and requeue
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Pod")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

type ConfigmapSetup struct {
	FirstName string
	LastName  string
	Avatar    string
}

func (r *SpeakerReconciler) deployConfigmap(s *horadek8sv1.Speaker) *corev1.ConfigMap {
	avatar := s.Spec.Avatar
	firstName := s.Spec.FirstName
	lastName := s.Spec.LastName
	cm_setup := ConfigmapSetup{firstName, lastName, avatar}

	indexTemplate := `
			<html>
				<head></head>
				<body>
					<div style="align-content: center;text-align: center;">
						<h1>Golang</h1>
						<img src="{{ .Avatar}}" alt="avatar" style="border-radius: 50%;">
						<h1>{{ .FirstName}} {{ .LastName}}</h1>
					</div>
				</body>

			</html>
			`
	t, err := template.New("index").Parse(indexTemplate)

	buf := &bytes.Buffer{}
	err = t.Execute(buf, cm_setup)
	if err != nil {
		panic(err)
	}

	// indexContent := fmt.Sprintf(indexTemplate, avatar, firstName, lastName)

	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Data: map[string]string{
			"index.html": buf.String(),
		},
	}
	return cm
}

func (r *SpeakerReconciler) deploySpeaker(s *horadek8sv1.Speaker) *corev1.Pod {

	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{{
				Image: "nginx",
				Name:  "nginx",
				Ports: []corev1.ContainerPort{{
					ContainerPort: 80,
					Name:          "http",
				}},
				VolumeMounts: []corev1.VolumeMount{{
					Name:      "config",
					MountPath: "/usr/share/nginx/html/index.html",
					SubPath:   "index.html",
				}},
			}},
			Volumes: []corev1.Volume{{
				Name: "config",
				VolumeSource: corev1.VolumeSource{
					ConfigMap: &corev1.ConfigMapVolumeSource{
						LocalObjectReference: corev1.LocalObjectReference{
							Name: s.Name,
						},
					},
				},
			}},
		},
	}
	return pod
}

func (r *SpeakerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&horadek8sv1.Speaker{}).
		Complete(r)
}
