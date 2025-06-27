import http from "k6/http";

const url = "http://gateway-clusterip";

export default function () {
  const namespaces = [fakeName(), fakeName(), fakeName()];

  for (const namespace of namespaces) {
    const issueTokenReq = { namespace: namespace };
    let res = http.post(url, JSON.stringify(issueTokenReq), {
      headers: { "Content-Type": "application/json" },
    });
    const token = res.json().token;

    http.get("http://gateway-clusterip/deployments", {
      headers: { Authorization: `Bearer ${token}` },
    });
    http.get("http://gateway-clusterip/pods", {
      headers: { Authorization: `Bearer ${token}` },
    });
    http.get("http://gateway-clusterip/services", {
      headers: { Authorization: `Bearer ${token}` },
    });
    http.get("http://gateway-clusterip/images", {
      headers: { Authorization: `Bearer ${token}` },
    });
    http.get("http://gateway-clusterip/flavors", {
      headers: { Authorization: `Bearer ${token}` },
    });
  }
}

function fakeName() {
  return Math.random().toString(32).substring(2);
}
