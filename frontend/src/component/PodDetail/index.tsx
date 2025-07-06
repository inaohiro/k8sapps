import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { useSetAtom } from "jotai";
import { setPageAtom } from "../../store/store";

interface Pod {
  id: string;
  name: string;
  status: string;
  image: string;
  created_at: string;
}

interface Props {
  id: string;
}

export function PodDetail({ id }: Props) {
  const [pod, setPod] = useState<Pod | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();
  const setPage = useSetAtom(setPageAtom);

  useEffect(() => {
    setLoading(true);
    setError(null);
    fetch(`/api/pods/${id}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch pod detail");
        return res.json();
      })
      .then((data) => {
        setPod(data);
        setLoading(false);
      })
      .catch((err: Error) => {
        setError(err.message);
        setLoading(false);
      });
  }, [id, token]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!pod) return <div>Not found</div>;

  return (
    <div>
      <button onClick={() => setPage({ type: "pods-list" })} style={{ marginBottom: 16 }}>
        ← 一覧に戻る
      </button>
      <h2>Pod詳細</h2>
      <table>
        <tbody>
          <tr>
            <th>ID</th>
            <td>{pod.id}</td>
          </tr>
          <tr>
            <th>Name</th>
            <td>{pod.name}</td>
          </tr>
          <tr>
            <th>Status</th>
            <td>{pod.status}</td>
          </tr>
          <tr>
            <th>Image</th>
            <td>{pod.image}</td>
          </tr>
          <tr>
            <th>Created At</th>
            <td>{pod.created_at}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
