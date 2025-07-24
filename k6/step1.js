import { sleep } from 'k6';
import { fakeName, retry } from "./common.js";

const url = "http://gateway:8080/api";

/**
 * 複数の namespace のリソース一覧を取得する
 * もしリソースがあれば、詳細取得も実行する
 */
export default function () {
  const namespace = fakeName();

  // token 発行
  const issueTokenReq = { namespace: namespace };
  let res = retry(
    "post",
    `${url}/tokens`,
    {
      headers: { "Content-Type": "application/json" },
    },
    JSON.stringify(issueTokenReq)
  );

  // レスポンスから token 取得
  const token = res.json().token;

  // これ以降で必要な Authorization header
  const headers = {
    headers: { Authorization: `Bearer ${token}` },
  };

  res = retry("get", `${url}/deployments`, headers);
  if (res.body.length > 100) {
    let name = res.body[0]["name"];
    retry("get", `${url}/deployments/${name}`, headers);
  }

  res = retry("get", `${url}/pods`, headers);
  if (res.body.length > 100) {
    let name = res.body[0]["name"];
    retry("get", `${url}/pods/${name}`, headers);
  }

  res = retry("get", `${url}/services`, headers);
  if (res.body.length > 100) {
    let name = res.body[0]["name"];
    retry("get", `${url}/services/${name}`, headers);
  }

  retry("get", `${url}/images`, headers);
  retry("get", `${url}/flavors`, headers);

  // 終わったら namespace を消す
  retry("del", `${url}/namespace/${namespace}`, headers)
}
