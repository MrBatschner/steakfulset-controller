// This file is licensed under the Apache Software License, v.2.0 except as
// noted otherwise in the LICENSE file
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	"context"
	"fmt"
	"math/rand"

	"k8s.io/apimachinery/pkg/runtime"
	krand "k8s.io/apimachinery/pkg/util/rand"
	ref "k8s.io/client-go/tools/reference"

	ctrl "sigs.k8s.io/controller-runtime"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	foodv1alpha1 "k8s.training/steakulset-controller/api/v1alpha1"
)

var (
	steakOwnerKey = ".metadata.controller"
	apiGVStr      = foodv1alpha1.GroupVersion.String()
)

// SteakfulSetReconciler reconciles a SteakfulSet object
type SteakfulSetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=food.k8s.training,resources=steakfulsets,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=food.k8s.training,resources=steakfulsets/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=food.k8s.training,resources=steakfulsets/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.1/pkg/reconcile
func (r *SteakfulSetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	var steakfulSet foodv1alpha1.SteakfulSet

	if err := r.Get(ctx, req.NamespacedName, &steakfulSet); err != nil {
		if client.IgnoreNotFound(err) == nil {
			logger.V(1).Info(fmt.Sprintf("SteakfulSet %s got deleted", steakfulSet.Name))
			return ctrl.Result{}, nil
		}

		logger.Error(err, "unable to fetch SteakfulSet")
		return ctrl.Result{}, err
	}

	logger.V(1).Info(fmt.Sprintf("Reconciling SteakfulSet %s...", steakfulSet.Name))

	// get a list of Steaks that have been cooked so far by this SteakfulSet
	var steakList foodv1alpha1.SteakList

	if err := r.List(ctx, &steakList, client.InNamespace(req.Namespace), client.MatchingFields{steakOwnerKey: req.Name}); client.IgnoreNotFound(err) != nil {
		logger.Error(err, "unable to list child Steaks")
		return ctrl.Result{}, err
	}

	// check if we have enough Steaks ready
	actualSteakCount := len(steakList.Items)
	desiredSteakCount := steakfulSet.Spec.Guests

	logger.V(1).Info(fmt.Sprintf("Found %d steaks belonging to this SteakfulSet, should have %d Steaks", actualSteakCount, desiredSteakCount))

	if actualSteakCount > desiredSteakCount && actualSteakCount > 0 {
		// we have too many Steaks, we need to discard some...
		for actualSteakCount > desiredSteakCount {
			steak := steakList.Items[actualSteakCount-1]

			logger.V(1).Info(fmt.Sprintf("Discarding charred Steak %s...", steak.Name))

			if err := r.Delete(ctx, &steak, client.PropagationPolicy(metav1.DeletePropagationForeground)); client.IgnoreNotFound(err) != nil {
				logger.Error(err, "unable to delete charred Steak", "steak", steak)
				return ctrl.Result{}, err
			}

			steakList.Items = steakList.Items[:actualSteakCount-1]
			actualSteakCount -= 1
		}
	} else {
		// we have too little Steaks, we need to cook some...
		for actualSteakCount < desiredSteakCount {
			steak := foodv1alpha1.Steak{
				ObjectMeta: metav1.ObjectMeta{
					Labels:      make(map[string]string),
					Annotations: make(map[string]string),
					Name:        fmt.Sprintf("%s-%s", steakfulSet.Name, krand.String(6)),
					Namespace:   steakfulSet.Namespace,
				},
				Spec: *steakfulSet.Spec.Steak.Spec.DeepCopy(),
			}

			// a Steak never has the desired weight so we introduce some randomness
			weightVariance := steak.Spec.Weight / 40
			steak.Spec.Weight = steak.Spec.Weight - weightVariance + rand.Intn(weightVariance*2)

			logger.V(1).Info(fmt.Sprintf("Cooking new %s Steak %s with %d grams...", steak.Spec.Variant, steak.Name, steak.Spec.Weight))

			if err := ctrl.SetControllerReference(&steakfulSet, &steak, r.Scheme); err != nil {
				logger.Error(err, "unable to set Steak owner reference")
				return ctrl.Result{}, err
			}

			if err := r.Create(ctx, &steak); err != nil {
				logger.Error(err, "failed to cook delicious Steak")
				return ctrl.Result{}, err
			}

			steakList.Items = append(steakList.Items, steak)
			actualSteakCount += 1
		}
	}

	// update the list of Steak references for the SteakfulSet's status
	steakfulSet.Status.SteaksServed = nil

	for _, steak := range steakList.Items {
		steakRef, err := ref.GetReference(r.Scheme, &steak)
		if err != nil {
			logger.Error(err, "unable to make reference to Steak", "steak", steak)
			continue
		}
		steakfulSet.Status.SteaksServed = append(steakfulSet.Status.SteaksServed, *steakRef)
	}

	if err := r.Status().Update(ctx, &steakfulSet); err != nil {
		logger.Error(err, "unable to update SteakfulSet status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *SteakfulSetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.Background(), &foodv1alpha1.Steak{}, steakOwnerKey, func(rawObj client.Object) []string {
		// grab the Steak object, extract the owner...
		steak := rawObj.(*foodv1alpha1.Steak)
		owner := metav1.GetControllerOf(steak)
		if owner == nil {
			return nil
		}
		// ...make sure it's a SteakfulSet...
		if owner.APIVersion != apiGVStr || owner.Kind != foodv1alpha1.SteakfulSetKind {
			return nil
		}

		// ...and if so, return it
		return []string{owner.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&foodv1alpha1.SteakfulSet{}).
		Owns(&foodv1alpha1.Steak{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: 1}).
		Complete(r)
}
