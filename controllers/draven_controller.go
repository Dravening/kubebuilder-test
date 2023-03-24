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

package controllers

import (
	"context"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"

	batchv1 "kubebuilder-test/api/v1"
)

// DravenReconciler reconciles a Draven object
type DravenReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=batch.my.domain,resources=dravens,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=batch.my.domain,resources=dravens/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=batch.my.domain,resources=dravens/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Draven object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *DravenReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logDraven := log.FromContext(ctx)
	logDraven.Info("reconciling foo custom resource")

	var foo batchv1.Draven
	if err := r.Get(ctx, req.NamespacedName, &foo); err != nil {
		logDraven.Error(err, "unable to list pods")
	}

	var podList corev1.PodList
	var friendFound bool
	if err := r.List(ctx, &podList); err != nil {
		logDraven.Error(err, "unable to list pods")
	} else {
		for _, item := range podList.Items {
			if item.GetName() == foo.Spec.Name {
				logDraven.Info("pod linked to a foo custom resource found", "name", item.GetName())
				friendFound = true
			}
		}
	}

	if friendFound {
		foo.Status.Happy = "发现了Draven"
	}
	if err := r.Status().Update(ctx, &foo); err != nil {
		logDraven.Error(err, "unable to update foo's happy status", "status", friendFound)
		return ctrl.Result{}, err
	}
	logDraven.Info("foo's happy status updated", "status", friendFound)
	logDraven.Info("foo custom resource reconciled")

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *DravenReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&batchv1.Draven{}).
		Watches(
			&source.Kind{Type: &corev1.Pod{}},
			handler.EnqueueRequestsFromMapFunc(r.mapPodsReqToFooReq),
		).
		Complete(r)
}

func (r *DravenReconciler) mapPodsReqToFooReq(obj client.Object) []reconcile.Request {
	ctx := context.Background()
	logDraven := log.FromContext(ctx)

	// List all the Foo custom resource
	req := []reconcile.Request{}
	var list batchv1.DravenList
	if err := r.Client.List(context.TODO(), &list); err != nil {
		logDraven.Error(err, "unable to list foo custom resources")
	} else {
		// Only keep Foo custom resources related to the Pod that triggered the reconciliation request
		for _, item := range list.Items {
			if item.Spec.Name == obj.GetName() {
				req = append(req, reconcile.Request{
					NamespacedName: types.NamespacedName{Name: item.Name, Namespace: item.Namespace},
				})
				logDraven.Info("pod linked to a foo custom resource issued an event", "name", obj.GetName())
			}
		}
	}
	return req
}
