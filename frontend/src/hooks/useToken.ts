import { useState, useCallback } from 'react';

function getCookie(name: string): string | null {
  const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
  return match ? decodeURIComponent(match[2]) : null;
}

export function useToken() {
  const [token, setToken] = useState(() => getCookie('token'));

  // cookieのtokenを再取得するための関数
  const refresh = useCallback(() => {
    setToken(getCookie('token'));
  }, []);

  return { token, hasToken: !!token, refresh };
}
