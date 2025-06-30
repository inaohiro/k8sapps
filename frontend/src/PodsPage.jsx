import React, { useEffect, useState } from 'react';

const PodsPage = () => {
  const [pods, setPods] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    fetch(`/api/pods`)
      .then((res) => {
        if (!res.ok) throw new Error('Failed to fetch pods');
        return res.json();
      })
      .then((data) => {
        setPods(data);
        setLoading(false);
      })
      .catch((err) => {
        setError(err.message);
        setLoading(false);
      });
  }, []);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;

  return (
    <div>
      <h1>Pods</h1>
      <ul>
        {pods.map((pod) => (
          <li key={pod.id}>
            <strong>{pod.name}</strong> (Status: {pod.status}, Node: {pod.node})
          </li>
        ))}
      </ul>
    </div>
  );
};

export default PodsPage;
