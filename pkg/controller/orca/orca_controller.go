package orca

import (
	"context"

	appv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"
	tufinv1alpha1 "github.com/tufin/orca-operator/pkg/apis/tufin/v1alpha1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"

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
)

var log = logf.Log.WithName("controller_orca")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Orca Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
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
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// service accounts
	kiteNameLabel := GetLabels(name + "=" + kite)

	kiteSA := GetServiceAccount(kite, instance.Namespace, kiteNameLabel)
	conntrackSA := GetServiceAccount(conntrack, instance.Namespace, GetLabels(name+"="+conntrack))

	kiteCR := GetClusterRole(kite, kiteNameLabel,
		GetPolicyRule([]string{VerbAll}, []string{"networking.k8s.io"}, []string{"networkpolicies"}),
		GetPolicyRule([]string{VerbAll}, []string{"config.istio.io"}, []string{VerbAll}),
		GetPolicyRule([]string{VerbAll}, []string{"istio.io"}, []string{"istioconfigs", "istioconfigs.istio.io"}),
		GetPolicyRule([]string{VerbAll}, []string{"apiextensions.k8s.io"}, []string{"customresourcedefinitions"}),
		GetPolicyRule([]string{VerbAll}, []string{"extensions"}, []string{"thirdpartyresources", "thirdpartyresources.extensions", ResourceIngresses, "ingresses/status"}),
		GetPolicyRule([]string{VerbAll}, []string{""}, []string{ResourceConfigMaps}),
		GetPolicyRule([]string{VerbGet, VerbList, VerbWatch}, []string{""}, []string{ResourceEndpoints, ResourceNamespaces, ResourceNodes, ResourcePods, ResourceServices, ResourceSecrets}),
	)

	conntrackCR := GetClusterRole(conntrack, GetLabels(name+"="+conntrack), GetPolicyRule([]string{VerbAll}, []string{VerbAll}, []string{VerbAll}))

	kiteCRB := GetClusterRoleBindig(kite, kiteSA, kiteCR)
	conntrackCRB := GetClusterRoleBindig(conntrack, conntrackSA, conntrackCR)

	deployment := getKiteDeployment(instance)
	service := getKiteService(instance)
	daemonset := getConntrackDaemonset(instance)

	reconcileResult, err := r.createServiceAccount(instance, kiteSA, request)
	if err != nil {
		return reconcileResult, err
	}

	reconcileResult, err = r.createKiteService(instance, service, request)
	if err != nil {
		return reconcileResult, err
	}

	reconcileResult, err = r.createServiceAccount(instance, conntrackSA, request)
	if err != nil {
		return reconcileResult, err
	}

	reconcileResult, err = r.createClusterRole(instance, kiteCR, request)
	if err != nil {
		return reconcileResult, err
	}

	reconcileResult, err = r.createClusterRole(instance, conntrackCR, request)
	if err != nil {
		return reconcileResult, err
	}

	reconcileResult, err = r.createClusterRoleBinding(instance, kiteCRB, request)
	if err != nil {
		return reconcileResult, err
	}

	reconcileResult, err = r.createClusterRoleBinding(instance, conntrackCRB, request)
	if err != nil {
		return reconcileResult, err
	}

	reconcileResult, err = r.createKiteDeployment(instance, deployment, request)
	if err != nil {
		return reconcileResult, err
	}

	reconcileResult, err = r.createConntrackDaemonset(instance, daemonset, request)
	return reconcileResult, err
}

func (r *ReconcileOrca) createKiteDeployment(instance *appv1alpha1.Orca, deployment *appsv1.Deployment, request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Deploying Kite")

	// Set Kite instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, deployment, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the Deployment already exists
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: deployment.Name, Namespace: deployment.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new deployment", "deployment.Namespace", deployment.Namespace, "deployment.Name", deployment.Name)
		err = r.client.Create(context.TODO(), deployment)
		if err != nil {
			return reconcile.Result{}, err
		}

		// deployment created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// deployment already exists - don't requeue
	reqLogger.Info("Skip reconcile: deployment already exists", "deployment.Namespace", found.Namespace, "deployment.Name", found.Name)
	return reconcile.Result{}, nil

}

func (r *ReconcileOrca) createKiteService(instance *appv1alpha1.Orca, service *corev1.Service, request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Deploying Kite Service")

	// Set Kite instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the Deployment already exists
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new deployment", "service.Namespace", service.Namespace, "service.Name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}

		// deployment created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// deployment already exists - don't requeue
	reqLogger.Info("Skip reconcile: service already exists", "service.Namespace", found.Namespace, "service.Name", found.Name)
	return reconcile.Result{}, nil

}

func (r *ReconcileOrca) createServiceAccount(instance *appv1alpha1.Orca, serviceAccount *corev1.ServiceAccount, request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Deploying Service Account: " + serviceAccount.Name)

	// Set Service Account instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, serviceAccount, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the Service Account already exists
	found := &corev1.ServiceAccount{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: serviceAccount.Name, Namespace: serviceAccount.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new service account", "serviceAccount.Namespace", serviceAccount.Namespace, "serviceAccount.Name", serviceAccount.Name)
		err = r.client.Create(context.TODO(), serviceAccount)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Service Account created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Service Account already exists - don't requeue
	reqLogger.Info("Skip reconcile: Service Account already exists", "serviceAccount.Namespace", found.Namespace, "serviceAccount.Name", found.Name)
	return reconcile.Result{}, nil

}

func (r *ReconcileOrca) createClusterRole(instance *appv1alpha1.Orca, clusterRole *rbacv1.ClusterRole, request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Deploying Cluster Role: " + clusterRole.Name)

	// Set clusterRole instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, clusterRole, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the clusterRole already exists
	found := &rbacv1.ClusterRole{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: clusterRole.Name, Namespace: clusterRole.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new clusterRole", "clusterRole.Namespace", clusterRole.Namespace, "clusterRole.Name", clusterRole.Name)
		err = r.client.Create(context.TODO(), clusterRole)
		if err != nil {
			return reconcile.Result{}, err
		}

		// clusterRole created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// clusterRole already exists - don't requeue
	reqLogger.Info("Skip reconcile: clusterRole already exists", "clusterRole.Namespace", found.Namespace, "clusterRole.Name", found.Name)
	return reconcile.Result{}, nil

}

func (r *ReconcileOrca) createClusterRoleBinding(instance *appv1alpha1.Orca, clusterRoleBinding *rbacv1.ClusterRoleBinding, request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Deploying Cluster Role Binding: " + clusterRoleBinding.Name)

	// Set clusterRoleBinding instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, clusterRoleBinding, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the clusterRoleBinding already exists
	found := &rbacv1.ClusterRoleBinding{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: clusterRoleBinding.Name, Namespace: ""}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new clusterRoleBinding", "clusterRoleBinding.Namespace", clusterRoleBinding.Namespace, "clusterRoleBinding.Name", clusterRoleBinding.Name)
		err = r.client.Create(context.TODO(), clusterRoleBinding)
		if err != nil {
			return reconcile.Result{}, err
		}

		// clusterRoleBinding created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// clusterRoleBinding already exists - don't requeue
	reqLogger.Info("Skip reconcile: clusterRoleBinding already exists", "clusterRoleBinding.Namespace", found.Namespace, "clusterRoleBinding.Name", found.Name)
	return reconcile.Result{}, nil

}

func (r *ReconcileOrca) createConntrackDaemonset(instance *appv1alpha1.Orca, daemonset *appsv1.DaemonSet, request reconcile.Request) (reconcile.Result, error) {

	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Deploying Conntrack")

	// Set Kite instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, daemonset, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if the Daemonset already exists
	found := &appsv1.DaemonSet{}
	err := r.client.Get(context.TODO(), types.NamespacedName{Name: daemonset.Name, Namespace: daemonset.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new daemonset", "daemonset.Namespace", daemonset.Namespace, "daemonset.Name", daemonset.Name)
		err = r.client.Create(context.TODO(), daemonset)
		if err != nil {
			return reconcile.Result{}, err
		}

		// daemonset created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// daemonset already exists - don't requeue
	reqLogger.Info("Skip reconcile: daemonset already exists", "daemonset.Namespace", found.Namespace, "deployment.Name", found.Name)
	return reconcile.Result{}, nil

}
