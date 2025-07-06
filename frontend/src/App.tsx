import { useAtom, useAtomValue, useSetAtom } from "jotai";
import React, { useEffect } from "react";
import { AppPage } from "./App.page";
import DeploymentCreate from "./component/DeploymentCreate";
import { DeploymentList } from "./component/DeploymentList";
import { PodList } from "./component/PodList";
import { ServiceList } from "./component/ServiceList";
import { useIssueToken } from "./hooks/useIssueToken";
import { pageAtom, setPageAtom } from "./store/store";
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

  // 初回レンダリング時に cookie から token を取得
  useEffect(() => {
    if (token !== null) {
      setPage({ type: "token-issue" });
      return;
    }
  }, []);

  // token があれば URL から表示するページを推測
  useEffect(() => {
    if (token === "" || token === null) return;

    const path = window.location.pathname;
    if (path.startsWith("/deployments")) {
      if (path === "/deployments/new") {
        setPage({ type: "deployments-create" });
      } else if (path.startsWith("/deployments/")) {
        const id = path.split("/")[2];
        setPage({ type: "deployments-detail", id });
      } else {
        setPage({ type: "deployments-list" });
      }
    } else if (path.startsWith("/services")) {
      if (path === "/services/new") {
        setPage({ type: "services-create" });
      } else if (path.startsWith("/services/")) {
        const id = path.split("/")[2];
        setPage({ type: "services-detail", id });
      } else {
        setPage({ type: "services-list" });
      }
    } else if (path.startsWith("/pods")) {
      if (path === "/pods/new") {
        setPage({ type: "pods-create" });
      } else if (path.startsWith("/pods/")) {
        const id = path.split("/")[2];
        setPage({ type: "pods-detail", id });
      } else {
        setPage({ type: "pods-list" });
      }
    } else {
      setPage({ type: "deployments-list" });
    }
  }, [token]);

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
