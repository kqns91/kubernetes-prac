# 完全ガイド

## chapter4

- Kubernetes リソースカテゴリ
    - Workload APIs
    - Service APIs
    - Config and Storage APIs
    - Cluster APIs
    - Metadata APIs

- kubectl
    - kubeconfig
        - 認証情報を設定する
            - clusters
            - users
            - contexts：clusters と users の組み合わせと namespace を指定したもの。
    - リソースの作成
        - create よりも apply を使う方が良い

### Workload APIs カテゴリ

#### Pod

- DNS 設定
    - ClusterFirst（デフォルト）：クラスタ内の DNS サーバーを優先する
    - None：外部の DNS サーバーを優先する
    - Default
    - ClusterFirstWithHostNet
    - 静的な DNS 設定：spec.hostAliases で設定する

#### ReplicaSet

- 指定した数の Pod が実行されていることを保証する
- ラベルによって監視する Pod を判断している

#### Deployment

- 複数の ReplicaSet を管理することで、ローリングアップデートやロールバックを行う

#### DaemonSet

- 全てのノードで Pod を実行する

#### StatefulSet

- Pod に一意な ID を割り当てる

#### Job

- 一度だけ処理を実行する Pod を作成する

#### CronJob

- 定期的に Job を実行する

### Service APIs カテゴリ

コンテナのサービスディスカバリや、クラスタの外部からもアクセス可能なエンドポイントを提供するためのリソース。
L4 ロードバランシングを呈する `Service` と、L7 ロードバランシングを提供する `Ingress` がある。

- Service
    - 受信したトラフィックを、複数の Pod へロードバランシングする。
    - サービスディスカバリ
        - 環境変数を利用したサービスディスカバリ
        - DNS A レコードを利用したサービスディスカバリ
        - DNS SRV レコードを利用したサービスディスカバリ
    - Node Local DNS Cache
        - Node に DNS キャッシュを持たせることで、クラスタ内の DNS クエリを高速化する。
    - ClusterIP Service
        - クラスタ内ロードバランサー。
        - ClusterIP 当ての通信は kube-proxy が Pod に転送を行う。
        - IP アドレスを静的に指定することも可能。
        - すでに作成する Service に対して ClusterIP を変更することはできない。
        - セッションアフィニティを有効にすることで、同じ Pod に対しての通信を維持することができる。
    - ExternalIP Service
        - type は ClusterIP。spec.externalIPs で Kubernetes Node の IP アドレスを指定する。
        - 特定の Kubernetes Node の IP アドレス:Port で受信したトラフィックをコンテナに転送する。
    - NodePort Service
        - 全ての Kubernetes Node の IP アドレス:NodePort で受信したトラフィックをコンテナに転送する。
        - 30000-32767 の範囲で指定する。
        - spec.externalTrafficPolicy の設定により、クラスタ外からのトラフィックを 別の Node の Pod に転送するかどうかを指定できる。
    - LoadBalancer Service
        - Kubernetes クラスタ外のロードバランサに外部疎通性のある仮想 IP を払い出す。
        - GCP/AWS/Azure/OpenStack などのプロバイダーで利用できるため、障害に強い。
        - デフォルトだとグローバルに公開されるため、spec.loadBalancerSourceRanges に接続を許可する送信元ネットワークを指定することでプロバイダのファイアウォール機能を利用したアクセス制御が可能。
    - Topology-aware Service routing
        - 通常 Service の転送先は Region や Availability Zone を考慮されない。ノード数が多すぎるとパフォーマンスが低下する。
        - 同一ノード、同一ゾーン、いずれかの Pod というような優先度指定ができる。
        - 送信元 IP アドレスは取得できない。


### Cluster APIs カテゴリ

- Namespace：仮想的なクラスタ
    - kube-sysmtem：クラスタのコンポーネント
    - kube-public：全てのユーザーが読み取り可能。ConfigMap など。
    - kube-node-lease：ノードのヘルスチェック
    - default：デフォルトの名前空間