# k8sapps

## 環境構築手順

manifest 以下に kubernetes のマニフェストファイルを用意しています  
`kubectl apply` を実行してください

```sh
git clone https://github.com/inaohiro/k8sapps.git

kubectl apply -f manifest/
```

## 動作確認

Service type: LoadBalancer で用意しているものが 3 つあります  
それぞれアクセスできることを確認してください

- ブラウザからアクセスできる画面
  - `http://<IP アドレス>:8080`
- Grafana
  - `http://<IP アドレス>:3000`
- VictoriaMetrics
  - `http://<IP アドレス>:8428`

## Troubleshooting

### ghcr.io からイメージをダウンロードする

コンテナイメージを ghcr.io から取得することもできます  
その場合、以下の修正が必要となります

## 1. secret.yaml の作成

```yaml
apiVersion: v1
data:
  .dockerconfigjson: <ここに base64 エンコードした ~/.docker/config.json を貼る>
kind: Secret
metadata:
  name: regcred
type: kubernetes.io/dockerconfigjson
```

```sh
kubectl apply -f secret.yaml
```

## 2. Deployment に imagePullSecret を追加、イメージ名を変更

spec.template.spec 以下に imagePullSecrets を追加してください
また、イメージ名を ghcr.io から始まるようにしてください
以下のような感じになります

```yaml
spec:
  template:
    spec:
      imagePullSecrets:
        - name: regcred
      containers:
        - name: ghcr.io/inaohiro/webapp:v1
```
