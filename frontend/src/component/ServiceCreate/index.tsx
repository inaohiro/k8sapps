import React, { useEffect, useState } from "react";
import { useToken } from "../../hooks/useToken";

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

const defaultPort: ServicePort = { port: 80, targetPort: 80, protocol: "TCP" };

export function ServiceCreate() {
  const [name, setName] = useState("");
  const [type, setType] = useState("ClusterIP");
  const [ports, setPorts] = useState<ServicePort[]>([{ ...defaultPort }]);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);
  const { token } = useToken();

  const handlePortChange = (idx: number, field: keyof ServicePort, value: string | number) => {
    setPorts((prev) =>
      prev.map((p, i) =>
        i === idx ? { ...p, [field]: field === "port" || field === "targetPort" ? Number(value) : value } : p
      )
    );
  };

  const addPort = () => setPorts((prev) => [...prev, { ...defaultPort }]);
  const removePort = (idx: number) => setPorts((prev) => prev.filter((_, i) => i !== idx));

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(false);
    try {
      const res = await fetch(`/api/services`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({ name, type, ports }),
      });
      if (!res.ok) throw new Error("Failed to create service");
      setSuccess(true);
      setName("");
      setType("ClusterIP");
      setPorts([{ ...defaultPort }]);
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div>
      <h1>新規 Service 作成</h1>
      <form onSubmit={handleSubmit} style={{ maxWidth: 500 }}>
        <div style={{ marginBottom: 12 }}>
          <label>
            名前<br />
            <input value={name} onChange={(e) => setName(e.target.value)} required />
          </label>
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>
            タイプ<br />
            <select value={type} onChange={(e) => setType(e.target.value)} required>
              <option value="ClusterIP">ClusterIP</option>
              <option value="NodePort">NodePort</option>
              <option value="LoadBalancer">LoadBalancer</option>
            </select>
          </label>
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>ポート設定</label>
          {ports.map((port, idx) => (
            <div key={idx} style={{ border: "1px solid #ccc", padding: 8, marginBottom: 8 }}>
              <div>
                <label>
                  名前 (任意)
                  <input
                    value={port.name || ""}
                    onChange={(e) => handlePortChange(idx, "name", e.target.value)}
                    style={{ marginLeft: 8 }}
                  />
                </label>
              </div>
              <div>
                <label>
                  ポート
                  <input
                    type="number"
                    value={port.port}
                    onChange={(e) => handlePortChange(idx, "port", e.target.value)}
                    required
                    style={{ marginLeft: 8 }}
                  />
                </label>
              </div>
              <div>
                <label>
                  ターゲットポート
                  <input
                    type="number"
                    value={port.targetPort}
                    onChange={(e) => handlePortChange(idx, "targetPort", e.target.value)}
                    required
                    style={{ marginLeft: 8 }}
                  />
                </label>
              </div>
              <div>
                <label>
                  プロトコル
                  <select
                    value={port.protocol}
                    onChange={(e) => handlePortChange(idx, "protocol", e.target.value)}
                    required
                    style={{ marginLeft: 8 }}
                  >
                    <option value="TCP">TCP</option>
                    <option value="UDP">UDP</option>
                  </select>
                </label>
              </div>
              {ports.length > 1 && (
                <button type="button" onClick={() => removePort(idx)} style={{ marginTop: 8 }}>
                  削除
                </button>
              )}
            </div>
          ))}
          <button type="button" onClick={addPort} style={{ marginTop: 8 }}>
            ポート追加
          </button>
        </div>
        <button type="submit">作成</button>
      </form>
      {error && <div style={{ color: "red" }}>Error: {error}</div>}
      {success && <div style={{ color: "green" }}>Service を作成しました</div>}
    </div>
  );
};
