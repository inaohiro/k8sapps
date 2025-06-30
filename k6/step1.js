import http from "k6/http";

const url = "http://nginx:8080/api";

export default function () {
  const namespaces = [fakeName(), fakeName(), fakeName()];

  for (const namespace of namespaces) {
    const issueTokenReq = { namespace: namespace };
    let res = http.post(url, JSON.stringify(issueTokenReq), {
      headers: { "Content-Type": "application/json" },
    });
    const token = res.json().token;
    const headers = {
      headers: { Authorization: `Bearer ${token}` },
    };

    http.get(`${url}/deployments`, headers);
    http.get(`${url}/pods`, headers);
    http.get(`${url}/services`, headers);
    http.get(`${url}/images`, headers);
    http.get(`${url}/flavors`, headers);
  }
}

function fakeName() {
  return Math.random().toString(32).substring(2);
}
