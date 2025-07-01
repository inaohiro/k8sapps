import { useAtom } from "jotai";
import React, { useEffect } from "react";
import { AppPage } from "./App.page";
import DeploymentCreatePage from "./component/DeploymentCreatePage";
import DeploymentsPage from "./component/DeploymentsPage";
import PodsPage from "./component/PodsPage";
import ServicesPage from "./component/ServicesPage";
import { useIssueToken } from "./hooks/useIssueToken";
import { pageAtom } from "./store/store";
import TokenIssuePage from "./component/TokenIssuePage";

export function App() {
  const { token, issueToken, loading, error } = useIssueToken();
  const [page, setPage] = useAtom(pageAtom);

  // 初回レンダリング時に cookie から token を取得
  // token があれば URL から表示するページを推測
  useEffect(() => {
    if (token !== null) {
      setPage({ type: "token-issue" });
      return;
    }

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
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }
  if (error) {
    return <div>Error: {error}</div>;
  }

  let content: React.ReactNode = null;
  switch (page.type) {
    case "deployments-list":
      content = <DeploymentsPage />;
      break;
    case "deployments-create":
      content = <DeploymentCreatePage />;
      break;
    case "services-list":
      content = <ServicesPage />;
      break;
    case "pods-list":
      content = <PodsPage />;
      break;
    // TODO: 他の画面（作成・詳細）も同様に分岐を追加
    case "token-issue":
      return <TokenIssuePage onTokenIssued={() => setPage({type: "deployments-list"})}/>
    default:
      content = <div>{page.type} Not implemented</div>;
  }

  return <AppPage content={content} handleClick={issueToken} />;
}
