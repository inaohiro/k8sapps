import React, { useState } from "react";

export function TokenIssuePage({
  issueToken,
  loading,
  onTokenIssued,
}: {
  issueToken: (namespace: string) => Promise<string | null>;
  loading: boolean;
  onTokenIssued: () => void;
}) {
  const [namespace, setNamespace] = useState("");
  const [err, setError] = useState<string | null>(null);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const token = await issueToken(namespace);
    if (token) {
      onTokenIssued();
    } else {
      setError("失敗しました");
    }
  };

  if (err) return <div>{err}</div>;

  return (
    <>
      <h2>Namespace を指定してください</h2>
      <form onSubmit={handleSubmit}>
        <label>
          Namespace:
          <input value={namespace} onChange={(e) => setNamespace(e.target.value)} required />
        </label>
        <button type="submit" disabled={loading}>
          発行
        </button>
      </form>
    </>
  );
}
