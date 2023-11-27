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

grpc サーバーを追加してみる。health エンドポイントを用意していないので、deployment の livenessProbe, readinessProbe は設定しない。

```sh
$ kubectl get po                                                                                            
NAME                                    READY   STATUS    RESTARTS   AGE
sample-service-go-5889644db8-8km84      1/1     Running   0          123m
sample-service-grpc-6db799f6cb-x5gjw    1/1     Running   0          2s
sample-service-js-589c4b65f4-vcv8m      1/1     Running   0          123m
sample-service-nginx-6d989844fc-rpccg   1/1     Running   0          123m
```

go サーバーに grpc にアクセスするエンドポイントを追加した。
/go を http://sample-service-go:8080/ にルーティングしていたが、クラスターの /go/sample にアクセスした場合に、go サーバーに /sample ではなく、/go/sample できてしまうため、http://sample-service-go:8080/go にルーティングするようにした。

nginx -> go -> grpc を実現した。

grpc の replicas を 5 にして、go サーバー経由でリクエストを送るが、リクエストが分散されていない。

```sh
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:09 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:19 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:19 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:29 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:29 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:39 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:39 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:49 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:49 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:59 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:48:59 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:49:09 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-krndv sample-service-helm 2023/11/23 13:49:09 Received request: /sample.service.HelloWorld/SayHello
```

リクエストを分散させてみる。

参考：https://christina04.hatenablog.com/entry/grpc-client-side-lb


grpc の client 生成時に、resolver を設定する。

```go
resolver.SetDefaultScheme("dns")
conn, err := grpc.Dial("sample-service-grpc:8080", grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingConfig": [{"round_robin":{}}]}`))
```

分散されなかった。。

ingress でできるかもしれない。

https://www.amazon.co.jp/Kubernetes%E3%81%AE%E7%9F%A5%E8%AD%98%E5%9C%B0%E5%9B%B3-%E2%80%94%E2%80%94-%E7%8F%BE%E5%A0%B4%E3%81%A7%E3%81%AE%E5%9F%BA%E7%A4%8E%E3%81%8B%E3%82%89%E6%9C%AC%E7%95%AA%E9%81%8B%E7%94%A8%E3%81%BE%E3%81%A7-%E9%9D%92%E5%B1%B1-%E7%9C%9F%E4%B9%9F/dp/4297135736/ref=sr_1_1?adgrpid=150201939336&hvadid=665617000696&hvdev=c&hvlocphy=1028853&hvnetw=g&hvqmt=e&hvrand=1998024573929357799&hvtargid=kwd-2118030279282&hydadcr=27490_14701104&jp-ad-ap=0&keywords=kubernetes+%E7%9F%A5%E8%AD%98%E5%9C%B0%E5%9B%B3&qid=1700757269&sr=8-1

のコラムより

> 9RPC通信の負荷分散においては、9RPCが使うHTTP2通信の特性から注 意しなければならない点があります。
従来のHTTP1.1では1つのサーバと並列して複数の通信をする際、複数の TCP接続を行っていました。TCP通信は接続の確立処理や誤り検出、再送機 能を備えており、接続本数が増えるにつれてサーバ側の消費する計算リソー スが大きくなります。
一方、HTTP2通信ではこのコストを削減するために、TCP接続を1つのサー バに対して1つにし、単一のTCP接続の中で並行して通信するようにしまし た。この手法はTCP通信の処理コストを抑える一方、ロードバランサが通信 を割り振る対象のサーバを増やしても既存のクライアントが同じサーバと通 信し続け、負荷が分散しないという問題を発生させます。
この問題に対してはいくつかの解決方法がありますが、HTTP2に対応した ロードバランサを使用するとクライアント側/サーバ側両方ともに変更を加
えることなく問題を解決できます。 HTTP2に対応したロードバランサを使用するとクライアントからロードバ
ランサにRPCを発行し、ロードバランサからバックエンドにそのRPCを伝 播させます。このときクライアントとバックエンドは直接TCP通信によって 接続されないため、先ほど挙げた負荷分散についての問題を克服できます。
Kubernetesの| n9ressにおいても使用しているロードバランサがHTTP2 に対応していれば、アノテーションやカスタムリソースを用いてHTTP2通信
のロードバランシングを設定できます。詳しい設定方法は各|n9ressコント ローラのドキュメントを参照してください。たとえばNGINX In9ress ControUerでは|n9ressリソースにnginx.ingress.kubernetes.io/backend-
protocol : "GRPC''アノテーションを付与してHTTP2機能を有効にします。

ingress が何者か分かってない。

Service は L4 ロードバランサーで、Ingress は L7 ロードバランサーなので、リソースとして分けられている。

ingress リソースと ingress controller がある。
ingress リソースが k8s に登録された際に、ingress controller がL7 ロードバランサーの設定や、Nginx の設定を変更してリロードを実施するなどの何らかの処理を行う。

ingress にも種類がある。以下はよく使われるもの。
- GKE Ingress
- Nginx Ingress

ちょっと Ingress の勉強途中だけど、別の方法で対応できそうぽい。

参考：https://techdozo.dev/grpc-load-balancing-on-kubernetes-using-headless-service/

> What is Headless Service ?
>> Luckily, Kubernetes allows clients to discover pod IPs through DNS lookups. Usually, when you perform a DNS lookup for a service, the DNS server returns a single IP — the service’s cluster IP. But if you tell Kubernetes you don’t need a cluster IP for your service (you do this by setting the clusterIP field to None in the service specification ), the DNS server will return the pod IPs instead of the single service IP. Instead of returning a single DNS A record, the DNS server will return multiple A records for the service, each pointing to the IP of an individual pod backing the service at that moment. Clients can therefore do a simple DNS A record lookup and get the IPs of all the pods that are part of the service. The client can then use that information to connect to one, many, or all of them.
>> Setting the clusterIP field in a service spec to None makes the service headless, as Kubernetes won’t assign it a cluster IP through which clients could connect to the pods backing it.

ClusterIP を使っているのが悪いかも。

```yaml
# 変更前
apiVersion: v1
kind: Service
metadata:
  name: "sample-service-grpc"
spec:
  type: {{ .Values.service.type }}
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      protocol: TCP
      name: http
  selector:
    app: "sample-service-grpc"
```

ClusterIP を None にする。  
これにより、DNSサーバーは単一のサービスIPの代わりにポッドIPを返すらしい。

```yaml
# 変更後
apiVersion: v1
kind: Service
metadata:
  name: "sample-service-grpc"
spec:
  clusterIP: None
  ports:
    - port: {{ .Values.service.port }}
      targetPort: http
      protocol: TCP
      name: http
  selector:
    app: "sample-service-grpc"
```

名前解決の結果もそうなってそう。

```sh
kubectl exec dnsutils -- nslookup sample-service-grpc
Server:         10.96.0.10
Address:        10.96.0.10#53

Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.117
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.114
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.115
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.113
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.116
```

```sh
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 07:30:45 Received request: /sample.service.HelloWorld/SayHello
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:30:45 | 200 |       2.098ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:30:45 | 200 |    2.363834ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 07:30:45 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-5n2d2 sample-service-helm 2023/11/24 07:30:55 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 07:30:55 Received request: /sample.service.HelloWorld/SayHello
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:30:55 | 200 |    7.400708ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:30:55 | 200 |   12.730125ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 07:31:05 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 07:31:05 Received request: /sample.service.HelloWorld/SayHello
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:31:05 | 200 |    1.613125ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:31:05 | 200 |    2.159209ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-grpc-6db799f6cb-5n2d2 sample-service-helm 2023/11/24 07:31:15 Received request: /sample.service.HelloWorld/SayHello
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:31:15 | 200 |    3.743666ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:31:15 | 200 |    4.533042ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 07:31:15 Received request: /sample.service.HelloWorld/SayHello
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:31:25 | 200 |    2.446458ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-go-765975f4-np9c7 sample-service-helm [GIN] 2023/11/24 - 07:31:25 | 200 |    5.501125ms |      10.244.0.1 | GET      "/go/sample"
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 07:31:25 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 07:31:25 Received request: /sample.service.HelloWorld/SayHello
```

ログ的にも良さそう。ClusterIP がないことによる他の影響を調べる必要はありそうだが、とりあえず分散はできた。

ちなみに、ClusterIP を None にしていても、grpc クライアント側のラウンドロビンの設定をなくすと分散されなかったので、それはそれとして必要そう。

---
(追記) この辺り理解できた。

Service は L4 ロードバランサーで、このロードバランサーの IP が ClusterIP になる。
Service 名を名前解決すると、ClusterIP が返ってくる。
Service 名にアクセスすると、このロードバランサーにアクセスがいき、ロードバランサーが自分に紐ついている Pod の1つにリクエストを送る。

gRPC は HTTP/2 ベースのプロトコルであり、1つの TCP 接続で複数のリクエストを送ることができる。(長期的接続)
つまり、一回接続した相手と一時やり取りを続けることになり、Service によってランダムに選ばれた 1つの Pod とのみ接続されていた。

ClusterIP を None にすると、Service は IP も持たないし、リクエストの分散も行わない。
DNS サーバーは、Service 名を名前解決すると、Service に紐ついている Pod の IP をすべて返す。
gRPC クライアント側は、すべての Pod と接続を確立することができる。
この状態で、gRPC クライアント側のラウンドロビンの設定をすると、リクエストを送る時に使用するコネクションを使い回してくれる。

---

さらに、grpc の Pod を増やしてみる。後から増えた分の Pod にも対応できるかの確認。10個にした。

```sh
kubectl exec dnsutils -- nslookup sample-service-grpc
Server:         10.96.0.10
Address:        10.96.0.10#53

Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.114
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.119
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.117
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.116
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.121
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.120
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.118
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.113
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.122
Name:   sample-service-grpc.default.svc.cluster.local
Address: 10.244.0.1150
```

名前解決は大丈夫そう。
でもリクエストが飛んでいるのは、5つの Pod だけだった。

```sh
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 08:38:48 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 08:38:58 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 08:38:58 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 08:39:08 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 08:39:08 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-5n2d2 sample-service-helm 2023/11/24 08:39:18 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 08:39:18 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 08:39:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 08:39:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 08:39:38 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-5n2d2 sample-service-helm 2023/11/24 08:39:38 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 08:39:48 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 08:39:48 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 08:39:58 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 08:39:58 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-5n2d2 sample-service-helm 2023/11/24 08:40:08 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 08:40:08 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 08:40:18 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 08:40:18 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 08:40:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-5n2d2 sample-service-helm 2023/11/24 08:40:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 08:40:38 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 08:40:38 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 08:40:48 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 08:40:48 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-5n2d2 sample-service-helm 2023/11/24 08:40:58 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 08:40:58 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 08:41:08 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 08:41:08 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-5n2d2 sample-service-helm 2023/11/24 08:41:18 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 08:41:18 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 08:41:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 08:41:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 08:41:38 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-8shf4 sample-service-helm 2023/11/24 08:41:38 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-nb5d2 sample-service-helm 2023/11/24 08:41:48 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-5n2d2 sample-service-helm 2023/11/24 08:41:48 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-sxjqn sample-service-helm 2023/11/24 08:41:58 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-qc9jx sample-service-helm 2023/11/24 08:41:58 Received request: /sample.service.HelloWorld/SayHello
```

これに対応したい。できないと auto scale に対応できない。

https://christina04.hatenablog.com/entry/grpc-client-side-lb#:~:text=%E3%81%8C%E3%81%82%E3%82%8A%E3%81%BE%E3%81%99%E3%80%82-,DNS%E3%81%AEresolve%E3%82%BF%E3%82%A4%E3%83%9F%E3%83%B3%E3%82%B0%E3%81%AF%EF%BC%9F,-Pod%E3%81%AF%E9%A0%BB%E7%B9%81

ここの話がそれっぽい。

再度名前解決が行われるのは、デフォルトだと 30 分後。ポーリングで行っているらしく、これも将来的には辞めようとしているらしい。サーバー側で、定期的にコネクションをリセットする実装を入れるべきらしい。

```go
grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionAge: 10 * time.Second,
		}))
```

interceptor を追加した。10秒は過剰なんだろうが、検証のため短めに設定。

スケールアウトの確認がしづらいので grpc の replicas を 1 に戻してから、再度増やしてみる。

```sh
kubectl get po                                                                                            
NAME                                    READY   STATUS    RESTARTS   AGE
dnsutils                                1/1     Running   0          19h
sample-service-go-765975f4-zzqcf        1/1     Running   0          76s
sample-service-grpc-6db799f6cb-mntmg    1/1     Running   0          76s
sample-service-js-589c4b65f4-vqz49      1/1     Running   0          76s
sample-service-nginx-6d989844fc-5cbms   1/1     Running   0          76s
```

リクエストを送って、ログをみとく。

```sh
 stern -A sample-service-grpc
+ default sample-service-grpc-6db799f6cb-mntmg › sample-service-helm
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:11:33 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:12:33 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:12:33 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:12:43 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:12:43 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:12:53 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:12:53 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:13:03 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:13:03 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:13:13 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:13:13 Received request: /sample.service.HelloWorld/SayHello
```

ここから 5 つに増やしてみる。

```sh
kubectl get po                                                                                            
NAME                                    READY   STATUS    RESTARTS   AGE
dnsutils                                1/1     Running   0          19h
sample-service-go-765975f4-zzqcf        1/1     Running   0          2m46s
sample-service-grpc-6db799f6cb-2n9cr    1/1     Running   0          5s
sample-service-grpc-6db799f6cb-6lp67    1/1     Running   0          5s
sample-service-grpc-6db799f6cb-mntmg    1/1     Running   0          2m46s
sample-service-grpc-6db799f6cb-p67t2    1/1     Running   0          5s
sample-service-grpc-6db799f6cb-vkf5k    1/1     Running   0          5s
sample-service-js-589c4b65f4-vqz49      1/1     Running   0          2m46s
sample-service-nginx-6d989844fc-5cbms   1/1     Running   0          2m46s
```

リクエストを送って、ログ。

```sh
+ default sample-service-grpc-6db799f6cb-2n9cr › sample-service-helm
+ default sample-service-grpc-6db799f6cb-p67t2 › sample-service-helm
+ default sample-service-grpc-6db799f6cb-vkf5k › sample-service-helm
+ default sample-service-grpc-6db799f6cb-6lp67 › sample-service-helm
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:14:03 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:14:03 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:14:13 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:14:13 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:14:23 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-p67t2 sample-service-helm 2023/11/24 09:14:23 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-6lp67 sample-service-helm 2023/11/24 09:14:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-vkf5k sample-service-helm 2023/11/24 09:14:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-2n9cr sample-service-helm 2023/11/24 09:14:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-p67t2 sample-service-helm 2023/11/24 09:14:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:14:28 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-6lp67 sample-service-helm 2023/11/24 09:14:29 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-vkf5k sample-service-helm 2023/11/24 09:14:29 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-vkf5k sample-service-helm 2023/11/24 09:14:29 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-2n9cr sample-service-helm 2023/11/24 09:14:29 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-vkf5k sample-service-helm 2023/11/24 09:14:33 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-6lp67 sample-service-helm 2023/11/24 09:14:33 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-mntmg sample-service-helm 2023/11/24 09:14:43 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-6lp67 sample-service-helm 2023/11/24 09:14:43 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-6lp67 sample-service-helm 2023/11/24 09:14:53 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-vkf5k sample-service-helm 2023/11/24 09:14:53 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-p67t2 sample-service-helm 2023/11/24 09:15:03 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-6lp67 sample-service-helm 2023/11/24 09:15:03 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-2n9cr sample-service-helm 2023/11/24 09:15:13 Received request: /sample.service.HelloWorld/SayHello
default sample-service-grpc-6db799f6cb-vkf5k sample-service-helm 2023/11/24 09:15:13 Received request: /sample.service.HelloWorld/SayHello
```

10秒で最初のコネクションが切断されて、再度名前解決が行われ、追加された Pod にもリクエストが飛ぶようになっている。
`MaxConnectionAge`の最適な値は、どれくらいだろうか。

grpc 側の処理に time.Sleep(15 * time.Second) を入れてみたが、10秒で切断されることはなさそう。
処理中のものは捌いてから、リセットしてくれてそう。流石にか。


