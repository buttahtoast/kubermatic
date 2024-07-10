/*
Copyright 2020 The Kubermatic Kubernetes Platform contributors.

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

package flatcar

import (
	"context"
	"fmt"

	kubermaticv1 "k8c.io/kubermatic/v2/pkg/apis/kubermatic/v1"
	userclustercontrollermanager "k8c.io/kubermatic/v2/pkg/controller/user-cluster-controller-manager"
	"k8c.io/kubermatic/v2/pkg/controller/user-cluster-controller-manager/flatcar/resources"
	nodelabelerapi "k8c.io/kubermatic/v2/pkg/controller/user-cluster-controller-manager/node-labeler/api"
	controllerutil "k8c.io/kubermatic/v2/pkg/controller/util"
	predicateutil "k8c.io/kubermatic/v2/pkg/controller/util/predicate"
	"k8c.io/kubermatic/v2/pkg/provider/kubernetes"
	"k8c.io/kubermatic/v2/pkg/resources/registry"
	"k8c.io/reconciler/pkg/reconciling"
	"sigs.k8s.io/controller-runtime/pkg/cache"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	ctrlruntimeclient "sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

const (
	// This controller is responsible for ensuring that the flatcar-linux-update-operator is installed when we have a healthy(running) flatcar
	// node in our cluster.
	ControllerName = "kkp-flatcar-controller"
)

type Reconciler struct {
	ctrlruntimeclient.Client
	seedClient        ctrlruntimeclient.Client
	namespace         string
	clusterIsPaused   userclustercontrollermanager.IsPausedChecker
	overwriteRegistry string
	updateWindow      kubermaticv1.UpdateWindow
	networkPolices    bool
	cache             cache.Cache
}

func Add(mgr manager.Manager, seedMgr manager.Manager, overwriteRegistry string, updateWindow kubermaticv1.UpdateWindow, clusterIsPaused userclustercontrollermanager.IsPausedChecker, networkPolices bool) error {
	reconciler := &Reconciler{
		Client:            mgr.GetClient(),
		overwriteRegistry: overwriteRegistry,
		updateWindow:      updateWindow,
		clusterIsPaused:   clusterIsPaused,
		networkPolices:    networkPolices,
	}

	reconciler.Client = mgr.GetClient()
	reconciler.seedClient = seedMgr.GetClient()
	reconciler.cache = mgr.GetCache()

	_, err := builder.ControllerManagedBy(mgr).
		Named(ControllerName).
		Watches(&corev1.Node{}, controllerutil.EnqueueConst(""), builder.WithPredicates(predicateutil.ByLabel(nodelabelerapi.DistributionLabelKey, nodelabelerapi.FlatcarLabelValue))).
		Build(reconciler)

	return err
}

func (r *Reconciler) Reconcile(ctx context.Context, _ reconcile.Request) (reconcile.Result, error) {
	paused, err := r.clusterIsPaused(ctx)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to check cluster pause status: %w", err)
	}
	if paused {
		return reconcile.Result{}, nil
	}

	cluster, err := r.getCluster(ctx)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to retrieve cluster: %w", err)
	}

	data := reconcileData{
		caCert:       caCert,
		userSSHKeys:  userSSHKeys,
		ccmMigration: r.ccmMigration || r.ccmMigrationCompleted,
		cluster:      cluster,
	}

	err = r.setupNetworkingData(cluster, &data)
	if err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to setup cluster networking data: %w", err)
	}

	var nodeList corev1.NodeList
	if err := r.List(ctx, &nodeList,
		ctrlruntimeclient.MatchingLabels{nodelabelerapi.DistributionLabelKey: nodelabelerapi.FlatcarLabelValue},
	); err != nil {
		return reconcile.Result{}, fmt.Errorf("failed to list nodes: %w", err)
	}

	// filter out any Flatcar nodes that are already being deleted
	var nodes []corev1.Node
	for _, node := range nodeList.Items {
		if node.ObjectMeta.DeletionTimestamp == nil {
			nodes = append(nodes, node)
		}
	}

	if len(nodes) == 0 {
		if err := r.cleanupUpdateOperatorResources(ctx); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to clean up UpdateOperator resources: %w", err)
		}
	} else {
		if err := r.reconcileUpdateOperatorResources(ctx); err != nil {
			return reconcile.Result{}, fmt.Errorf("failed to reconcile the UpdateOperator resources: %w", err)
		}
	}

	return reconcile.Result{}, nil
}

func (r *Reconciler) cleanupUpdateOperatorResources(ctx context.Context) error {
	return resources.EnsureAllDeleted(ctx, r.Client)
}

// reconcileUpdateOperatorResources deploys the FlatcarUpdateOperator
// https://github.com/kinvolk/flatcar-linux-update-operator
func (r *Reconciler) reconcileUpdateOperatorResources(ctx context.Context) error {
	saReconcilers := []reconciling.NamedServiceAccountReconcilerFactory{
		resources.OperatorServiceAccountReconciler(),
		resources.AgentServiceAccountReconciler(),
	}
	if err := reconciling.ReconcileServiceAccounts(ctx, saReconcilers, metav1.NamespaceSystem, r.Client); err != nil {
		return fmt.Errorf("failed to reconcile the ServiceAccounts: %w", err)
	}

	crReconcilers := []reconciling.NamedClusterRoleReconcilerFactory{
		resources.OperatorClusterRoleReconciler(),
		resources.AgentClusterRoleReconciler(),
	}
	if err := reconciling.ReconcileClusterRoles(ctx, crReconcilers, metav1.NamespaceNone, r.Client); err != nil {
		return fmt.Errorf("failed to reconcile the ClusterRoles: %w", err)
	}

	crbReconcilers := []reconciling.NamedClusterRoleBindingReconcilerFactory{
		resources.OperatorClusterRoleBindingReconciler(),
		resources.AgentClusterRoleBindingReconciler(),
	}
	if err := reconciling.ReconcileClusterRoleBindings(ctx, crbReconcilers, metav1.NamespaceNone, r.Client); err != nil {
		return fmt.Errorf("failed to reconcile the ClusterRoleBindings: %w", err)
	}

	depReconcilers := []reconciling.NamedDeploymentReconcilerFactory{
		resources.OperatorDeploymentReconciler(registry.GetImageRewriterFunc(r.overwriteRegistry), r.updateWindow),
	}
	if err := reconciling.ReconcileDeployments(ctx, depReconcilers, metav1.NamespaceSystem, r.Client); err != nil {
		return fmt.Errorf("failed to reconcile the Deployments: %w", err)
	}

	dsReconcilers := []reconciling.NamedDaemonSetReconcilerFactory{
		resources.AgentDaemonSetReconciler(registry.GetImageRewriterFunc(r.overwriteRegistry)),
	}
	if err := reconciling.ReconcileDaemonSets(ctx, dsReconcilers, metav1.NamespaceSystem, r.Client); err != nil {
		return fmt.Errorf("failed to reconcile the DaemonSets: %w", err)
	}

	if r.networkPolices {
		netReconcilers := []reconciling.NamedNetworkPolicyReconcilerFactory{
			resources.NetworkPolicyReconciler(data),
		}
		if err := reconciling.ReconcileNetworkPolicies(ctx, netReconcilers, metav1.NamespaceSystem, r.Client); err != nil {
			return fmt.Errorf("failed to reconcile the NetworkPolices: %w", err)
		}
	}

	return nil
}

func (r *Reconciler) getCluster(ctx context.Context) (*kubermaticv1.Cluster, error) {
	cluster, err := kubernetes.ClusterFromNamespace(ctx, r.seedClient, r.namespace)
	if err != nil {
		return nil, err
	}
	if cluster == nil {
		return nil, fmt.Errorf("no cluster exists for namespace %q", r.namespace)
	}

	return cluster, nil
}

func (r *Reconciler) setupNetworkingData(cluster *kubermaticv1.Cluster, data *reconcileData) (err error) {
	data.k8sServiceApiIP, err = resources.InClusterApiserverIP(cluster)
	if err != nil {
		return fmt.Errorf("failed to get Cluster Apiserver IP: %w", err)
	}

	// Reconcile kubernetes service endpoints, unless it is not supported or disabled in the apiserver override settings.
	data.reconcileK8sSvcEndpoints = true

	if cluster.Spec.ComponentsOverride.Apiserver.EndpointReconcilingDisabled != nil && *cluster.Spec.ComponentsOverride.Apiserver.EndpointReconcilingDisabled {
		// Do not reconcile if explicitly disabled.
		data.reconcileK8sSvcEndpoints = false
	}

	if cluster.Spec.ExposeStrategy == kubermaticv1.ExposeStrategyTunneling {
		data.k8sServiceEndpointAddress = r.tunnelingAgentIP.String()
		data.k8sServiceEndpointPort = int32(r.kasSecurePort)
	} else {
		data.k8sServiceEndpointAddress = cluster.Status.Address.IP
		data.k8sServiceEndpointPort = cluster.Status.Address.Port
	}
	data.ipFamily = cluster.Spec.ClusterNetwork.IPFamily
	return nil
}
