import React, { useState } from "react";
import { useToken } from "../../hooks/useToken";

interface Pod {
  id: string;
  name: string;
  status: string;
  image: string;
  created_at: string;
}

export function PodCreate() {
  const [name, setName] = useState("");
  const [image, setImage] = useState("");
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const { token } = useToken();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(false);
    try {
      const res = await fetch(`/api/pods`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({ name, image }),
      });
      if (!res.ok) throw new Error("Failed to create pod");
      setSuccess(true);
      setName("");
      setImage("");
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div>
      <h1>新規 Pod 作成</h1>
      <form onSubmit={handleSubmit} style={{ maxWidth: 400 }}>
        <div style={{ marginBottom: 12 }}>
          <label>
            名前<br />
            <input value={name} onChange={(e) => setName(e.target.value)} required />
          </label>
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>
            イメージ<br />
            <input value={image} onChange={(e) => setImage(e.target.value)} required />
          </label>
        </div>
        <button type="submit">作成</button>
      </form>
      {error && <div style={{ color: "red" }}>Error: {error}</div>}
      {success && <div style={{ color: "green" }}>Pod を作成しました</div>}
    </div>
  );
};
