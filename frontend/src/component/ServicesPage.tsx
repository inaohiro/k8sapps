import React, { useEffect, useState } from "react";
import { useToken } from "../hooks/useToken";

interface ServicePort {
  name?: string;
  port: number;
  targetPort: number;
  protocol: string;
}

interface Service {
  id: string;
  name: string;
  type: string;
  clusterIP: string;
  ports: ServicePort[];
  created_at: string;
}

const ServicesPage: React.FC = () => {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();

  useEffect(() => {
    fetch(`/api/services`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch services");
        return res.json();
      })
      .then((data: Service[]) => {
        setServices(data);
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
