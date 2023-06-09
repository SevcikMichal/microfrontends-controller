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
	"sync"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	microfrontendv1alpha1 "github.com/SevcikMichal/microfrontends-controller/api/v1alpha1"
	"github.com/SevcikMichal/microfrontends-controller/internal/model"
)

const webComponentFinalizer = "microfrontend.michalsevcik.dev/finalizer"

// WebComponentReconciler reconciles a WebComponent object
type WebComponentReconciler struct {
	client.Client
	Scheme               *runtime.Scheme
	Recorder             record.EventRecorder
	MicroFrontendConfigs *sync.Map
}

//+kubebuilder:rbac:groups=microfrontend.michalsevcik.dev,resources=webcomponents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=microfrontend.michalsevcik.dev,resources=webcomponents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=microfrontend.michalsevcik.dev,resources=webcomponents/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=events,verbs=create;patch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the WebComponent object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *WebComponentReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// Fetch the WebComponent instance
	// The purpose is check if the Custom Resource for the Kind WebComponent
	// is applied on the cluster if not we return nil to stop the reconciliation
	webComponent := &microfrontendv1alpha1.WebComponent{}
	err := r.Get(ctx, req.NamespacedName, webComponent)
	if err != nil {
		if apierrors.IsNotFound(err) {
			// If the custom resource is not found then, it usually means that it was deleted or not created
			// In this way, we will stop the reconciliation
			log.Info("WebComponent resource not found. Ignoring since object must be deleted!")
			return ctrl.Result{}, nil
		}

		// Error reading the object - requeue the request.
		log.Error(err, "Failed to get WebComponent!")
		return ctrl.Result{}, err
	}

	// Let's add a finalizer. Then, we can define some operations which should
	// occurs before the custom resource to be deleted.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !controllerutil.ContainsFinalizer(webComponent, webComponentFinalizer) {
		log.Info("Adding Finalizer for WebComponent.")
		if ok := controllerutil.AddFinalizer(webComponent, webComponentFinalizer); !ok {
			log.Error(err, "Failed to add finalizer into the custom resource!")
			return ctrl.Result{Requeue: true}, nil
		}

		if err = r.Update(ctx, webComponent); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer!")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// Check if the Memcached instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isWebComponentMarkedToBeDeleted := webComponent.GetDeletionTimestamp() != nil
	if isWebComponentMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(webComponent, webComponentFinalizer) {
			log.Info("Performing finalizer operations for the WebComponent before deleting the custom resource.")

			if err := r.doFinalizerOperationsForWebComponent(webComponent); err != nil {
				log.Error(err, "Failed to perform finalizer operations for the WebComponent!")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Get(ctx, req.NamespacedName, webComponent); err != nil {
				log.Error(err, "Failed to re-fetch WebComponent!")
				return ctrl.Result{}, err
			}

			log.Info("Removing Finalizer for WebComponent after successfully performing the operations.")
			if ok := controllerutil.RemoveFinalizer(webComponent, webComponentFinalizer); !ok {
				log.Error(err, "Failed to remove finalizer for WebComponent!")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Update(ctx, webComponent); err != nil {
				log.Error(err, "Failed to remove finalizer for WebComponent!")
				return ctrl.Result{}, err
			}

			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, nil
	}

	r.addWebComponentSpecToMicroFrontendConfigs(webComponent)

	// TODO: Maybe create a deployment and a service here eventually

	return ctrl.Result{}, nil
}

func (r *WebComponentReconciler) doFinalizerOperationsForWebComponent(webComponent *microfrontendv1alpha1.WebComponent) error {
	log := log.FromContext(context.Background())
	log.Info("Removing web component from frontend configs.")
	r.MicroFrontendConfigs.Delete(webComponent.ObjectMeta.UID)
	return nil
}

func (r *WebComponentReconciler) addWebComponentSpecToMicroFrontendConfigs(webComponent *microfrontendv1alpha1.WebComponent) {
	log := log.FromContext(context.Background())
	log.Info("Converting web component to frontend config.")
	frontendConfig := model.CreateFrontendConfigFromWebComponent(webComponent)
	log.Info("WebComponent converted to frontend config.", "UID", webComponent.ObjectMeta.UID)
	r.MicroFrontendConfigs.Store(webComponent.ObjectMeta.UID, frontendConfig)
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebComponentReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&microfrontendv1alpha1.WebComponent{}).
		Complete(r)
}
