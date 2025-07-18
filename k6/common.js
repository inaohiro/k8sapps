import { sleep } from 'k6';
import tempo from "./jslib.js";

const http = new tempo.Client({
  propagator: "w3c",
});

export function fakeName() {
  return "a" + Math.random().toString(32).substring(2);
}

export function retry(method, url, params, body, count) {
  if (count === undefined) {
    count = 0;
  }

  var res;
  if (method === "get") {
    res = http[method](url, params);
  } else if (method === "del") {
    res = http[method](url, null, params);
  } else {
    res = http[method](url, body, params);
  }
  if (res.status < 300) {
    return res;
  }

  // 最大 5 回まで
  if (count >= 5) {
    return res;
  }

  sleep(1);

  return retry(method, url, params, body, count+1)
}
