import React, { useEffect, useState } from 'react';
import { useToken } from '../hooks/useToken';

interface Pod {
  id: string;
  name: string;
  status: string;
  image: string;
  created_at: string;
}

const PodsPage: React.FC = () => {
  const [pods, setPods] = useState<Pod[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();

  useEffect(() => {
    fetch(`/api/pods`, {
      headers: {
        "Authorization": `Bearer ${token}`
      }
    })
      .then((res) => {
        if (!res.ok) throw new Error('Failed to fetch pods');
        return res.json();
      })
      .then((data: Pod[]) => {
        setPods(data);
        setLoading(false);
      })
      .catch((err: Error) => {
        setError(err.message);
        setLoading(false);
      });
  }, [token]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Pods</h1>
      <ul>
        {pods.map((pod) => (
          <li key={pod.id}>
            <strong>{pod.name}</strong> (Status: {pod.status}, Image: {pod.image}, Created At: {pod.created_at})
          </li>
        ))}
      </ul>
    </div>
  );
};

export default PodsPage;
