import step1 from "./step1.js";
import step2 from "./step2.js";
import step3 from "./step3.js";

export { step1, step2, step3 };

export const options = {
  scenarios: {
    step1: {
      executor: "per-vu-iterations",
      exec: "step1",
      vus: 5,
      iterations: 5,
      startTime: "0s",
      maxDuration: "20s",
    },
    step2: {
      executor: "per-vu-iterations",
      exec: "step2",
      vus: 5,
      iterations: 5,
      startTime: "0s",
      maxDuration: "20s",
    },
    step3: {
      executor: "per-vu-iterations",
      exec: "step3",
      vus: 1,
      iterations: 1,
      startTime: "20s",
    },
  },
};

export default function () {}
