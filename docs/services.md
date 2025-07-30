# サービス構成

## 全体構成

Grafana, Prometheus 等を含めた全体の構成です

<img src="./assets/otel_service_architecture.png" width=512>

- k6
  - 負荷テスト、パフォーマンステスト等で用いられるツールです
  - 設定した通りにリクエストを送信することができるため、利用しています
- Nginx
  - 外部からアクセスするときに Nginx を通してアプリケーションにリクエストが届きます
- Otel Collector
  - アプリケーションに計装した OpenTelemetry データを Otel Collector で集めます
  - その後 Tempo, Prometheus, VictoriaMetrics に送信しています
- Tempo
  - Grafana と連携することで Service Graph が表示できるらしい、と見て導入しています
  - 結果 traces_spanmetrics 系のメトリクスが使えるようになりました
- VictoriaMetrics
  - Prometheus 互換の時系列データベースです
  - 時間方向の合計 (累積値, running_sum) が Prometheus では出来なさそうだったので、使うことにしました
- Prometheus
  - Service Graph でエラー率がぱっと見で分かるのは Prometheus を使ったときのみでしたので、残しています
- Grafana
  - OSS なグラフ描画ツールといえば Grafana ! と思い使っています
  - Grafana 社製ツール (k6, Tempo 等) との連携もできとても便利です

## OpenTelemetry の測定対象の構成

<img src="./assets/grafana_service_graph.png" width=512>

OpenTelemetry の測定対象、という観点で登場するサービスは以下の 4 つです

- gateway
- auth
- webapp
- db

### gateway

アプリケーションへのリクエストは必ず gateway を経由するようにしています  
リクエストによって、トークン発行やトークン検証の実施、またはそのまま通過させるなどします

### auth

トークンの発行、検証をするサービスです  
Auth という名前ですがユーザ認証などはしていません

### webapp

Go の kubernetes クライアントである client-go を使って kubernetes を操作するアプリケーション本体です

### db

MySQL データベースです  
webapp は kubernetes の操作が主なため、利用は簡単な SELECT のみです
