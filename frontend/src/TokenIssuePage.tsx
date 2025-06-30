import React, { useState } from 'react';

interface TokenIssuePageProps {
  onTokenIssued: () => void;
}

const TokenIssuePage: React.FC<TokenIssuePageProps> = ({ onTokenIssued }) => {
  const [namespace, setNamespace] = useState('');
  const [error, setError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      const res = await fetch(`api/tokens`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ namespace }),
      });
      if (!res.ok) throw new Error('トークン発行に失敗しました');
      const data = await res.json();
      document.cookie = `token=${data.token}; path=/`;
      onTokenIssued();
    } catch (err: any) {
      setError(err.message);
    }
  };

  return (
    <div>
      <h2>トークン発行</h2>
      <form onSubmit={handleSubmit}>
        <label>
          Namespace:
          <input value={namespace} onChange={e => setNamespace(e.target.value)} required />
        </label>
        <button type="submit">発行</button>
      </form>
      {error && <div style={{color: 'red'}}>{error}</div>}
    </div>
  );
};

export default TokenIssuePage;
