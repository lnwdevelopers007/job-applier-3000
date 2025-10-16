export async function apiFetch(url: string, options: RequestInit = {}) {
	const res = await fetch(url, {
		...options,
		credentials: 'include',
		headers: {
			...(options.headers || {})
		}
	});

	return res;
}
