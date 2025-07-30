# グラフ, クエリ対照表 <!-- omit in toc -->

作成したいグラフと、それを実現するためのクエリ例の対応を表にしたものです

- [総リクエスト数](#総リクエスト数)
- [リクエスト数の遷移](#リクエスト数の遷移)
- [成功率](#成功率)
- [レイテンシ](#レイテンシ)

## 総リクエスト数

- running_sum\(sum\(rate\(http_client_request_duration_seconds_count\[$\_\_rate_interval])))
- running_sum\(sum\(rate\(http_server_request_duration_seconds_count\[$\_\_rate_interval])))
- running_sum\(sum\(rate\(traces_spanmetrics_calls_total\[$\_\_rate_interval]))) |

## リクエスト数の遷移

- sum by (http_request_method) (rate(http_client_request_duration_seconds_count\[$\_\_rate_interval]))
- sum by (http_request_method, service_name) (rate(http_server_request_duration_seconds_count\[$\_\_rate_interval]))
- sum by (span_name) (rate(traces_spanmetrics_calls_total{span_kind="SPAN_KIND_SERVER"}\[$\_\_rate_interval]))

## 成功率

- http_client_request_duration_seconds_count の場合
  ```
  # total とする
  running_sum(sum(rate(http_client_request_duration_seconds_count[$__rate_interval])))
  # success とする
  running_sum(sum(rate(http_client_request_duration_seconds_count{http_response_status_code=~"2.."}[$__rate_interval])))
  # Expression
  $success / $total * 100
  ```
- http_server_request_duration_seconds_count

  ```
  # total とする
  running_sum(sum(rate(http_server_request_duration_seconds_count[$__rate_interval])))
  # success とする
  running_sum(sum(rate(http_server_request_duration_seconds_count{http_response_status_code!~"(4..)|(5..)"}[$__rate_interval])))
  # Expression
  $success / $total * 100
  ```

- traces_spanmetrics_calls_total の場合

  ```
  # total とする
  running_sum(sum(rate(traces_spanmetrics_calls_total{span_kind="SPAN_KIND_SERVER"}[$__rate_interval])))
  # error とする
  running_sum(sum(rate(traces_spanmetrics_calls_total{span_kind="SPAN_KIND_SERVER", status_code="STATUS_CODE_ERROR"}[$__rate_interval])))
  # Expression
  ($total - $error) / $total * 100
  ```

## レイテンシ

> http_client_request_duration_seconds_bucket, http_server_request_duration_seconds_bucket の場合、値がちょっと違います  
> 何か間違っているかもです ...

- histogram_quantile(.5, sum by (le) (rate(http_client_request_duration_seconds_bucket[$__rate_interval])))
- histogram_quantile(.75, sum by (le) (rate(http_server_request_duration_seconds_bucket[$__rate_interval])))
- histogram_quantile(.95, sum by (le, span_name) (rate(traces_spanmetrics_latency_bucket{span_kind="SPAN_KIND_SERVER"}[$__rate_interval])))
