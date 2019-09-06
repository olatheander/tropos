package kubernetes

import (
	"bytes"
	"fmt"
	"io"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/portforward"
	"k8s.io/client-go/tools/remotecommand"
	"k8s.io/client-go/transport/spdy"
	"net/http"
	"net/url"
	"strings"
	args "tropos/pkg/args"
)

func getDeployment(kube *args.Kubernetes) (*appsv1.Deployment, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kube.Config)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	namespace := getNamespace(kube)
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)
	return deploymentsClient.Get(kube.DeploymentName, metav1.GetOptions{})
}

//NewDeployment create a new deployment
func NewDeployment(kube *args.Kubernetes) (*appsv1.Deployment, error) {
	result, err := getDeployment(kube)
	if err == nil {
		fmt.Printf("Found existing deployment %s, reusing it.\n", kube.DeploymentName)
		return result, err
	}

	config, err := clientcmd.BuildConfigFromFlags("", kube.Config)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	namespace := getNamespace(kube)
	deploymentsClient := clientSet.AppsV1().Deployments(namespace)

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: kube.DeploymentName,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": kube.DeploymentName,
				},
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": kube.DeploymentName,
					},
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  kube.DeploymentName,
							Image: kube.Image,
							//ImagePullPolicy: apiv1.PullNever,	//TODO: remove, only added for Minikube development.
							Ports: []apiv1.ContainerPort{
								{
									Name:          "ssh",
									Protocol:      apiv1.ProtocolTCP,
									ContainerPort: kube.ContainerPort,
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
	result, err = deploymentsClient.Create(deployment)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return result, err
}

func getNamespace(kube *args.Kubernetes) string {
	var namespace string
	if kube.Namespace != "" {
		namespace = kube.Namespace
	} else {
		namespace = apiv1.NamespaceDefault
	}
	return namespace
}

//SwapDeployment swap out an existing deployment.
func SwapDeployment() (*appsv1.Deployment, error) {
	return nil, nil
}

// Get just a single pod of the deployment or fail if not a single pod deployment.
func getDeploymentPod(kube *args.Kubernetes, deployment *appsv1.Deployment) (*apiv1.Pod, error) {
	pods, err := getDeploymentPods(kube, deployment)
	if err != nil {
		panic(err)
	}
	if len(pods.Items) != 1 {
		return nil, fmt.Errorf("expected only a single pod in the deployment, found %d", len(pods.Items))
	}
	pod := pods.Items[0]
	return &pod, nil
}

//getDeploymentPods get all pods belonging to the deployment.
func getDeploymentPods(kube *args.Kubernetes, deployment *appsv1.Deployment) (*apiv1.PodList, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kube.Config)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	labelMap, _ := metav1.LabelSelectorAsMap(deployment.Spec.Selector)
	listOptions := metav1.ListOptions{LabelSelector: labels.SelectorFromSet(labelMap).String()}
	pods, err := clientSet.CoreV1().Pods(getNamespace(kube)).List(listOptions)
	return pods, err
}

//PortForward set up port-forward to the deployment.
func PortForward(kube *args.Kubernetes, deployment *appsv1.Deployment) error {
	config, err := clientcmd.BuildConfigFromFlags("", kube.Config)
	if err != nil {
		panic(err)
	}

	pod, err := getDeploymentPod(kube, deployment)
	if err != nil {
		panic(err)
	}

	roundTripper, upgrader, err := spdy.RoundTripperFor(config)
	if err != nil {
		panic(err)
	}

	path := fmt.Sprintf("/api/v1/namespaces/%s/pods/%s/portforward",
		deployment.GetNamespace(),
		pod.GetName())
	hostIP := strings.TrimLeft(config.Host, "https://")
	serverURL := url.URL{Scheme: "https", Path: path, Host: hostIP}

	dialer := spdy.NewDialer(upgrader, &http.Client{Transport: roundTripper}, http.MethodPost, &serverURL)

	stopChan, readyChan := make(chan struct{}, 1), make(chan struct{}, 1)
	out, errOut := new(bytes.Buffer), new(bytes.Buffer)
	ports := []string{fmt.Sprintf("%d:%d", kube.HostPort, kube.ContainerPort)}

	forwarder, err := portforward.New(dialer, ports, stopChan, readyChan, out, errOut)
	if err != nil {
		panic(err)
	}

	go func() {
		for range readyChan { // Kubernetes will close this channel when it has something to tell us.
		}
		if len(errOut.String()) != 0 {
			panic(errOut.String())
		} else if len(out.String()) != 0 {
			fmt.Println(out.String())
			go func() {
				fmt.Println("Mounting working directory in pod.")
				//TODO: this is just a dummy test. Should set up sshfs like in https://superuser.com/questions/616182/how-to-mount-local-directory-to-remote-like-sshfs
				output, stderr, err := Exec("ls -l /", kube, deployment, nil)
				fmt.Println(output)
				fmt.Println(stderr)
				if err != nil {
					panic(err)
				}
			}()
		}
	}()

	if err = forwarder.ForwardPorts(); err != nil { // Locks until stopChan is closed.
		panic(err)
	}

	return nil
}

// Exec execute the specified command in the Pod
func Exec(command string, kube *args.Kubernetes, deployment *appsv1.Deployment, stdin io.Reader) (string, string, error) {
	config, err := clientcmd.BuildConfigFromFlags("", kube.Config)
	if err != nil {
		panic(err)
	}

	pod, err := getDeploymentPod(kube, deployment)
	if err != nil {
		panic(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	req := clientSet.CoreV1().RESTClient().Post().
		Resource("pods").
		Name(pod.Name).
		Namespace(deployment.Namespace).
		SubResource("exec")
	scheme := runtime.NewScheme()
	if err := apiv1.AddToScheme(scheme); err != nil {
		return "", "", fmt.Errorf("error adding to scheme: %v", err)
	}

	//TODO: if multiple containers in Pod and kube.containerName is not set or not matching existing container => fail.
	parameterCodec := runtime.NewParameterCodec(scheme)
	req.VersionedParams(&apiv1.PodExecOptions{
		Command:   strings.Fields(command),
		Container: kube.ContainerName,
		Stdin:     stdin != nil,
		Stdout:    true,
		Stderr:    true,
		TTY:       false,
	}, parameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", "", fmt.Errorf("error while creating Executor: %v", err)
	}

	var stdout, stderr bytes.Buffer
	err = exec.Stream(remotecommand.StreamOptions{
		Stdin:  stdin,
		Stdout: &stdout,
		Stderr: &stderr,
		Tty:    false,
	})
	if err != nil {
		return "", "", fmt.Errorf("error in Stream: %v", err)
	}

	return stdout.String(), stderr.String(), nil
}

//DeleteDeployment Delete the newly (non-swapped) deployment.
func DeleteDeployment(*appsv1.Deployment) error {
	return nil
}

//RestoreDeployment restore a swapped out deployment.
func RestoreDeployment(*appsv1.Deployment) error {
	return nil
}

func int32Ptr(i int32) *int32 { return &i }

func newTrue() *bool {
	b := true
	return &b
}
