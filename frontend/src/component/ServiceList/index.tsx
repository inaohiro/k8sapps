import { useSetAtom } from "jotai";
import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { pageAtom } from "../../store/store";

interface Service {
  id: string;
  name: string;
  type: string;
  clusterIP: string;
  ports: string;
}

export function ServiceList() {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();
  const setPageState = useSetAtom(pageAtom);

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
      .then((data) => {
        setServices(data);
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
      <h1>Services</h1>
      <button
        onClick={() => {
          setPageState({ type: "services-create" });
        }}
        style={{ marginBottom: "1em" }}
      >
        + 新規 Service 作成
      </button>
      <ul>
        {services.map((svc) => (
          <li key={svc.id}>
            <strong>{svc.name}</strong> (Type: {svc.type}, ClusterIP: {svc.clusterIP}, Ports: {svc.ports})
          </li>
        ))}
      </ul>
    </div>
  );
};
