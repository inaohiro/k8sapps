import React, { useEffect, useState } from 'react';

interface Deployment {
  id: string;
  name: string;
  image: string;
  status: string;
}

const DeploymentsPage: React.FC = () => {
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    fetch(`/api/deployments`)
      .then((res) => {
        if (!res.ok) throw new Error('Failed to fetch deployments');
        return res.json();
      })
      .then((data) => {
        setDeployments(data);
        setLoading(false);
      })
      .catch((err: Error) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Deployments</h1>
      <button
        onClick={() => {
          window.location.href = '/deployments/new';
        }}
        style={{ marginBottom: '1em' }}
      >
        + 新規 Deployment 作成
      </button>
      <ul>
        {deployments.map((dep) => (
          <li key={dep.id}>
            <strong>{dep.name}</strong> (Image: {dep.image}, Status: {dep.status})
          </li>
        ))}
      </ul>
    </div>
  );
};

export default DeploymentsPage;
