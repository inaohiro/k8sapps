import { retry } from "./common.js";

const url = "http://gateway:8080/api";

export default function () {
  retry("del", `${url}/namespace/_all`);
}
