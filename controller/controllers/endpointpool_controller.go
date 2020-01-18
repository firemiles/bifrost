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
	"context"
	"fmt"

	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/rand"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	networkv1 "github.com/firemiles/bifrost/controller/api/v1"
)

// EndpointPoolReconciler reconciles a EndpointPool object
type EndpointPoolReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=network.crd.firemiles.top,resources=endpointpools,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=network.crd.firemiles.top,resources=endpointpools/status,verbs=get;update;patch

func (r *EndpointPoolReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("endpointpool", req.NamespacedName)

	var endpointPool networkv1.EndpointPool
	if err := r.Get(ctx, req.NamespacedName, &endpointPool); err != nil {
		log.Error(err, "unable to fetch EndpointPool")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	var endpoints networkv1.EndpointList
	if err := r.List(ctx, &endpoints, client.MatchingFields{ownerKey: req.Name}); err != nil {
		log.Error(err, "unable to list endpoints in pool %s/%s", endpointPool.Namespace, endpointPool.Name)
		return ctrl.Result{}, err
	}

	constructEndpointForEndpointPool := func(endpointPool *networkv1.EndpointPool) (*networkv1.Endpoint, error) {
		name := fmt.Sprintf("%s-%s", endpointPool.Name, rand.String(6))

		endpoint := &networkv1.Endpoint{
			ObjectMeta: metav1.ObjectMeta{
				Labels:      make(map[string]string),
				Annotations: make(map[string]string),
				Name:        name,
			},
			Spec: networkv1.EndpointSpec{
				IPs: []networkv1.FixedIP{
					{Subnet: "xxxx", IP: "xxxx"},
				},
			},
			Status: networkv1.EndpointStatus{},
		}
		if err := ctrl.SetControllerReference(endpointPool, endpoint, r.Scheme); err != nil {
			return nil, err
		}
		return endpoint, nil
	}

	count := 0
	for i := len(endpoints.Items); i < endpointPool.Spec.PoolSize; i++ {
		ep, err := constructEndpointForEndpointPool(&endpointPool)
		if err != nil {
			log.Error(err, "construct endpoint for pool %s/%s failed", endpointPool.Namespace, endpointPool.Name)
			// don't bother requeuing until we get a change to the spec
			return ctrl.Result{}, nil
		}
		if err := r.Create(ctx, ep); err != nil {
			log.Error(err, "unable to create Endpoint for EndpointPool", "endpoint", ep)
			return ctrl.Result{}, err
		}
		count += 1
	}

	endpointPool.Status.AvailableEndpoints = len(endpoints.Items) + count
	endpointPool.Status.Phase = networkv1.PhaseAvailable

	if err := r.Status().Update(ctx, &endpointPool); err != nil {
		log.Error(err, "unable to update EndpointPool status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *EndpointPoolReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(&networkv1.Endpoint{}, ownerKey, func(rawObj runtime.Object) []string {
		endpoint := rawObj.(*networkv1.Endpoint)
		owner := metav1.GetControllerOf(endpoint)
		if owner == nil {
			return nil
		}
		if owner.APIVersion != apiGVStr || owner.Kind != "EndpointPool" {
			return nil
		}
		return []string{owner.Name}

	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&networkv1.EndpointPool{}).
		Owns(&networkv1.Endpoint{}).
		Complete(r)
}
