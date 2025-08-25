export const load = async ({ fetch }) => {
	const res = await fetch('/health');
	const text = await res.text();
	console.log('Raw response from backend:', text);
	try {
		const data = JSON.parse(text);
		return { health: data };
	} catch (err) {
		console.error('Failed to parse JSON:', err);
		return { health: false };
	}
};
