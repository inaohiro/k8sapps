# トラブルシューティング

## namespace が大量に存在する

k6 でリクエストを送信するとき、namespace を作っています  
k6 のシナリオ実行後消すようにしていますが、消えずに残っている場合以下を実行してください

```sh
kubectl get ns -o name | xargs kubectl delete
```
