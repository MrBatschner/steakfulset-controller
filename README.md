# SteakfulSet controller :cut_of_meat:

In Kubernetes, we are talking a lot about StatefulSets. But what about SteakfulSets? Sounds delicious to you? Read on...

## Description

This repository implements a small and simple Kubernetes controller that reconciles on the custom resource `SteakfulSet` and creates `Steaks` - another custom resource - out of it. It is intended to give a you an idea about what it takes to implement your own Kubernetes operator - something you might come across doing at some point in time.

This repository is heavily based on [Kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) and the [Kubebuilder book](https://book.kubebuilder.io/). In fact, Kubebuilder has been used to create the skeleton of this repository, to provide all the generators and to design the CRD layout.

## File layout

### The custom resources

The custom resources `SteakfulSet` and `Steak` are defined in Go source code within the files inside the [`api/v1alpha`](api/v1alpha1/) folder. Inside the Go source files, we use some Kubebuilder decorators to describe some spefic of the custom resources. Then, the following command will invoke Kubebuilder's generator to create the [custom resource definition](https://kubernetes.io/docs/tasks/extend-kubernetes/custom-resources/custom-resource-definitions/) files:

    make manifests
    
The resulting CRDs can be found in the [`config/crd`](config/crd/) directory.

### The controller

The controller consists of only two files: the controller is set up and initialized in [`main.go`](main.go). A quick overview about what happens in here is given in [chapter 1.7.1](https://book.kubebuilder.io/cronjob-tutorial/main-revisited.html) of the Kubebuilder book - in fact, `main.go` only contains generated code.

The real fun happens in [`controllers/steakfulset_controller.go`](controllers/steakfulset_controller.go). In here, the reconcilation logic for `SteakfulSets` is implemented. Have a look at the comments of the `Reconcile()` method to get an idea what happens.


## Getting Started

You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.
**Note:** Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Helm Chart

To install the SteakfulSet controller into your cluster with the help of Helm, simply look into the [`charts`](charts) directory. Install it with:

```
kubectl create ns steakfulset-controller
helm install -n steakfulset-controller bbq charts/steakfulset-controller
```

To fry your first Steaks, just `kubectl apply -f config/samples/food_v1alpha1_steakfulset.yaml`.

### Running manually on the cluster

1. Install the CRDs into the cluster:

    ```sh
    make install
    ```

1. Build and push your image to the location specified by `IMG`:

    ```sh
    make docker-build docker-push IMG=<some-registry>/steakfulset-controller:tag
    ```

1. Deploy the controller to the cluster with the image specified by `IMG`:

    ```sh
    make deploy IMG=<some-registry>/steakfulset-controller:tag
    ```

### Uninstall CRDs
To delete the CRDs from the cluster:

```sh
make uninstall
```

### Undeploy controller
UnDeploy the controller from the cluster:

```sh
make undeploy
```

## More on Kubernetes controllers?

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)
