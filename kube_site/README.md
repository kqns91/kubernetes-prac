# 現場で使えるKubernetes

## chapter01

### アーキテクチャ

- Cluster
  - Control Plane
    - kube-apiserver
      - 司令塔
    - kube-scheduler
      - 配置可能なnodeを探し、リソースを配置するのにふさわしいnodeを決定する
    - etcd
      - キーバリューストア。構成を保存する。
    - kube-controller-manager
      - 複数のコントローラーからなる。nodeで実行されるワークロードをの制御を行う。
    - cloud-controller-manager
      - AWSやGCPなどのサービスと連携する。
  - Nodes
    - kubelet
      - node上で動作するエージェント。kube-apiserverからリクエストを受け取り定義通りにpodを起動。
    - kube-proxy
      - cluster内外のネットワークトラフィックの制御を行う。
    - container-runtime
      - nodeでpodを実行するためのコンテナ実行環境。

### 環境構築

書籍では普通にkubernetesを使っているが、Minikubeで進める。

#### Minikubeインストール
- [公式ドキュメント](https://minikube.sigs.k8s.io/docs/start/)に従ってインストール。

```
brew install minikube
```

- クラスター開始

```
minikube start
```

- クラスター操作

```
kubectl get po -A
```

- ダッシュボードを開く

うわあ、いつも見るやつ~

```
minikube dashboard
```

- 停止

```
minikube stop
```

- 削除

```
minikube delete --all
```

### ワークロードリソース

#### Pod
- 1つ以上のコンテナで構成されるワークロードの最小単位。各podのCPUとメモリは共有される。コンテナはIPを持たず、localhostとportで名前解決できる。 
- マニュフェストを作成しkubectlでPod作成、確認、コンテナ接続などできる。

#### Replicaset
- マニフェスト記載のレプリカ数に合わせるようにPodの数を維持する。
- Podを削除しても勝手に立ち上がる。

#### Deployment
- 実際の開発や運用現場では基本PodやReplicaSetを直接定義して操作しない。
- PodやReplicaSetを管理し、コンテナアプリケーションのデプロイに関する機能を提供する。
- マニュフェストを更新して`kubectl apply`するとロールアウトされる。

#### Job
- 複数のPodを実行できる。
- Podの正常な終了を追跡することが目的。
- 並列化や、Pod内のプロセス途中で失敗した際のリトライ回数の設定が可能。

#### CronJob
- 定期的にjobを作成する。
- 分、時間、日、月、曜日で指定。

#### DemonSet
- Nodeで1つのPodのコピーが稼働することを保証する。
- ClusterにNodeが追加・削除された場合、Podも追加・削除される。
- 基本的に1つのNodeに1つずつPodを実行する。
- Clusterレベルの監視やロギングに使われがち。

#### StatefulSet
- 安定的にネットワークIDをもち、永続ストレージを提供する。
- 安定的なネットワークID
  - Podのレプリカ数を指定できる。 
  - replicasetやPodが削除されても、Podに割り当てられる名前は変わらない。
- 永続的なストレージの提供
  - RersistentVolumeリソースを使ってそれぞれのPodに個別の永続ストレージを提供する。

### ネットワーキング
#### Service
- Nodeに作成されたPodのIPは基本的に不変ではないため、名前解決が必要。
- Pod間通信での名前解決はServiceが必要。
- ServiceにはTypeが色々ある。

#### Service Type
##### ClusterIP
- デフォルト。
- Cluster内部のIPで内部のみに公開。

##### NodePort
- Clusterの各Nodeにポートを割り当て外部からの通信をCluster内部のPodに転送する。

##### LoadBalancer
- クラウドプロバイダのロードバランシングサービスを使用して、Cluster外部に公開する。
- EKSでLoadBalancerのServiceを作成するとAWSのCloud Provider Load Balancer ControllerによってCLBまたはNLBが作成される。
- L4で通信する。

#### Ingress
- LoadBalancerと同じくPodをCluster外部に公開する。
- IngressはL7で通信する。
- Ingressコントローラーを作成しないと、IngressリソースをClusterに作成しただけでは使えない。
- Ingressコントローラーは環境やプロバイダに合わせて選択する必要がある。

### Storage/Configリソース

#### Volume
- Node内部のローカルストレージ、Amazon Elastic Block Store、Google Compute Engineなど外部ストレージを利用可能。

#### Persistent Volume / PersistentVolumeClaim
- Podとは切り離して単体で作成可能。
- PersistentVolumeClaimでPodとPersisntentVolumeを紐づける。
- PersisntentVolumeへのアクセスモード設定が可能。

#### ConfigMap
- キーバリュー形式でデータをClusterに登録できる。
- 1MB未満。

#### Secret
- configMapは平文、secretはbase64でエンコードされる。
- secretリソースはgitで管理しないか、RBACを使うなどすべき。


## chapter02

### Custom Resources
- CRDでユーザーが独自にリソースを定義し、CRでそのリソースを作成できる。
- Custom controllersで定期的に監視し、ユーザーが定義したあるべき状態を維持する。

### Kustomize
- kubernetesのマニュフェスト管理ツール。
- 開発、テスト、本番環境など複数の環境の一元管理が目的。
- 共通リソースと環境固有のリソースに分けて作成。
- kustomize cliかv1.14以上のkubectlで操作。

### Helm
- kubernetes専用パッケージマネージャー。
- chartを操作することでClusterへのリソースの作成、更新、削除、ロールバックなどが可能。
- Chartによって作成されたインスタンスはリリースと呼ばれる。

```
brew install helm
```

- ClusterにReleaseを作成手順

```
// Chartリポジトリの追加
helm repo add eks ${HELM_REPOSITORY}

// リポジトリの更新
helm repo update

// ClusterへのChartのインストール
helm install
```

