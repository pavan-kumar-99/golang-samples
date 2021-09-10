/*
Copyright 2021.

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
	"fmt"
	"time"

	acmautosynciov1beta1 "pavan-kumar-99/k8s-acm-autosync/api/v1beta1"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// AcmAutoSyncReconciler reconciles a AcmAutoSync object
type AcmAutoSyncReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	recorder record.EventRecorder
}

//+kubebuilder:rbac:groups=acm-autosync.io,resources=acmautosyncs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=acm-autosync.io,resources=acmautosyncs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=acm-autosync.io,resources=acmautosyncs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the AcmAutoSync object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *AcmAutoSyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logs := log.FromContext(ctx)
	var acmas acmautosynciov1beta1.AcmAutoSync
	var acmasSecret corev1.Secret
	if err := r.Get(ctx, req.NamespacedName, &acmas); err != nil {
		logs.Info("AcmAutoSync Not found", "Error", err)
		// fmt.Println(acmas.Spec.SecretName)
		custom := types.NamespacedName{
			Name:      acmas.Spec.SecretName,
			Namespace: acmas.Namespace,
		}
		fmt.Println(custom.Name, custom.Namespace)
		if err = r.Get(ctx, custom, &acmasSecret); err == nil {
			logs.Info("Found Secret", "Secret Details", custom)
			return r.RemoveSecret(ctx, &acmasSecret, logs)
		}
		if err = r.Get(ctx, custom, &acmasSecret); err != nil {
			logs.Info("Not Found Secret", "Secret Details", custom)
			return r.RemoveSecret(ctx, &acmasSecret, logs)
		}
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	custom := types.NamespacedName{
		Name:      acmas.Spec.SecretName,
		Namespace: acmas.Namespace,
	}
	if err := r.Get(ctx, custom, &acmasSecret); err == nil {
		logs.Info("Found Secret", "Secret Details", custom)
		return r.RemoveSecret(ctx, &acmasSecret, logs)
	}
	if err := r.Get(ctx, custom, &acmasSecret); err != nil {
		logs.Info("unable to fetch Secret for AcmAutoSync", "AcmAutoSync", custom)
		return r.CreateSecret(ctx, req, acmas, logs)
	}
	if err := r.Get(ctx, custom, &acmasSecret); err == nil {
		logs.Info("Found Secret", "AcmAutoSync", custom)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	return ctrl.Result{}, nil
}

func (r *AcmAutoSyncReconciler) RemoveSecret(ctx context.Context, delsecname *corev1.Secret, log logr.Logger) (ctrl.Result, error) {
	name := delsecname.Name
	if err := r.Delete(ctx, delsecname); err != nil {
		log.Error(err, "unable to delete secret for secret", "AcmAutoSync", delsecname.Name)
		return ctrl.Result{}, err
	}
	log.Info("Deleted the secret for", "AcmAutoSync", name)
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	//return ctrl.Result{}, nil
}

func (r *AcmAutoSyncReconciler) CreateSecret(ctx context.Context, req ctrl.Request, acmas acmautosynciov1beta1.AcmAutoSync, log logr.Logger) (ctrl.Result, error) {
	secret := &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      acmas.Spec.SecretName,
			Namespace: acmas.Namespace,
		},
		StringData: map[string]string{
			"username": "Pavan",
			"Password": "Password123@",
		},
		Type: corev1.SecretType("Opqaue"),
	}
	if err := r.Create(ctx, secret); err != nil {
		log.Error(err, "unable to create secret for AcmAutoSync", "AcmAutoSync", secret)
		return ctrl.Result{}, err
	}
	log.Info("created secret for AcmAutoSync", "AcmAutoSync", secret)
	r.UpdateStatus(ctx, req, acmas, secret)
	r.recorder.Event(&acmas, corev1.EventTypeNormal, "Created", "Created Secreet")
	return ctrl.Result{RequeueAfter: 10 * time.Second}, nil
	//return ctrl.Result{}, nil
}

func (r *AcmAutoSyncReconciler) UpdateStatus(ctx context.Context, req ctrl.Request, acmas acmautosynciov1beta1.AcmAutoSync, secname *corev1.Secret) {
	// Here we are only maintaining the current replica count for the status, more can be done.
	fmt.Println("Updating status........")
	fmt.Println(secname.CreationTimestamp)
	acmas.Status.SecretCreationTime = secname.CreationTimestamp
	fmt.Println(acmas.Status.SecretCreationTime)
}

// SetupWithManager sets up the controller with the Manager.
func (r *AcmAutoSyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.recorder = mgr.GetEventRecorderFor("AcmAutoSync")
	return ctrl.NewControllerManagedBy(mgr).
		For(&acmautosynciov1beta1.AcmAutoSync{}).
		Complete(r)
}
