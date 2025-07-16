import { check } from "k6";
import tempo from "./jslib.js";

const http = new tempo.Client({
  propagator: "w3c",
});

const url = "http://gateway:8080/api";

/**
 * 複数の namespace のリソース一覧を取得する
 * もしリソースがあれば、詳細取得も実行する
 */
export default function () {
  const namespaces = [fakeName(), fakeName(), fakeName()];

  for (const namespace of namespaces) {
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
    check(res, { "response code was 200": (res) => res.status == 200 });
    if (res.body.length > 100) {
      let name = res.body[0]["name"];
      retry("get", `${url}/deployments/${name}`, headers);
    }

    res = retry("get", `${url}/pods`, headers);
    check(res, { "response code was 200": (res) => res.status == 200 });
    if (res.body.length > 100) {
      let name = res.body[0]["name"];
      retry("get", `${url}/pods/${name}`, headers);
    }

    res = retry("get", `${url}/services`, headers);
    check(res, { "response code was 200": (res) => res.status == 200 });
    if (res.body.length > 100) {
      let name = res.body[0]["name"];
      retry("get", `${url}/services/${name}`, headers);
    }

    res = retry("get", `${url}/images`, headers);
    check(res, { "response code was 200": (res) => res.status == 200 });

    res = retry("get", `${url}/flavors`, headers);
    check(res, { "response code was 200": (res) =>ores.status == 200 });

    // 終わったら namespace を消す
    retry("del", `${url}/namespace/${namespace}`, headers)
  }
}

function fakeName() {
  return Math.random().toString(32).substring(2);
}

function retry(method, url, params, body) {
  var res;
  for (let retries = 3; retries > 0; retries--) {
    if (method === "get" || method === "del") {
      res = http[method](url, params);
    } else {
      res = http[method](url, body, params);
    }
    if (res.status < 300) {
      return res;
    }
  }

  return res;
}
