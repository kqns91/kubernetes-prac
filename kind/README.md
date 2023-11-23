# kubernetes ローカル構築

## kubernetes in docker

ローカルマシンに kubernetes 環境を構築する。
Docker コンテナを複数個起動し、そのコンテナを kubernetes node として利用することで、複数台構成の kubernetes クラスタを構成する。

[kindファイル](./kind.yaml)を作成し、以下のコマンドで起動する。

```cmd
$kind create cluster --config ./kind.yaml
Creating cluster "kind" ..00.
 ✓ Ensuring node image (kindest/node:v1.28.0) 🖼
 ✓ Preparing nodes 📦 📦 📦 📦 📦 📦  
 ✓ Configuring the external load balancer ⚖️ 
 ✓ Writing configuration 📜 
 ✓ Starting control-plane 🕹️ 
 ✓ Installing CNI 🔌 
 ✓ Installing StorageClass 💾 
 ✓ Joining more control-plane nodes 🎮 
 ✓ Joining worker nodes 🚜 
Set kubectl context to "kind-kind"
You can now use your cluster with:

kubectl cluster-info --context kind-kind

Have a nice day! 👋
```

kubectl 使える。node 確認。

```cmd
$kubectl get nodes                       
NAME                  STATUS   ROLES           AGE     VERSION
kind-control-plane    Ready    control-plane   2m46s   v1.28.0
kind-control-plane2   Ready    control-plane   2m30s   v1.28.0
kind-control-plane3   Ready    control-plane   101s    v1.28.0
kind-worker           Ready    <none>          96s     v1.28.0
kind-worker2          Ready    <none>          96s     v1.28.0
kind-worker3          Ready    <none>          96s     v1.28.0
```

docker コンテナが起動している。

```cmd
docker ps
CONTAINER ID   IMAGE                                COMMAND                  CREATED         STATUS         PORTS                       NAMES
2494de7576be   kindest/haproxy:v20230606-42a2262b   "haproxy -W -db -f /…"   3 minutes ago   Up 3 minutes   127.0.0.1:51625->6443/tcp   kind-external-load-balancer
1caef72990cc   kindest/node:v1.28.0                 "/usr/local/bin/entr…"   4 minutes ago   Up 4 minutes                               kind-worker2
159d36d6805f   kindest/node:v1.28.0                 "/usr/local/bin/entr…"   4 minutes ago   Up 4 minutes   127.0.0.1:51627->6443/tcp   kind-control-plane
f07d8ebf8b88   kindest/node:v1.28.0                 "/usr/local/bin/entr…"   4 minutes ago   Up 4 minutes   127.0.0.1:51626->6443/tcp   kind-control-plane3
8aa8faa13205   kindest/node:v1.28.0                 "/usr/local/bin/entr…"   4 minutes ago   Up 4 minutes   127.0.0.1:51624->6443/tcp   kind-control-plane2
83cacb59bdad   kindest/node:v1.28.0                 "/usr/local/bin/entr…"   4 minutes ago   Up 4 minutes                               kind-worker
95882cd5392c   kindest/node:v1.28.0                 "/usr/local/bin/entr…"   4 minutes ago   Up 4 minutes                               kind-worker3
```

## kind に関する記事を試す

https://faun.pub/local-kubernetes-with-kind-helm-and-a-sample-service-4755e3e6eff4 に沿ってやってみる


mac ユーザー用の config を作成してクラスター作成。

```
kind create cluster --name local-dev --config k8s-cluster-config.yaml
```

```
$ kind get clusters
local-dev
```

```
helm create sample-service-helm

helm install sample-service --dry-run --debug ./sample-service-helm
```

```
$ helm install sample-service ./sample-service-helm --set service.type=NodePort --set service.nodePort=31234
NAME: sample-service
LAST DEPLOYED: Thu Nov 23 11:55:18 2023
NAMESPACE: default
STATUS: deployed
REVISION: 1
NOTES:
1. Get the application URL by running these commands:
  export NODE_PORT=$(kubectl get --namespace default -o jsonpath="{.spec.ports[0].nodePort}" services sample-service-sample-service-helm)
  export NODE_IP=$(kubectl get nodes --namespace default -o jsonpath="{.items[0].status.addresses[0].address}")
  echo http://$NODE_IP:$NODE_PORT
```

```
$ kubectl get po
NAME                                                READY   STATUS             RESTARTS   AGE
sample-service-sample-service-helm-c9ff9466-t4xrx   0/1     InvalidImageName   0          6s
```

ローカルの docker image を load する

```
kind load --name local-dev docker-image sample-service:latest
```

deployment.yaml を以下になるようにする

```yaml
containers:
     - name: sample-service
       image: sample-service:latest
       imagePullPolicy: IfNotPresent
```

```
helm delete sample-service 
helm install sample-service ./sample-service-helm --set service.type=NodePort --set service.nodePort=31234
```

```
kubectl get po                                                                                            
NAME                                                  READY   STATUS    RESTARTS   AGE
sample-service-sample-service-helm-6854d8b88b-vjlrt   1/1     Running   0          3s
```

できた。追加でちょっと遊んでみる。

go サーバーの deployment を追加した。

```sh
$ kubectl get po
NAME                                                  READY   STATUS    RESTARTS   AGE
sample-service-go-5599df85bc-ctfw4                    1/1     Running   0          21s
sample-service-sample-service-helm-6854d8b88b-cdgwj   1/1     Running   0          21s
```

nginx を追加したけど、CrashLoopBackOff になった。

```sh
$ kubectl get po
NAME                                                  READY   STATUS             RESTARTS        AGE
sample-service-go-5599df85bc-wrfv5                    1/1     Running            0               24m
sample-service-nginx-96f57b765-89pc5                  0/1     CrashLoopBackOff   9 (3m20s ago)   24m
sample-service-sample-service-helm-6854d8b88b-lkkzf   1/1     Running            0               24m
```

```sh
kubectl describe po sample-service-nginx
Name:             sample-service-nginx-96f57b765-89pc5
Namespace:        default
Priority:         0
Service Account:  sample-service-sample-service-helm
Node:             local-dev-control-plane/172.22.0.2
Start Time:       Thu, 23 Nov 2023 17:06:43 +0900
Labels:           app.kubernetes.io/instance=sample-service
                  app.kubernetes.io/name=sample-service-helm
                  pod-template-hash=96f57b765
Annotations:      <none>
Status:           Running
IP:               10.244.0.32
IPs:
  IP:           10.244.0.32
Controlled By:  ReplicaSet/sample-service-nginx-96f57b765
Containers:
  sample-service-helm:
    Container ID:   containerd://b987de03f1ef8e154ca5c09070dc5c2096e92c24d124bfbb474ec04386771ecc
    Image:          sample-service-nginx:latest
    Image ID:       docker.io/library/import-2023-11-23@sha256:d1edc26cdb9e57eddc30cd7a9213a9898d5f2197061b722d4af76727e4de6cad
    Port:           8080/TCP
    Host Port:      0/TCP
    State:          Waiting
      Reason:       CrashLoopBackOff
    Last State:     Terminated
      Reason:       Error
      Exit Code:    1
      Started:      Thu, 23 Nov 2023 17:37:58 +0900
      Finished:     Thu, 23 Nov 2023 17:37:58 +0900
    Ready:          False
    Restart Count:  11
    Liveness:       http-get http://:http/ delay=0s timeout=1s period=10s #success=1 #failure=3
    Readiness:      http-get http://:http/ delay=0s timeout=1s period=10s #success=1 #failure=3
    Environment:    <none>
    Mounts:
      /var/run/secrets/kubernetes.io/serviceaccount from kube-api-access-tk7j4 (ro)
Conditions:
  Type              Status
  Initialized       True 
  Ready             False 
  ContainersReady   False 
  PodScheduled      True 
Volumes:
  kube-api-access-tk7j4:
    Type:                    Projected (a volume that contains injected data from multiple sources)
    TokenExpirationSeconds:  3607
    ConfigMapName:           kube-root-ca.crt
    ConfigMapOptional:       <nil>
    DownwardAPI:             true
QoS Class:                   BestEffort
Node-Selectors:              <none>
Tolerations:                 node.kubernetes.io/not-ready:NoExecute op=Exists for 300s
                             node.kubernetes.io/unreachable:NoExecute op=Exists for 300s
Events:
  Type     Reason     Age                    From               Message
  ----     ------     ----                   ----               -------
  Normal   Scheduled  32m                    default-scheduler  Successfully assigned default/sample-service-nginx-96f57b765-89pc5 to local-dev-control-plane
  Normal   Pulled     30m (x5 over 32m)      kubelet            Container image "sample-service-nginx:latest" already present on machine
  Normal   Created    30m (x5 over 32m)      kubelet            Created container sample-service-helm
  Normal   Started    30m (x5 over 32m)      kubelet            Started container sample-service-helm
  Warning  BackOff    2m16s (x150 over 32m)  kubelet            Back-off restarting failed container sample-service-helm in pod sample-service-nginx-96f57b765-89pc5_default(f1cefa51-cd68-4044-af3a-b72d8d25f675)
```

go サーバーと js サーバーの Service 設定ができてなかった。
deployment と service は、selector によって紐づけられるらしい。
deployment の spec.selector.matchLabels と service の selector の設定を一致させる。

```yaml
$ kubectl get po
NAME                                    READY   STATUS    RESTARTS   AGE
sample-service-go-5889644db8-8km84      1/1     Running   0          6m47s
sample-service-js-589c4b65f4-vcv8m      1/1     Running   0          6m47s
sample-service-nginx-6d989844fc-rpccg   1/1     Running   0          6m47s
```

/go は go サーバー、/js は js サーバー、/ は nginx が返すようにできた。


メトリクスサーバーを入れてみる。

参考：https://qiita.com/dingtianhongjie/items/a8ddc2d7f7b57291a13e#%E3%82%B3%E3%83%BC%E3%83%89%E3%81%AE%E3%83%80%E3%82%A6%E3%83%B3%E3%83%AD%E3%83%BC%E3%83%89

```sh
helm repo add metrics-server https://kubernetes-sigs.github.io/metrics-server
```

```sh
$ helm upgrade -n kube-system --install -f values.yaml metrics-server metrics-server/metrics-server
Release "metrics-server" does not exist. Installing it now.
NAME: metrics-server
LAST DEPLOYED: Thu Nov 23 17:54:04 2023
NAMESPACE: kube-system
STATUS: deployed
REVISION: 1
TEST SUITE: None
NOTES:
***********************************************************************
* Metrics Server                                                      *
***********************************************************************
  Chart version: 3.11.0
  App version:   0.6.4
  Image tag:     registry.k8s.io/metrics-server/metrics-server:v0.6.4
***********************************************************************
```

CPU, Memory の使用量を確認できるようになった。

```sh
kubectl top node
NAME                      CPU(cores)   CPU%   MEMORY(bytes)   MEMORY%   
local-dev-control-plane   204m         2%     1014Mi          12%       
```

```sh
$ kubectl top pod   
NAME                                    CPU(cores)   MEMORY(bytes)   
sample-service-go-5889644db8-8km84      1m           5Mi             
sample-service-js-589c4b65f4-vcv8m      1m           8Mi             
sample-service-nginx-6d989844fc-rpccg   1m           1Mi             
```