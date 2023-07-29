# kind

kubernetes in docker

- ローカル kubernetes 環境。
- Dockerコンテナを複数個起動し、そのコンテナを kubernetes node として利用することで、複数台構成の kubernetes クラスタを構成する。

[kindファイル](./kind.yaml)を作成し、以下のコマンドで起動する。

```cmd
kind create cluster --config kind.yaml --name kindcluster
```

