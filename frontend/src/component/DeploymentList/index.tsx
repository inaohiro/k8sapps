import { useSetAtom } from "jotai";
import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { pageAtom, setPageAtom } from "../../store/store";

interface Deployment {
  id: string;
  name: string;
  image: string;
  status: string;
}

export function DeploymentList() {
  const [deployments, setDeployments] = useState<Deployment[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();
  const setPage = useSetAtom(setPageAtom);

  useEffect(() => {
    fetch(`/api/deployments`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch deployments");
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
          setPage({ type: "deployments-create" });
        }}
        style={{ marginBottom: "1em" }}
      >
        + 新規 Deployment 作成
      </button>
      <ul>
        {deployments.map((dep) => (
          <li key={dep.id} style={{ display: "flex", alignItems: "center", gap: 8 }}>
            <a
              href="#"
              onClick={e => {
                e.preventDefault();
                setPage({ type: "deployments-detail", id: dep.id });
              }}
              style={{ fontWeight: "bold", cursor: "pointer", color: "#1976d2", textDecoration: "underline" }}
            >
              {dep.name}
            </a>
            <span>(Image: {dep.image}, Status: {dep.status})</span>
            <button
              style={{ marginLeft: 8, color: "#fff", background: "#d32f2f", border: "none", borderRadius: 4, padding: "2px 8px", cursor: "pointer" }}
              onClick={async (e) => {
                e.stopPropagation();
                setLoading(true);
                setError(null);
                try {
                  const res = await fetch(`/api/deployments/${dep.id}`, {
                    method: "DELETE",
                    headers: { Authorization: `Bearer ${token}` },
                  });
                  if (!res.ok) throw new Error("Failed to delete deployment");
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
};
