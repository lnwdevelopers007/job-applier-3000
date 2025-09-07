<div class="p-2 bg-slate-50">
  <h1 class="text-2xl font-bold text-gray-900">
    My Applications
  </h1>
  <p class="my-2 text-base text-gray-600">
    Track and manage all your job applications in one place
  </p>
  <div class="mt-6 grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Total Applications</h2>
      <p class="text-gray-800 text-xl font-bold">24</p>
      <p class="text-green-500 text-sm">↑ 12% from last month</p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">In Review</h2>
      <p class="text-gray-800 text-xl font-bold">8</p>
      <p class="text-green-500 text-sm">↑ 2 new this week</p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Offer Received</h2>
      <p class="text-gray-800 text-xl font-bold">3</p>
      <p class="text-green-500 text-sm">↑ 2 pending response</p>
    </div>

    <div class="p-4 bg-white rounded-lg shadow">
      <h2 class="text-gray-600">Response Rate</h2>
      <p class="text-gray-800 text-xl font-bold">67%</p>
      <p class="text-green-500 text-sm">↑ Above average</p>
    </div>
  </div>
  <div class="flex justify-between py-4">
    <div class="flex gap-2">
      <button class="px-3 py-1 bg-white rounded-full text-sm border-1 border-gray-300">All</button>
      <button class="px-3 py-1 bg-white rounded-full text-sm border-1 border-gray-300">Applied</button>
      <button class="px-3 py-1 bg-white rounded-full text-sm border-1 border-gray-300">In Review</button>
      <button class="px-3 py-1 bg-white rounded-full text-sm border-1 border-gray-300">Accepted</button>
      <button class="px-3 py-1 bg-white rounded-full text-sm border-1 border-gray-300">Rejected</button>
    </div>
    <button class="px-3 py-1 bg-white rounded-lg font-semibold text-sm flex item-center gap-1 border-1 border-gray-300 "><Download class="w-5 h-5" /> Export</button>
  </div>

  <div class="p-4 w-full">
    <div class="bg-white shadow-md rounded-lg overflow-hidden border border-gray-300">
      <div class="p-3 bg-gray-100 border-b border-gray-300">
        <h1 class=" px-2 font-bold">Recent Applications</h1>
      </div>

      {#each applications as app, i (i)}
        <div class={`flex items-center justify-between p-4 border-b border-gray-300 ${i === applications.length - 1 ? 'border-b-0' : ''}`}>
          <!-- Left part block -->
          <div class="flex items-center gap-4">
            <img src={app.companyLogo} alt="Company Logo" class="w-10 h-10 rounded-full object-cover" />
            <div class="flex flex-col">
              <span class="font-semibold">{app.jobTitle}</span>
              <div class="flex items-center gap-1.5 text-gray-500 text-sm">
                <span>{app.companyName}</span>
                <MapPin class="h-4 w-4 pl-1 text-black" />
                <span>{app.location}</span>
              </div>
            </div>
          </div>

          <!-- Right part block -->
          <div class="flex items-center gap-4">
            <div class="flex flex-col items-end text-right">
              <span class="text-xs font-bold text-gray-600">APPLIED</span>
              <span class="text-sm text-gray-500">{app.daysAgo} days ago</span>
            </div>

            <span
              class={`badge text-xs font-semibold p-2 rounded-full border ${
                app.status === 'IN REVIEW'
                  ? 'bg-yellow-100 text-yellow-800 border-yellow-300'
                  : app.status === 'ACCEPTED'
                  ? 'bg-green-100 text-green-800 border-green-300'
                  : app.status === 'APPLIED'
                  ? 'bg-blue-100 text-blue-800 border-blue-300'
                  : app.status === 'REJECTED'
                  ? 'bg-red-100 text-red-800 border-red-300'
                  : 'badge-outline'
              }`}
            >
              {app.status}
            </span>

            <div class="flex items-center gap-2">
              <button class="btn btn-sm border border-gray-300 p-1 rounded-lg btn-ghost">
                <Eye class="w-5 h-5" />
              </button>
              <button class="btn btn-sm border border-gray-300 p-1 rounded-lg btn-ghost">
                <EllipsisVertical class="w-5 h-5" />
              </button>
            </div>
          </div>
        </div>
      {/each}
    </div>
  </div>

  <div class="mt-8">
    <h2 class="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h2>
    <div class="relative">
      <div
        class="absolute left-4 bg-gray-300 w-px"
        style="
          top: calc(2rem); 
          bottom: calc(2rem);
        "
      ></div>
      {#each activities as activity, i (i)}
        <div class="relative mb-8 last:mb-0">
          <span
            class="absolute left-4 top-1/2 w-5 h-5 rounded-full bg-white border-4 border-green-500 -translate-x-1/2 -translate-y-1/2"
          ></span>

          <div class="ml-10 bg-white shadow p-3 rounded-lg">
            <p class="text-xs text-gray-500 uppercase tracking-wide">{activity.time}</p>
            <p class="text-sm font-semibold text-gray-800">{activity.title}</p>
            <p class="text-sm text-gray-600">{activity.description}</p>
          </div>
        </div>
      {/each}
    </div>
  </div>

</div>
<script>
  import { Download, MapPin, Eye, EllipsisVertical } from 'lucide-svelte';
  let activities = [
    {
      time: "Today, 10:30 AM",
      title: "Interview Scheduled",
      description:
        "Video interview with Agoda for Software Developer position scheduled for Dec 28, 2:00 PM",
    },
    {
      time: "Yesterday, 3:45 PM",
      title: "Application Viewed",
      description:
        "Your application for Google Software Engineer position was viewed by the hiring team",
    },
    {
      time: "Dec 20, 11:00 AM",
      title: "Offer Received",
      description:
        "Microsoft extended an internship offer for Data Engineer position",
    },
  ];
  let applications = [
    {
      id: 1,
      companyLogo: 'https://logo.clearbit.com/google.com',
      jobTitle: 'Frontend Developer',
      companyName: 'Google',
      location: 'Bangkok, Thailand',
      daysAgo: 3,
      status: 'IN REVIEW',
    },
    {
      id: 2,
      companyLogo: 'https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRbx0J4R5cqy54hvPd9lNP8ywO5kGWl5JlO9A&s',
      jobTitle: 'Backend Engineer',
      companyName: 'Agoda',
      location: 'Beijing, China',
      daysAgo: 5,
      status: 'ACCEPTED',
    },
    {
      id: 3,
      companyLogo: 'https://logo.clearbit.com/kasikornbank.com',
      jobTitle: 'UI/UX Designer',
      companyName: 'Kasikorn Bank',
      location: 'Chiang Mai, Thailand',
      daysAgo: 1,
      status: 'REJECTED',
    },
  ];
</script>