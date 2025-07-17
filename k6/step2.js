import { sleep } from 'k6';
import { check } from "k6";
import tempo from "./jslib.js";

const http = new tempo.Client({
  propagator: "w3c",
});

const url = "http://gateway:8080/api";

/**
 * Deployment, Service を作成する
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
      headers: {
        Authorization: `Bearer ${token}`,
        "Content-Type": "application/json",
      },
    };

    // image 一覧
    res = retry("get", `${url}/images`, headers);
    let images = res.json();
    let index = Math.floor(Math.random() * images.length);
    let image_name = images[index].name;

    // 作成時のリクエストボディ
    const name = fakeName();
    const createDeployment = { name: name, image: image_name };
    const createService = {
      name: name,
      type: "ClusterIP",
      ports: [{ port: 80, targetPort: 80, protocol: "TCP" }],
    };

    res = retry(
      "post",
      `${url}/deployments`,
      headers,
      JSON.stringify(createDeployment)
    );
    check(res, { "response code was 200": (res) => res.status == 200 });
    res = retry("get", `${url}/deployments/${name}`, headers);
    check(res, { "response code was 200": (res) => res.status == 200 });
    res = retry(
      "post",
      `${url}/services`,
      headers,
      JSON.stringify(createService)
    );
    check(res, { "response code was 200": (res) => res.status == 200 });
    res = retry("get", `${url}/services/${name}`, headers);
    check(res, { "response code was 200": (res) => res.status == 200 });

    sleep(5);

    retry("del", `${url}/deployments/${name}`, headers)
    retry("del", `${url}/services/${name}`, headers)

    sleep(10);

    // 終わったら namespace を消す
    retry("del", `${url}/namespace/${namespace}`, headers)
  }
}

function fakeName() {
  return Math.random().toString(32).substring(2);
}

function retry(method, url, params, body, count) {
  if (count === undefined) {
    count = 0;
  }

  var res;
  if (method === "get" || method === "del") {
    res = http[method](url, params);
  } else {
    res = http[method](url, body, params);
  }
  if (res.status < 300) {
    return res;
  }

  // 最大 5 回まで
  if (count > 5) {
    return res;
  }

  sleep(1 + (2**count) * Math.round(1 + Math.random()));

  return retry(method, url, params, body, count+1)
}
