import { writable } from 'svelte/store';

interface JobSearchState {
	query: string;
	shouldFetch: boolean;
}

function createJobSearchStore() {
	const { subscribe, set, update } = writable<JobSearchState>({
		query: '',
		shouldFetch: false
	});

	return {
		subscribe,
		setSearchQuery: (query: string) => {
			set({ query, shouldFetch: true });
		},
		clearFetchFlag: () => {
			update(state => ({ ...state, shouldFetch: false }));
		},
		reset: () => {
			set({ query: '', shouldFetch: false });
		}
	};
}

export const jobSearchStore = createJobSearchStore();