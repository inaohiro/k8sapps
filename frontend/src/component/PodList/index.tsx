import { useSetAtom } from "jotai";
import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { pageAtom, setPageAtom } from "../../store/store";

interface Pod {
  id: string;
  name: string;
  status: string;
  image: string;
  created_at: string;
}

export function PodList() {
  const [pods, setPods] = useState<Pod[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();
  const setPage = useSetAtom(setPageAtom);

  useEffect(() => {
    fetch(`/api/pods`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch pods");
        return res.json();
      })
      .then((data) => {
        setPods(data);
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
      <h1>Pods</h1>
      <button
        onClick={() => {
          setPage({ type: "pods-create" });
        }}
        style={{ marginBottom: "1em" }}
      >
        + 新規 Pod 作成
      </button>
      <ul>
        {pods.map((pod) => (
          <li key={pod.id} style={{ display: "flex", alignItems: "center", gap: 8 }}>
            <a
              href="#"
              onClick={(e) => {
                e.preventDefault();
                setPage({ type: "pods-detail", id: pod.name });
              }}
              style={{ fontWeight: "bold", cursor: "pointer", color: "#1976d2", textDecoration: "underline" }}
            >
              {pod.name}
            </a>
            <span>
              (Status: {pod.status}, Image: {pod.image})
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
                  const res = await fetch(`/api/pods/${pod.name}`, {
                    method: "DELETE",
                    headers: { Authorization: `Bearer ${token}` },
                  });
                  if (!res.ok) throw new Error("Failed to delete pod");
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
