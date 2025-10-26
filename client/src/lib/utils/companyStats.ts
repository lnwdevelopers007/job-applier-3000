import { isAuthenticated } from '$lib/utils/auth';

interface CompanyStats {
	activeJobs: number;
	totalApplicants: number;
	pendingReview: number;
	offersAccepted: number;
	trend: {
		activeJobs: number;
		totalApplicants: number;
		pendingReview: number;
		offersAccepted: number;
	};
}

export async function getCompanyAnalytics(companyID: string): Promise<CompanyStats> {
	if (!isAuthenticated()) throw new Error('Not authenticated');
	if (!companyID) throw new Error('Invalid companyID');

	const jobsRes = await fetch(`/jobs/query?companyID=${companyID}`);
	if (!jobsRes.ok) throw new Error('Failed to fetch company jobs');
	const jobs = await jobsRes.json();

	if (!Array.isArray(jobs) || jobs.length === 0) {
		return emptyStats();
	}

	const now = new Date();
	const thisMonthStart = new Date(now.getFullYear(), now.getMonth(), 1);
	const prevMonthStart = new Date(now.getFullYear(), now.getMonth() - 1, 1);
	const prevMonthEnd = new Date(now.getFullYear(), now.getMonth(), 0);

	const todayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate());
	const yesterdayStart = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1);
	const yesterdayEnd = new Date(now.getFullYear(), now.getMonth(), now.getDate());

	// Counters
	const total = { activeJobs: 0, totalApplicants: 0, pendingReview: 0, offersAccepted: 0 };
	const currentPeriod = { activeJobs: 0, totalApplicants: 0, pendingReview: 0, offersAccepted: 0 };
	const previousPeriod = { activeJobs: 0, totalApplicants: 0, pendingReview: 0, offersAccepted: 0 };

	for (const job of jobs) {
		const openAt = new Date(job.postOpenDate || job.createdAt);
		const closeAt = job.applicationDeadline ? new Date(job.applicationDeadline) : null;

		// Active job
		if (openAt <= now && (!closeAt || closeAt >= now)) total.activeJobs++;

		// Previous month active
		if (openAt <= prevMonthEnd && (!closeAt || closeAt >= prevMonthStart))
			previousPeriod.activeJobs++;
		if (openAt <= now && (!closeAt || closeAt >= thisMonthStart)) currentPeriod.activeJobs++;

		const jobID = job.id || job._id;
		if (!jobID) continue;

		const applyRes = await fetch(`/apply?jobID=${jobID}`);
		if (!applyRes.ok) continue;
		const applications = await applyRes.json();
		if (!Array.isArray(applications)) continue;

		for (const app of applications) {
			const createdAt = new Date(app.jobApplication?.createdAt || app.createdAt);
			const status = app.jobApplication?.status?.toLowerCase();

			// Lifetime metrics
			total.totalApplicants++;
			if (status === 'pending') total.pendingReview++;
			if (status === 'accepted') total.offersAccepted++;

			// Monthly metrics
			if (createdAt >= thisMonthStart && createdAt <= now) {
				currentPeriod.totalApplicants++;
				if (status === 'accepted') currentPeriod.offersAccepted++;
			}
			if (createdAt >= prevMonthStart && createdAt <= prevMonthEnd) {
				previousPeriod.totalApplicants++;
				if (status === 'accepted') previousPeriod.offersAccepted++;
			}

			// Daily metrics
			if (status === 'pending') {
				if (createdAt >= todayStart && createdAt <= now) currentPeriod.pendingReview++;
				else if (createdAt >= yesterdayStart && createdAt < yesterdayEnd)
					previousPeriod.pendingReview++;
			}
		}
	}

	const trend = {
		activeJobs: currentPeriod.activeJobs - previousPeriod.activeJobs,
		totalApplicants: currentPeriod.totalApplicants - previousPeriod.totalApplicants,
		pendingReview: currentPeriod.pendingReview - previousPeriod.pendingReview,
		offersAccepted: currentPeriod.offersAccepted - previousPeriod.offersAccepted
	};

	return { ...total, trend };
}

function emptyStats(): CompanyStats {
	return {
		activeJobs: 0,
		totalApplicants: 0,
		pendingReview: 0,
		offersAccepted: 0,
		trend: { activeJobs: 0, totalApplicants: 0, pendingReview: 0, offersAccepted: 0 }
	};
}
