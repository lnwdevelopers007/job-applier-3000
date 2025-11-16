export interface JobseekerStats {
	totalApplications: number;
	inReview: number;
	offerReceived: number;
	responseRate: number; // percentage
	trend: {
		totalApplications: number;
		inReview: number;
		offerReceived: number;
		responseRate: number;
	};
}

function emptyStats(): JobseekerStats {
	return {
		totalApplications: 0,
		inReview: 0,
		offerReceived: 0,
		responseRate: 0,
		trend: {
			totalApplications: 0,
			inReview: 0,
			offerReceived: 0,
			responseRate: 0
		}
	};
}

export async function getJobseekerStats(userID: string): Promise<JobseekerStats> {
	// Remove authentication check - let the backend handle it
	if (!userID) {
		console.error('getJobseekerStats: Invalid userID');
		throw new Error('Invalid userID');
	}

	try {
		const res = await fetch(`/apply/query?applicantID=${userID}`, {
			credentials: 'include'
		});
		
		if (!res.ok) {
			// Handle authentication errors specifically
			if (res.status === 401) {
				console.error('Not authenticated - redirecting to login');
				// Optionally redirect to login page
				// window.location.href = '/auth/google';
				throw new Error('Not authenticated');
			}
			
			console.error(`Failed to fetch applications: ${res.status} ${res.statusText}`);
			const errorText = await res.text();
			console.error('Error response:', errorText);
			throw new Error(`Failed to fetch applications: ${res.status}`);
		}

		const applications = await res.json();
		
		if (!Array.isArray(applications)) {
			console.error('Expected array but got:', typeof applications);
			return emptyStats();
		}
		
		if (applications.length === 0) {
			return emptyStats();
		}

		const now = new Date();

		// Time boundaries
		const lastMonthStart = new Date(now.getFullYear(), now.getMonth() - 1, 1);
		const lastMonthEnd = new Date(now.getFullYear(), now.getMonth(), 0);

		const weekAgo = new Date();
		weekAgo.setDate(now.getDate() - 7);

		// Counters
		let totalApplications = 0;
		let inReview = 0;
		let offerReceived = 0;
		let responded = 0;

		// Trend counters
		let totalApplicationsPrev = 0;
		let inReviewPrev = 0;
		let offerReceivedPrev = 0;
		let respondedPrev = 0;

		for (const app of applications) {
			const createdAt = new Date(app.jobApplication?.createdAt || app.createdAt);
			const status = (app.jobApplication?.status || app.status || 'pending').toLowerCase();

			totalApplications++;
			if (status === 'pending') inReview++;
			if (status === 'accepted') offerReceived++;
			if (status === 'accepted' || status === 'rejected') responded++;

			// Trends
			if (createdAt >= lastMonthStart && createdAt <= lastMonthEnd) {
				totalApplicationsPrev++;
			}

			if (createdAt >= weekAgo && createdAt < now) {
				if (status === 'pending') inReviewPrev++;
				if (status === 'accepted') offerReceivedPrev++;
				if (status === 'accepted' || status === 'rejected') respondedPrev++;
			}
		}

		const responseRate =
			totalApplications > 0 ? Math.round((responded / totalApplications) * 100) : 0;
		const responseRatePrev =
			totalApplicationsPrev > 0 ? Math.round((respondedPrev / totalApplicationsPrev) * 100) : 0;

		return {
			totalApplications,
			inReview,
			offerReceived,
			responseRate,
			trend: {
				totalApplications: totalApplications - totalApplicationsPrev,
				inReview: inReview - inReviewPrev,
				offerReceived: offerReceived - offerReceivedPrev,
				responseRate: responseRate - responseRatePrev
			}
		};
	} catch (err) {
		console.error('Error computing jobseeker stats:', err);
		return emptyStats();
	}
}
