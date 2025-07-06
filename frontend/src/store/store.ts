import { atom } from "jotai";

export type PageState =
  | { type: "token-issue" }
  | { type: "deployments-list" }
  | { type: "deployments-create" }
  | { type: "deployments-detail"; id: string }
  | { type: "services-list" }
  | { type: "services-create" }
  | { type: "services-detail"; id: string }
  | { type: "pods-list" }
  | { type: "pods-create" }
  | { type: "pods-detail"; id: string };

const baseAtom = atom<PageState>({ type: "deployments-list" });
export const pageAtom = atom((get) => get(baseAtom));

// URL の更新、Atom の更新
export const setPageAtom = atom(null, (get, set, state: PageState) => {
  const page = get(baseAtom);
  set(baseAtom, state);

  // URL の更新
  const path = `/${state.type.replace("-", "/")}`;
  if (state.type === "deployments-detail" || state.type === "services-detail" || state.type === "pods-detail") {
    window.history.pushState({}, "", `${path}/${(state as { id: string }).id}`);
  } else {
    window.history.pushState({}, "", path);
  }
});
