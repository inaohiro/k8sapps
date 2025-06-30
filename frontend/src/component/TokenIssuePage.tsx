import React, { useState } from 'react';
import { useIssueToken } from '../hooks/useIssueToken';

interface TokenIssuePageProps {
  onTokenIssued: () => void;
}

const TokenIssuePage: React.FC<TokenIssuePageProps> = ({ onTokenIssued }) => {
  const [namespace, setNamespace] = useState('');
  const { issueToken, loading, error } = useIssueToken();
  const [localError, setLocalError] = useState('');

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLocalError('');
    const token = await issueToken(namespace);
    if (token) {
      onTokenIssued();
    } else {
      setLocalError(error || '');
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
        <button type="submit" disabled={loading}>発行</button>
      </form>
      {(localError || error) && <div style={{color: 'red'}}>{localError || error}</div>}
    </div>
  );
};

export default TokenIssuePage;
