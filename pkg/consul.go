package pkg

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ConsulSecret(name string, namespace string) *v1.Secret {

	workloadName := name + "-gossip-key"
	gossip_keySecret := &v1.Secret{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Secret",
			APIVersion: "corev1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName,
			Namespace: namespace,
		},
		Data: map[string][]uint8{"gossip-key": {
			73,
			71,
			48,
			71,
			49,
			119,
			55,
			103,
			56,
			65,
			110,
			88,
			48,
			48,
			55,
			113,
			65,
			48,
			88,
			73,
			106,
			49,
			50,
			70,
		},
		},
		Type: v1.SecretType("Opaque"),
	}

	return gossip_keySecret
}

func ConsulService(name, namespace string) *v1.Service {
	workloadName := name + "-consul"
	consulService := &v1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "corev1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:        workloadName,
			Namespace:   namespace,
			Labels:      map[string]string{"app": "consul"},
			Annotations: map[string]string{"service.alpha.kubernetes.io/tolerate-unready-endpoints": "true"},
		},
		Spec: v1.ServiceSpec{
			Ports: []v1.ServicePort{
				{
					Name: "http",
					Port: 8500,
				},
				{
					Name: "rpc",
					Port: 8400,
				},
				{
					Name:     "tcp-serflan",
					Protocol: v1.Protocol("TCP"),
					Port:     8301,
				},
				{
					Name:     "udp-serflan",
					Protocol: v1.Protocol("UDP"),
					Port:     8301,
				},
				{
					Name:     "tcp-serfwan",
					Protocol: v1.Protocol("TCP"),
					Port:     8302,
				},
				{
					Name:     "udp-serfwan",
					Protocol: v1.Protocol("UDP"),
					Port:     8302,
				},
				{
					Name: "server",
					Port: 8300,
				},
				{
					Name: "consuldns",
					Port: 8600,
				},
			},
			Selector:  map[string]string{"app": "consul"},
			ClusterIP: "None",
		},
	}

	return consulService
}

func ConsulSts(name string, namespace string, size int32, storageSize string, storageClassName *string) (*appsv1.StatefulSet, error) {

	var fsGroup int64 = 1000
	quantity, err := resource.ParseQuantity(storageSize)
	if err != nil {
		return nil, err
	}
	workloadName := name + "-consul"
	consulStatefulSet := &appsv1.StatefulSet{
		TypeMeta: metav1.TypeMeta{
			Kind:       "StatefulSet",
			APIVersion: "apps/appsv1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      workloadName,
			Namespace: namespace,
		},
		Spec: appsv1.StatefulSetSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{MatchLabels: map[string]string{
				"app": "consul",
			},
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   workloadName,
					Labels: map[string]string{"app": "consul"},
				},
				Spec: v1.PodSpec{
					Volumes: []v1.Volume{
						{
							Name: "gossip-key",
							VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{
								SecretName: name + "-gossip-key",
							},
							},
						},
					},
					Containers: []v1.Container{
						{
							Name:  workloadName,
							Image: "consul:1.11.3",
							Command: []string{
								"/bin/sh",
								"-ec",
								`IP=$(hostname -i)
if [ -e /etc/consul/secrets/gossip-key ]; then
  echo "{\"encrypt\": \"$(base64 /etc/consul/secrets/gossip-key)\"}" > /etc/consul/encrypt.json
  GOSSIP_KEY="-config-file /etc/consul/encrypt.json"
fi
for i in $(seq 0 $((${INITIAL_CLUSTER_SIZE} - 1))); do
    while true; do
        echo "Waiting for ${PETSET_NAME}-${i}.${PETSET_NAME} to come up"
        ping -W 1 -c 1 ${PETSET_NAME}-${i}.${PETSET_NAME}.${PETSET_NAMESPACE}.svc.cluster.local > /dev/null && break
        sleep 1s
    done
done
PEERS=""
for i in $(seq 0 $((${INITIAL_CLUSTER_SIZE} - 1))); do
    PEERS="${PEERS}${PEERS:+ } -retry-join $(ping -c 1 ${PETSET_NAME}-${i}.${PETSET_NAME}.${PETSET_NAMESPACE}.svc.cluster.local | awk -F'[()]' '/PING/{print $2}')"
done
exec /bin/consul agent \
  -data-dir=/var/lib/consul \
  -server \
  -ui \
  -bootstrap-expect=${INITIAL_CLUSTER_SIZE} \
  -bind=0.0.0.0 \
  -advertise=${IP} \
  ${PEERS} \
  ${GOSSIP_KEY} \
  -client=0.0.0.0
`,
							},
							Ports: []v1.ContainerPort{
								{
									Name:          "http",
									ContainerPort: 8500,
								},
								{
									Name:          "rpc",
									ContainerPort: 8400,
								},
								{
									Name:          "tcp-serflan",
									ContainerPort: 8301,
									Protocol:      v1.Protocol("TCP"),
								},
								{
									Name:          "udp-serflan",
									ContainerPort: 8301,
									Protocol:      v1.Protocol("UDP"),
								},
								{
									Name:          "tcp-serfwan",
									ContainerPort: 8302,
									Protocol:      v1.Protocol("TCP"),
								},
								{
									Name:          "udp-serfwan",
									ContainerPort: 8302,
									Protocol:      v1.Protocol("UDP"),
								},
								{
									Name:          "server",
									ContainerPort: 8300,
								},
								{
									Name:          "consuldns",
									ContainerPort: 8600,
								},
							},
							Env: []v1.EnvVar{
								{
									Name:  "INITIAL_CLUSTER_SIZE",
									Value: fmt.Sprintf("%d", size),
								},
								{
									Name:  "PETSET_NAME",
									Value: workloadName,
								},
								{
									Name: "PETSET_NAMESPACE",
									ValueFrom: &v1.EnvVarSource{
										FieldRef: &v1.ObjectFieldSelector{FieldPath: "metadata.namespace"},
									},
								},
							},
							VolumeMounts: []v1.VolumeMount{
								{
									Name:      "datadir",
									MountPath: "/var/lib/consul",
								},
								{
									Name:      "gossip-key",
									ReadOnly:  true,
									MountPath: "/etc/consul/secrets",
								},
							},
							ImagePullPolicy: v1.PullPolicy("Always"),
						},
					},
					SecurityContext: &v1.PodSecurityContext{FSGroup: &fsGroup},
				},
			},
			VolumeClaimTemplates: []v1.PersistentVolumeClaim{
				{
					ObjectMeta: metav1.ObjectMeta{
						Name: "datadir",
					},
					Spec: v1.PersistentVolumeClaimSpec{
						AccessModes: []v1.PersistentVolumeAccessMode{v1.PersistentVolumeAccessMode("ReadWriteOnce")},
						Resources:   v1.ResourceRequirements{Requests: v1.ResourceList{v1.ResourceName("storage"): quantity}},
					},
				},
			},
			ServiceName: workloadName,
		},
	}

	if storageClassName != nil {
		consulStatefulSet.Spec.VolumeClaimTemplates[0].Spec.StorageClassName = storageClassName
	}

	return consulStatefulSet, nil
}
