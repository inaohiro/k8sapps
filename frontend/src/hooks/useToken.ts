import { useState, useCallback } from 'react';

function getCookie(name: string): string | null {
  const match = document.cookie.match(new RegExp('(^| )' + name + '=([^;]+)'));
  return match ? decodeURIComponent(match[2]) : null;
}

export function useToken() {
  const token = getCookie('token');
  const setToken = (token: string) => {
    document.cookie = `token=${token}; path=/`;
  }

  return { token, setToken } as const;
}
