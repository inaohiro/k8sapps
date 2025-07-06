import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { useSetAtom } from "jotai";
import { setPageAtom } from "../../store/store";

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

interface Props {
  id: string;
}

export function ServiceDetail({ id }: Props) {
  const [service, setService] = useState<Service | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();
  const setPage = useSetAtom(setPageAtom);

  useEffect(() => {
    setLoading(true);
    setError(null);
    fetch(`/api/services/${id}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch service detail");
        return res.json();
      })
      .then((data) => {
        setService(data);
        setLoading(false);
      })
      .catch((err: Error) => {
        setError(err.message);
        setLoading(false);
      });
  }, [id, token]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!service) return <div>Not found</div>;

  return (
    <div>
      <button onClick={() => setPage({ type: "services-list" })} style={{ marginBottom: 16 }}>
        ← 一覧に戻る
      </button>
      <h2>Service詳細</h2>
      <table>
        <tbody>
          <tr>
            <th>ID</th>
            <td>{service.id}</td>
          </tr>
          <tr>
            <th>Name</th>
            <td>{service.name}</td>
          </tr>
          <tr>
            <th>Type</th>
            <td>{service.type}</td>
          </tr>
          <tr>
            <th>Cluster IP</th>
            <td>{service.clusterIP}</td>
          </tr>
          <tr>
            <th>Created At</th>
            <td>{service.created_at}</td>
          </tr>
          <tr>
            <th>Ports</th>
            <td>
              <table style={{ borderCollapse: "collapse", width: "100%" }}>
                <thead>
                  <tr>
                    <th>Name</th>
                    <th>Port</th>
                    <th>Target Port</th>
                    <th>Protocol</th>
                  </tr>
                </thead>
                <tbody>
                  {service.ports.map((port, i) => (
                    <tr key={i}>
                      <td>{port.name || "-"}</td>
                      <td>{port.port}</td>
                      <td>{port.targetPort}</td>
                      <td>{port.protocol}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
