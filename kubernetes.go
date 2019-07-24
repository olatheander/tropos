package main

import (
	"fmt"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func NewDeployment(kube *Kubernetes) (*appsv1.Deployment, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kube.config)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	var namespace string
	if kube.namespace != "" {
		namespace = kube.namespace
	} else {
		namespace = apiv1.NamespaceDefault
	}
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: kube.deploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": kube.deploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": kube.deploymentName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  kube.deploymentName,
							Image: kube.image,
							Ports: []apiv1.ContainerPort{
								{
									Name:          "ssh",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: kube.containerPort,
								},
							},
							SecurityContext: &apiv1.SecurityContext{
								Privileged: newTrue(),
								Capabilities: &apiv1.Capabilities{
									Add: []apiv1.Capability{
										"SYS_ADMIN",
									},
								},
							},
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return result, err
}

func SwapDeployment() (*appsv1.Deployment, error) {
	return nil, nil
}

func int32Ptr(i int32) *int32 { return &i }

func newTrue() *bool {
	b := true
	return &b
}
