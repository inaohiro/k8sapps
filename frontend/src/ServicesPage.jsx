import React, { useEffect, useState } from 'react';

const ServicesPage = () => {
  const [services, setServices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

  useEffect(() => {
    const apiUrl = import.meta.env.VITE_API_URL;
    fetch(`${apiUrl}/services`)
      .then((res) => {
        if (!res.ok) throw new Error('Failed to fetch services');
        return res.json();
      })
      .then((data) => {
        setServices(data);
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
      <h1>Services</h1>
      <ul>
        {services.map((svc) => (
          <li key={svc.id}>
            <strong>{svc.name}</strong> (Type: {svc.type}, ClusterIP: {svc.clusterIP})
          </li>
        ))}
      </ul>
    </div>
  );
};

export default ServicesPage;
