/*
Copyright 2022.

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
	"github.com/hipster-labs/jhipster-operator/pkg"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	k8sv1alpha1 "github.com/hipster-labs/jhipster-operator/api/v1alpha1"
)

// JHipsterSetupReconciler reconciles a JHipsterSetup object
type JHipsterSetupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=k8s.jhipster.tech,resources=jhipstersetups,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=k8s.jhipster.tech,resources=jhipstersetups/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=k8s.jhipster.tech,resources=jhipstersetups/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// the JHipsterSetup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.1/pkg/reconcile
func (r *JHipsterSetupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	setup := &k8sv1alpha1.JHipsterSetup{}
	err := r.Get(ctx, req.NamespacedName, setup)

	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}

		return ctrl.Result{}, err
	}

	logger.Info("reconciling", "setup", setup)
	if setup.Spec.ServiceDiscoveryType == k8sv1alpha1.Consul {
		return r.ensureConsulResources(ctx, setup)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *JHipsterSetupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&k8sv1alpha1.JHipsterSetup{}).
		Owns(&k8sv1alpha1.JHipsterApplication{}).
		Owns(&appsv1.StatefulSet{}).
		Owns(&v1.Secret{}).
		Owns(&v1.Service{}).
		Owns(&v1.ConfigMap{}).
		Complete(r)
}

func (r *JHipsterSetupReconciler) ensureConsulResources(ctx context.Context, setup *k8sv1alpha1.JHipsterSetup) (ctrl.Result, error) {
	// ensure the gossipKey
	gossipKey := &v1.Secret{}
	gossipKeyTpl, err := pkg.ConsulSecret(setup.Name, setup.Namespace, pkg.RandSeq(20))
	if err != nil {
		return ctrl.Result{}, err
	}
	err = r.Get(ctx, types.NamespacedName{
		Namespace: gossipKeyTpl.Namespace,
		Name:      gossipKeyTpl.Name,
	}, gossipKey)

	if err != nil {
		if errors.IsNotFound(err) {
			ctrl.SetControllerReference(setup, gossipKeyTpl, r.Scheme)
			err = r.Create(ctx, gossipKeyTpl)

			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}
	}

	// todo, reconcile if user changed

	consulService := &v1.Service{}
	consulServiceTpl, err := pkg.ConsulService(setup.Name, setup.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Get(ctx, types.NamespacedName{
		Namespace: consulServiceTpl.Namespace,
		Name:      consulServiceTpl.Name,
	}, consulService)

	if err != nil {
		if errors.IsNotFound(err) {
			ctrl.SetControllerReference(setup, consulServiceTpl, r.Scheme)
			err = r.Create(ctx, consulServiceTpl)

			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}

		return ctrl.Result{}, err
	}

	// todo, reconcile if user changed
	consulSts := &appsv1.StatefulSet{}
	consulStsTpl, err := pkg.ConsulSts(setup.Name, setup.Namespace, setup.Spec.StorageClassName)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Get(ctx, types.NamespacedName{
		Namespace: consulStsTpl.Namespace,
		Name:      consulStsTpl.Name,
	}, consulSts)

	if err != nil {
		if errors.IsNotFound(err) {
			ctrl.SetControllerReference(setup, consulStsTpl, r.Scheme)
			err = r.Create(ctx, consulStsTpl)

			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}

		return ctrl.Result{}, err
	}

	// todo, reconcile if user changed
	consulConfigLoaderDeployment := &appsv1.Deployment{}
	consulConfigLoaderDeploymentTpl, err := pkg.ConsulConfigLoader(setup.Name, setup.Namespace)
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Get(ctx, types.NamespacedName{
		Namespace: consulConfigLoaderDeploymentTpl.Namespace,
		Name:      consulConfigLoaderDeploymentTpl.Name,
	}, consulConfigLoaderDeployment)

	if err != nil {

		if errors.IsNotFound(err) {
			ctrl.SetControllerReference(setup, consulConfigLoaderDeploymentTpl, r.Scheme)
			err = r.Create(ctx, consulConfigLoaderDeploymentTpl)

			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}

		return ctrl.Result{}, err
	}
	applicationConfigMap := &v1.ConfigMap{}
	applicationConfigMapTpl, err := pkg.ConsulApplicationConfig(setup.Name, setup.Namespace, pkg.RandSeq(64))
	if err != nil {
		return ctrl.Result{}, err
	}

	err = r.Get(ctx, types.NamespacedName{
		Namespace: applicationConfigMapTpl.Namespace,
		Name:      applicationConfigMapTpl.Name,
	}, applicationConfigMap)

	if err != nil {

		if errors.IsNotFound(err) {
			ctrl.SetControllerReference(setup, applicationConfigMapTpl, r.Scheme)
			err = r.Create(ctx, applicationConfigMapTpl)

			if err != nil {
				return ctrl.Result{}, err
			}

			return ctrl.Result{Requeue: true}, nil
		}

		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}
