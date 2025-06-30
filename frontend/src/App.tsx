import React from 'react';
import { atom, useAtom } from 'jotai';
import { useToken } from './hooks/useToken';
import TokenIssuePage from './TokenIssuePage';
import DeploymentsPage from './DeploymentsPage';
import ServicesPage from './ServicesPage';
import PodsPage from './PodsPage';
import DeploymentCreatePage from './DeploymentCreatePage';
import { AppPage } from './App.page';
import { pageAtom } from './store/store';


export function App() {
  const { hasToken, refresh } = useToken();
  const [page, setPage] = useAtom(pageAtom);

  const handleTokenIssued = () => {
    refresh();
    setPage({ type: 'deployments-list' });
  };

  const handleReissueClick = () => {
    setPage({ type: 'token-issue' });
  };

  if (!hasToken || page.type === 'token-issue') {
    return <TokenIssuePage onTokenIssued={handleTokenIssued} />;
  }

  let content: React.ReactNode = null;
  switch (page.type) {
    case 'deployments-list':
      content = <DeploymentsPage />;
      break;
    case 'deployments-create':
      content = <DeploymentCreatePage />;
      break;
    case 'services-list':
      content = <ServicesPage />;
      break;
    case 'pods-list':
      content = <PodsPage />;
      break;
    // TODO: 他の画面（作成・詳細）も同様に分岐を追加
    default:
      content = <div>Not implemented</div>;
  }

  return <AppPage content={content} handleReissueToken={handleReissueClick} />
}
