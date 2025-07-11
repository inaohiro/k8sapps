import { useSetAtom } from "jotai";
import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { pageAtom, setPageAtom } from "../../store/store";

interface Service {
  id: string;
  name: string;
  type: string;
  clusterIP: string;
  ports: ServicePort[];
}
type ServicePort = {
  port: string;
  targetPort: string;
  protocol: null | string;
};

export function ServiceList() {
  const [services, setServices] = useState<Service[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();
  const setPage = useSetAtom(setPageAtom);

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
          setPage({ type: "services-create" });
        }}
        style={{ marginBottom: "1em" }}
      >
        + 新規 Service 作成
      </button>

      <div className="relative overflow-x-auto">
        <table className="table-auto">
          <thead>
            <tr>
              <th>Name</th>
              <th>Type</th>
              <th>IP</th>
              <th>Ports</th>
              <th>Action</th>
            </tr>
          </thead>
          <tbody>
            {services.map((svc) => (
              <tr key={svc.id}>
                <th scope="row">
                  <a
                    href="#"
                    onClick={(e) => {
                      e.preventDefault();
                      setPage({ type: "services-detail", id: svc.name });
                    }}
                    style={{ fontWeight: "bold", cursor: "pointer", color: "#1976d2", textDecoration: "underline" }}
                  >
                    {svc.name}
                  </a>
                </th>
                <td>{svc.type}</td>
                <td>{svc.clusterIP}</td>
                <td>
                  {svc.ports.map((port) => (
                    <>
                      <span>
                        Port: {port.port}, TargetPort: {port.targetPort}, Protocol: {port.protocol ?? "tcp"}
                      </span>
                      <br />
                    </>
                  ))}
                </td>
                <td>
                  <button
                    style={{
                      marginLeft: 8,
                      color: "#fff",
                      background: "#d32f2f",
                      border: "none",
                      borderRadius: 4,
                      padding: "2px 8px",
                      cursor: "pointer",
                    }}
                    onClick={async (e) => {
                      e.stopPropagation();
                      setLoading(true);
                      setError(null);
                      try {
                        const res = await fetch(`/api/services/${svc.name}`, {
                          method: "DELETE",
                          headers: { Authorization: `Bearer ${token}` },
                        });
                        if (!res.ok) throw new Error("Failed to delete service");
                      } catch (err: any) {
                        setError(err.message);
                      } finally {
                        setLoading(false);
                      }
                    }}
                  >
                    削除
                  </button>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
