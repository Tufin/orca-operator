package orca

import (
	"context"

	appv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"
	tufinv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"fmt"
	"reflect"
)

var log = logf.Log.WithName("controller_orca")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

const (
	StatusUnknown  = "Unknown"
	StatusCreating = "Creating"
	StatusReady    = "Ready"
	StatusFailed   = "Failed"
)

func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileOrca{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("orca-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Orca
	err = c.Watch(&source.Kind{Type: &tufinv1alpha1.Orca{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner Orca
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &tufinv1alpha1.Orca{},
	})
	if err != nil {
		return err
	}

	return nil
}

var _ reconcile.Reconciler = &ReconcileOrca{}

// ReconcileOrca reconciles a Orca object
type ReconcileOrca struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

type ResourceRequest struct {
	Required       metav1.Object
	RequiredStruct runtime.Object
}

// Reconcile reads that state of the cluster for a Orca object and makes changes based on the state read
// and what is in the Orca.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileOrca) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Orca")

	// Fetch the Orca instance
	instance := &tufinv1alpha1.Orca{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		reqLogger.Error(err, "failed to fetch CRD")
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}

		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	if instance.Status.Ready != StatusReady && instance.Status.Ready != "" {
		return reconcile.Result{}, nil
	}

	if instance.Spec.Namespace == "" {
		instance.Namespace = "default"
	} else {
		instance.Namespace = instance.Spec.Namespace
	}

	kiteDeployment := getKiteDeployment(instance)
	kiteService := getKiteService(instance)
	conntrackDaemonset := getConntrackDaemonset(instance)

	reconcileResult, err := r.createResourceArray(instance,
		ResourceRequest{Required: kiteDeployment, RequiredStruct: &appsv1.Deployment{}},
		ResourceRequest{Required: kiteService, RequiredStruct: &corev1.Service{}},
		ResourceRequest{Required: conntrackDaemonset, RequiredStruct: &appsv1.DaemonSet{}},
	)

	reqLogger.Info("Orca was successfully deployed in the cluster!")
	return reconcileResult, nil
}

func (r *ReconcileOrca) createResourceArray(instance *appv1alpha1.Orca, resources ...ResourceRequest) (reconcile.Result, error) {

	var reconcileResult reconcile.Result
	var err error
	instance.Status.Ready = StatusCreating

	for _, resourceRequest := range resources {
		reconcileResult, err = r.createResource(instance, resourceRequest.Required, resourceRequest.RequiredStruct)
		if err != nil {
			r.UpdateStatus(instance, StatusFailed)
			return reconcileResult, err
		}
	}

	r.UpdateStatus(instance, StatusReady)

	return reconcile.Result{}, nil
}

func (r *ReconcileOrca) UpdateStatus(instance *appv1alpha1.Orca, status string) error {

	var err error

	instance.Status.Ready = status
	err = r.client.Status().Update(context.TODO(), instance)

	return err
}

func (r *ReconcileOrca) createResource(instance *appv1alpha1.Orca, required metav1.Object, requiredStruct runtime.Object) (reconcile.Result, error) {

	reqLogger := log.WithValues("Kind", fmt.Sprintf("%T", requiredStruct), "Namespace", required.GetNamespace(), "Resource Name", required.GetName())
	ns := required.GetNamespace()

	if instance.Status.Ready != StatusCreating {
		return reconcile.Result{}, nil
	}

	if err := controllerutil.SetControllerReference(instance, required, r.scheme); err != nil {
		reqLogger.Error(err, "Failed to set the operator as the resource owner")
		return reconcile.Result{}, err
	}

	err := r.client.Get(context.TODO(), types.NamespacedName{Name: required.GetName(), Namespace: ns}, requiredStruct)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating Resource...")
		err = r.client.Create(context.TODO(), required.(runtime.Object))
		if err != nil {
			reqLogger.Error(err, "Resource creation failed")
			return reconcile.Result{}, err
		}

		reqLogger.Info("Resource created successfully")
		return reconcile.Result{}, nil
	} else if err != nil {

		return reconcile.Result{}, err
	} else {
		reqLogger.Info("Resource already exists, trying to update...")

		if reflect.DeepEqual(required, requiredStruct.(metav1.Object)) {
			reqLogger.Info("Resource is already up to date")
			return reconcile.Result{}, nil
		}

		err = r.client.Delete(context.TODO(), requiredStruct)
		if err != nil {
			reqLogger.Error(err, "Resource update failed, deletion failed")
			return reconcile.Result{}, err
		}
		err = r.client.Create(context.TODO(), required.(runtime.Object))
		if err != nil {
			reqLogger.Error(err, "Resource update failed, creation failed")
			return reconcile.Result{}, err
		}

		reqLogger.Info("Resource update succeeded")
	}

	return reconcile.Result{}, nil
}
