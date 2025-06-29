import React, { useState } from 'react';

const DeploymentCreatePage: React.FC = () => {
  const [name, setName] = useState('');
  const [image, setImage] = useState('');
  const [status, setStatus] = useState('');
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError(null);
    setSuccess(false);
    try {
      const apiUrl = import.meta.env.VITE_API_URL;
      const res = await fetch(`${apiUrl}/deployments`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name, image, status }),
      });
      if (!res.ok) throw new Error('Failed to create deployment');
      setSuccess(true);
      setName('');
      setImage('');
      setStatus('');
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div>
      <h1>新規 Deployment 作成</h1>
      <form onSubmit={handleSubmit} style={{ maxWidth: 400 }}>
        <div style={{ marginBottom: 12 }}>
          <label>名前<br />
            <input value={name} onChange={e => setName(e.target.value)} required />
          </label>
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>イメージ<br />
            <input value={image} onChange={e => setImage(e.target.value)} required />
          </label>
        </div>
        <div style={{ marginBottom: 12 }}>
          <label>ステータス<br />
            <input value={status} onChange={e => setStatus(e.target.value)} />
          </label>
        </div>
        <button type="submit">作成</button>
      </form>
      {error && <div style={{ color: 'red' }}>Error: {error}</div>}
      {success && <div style={{ color: 'green' }}>Deployment を作成しました</div>}
    </div>
  );
};

export default DeploymentCreatePage;
