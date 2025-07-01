import { useCallback, useState } from "react";

function getCookie(name: string): string | null {
  return (
    document.cookie.split("; ").reduce((r, v) => {
      const parts = v.split("=");
      return parts[0] === name ? decodeURIComponent(parts[1]) : r;
    }, "") || null
  );
}

function setCookie(name: string, value: string) {
  document.cookie = `${name}=${encodeURIComponent(value)}; path=/`;
}
function removeCookie(name: string) {
  document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 GMT; path=/`;
}

export function useIssueToken() {
  const [item, setItem] = useState(() => {
    return getCookie("token");
  });
  const removeItem = useCallback(() => {
    setItem(null);
    removeCookie("token");
  }, []);
  const updateItem = useCallback((value: string) => {
    setItem(value);
    setCookie("token", value);
  }, []);

  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const issueToken = async (namespace: string): Promise<string | null> => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch("api/tokens", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ namespace }),
      });
      if (!res.ok) throw new Error("トークン発行に失敗しました");
      const data = await res.json();
      //   document.cookie = `token=${data.token}; path=/`;
      updateItem(data.token);
      setLoading(false);
      return data.token;
    } catch (err: any) {
      setError(err.message);
      setLoading(false);
      return null;
    }
  };

  return { token: item, issueToken, loading, error, logout: removeItem };
}
