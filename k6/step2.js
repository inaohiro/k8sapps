import { sleep } from 'k6';
import { fakeName, retry } from "./common.js";

const url = "http://gateway:8080/api";

/**
 * Deployment, Service を作成する
 */
export default function () {
  const namespace = fakeName()

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

  retry(
    "post",
    `${url}/deployments`,
    headers,
    JSON.stringify(createDeployment)
  );
  retry("get", `${url}/deployments/${name}`, headers);
  res = retry(
    "post",
    `${url}/services`,
    headers,
    JSON.stringify(createService)
  );
  if (res.status !== 500) {
    retry("get", `${url}/services/${name}`, headers);
    retry("del", `${url}/services/${name}`, headers);
  }

  retry("del", `${url}/deployments/${name}`, headers)
  while(true) {
    res = retry("get", `${url}/deployments/${name}`, headers, null, 10)
    if (res.status == 404) {
      break
    }
    sleep(1)
  }

  // 終わったら namespace を消す
  retry("del", `${url}/namespace/${namespace}`, headers, null, -10)
}
