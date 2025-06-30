import { atom } from "jotai";

export type PageState =
  | { type: 'token-issue' }
  | { type: 'deployments-list' }
  | { type: 'deployments-create' }
  | { type: 'deployments-detail'; id: string }
  | { type: 'services-list' }
  | { type: 'services-create' }
  | { type: 'services-detail'; id: string }
  | { type: 'pods-list' }
  | { type: 'pods-create' }
  | { type: 'pods-detail'; id: string };

export const pageAtom = atom<PageState>({ type: 'deployments-list' });
