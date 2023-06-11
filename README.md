# microfrontends-controller
This repository is a (GO) re-implementation of the [Kubernetes Controller] pattern over custom resources specifying front-end web components to be dynamically integrated into a user interface application shell, using operator-sdk. It is heavily inspired by https://github.com/milung/ufe-controller. The project is of educational nature.

This is an experimental concept design of micro-frontends architecture, considering declarative definition of micro-frontends as part of the Kubernetes API custom resource definitions, and leveraging the web components technology. This enables us to approach the development of particular micro-frontends in a similar way as is done with the development of cloud-native microservices.

## Description
The original ufe-controller was implemented in Prolog using the Kubernetes (k8s) client to read the Custom Resources (CRs) called web components. It then aggregates the CRs to a frontend config structure which is then served via a REST endpoint. This config is then used to inject the web components to a frontend application container.

The aim of this project is to use operator-sdk and Golang to recreate what was done in the ufe-controller. It was created as a learning experience to understand how the k8s operators work and how to implement one.

There are two main components in the application: WebComponentController and an API server serving the fe-config.

WebComponentController is reconciling the WebComponent CRs. Currently, it is applying the following logic on a WebComponent:

- Check if the WebComponent still exists.
- If WebComponent does not contain finalizers it adds one so that the controller does its cleanup before the k8s removes the WebComponent.
- It converts the WebComponent to internal MicroFrontendConfig model which in turn is converted to TransferObject and stored in memory to be served through the fe-config endpoint.
- If WebComponent is set to be deleted it removes it also from the in-memory storage so it is no longer a part of fe-config response.

The API server simply takes the last snapshot of the in-memory storage and returns it to the frontend.

### Migration TODO list:
Here is a list of functionality that is provided by the original ufe-controller and whit information of what is implemented in the go alternative:
- [x] An operator observing specific CRs either on specific namespaces or in all
- [x] REST endpoint serving the MicroFrontendConfiguration as JSON (`/fe-config`)
- [x] REST endpoint serving the MicroFrontendConfiguration as JavaScript
- TBD

### Configuration
You can use environment variables to configure the following parameters:
| Env. Variable | Default Value | Description |
|- |- |- |
|BASE_URL| / |Base URL of the server, all absolute links are prefixed with this address|
|OBSERVE_NAMESPACES||Comma separated list of namespaces in which to look for webcomponents to be served by this instance|
|USER_ID_HEADER|x-forwarded-email|incomming request`s header name (lowercase) specifying the user identifier, typically email|
|USER_NAME_HEADER|x-forwarded-user|incomming request`s header name (lowercase) specifying the user name|
|USER_ROLES_HEADER|x-forwarded-groups|incomming request`s header name (lowercase) specifying the list of user roles (or groups)|
|WEBCOMPONENTS_SELECTOR||comma separate list of key-value pairs, used to filter WebComponent resources handled by this controller|

## Getting Started
You’ll need a Kubernetes cluster to run against. You can use [KIND](https://sigs.k8s.io/kind) to get a local cluster for testing, or run against a remote cluster.

> Your controller will automatically use the current context in your kubeconfig file (i.e. whatever cluster `kubectl cluster-info` shows).

### Running locally
> This will still require a cluster where the CRs will be stored only the processing part will run from your code. Allowing debuging etc.
1. Install CRD into your local cluster
```sh
kubectl apply -k config/crd
```

2. Install Instances of Custom Resources:

```sh
kubectl apply -k config/samples/
```

3. Or if you want to test the namespace observation you can apply namespaced version of the samples:
```sh
kubect apply -k config/samples/namespaced
```

3. Run a debug session using included debug configurations in `.vscode/launch.json` there is one for observing all namespaces and one where you can limit the observed namespaces.

4. Call `localhost:10000\fe-config` to get the MicroFrontendConfigurations created from the CRs.
### Delete CRs
To delete the CRs from the cluster and see the delete logic in the controller:
```sh
kubectl delete -k config/samples
```
or if you used the namespaced sample:
```sh
kubectl delete -k config/samples/namespaced
```
### Uninstall CRDs
To delete the CRDs from the cluster:
```sh
kubectl delete -k config/crd
```

### Running on the cluster
. Install CRD into your local cluster
```sh
kubectl apply -k config/crd
```

2. Install Instances of Custom Resources:

```sh
kubectl apply -k config/samples/
```

3. Or if you want to test the namespace observation you can apply namespaced version of the samples:
```sh
kubect apply -k config/samples/namespaced
```

4. Build and push your image to the location specified by `IMG`:

```sh
make docker-build docker-push IMG=<some-registry>/microfrontends-controller:tag
```

3. Deploy the controller to the cluster with the image specified by `IMG`:

```sh
make deploy IMG=<some-registry>/microfrontends-controller:tag
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

### How it works
This project aims to follow the Kubernetes [Operator pattern](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/).

It uses [Controllers](https://kubernetes.io/docs/concepts/architecture/controller/),
which provide a reconcile function responsible for synchronizing resources until the desired state is reached on the cluster.

### Modifying the API definitions
If you are editing the API definitions, generate the manifests such as CRs or CRDs using:

```sh
make manifests
```

**NOTE:** Run `make --help` for more information on all potential `make` targets

More information can be found via the [Kubebuilder Documentation](https://book.kubebuilder.io/introduction.html)

## License

Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.

