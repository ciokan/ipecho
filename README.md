This container prints the ip address of the 
user accessing it via GET on port 8080.

`docker run -p 8080:8080 ciokan/ipecho`

I use it on Kubernetes deployments to make 
sure the ip address of the visitor is properly 
forwarded and available to the services 
running on our clusters.

The code is written in Go and uses very little
resources. Example deployment on kubernetes:

```
kind: Deployment
spec:
    template:
        spec:
            containers:
                -
                    image: 'ciokan/ipecho'
                    imagePullPolicy: IfNotPresent
                    ports:
                        - containerPort: 8080
                    name: ipecho
                    resources:
                        requests:
                            cpu: 20m
                            memory: 64Mi
        metadata:
            labels:
                component: ipecho
    replicas: 1
apiVersion: extensions/v1beta1
metadata:
    labels:
        component: ipecho
    name: ipecho

---
kind: Service
spec:
    type: NodePort
    ports:
        -
            name: http
            protocol: TCP
            port: 8080
            targetPort: 8080
            nodePort: 30341
    selector:
        component: ipecho
apiVersion: v1
metadata:
    name: ipecho
```