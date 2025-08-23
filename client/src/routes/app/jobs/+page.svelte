<div class="flex h-screen">
  <main class="flex-1 p-0.5 space-y-6 w-full min-w-0">
    <div class="flex items-center gap-3">
      <!-- Search bar -->
      <div class="flex-1 flex items-center bg-white rounded-lg shadow px-3 py-1">
        <Search class="w-5 h-5 text-gray-500" />
        <input
          type="text"
          placeholder="Search jobs..."
          class="flex-1 ml-2 outline-none border-none"
        />
      </div>
      <div class="flex gap-2">
        <button class="px-3 py-1 bg-white rounded-full text-sm">Full-time</button>
        <button class="px-3 py-1 bg-white rounded-full text-sm">1 day ago</button>
        <button class="px-3 py-1 bg-white rounded-full text-sm">AI</button>
        <button class="px-3 py-1 bg-white rounded-full text-sm">Remote</button>
        <button class="px-3 py-1 bg-white rounded-full text-sm flex items-center"><Filter class="h-4 w-4 mr-1" /> Filters</button>
      </div>
    </div>

    <div class="grid grid-cols-2 gap-6 w-full">
      <!--Jobs List -->
      <section class="col-span-1 flex flex-col space-y-4">
        <div class="flex-1 overflow-y-auto w-full space-y-3 p-2">
          {#each jobs as job}
          <button
            on:click={() => (selectedJob = job)}
            class="w-full flex flex-col cursor-pointer bg-white rounded-lg shadow ring-offset-2 hover:ring-2 hover:ring-green-300 {selectedJob.id === job.id ? 'ring-2 ring-green-500' : ''}"
          >
            <div class="flex items-start gap-3 p-2">
              <div class="mt-2">
              <img src={job.logo} alt={job.company} class="w-12 h-12 rounded-full object-cover flex-shrink-0" />
              </div>
              <div class="flex flex-col text-left flex-1">
                <div class="font-semibold">{job.title}</div>
                <div class="text-sm text-gray-600">{job.company}</div>
                <div class="text-sm text-black flex items-center gap-2 mt-1">
                  <MapPin class="w-4 h-4 mt-0.5" />
                  {job.location}
                </div>
              </div>
            </div>
            <div class="flex flex-wrap gap-2 mt-2 p-2">
              {#each job.tags as tag}
                <span class="px-2 py-1 bg-gray-100 rounded-full text-sm">{tag}</span>
              {/each}
            </div>
            <div class="text-xs text-left mx-3 pb-2 text-gray-500 mt-1">{job.posted}</div>
          </button>

          {/each}
        </div>
      </section>

      <!--Job Detail -->
      <section class="col-span-1 bg-white p-6 rounded-lg shadow space-y-2">
        {#if selectedJob}
          <div class="mb-2">
            <div class="flex items-start gap-3 px-2">
              <img src={selectedJob.logo} alt={selectedJob.company} class="w-4 h-4 rounded-full object-cover flex-shrink-0 mt-0.5" />
              <p class="text-sm font-semibold text-gray-600">{selectedJob.company}</p>
            </div>
            <h2 class="text-xl font-bold">{selectedJob.title}</h2>
            <div class="text-sm text-black flex items-center gap-2 my-2">
                <MapPin class="w-4 h-4 mt-0.5" />
                {selectedJob.location}
            </div>
            <p class="text-sm text-gray-500 my-2">{selectedJob.posted}</p>
          </div>
          <div class="flex gap-2">
            {#each selectedJob.tags as tag}
              <span class="px-2 py-1 bg-gray-100 rounded-full text-sm">{tag}</span>
            {/each}
          </div>
          <div class="flex gap-2 mb-4">
            <button class="px-4 py-2 bg-green-600 text-white text-sm rounded-lg">Apply</button>
            <button class="px-4 py-2 bg-yellow-400 text-sm rounded-lg">Bookmark</button>
          </div>
          <h3 class="font-semibold text-lg">Job Description</h3>
          <div class="space-y-4">
            <p>{selectedJob.description}</p>
          </div>
        {/if}
      </section>
    </div>
  </main>
</div>

<script>
  import { Search, Filter, MapPin } from 'lucide-svelte';

  let jobs = [
    {
      id: 1,
      title: "Software Engineer, AI Acceleration",
      company: "Google",
      location: "Singapore",
      type: "Full-time",
      tags: ["Remote", "AI"],
      posted: "1 day ago",
      description: "Google wants to have more software engineer for researching and improving AI abilities.",
      logo: "https://logo.clearbit.com/google.com",
    },
    {
      id: 2,
      title: "Software Developer, Back End",
      company: "Agoda",
      location: "Bangkok, Thailand",
      type: "Full-time",
      tags: ["Backend", "Node.js"],
      posted: "2 days ago",
      description: "Backend-focused developer role which is responsible for creating API for our applications and managing database.",
      logo: "https://encrypted-tbn0.gstatic.com/images?q=tbn:ANd9GcRbx0J4R5cqy54hvPd9lNP8ywO5kGWl5JlO9A&s",
    },
    {
      id: 3,
      title: "Fullstack Web Developer",
      company: "Kasikorn Bank",
      location: "Bangkok, Thailand",
      type: "Full-time",
      tags: ["React", "Node.js"],
      posted: "5 days ago",
      description: "Work across frontend and backend...",
      logo: "https://logo.clearbit.com/kasikornbank.com",
    },
  ];

  let selectedJob = jobs[0];
</script>