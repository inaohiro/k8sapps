import React from 'react';
import { atom, useAtom } from 'jotai';
import { useToken } from './hooks/useToken';
import TokenIssuePage from './TokenIssuePage';
import DeploymentsPage from './DeploymentsPage';
import ServicesPage from './ServicesPage';
import PodsPage from './PodsPage';
import DeploymentCreatePage from './DeploymentCreatePage';

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

  return (
    <div className="min-h-screen bg-gray-50">
      <div className="flex justify-end p-4 bg-white shadow">
        <button
          onClick={handleReissueClick}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
        >
          トークン再発行
        </button>
      </div>
      <div className="max-w-3xl mx-auto mt-8">
        {content}
      </div>
    </div>
  );
}
