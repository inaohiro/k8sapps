import step1 from "./step1.js";
// import step2 from "./step2";
// import step3 from "./step3";

// export { step1, step2, step3 };
export { step1 };

export const options = {
  scenarios: {
    step1: {
      executor: "per-vu-iterations",
      exec: "step1",
      vus: 1,
      iterations: 1,
      startTime: "0s",
    },
    // step2: {
    //   executor: "per-vu-iterations",
    //   exec: "step2",
    //   vus: 1,
    //   iterations: 1,
    //   startTime: "0s",
    // },
    // step3: {
    //   executor: "per-vu-iterations",
    //   exec: "step3",
    //   vus: 1,
    //   iterations: 1,
    //   startTime: "0s",
    // },
  },
};

export default function () {}
