export async function apiFetch(url: string, options: RequestInit = {}) {
	// Since we're using httpOnly cookies, no need to manually handle tokens
	// The browser will automatically include cookies with requests
	let res = await fetch(url, {
		...options,
		credentials: 'include', // Ensure cookies are included in the request
		headers: {
			...(options.headers || {})
		}
	});

	// If unauthorized, the server will automatically try to refresh via cookies
	// If that fails, the user will be redirected to login by hooks.server.ts
	return res;
}
