# つくって、壊して、直して学ぶ Kubernetes 入門

## k8s クラスタ上にシンプルアプリケーションをデプロイする

ローカルに k8s クラスタを構築するために kind を使用する

```sh
kind create cluster --image=kindest/node:v1.29.0
```

ローカルでビルドした docker イメージを使用するには、kind クラスタにイメージをロードし、ImagePullPolicy を IfNotPresent に設定する必要がある。

```sh
# ローカルでビルドした docker イメージを kind クラスタにロードする
kind load docker-image {image-name:tag}
```

```sh
# マニフェストファイルを適用する
kubectl apply -f {manufest-file}
```

## トラブルシューティング

- ガイド
  - pod の status を確認する
  - describe でイベントを確認する
  - image を確認する
  - describe で Probe の結果を確認する
  - ログを確認する
- デバッグ用のサイドカーコンテナ
  - kubectl debug コマンド
- マニフェストをその場で修正して適用する
  - [非推奨] kubectl edit コマンド

プラグイン
- kubens：namespace を切り替える
- kubectx：cluster を切り替える

## kubernetes のアーキテクチャ

![](./スクリーンショット%202024-06-22%2023.20.35.png)

## オブザーバビリティとモニタリング

### ログ

- コンテナの標準出力/エラーの内容をコンテナnのログとして収集する仕組みがある
- Pod が Node から削除されるとログも消えるため、永続化が必要
- ログを外部に転送するために Fluentd や FluentBit などの OSS が利用される

### メトリクス

- メトリクスを収集するために Prometheus が利用される
- 独自の PromQL というクエリ言語を利用して収集したメトリクスを参照する

### トレース

- リクエストが通る系路上の全てで Traces/Span が導入できている必要がある
- トレースの実装を簡易化するために OpenTelemetry が利用される
- トレースを収集するために Jaeger, Grafana Tempo が利用される

### 構築

```sh
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo update
kubectl create namespace monitoring
helm install kube-prometheus-stack -n monitoring prometheus-community/kube-prometheus-stack
```

```sh
kc get pod
NAME                                                        READY   STATUS    RESTARTS   AGE
alertmanager-kube-prometheus-stack-alertmanager-0           2/2     Running   0          47s
kube-prometheus-stack-grafana-76858ff8dd-ldlcn              3/3     Running   0          59s
kube-prometheus-stack-kube-state-metrics-7f6967956d-brggf   1/1     Running   0          59s
kube-prometheus-stack-operator-7fcbcf6c8f-g799k             1/1     Running   0          59s
kube-prometheus-stack-prometheus-node-exporter-x7pzc        1/1     Running   0          59s
prometheus-kube-prometheus-stack-prometheus-0               2/2     Running   0          46s
```

