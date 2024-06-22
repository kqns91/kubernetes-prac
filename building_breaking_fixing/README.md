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
