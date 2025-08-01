import { useState } from "react";

export function GlobalHeader({ issueToken }: { issueToken: (namespace: string) => void }) {
  const [namespace, setNamespace] = useState("");

  return (
    <div>
      <div className="flex justify-end p-4 bg-white shadow gap-2">
        <input
          type="text"
          placeholder="namespace"
          value={namespace}
          onChange={(e) => setNamespace(e.target.value)}
          className="px-2 py-1 border rounded"
          style={{ minWidth: 120 }}
        />
        <button
          onClick={() => issueToken(namespace)}
          className="px-4 py-2 bg-blue-600 text-white rounded hover:bg-blue-700 transition"
        >
          トークン再発行
        </button>
      </div>
    </div>
  );
}
