import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";
import { useSetAtom } from "jotai";
import { setPageAtom } from "../../store/store";

interface Deployment {
  id: string;
  name: string;
  image: string;
  status: string;
  created_at: string;
}

interface Props {
  id: string;
}

export function DeploymentDetail({ id }: Props) {
  const [deployment, setDeployment] = useState<Deployment | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const { token } = useToken();
  const setPage = useSetAtom(setPageAtom);

  useEffect(() => {
    setLoading(true);
    setError(null);
    fetch(`/api/deployments/${id}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((res) => {
        if (!res.ok) throw new Error("Failed to fetch deployment detail");
        return res.json();
      })
      .then((data) => {
        setDeployment(data);
        setLoading(false);
      })
      .catch((err: Error) => {
        setError(err.message);
        setLoading(false);
      });
  }, [id, token]);

  if (loading) return <div>Loading...</div>;
  if (error) return <div>Error: {error}</div>;
  if (!deployment) return <div>Not found</div>;

  return (
    <div>
      <button onClick={() => setPage({ type: "deployments-list" })} style={{ marginBottom: 16 }}>
        ← 一覧に戻る
      </button>
      <h2>Deployment詳細</h2>
      <table>
        <tbody>
          <tr>
            <th>ID</th>
            <td>{deployment.id}</td>
          </tr>
          <tr>
            <th>Name</th>
            <td>{deployment.name}</td>
          </tr>
          <tr>
            <th>Image</th>
            <td>{deployment.image}</td>
          </tr>
          <tr>
            <th>Status</th>
            <td>{deployment.status}</td>
          </tr>
          <tr>
            <th>Created At</th>
            <td>{deployment.created_at}</td>
          </tr>
        </tbody>
      </table>
    </div>
  );
}
