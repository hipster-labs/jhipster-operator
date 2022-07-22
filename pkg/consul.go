package pkg

import (
	"encoding/base64"
	"errors"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
)

const ConsulTplPath = "config/workload-templates/consul/"

func ConsulSecret(name, namespace, gossipKey string) (*v1.Secret, error) {
	gossipKeySecret := &v1.Secret{}
	vars := map[string]string{
		"Name":            name,
		"Namespace":       namespace,
		"GossipKeyBase64": base64.StdEncoding.EncodeToString([]byte(gossipKey)),
	}
	obj, groupVersionKind, err := GetK8sObjectFromTemplate(ConsulTplPath+"secret.yaml", vars)
	if err != nil {
		return nil, err
	}

	if groupVersionKind.Kind == "Secret" {
		gossipKeySecret = obj.(*v1.Secret)

		return gossipKeySecret, nil
	}

	return nil, errors.New("error while parsing secret file")

}

func ConsulService(name, namespace string) (*v1.Service, error) {
	consulService := &v1.Service{}
	vars := map[string]string{
		"Name":      name,
		"Namespace": namespace,
	}

	obj, groupVersionKind, err := GetK8sObjectFromTemplate(ConsulTplPath+"service.yaml", vars)
	if err != nil {
		return nil, err
	}

	if groupVersionKind.Kind == "Service" {
		consulService = obj.(*v1.Service)

		return consulService, nil
	}

	return nil, errors.New("error during parsing consul service file")
}

func ConsulSts(name string, namespace string, storageClassName *string) (*appsv1.StatefulSet, error) {
	consulStatefulSet := &appsv1.StatefulSet{}
	vars := map[string]string{
		"Name":      name,
		"Namespace": namespace,
	}
	obj, groupVersionKind, err := GetK8sObjectFromTemplate(ConsulTplPath+"sts.yaml", vars)
	if err != nil {
		return nil, err
	}

	if groupVersionKind.Kind == "StatefulSet" {
		consulStatefulSet = obj.(*appsv1.StatefulSet)
		if storageClassName != nil {
			consulStatefulSet.Spec.VolumeClaimTemplates[0].Spec.StorageClassName = storageClassName
		}

		return consulStatefulSet, nil
	}

	return nil, errors.New("error during processing STS file")
}

func ConsulConfigLoader(name, namespace string) (*appsv1.Deployment, error) {
	consulConfigLoader := &appsv1.Deployment{}
	vars := map[string]string{
		"Name":      name,
		"Namespace": namespace,
	}

	obj, groupVersionKind, err := GetK8sObjectFromTemplate(ConsulTplPath+"config-loader-deployment.yaml", vars)

	if err != nil {
		return nil, err
	}

	if groupVersionKind.Kind == "Deployment" {
		consulConfigLoader = obj.(*appsv1.Deployment)

		return consulConfigLoader, nil
	}

	return nil, errors.New("error parsing consul config loader deployment")
}

func ConsulApplicationConfig(name, namespace, jwtSecret string) (*v1.ConfigMap, error) {
	applicationConfig := &v1.ConfigMap{}

	vars := map[string]string{
		"Name":            name,
		"Namespace":       namespace,
		"JwtSecretBase64": base64.StdEncoding.EncodeToString([]byte(jwtSecret)),
	}

	obj, groupVersionKind, err := GetK8sObjectFromTemplate(ConsulTplPath+"application-config.yaml", vars)

	if err != nil {
		return nil, err
	}

	if groupVersionKind.Kind == "ConfigMap" {
		applicationConfig = obj.(*v1.ConfigMap)

		return applicationConfig, err
	}

	return nil, errors.New("error with parsing application config")
}
