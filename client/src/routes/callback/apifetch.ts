// utils/api.ts

async function apiFetch(url: string, options: RequestInit = {}) {
  let token = localStorage.getItem("access_token");

  // Try the request with the current token
  let res = await fetch(url, {
    ...options,
    headers: {
      ...(options.headers || {}),
      Authorization: `Bearer ${token}`,
    },
  });

  // If unauthorized, try refreshing the token
  if (res.status === 401) {
    const refreshRes = await fetch("http://localhost:8080/auth/refresh", {
      method: "POST",
      credentials: "include", // important if backend sets httpOnly cookies
    });

    if (!refreshRes.ok) {
      throw new Error("Failed to refresh token");
    }

    const { access_token } = await refreshRes.json();

    // Save new token
    localStorage.setItem("access_token", access_token);

    // Retry original request with new token
    res = await fetch(url, {
      ...options,
      headers: {
        ...(options.headers || {}),
        Authorization: `Bearer ${access_token}`,
      },
    });
  }

  return res;
}

export default apiFetch;
