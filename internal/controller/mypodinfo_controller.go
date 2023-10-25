/*
Copyright 2023.

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

package controller

import (
	"context"
	//"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	//"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	appv1 "cproject.domain/MyPodinfo/api/v1"
)

// MyPodinfoReconciler reconciles a MyPodinfo object
type MyPodinfoReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// convert int32 to *int32, required for replicas
func Int32(v int32) *int32 {
	return &v
}

// create an rc
func (r *MyPodinfoReconciler) CreateRC(name string, ctx context.Context, podInfo appv1.MyPodinfo) error {

	createRc := &corev1.ReplicationController{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
		},
		Spec: corev1.ReplicationControllerSpec{
			Replicas: Int32(int32(podInfo.Spec.ReplicaCount)),
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "podinfo",
					Labels: map[string]string{"app": name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: podInfo.Spec.Image.Image,
							Name:  podInfo.Spec.Image.Name,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									HostPort:      8082,
								},
							},
						},
					},
				},
			},
		},
	}

	if err := r.Client.Create(ctx, createRc); err != nil {
		log.Log.Error(err, "Unable to create new-updated pod")
		return err
	}
	return nil
}

// update an rc
func (r *MyPodinfoReconciler) UpdateRC(name string, ctx context.Context, podInfo appv1.MyPodinfo) {
	r.DeleteRC(name, ctx)

	createRc := &corev1.ReplicationController{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
		},
		Spec: corev1.ReplicationControllerSpec{
			Replicas: Int32(int32(podInfo.Spec.ReplicaCount)),
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "podinfo",
					Labels: map[string]string{"app": name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: podInfo.Spec.Image.Image,
							Name:  podInfo.Spec.Image.Name,
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 8080,
									HostPort:      8082,
								},
							},
						},
					},
				},
			},
		},
	}

	if err := r.Client.Create(ctx, createRc); err != nil {
		log.Log.Error(err, "Unable to create new-updated pod")
	}
}

// deletes an rc by setting replica to 0 then calling Delete()
func (r *MyPodinfoReconciler) DeleteRC(name string, ctx context.Context) {
	deleteRc := &corev1.ReplicationController{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "default",
			Name:      name,
		},
		Spec: corev1.ReplicationControllerSpec{
			Replicas: Int32(0),
			Selector: map[string]string{"app": "podinfo"},
			Template: &corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   "podinfo",
					Labels: map[string]string{"app": name},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Image: "nginx",
							Name:  "nginx",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 9898,
									HostPort:      8082,
								},
							},
						},
					},
				},
			}},
	}
	// must update replicas to 0 before deleting
	if uerr := r.Update(ctx, deleteRc); uerr != nil {
		log.Log.Error(uerr, "Unable to update podInfo")
	}
	// delete rc
	if derr := r.Delete(ctx, deleteRc); derr != nil {
		log.Log.Error(derr, "Unable to delete podInfo")
	}
}

//+kubebuilder:rbac:groups=app.cproject.domain,resources=mypodinfoes,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=app.cproject.domain,resources=mypodinfoes/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=app.cproject.domain,resources=mypodinfoes/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the MyPodinfo object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.0/pkg/reconcile
func (r *MyPodinfoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)
	var podInfo appv1.MyPodinfo

	// Attempt to get, if err, deleted has been called
	if err := r.Get(ctx, req.NamespacedName, &podInfo); err != nil {

		r.DeleteRC("podinfo", ctx)
		log.Log.Error(err, "Unable to fetch podInfo")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	/*
		if berr := r.Client.Update(ctx, createRc); berr != nil {
			log.Log.Error(berr, "Unable to update podInfo")
		}
	*/

	// Attempt to create, if err, rc already exists, update instead
	if err := r.CreateRC("podinfo", ctx, podInfo); err != nil {
		log.Log.Error(err, "Unable to create second round pod")
		r.UpdateRC("podinfo", ctx, podInfo)
		// update
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	log.Log.Info("Inside Reconcicle Loop", "pod", req.NamespacedName)
	log.Log.Info("Replication count", "t", podInfo.Spec.ReplicaCount)

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *MyPodinfoReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appv1.MyPodinfo{}).
		Complete(r)
}
