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

### Cluster APIs カテゴリ

- Namespace：仮想的なクラスタ
    - kube-sysmtem：クラスタのコンポーネント
    - kube-public：全てのユーザーが読み取り可能。ConfigMap など。
    - kube-node-lease：ノードのヘルスチェック
    - default：デフォルトの名前空間