import { useAtom, useAtomValue, useSetAtom } from "jotai";
import React, { useEffect } from "react";
import { AppPage } from "./App.page";
import DeploymentCreate from "./component/DeploymentCreate";
import { DeploymentList } from "./component/DeploymentList";
import { PodList } from "./component/PodList";
import { ServiceList } from "./component/ServiceList";
import { useIssueToken } from "./hooks/useIssueToken";
import { pageAtom, PageState, setPageAtom } from "./store/store";
import { TokenIssuePage } from "./component/TokenIssuePage";
import { GlobalHeader } from "./GlobalHeader.page";
import { Page } from "./Page.page";
import { ServiceCreate } from "./component/ServiceCreate";
import { PodCreate } from "./component/PodCreate";
import { DeploymentDetail } from "./component/DeploymentDetail";
import { ServiceDetail } from "./component/ServiceDetail";
import { PodDetail } from "./component/PodDetail";

export function App() {
  const { token, issueToken, loading, error } = useIssueToken();
  const page = useAtomValue(pageAtom);
  const setPage = useSetAtom(setPageAtom);

  // token, path から表示するページを決定
  useEffect(() => {
    setPage(getPage(token, window.location.pathname))
  }, [])

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    return <div>Error: {error}</div>;
  }

  let content: React.ReactNode = null;
  switch (page.type) {
    case "deployments-list":
      content = <DeploymentList />;
      break;
    case "deployments-create":
      content = <DeploymentCreate />;
      break;
    case "deployments-detail":
      content = <DeploymentDetail id={page.id} />;
      break;
    case "services-list":
      content = <ServiceList />;
      break;
    case "services-create":
      content = <ServiceCreate />;
      break;
    case "services-detail":
      content = <ServiceDetail id={page.id} />;
      break;
    case "pods-list":
      content = <PodList />;
      break;
    case "pods-create":
      content = <PodCreate />;
      break;
    case "pods-detail":
      content = <PodDetail id={page.id} />;
      break;
    case "token-issue":
      return (
        <TokenIssuePage
          issueToken={issueToken}
          loading={loading}
          onTokenIssued={() => setPage({ type: "deployments-list" })}
        />
      );
  }

  return (
    <>
      <AppPage>
        <GlobalHeader issueToken={issueToken} setPage={setPage} />
        <Page>{content}</Page>
      </AppPage>
    </>
  );
}

function getPage(token: string|null, path: string): PageState {
  if (token === "" || token === null || token === undefined || token === "undefined") {
    return {type: "token-issue"}
  }

  const elements = path.split("/")
  if (elements.length === 0) {
    return {type: "deployments-list"};
  }

  if (elements.length === 1) {
    if (elements[0] === "services") return {type: "services-list"}
    if (elements[0] === "pods") return {type: "pods-list"}
    return {type: "deployments-list"}
  }

  if (elements.length === 2) {
    if (path === "/deployments/create") return {type: "deployments-create"}
    if (path === "/services/create") return {type: "services-create"}
    if (path === "/pods/create") return {type: "pods-create"}
    return {type: "deployments-list"}
  }

  if (elements.length === 3) {
    if (path.startsWith("/deployments/detail")) return {type: "deployments-detail", id: elements[2]}
    if (path.startsWith("/services/detail")) return {type: "services-detail", id: elements[2]}
    if (path.startsWith("/pods/detail")) return {type: "pods-detail", id: elements[2]}
    return {type: "deployments-list"}
  }

  return {type: "deployments-list"}
}
