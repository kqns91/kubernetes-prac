# クラウドネイティブな監視システムの構築 Prometheus 実践ガイド

## Prometheus と監視の基本

### Prometheus の概要

メトリクスの収集
- Pull 型を採用しており、監視エージェントではなくエクスポーターを使ってメトリクスを収集する
- Pull 型は監視サーバーがクライアント IP などを知る必要があり、サービスディスカバリーが用意されている
- 特定のソフトウェアごとにエクスポーターが用意されている（Node exporter, nginx exporter, ...）

アーキテクチャ

![アーキテクチャ](./images/スクリーンショット%202024-06-23%2023.27.53.png)

### Prometheus の実践

式ブラウザ
- Prometheus の Web UI にアクセスして最初に表示される画面。デバッグ用途で利用される。
- Prometheus で PromQL を用いたクエリを実行できる。

```sh
kind create cluster --image=kindest/node:v1.29.0
```

`helmfile.yaml` と `values.yaml` を用意。

```sh
# helmfile.yaml
repositories:
  - name: prometheus-community
    url: https://prometheus-community.github.io/helm-charts

releases:
  - name: kube-prometheus-stack
    namespace: prometheus
    chart: prometheus-community/kube-prometheus-stack
    version: 60.3.0
    values:
      - values.yaml
```

```sh
# values.yaml
```

kube-prometheus-stack に用意されているカスタムリソースが nod found と言われるため、apply せずに sync する。

```sh
helmfile sync
```

```sh
$ kc get pod   
NAME                                                        READY   STATUS    RESTARTS   AGE
alertmanager-kube-prometheus-stack-alertmanager-0           2/2     Running   0          110s
kube-prometheus-stack-grafana-76858ff8dd-xwz7s              3/3     Running   0          110s
kube-prometheus-stack-kube-state-metrics-7f6967956d-lsv4z   1/1     Running   0          110s
kube-prometheus-stack-operator-7fcbcf6c8f-7clld             1/1     Running   0          110s
kube-prometheus-stack-prometheus-node-exporter-hfd6m        1/1     Running   0          110s
prometheus-kube-prometheus-stack-prometheus-0               2/2     Running   0          110s
```

k9s でポートフォワードして、Prometheus と Grafana にアクセス。

![](./images/スクリーンショット%202024-06-25%200.07.31.png)
![](./images/スクリーンショット%202024-06-25%200.07.04.png)
