import { useSetAtom } from "jotai";
import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { pageAtom, setPageAtom } from "../../store/store";

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
      <ul>
        {services.map((svc) => (
          <li key={svc.id} style={{ display: "flex", alignItems: "center", gap: 8 }}>
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
            <span>
              (Type: {svc.type}, ClusterIP: {svc.clusterIP}, Ports: {svc.ports})
            </span>
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
          </li>
        ))}
      </ul>
    </div>
  );
}
