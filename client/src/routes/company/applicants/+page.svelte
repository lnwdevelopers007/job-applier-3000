<div class="p-2 bg-gray-100">
  <h1 class="text-2xl font-bold text-gray-900">
    All Applicants
  </h1>
  <p class="my-3 text-base text-gray-600">
    Manage candidates across all job postings
  </p>
      <div class="flex-1 flex items-center bg-white rounded-lg shadow px-3 py-1">
        <Search class="w-5 h-5 text-gray-500" />
        <input
          type="text"
          placeholder="Search candidates by name, email, skills or position..."
          class="flex-1 ml-2 outline-none border-none text-sm"
        />
      </div>

      <div class="flex gap-2 my-2">      
      <select class="px-7 py-1 bg-white border border-gray-200 rounded text-sm">
        <option>All Jobs</option>
        <option>AI Researcher</option>
        <option>Developer</option>
        <option>Tester</option>
      </select>

      <select class="px-7 py-1 bg-white border border-gray-200 rounded text-sm">
        <option>All Statuses</option>
        <option>Accepted</option>
        <option>Pending</option>
        <option>Rejected</option>
      </select>
      <div class="border-l border-gray-300 h-6 py-1"></div>
        <button class="px-3 py-1 bg-white border-1 border-gray-200 rounded-full text-sm">All</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">New (34)</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">High Match</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">Recent Grads</button>
        <button class="px-3 py-1 bg-white border-1 border-gray-200  rounded-full text-sm">Experienced</button>
      </div>

  <div class="grid grid-cols-3 gap-4 mt-4">
    <div class="col-span-1 bg-white p-4 rounded-lg shadow">
      <h2 class="text-lg font-semibold text-gray-800">All Candidates</h2>
      <div class="flex gap-2 my-2">
        <button class="px-2 py-1 bg-white border-1 border-gray-200 rounded-full text-xs">All</button>
        <button class="px-2 py-1 bg-white border-1 border-gray-200  rounded-full text-xs">New</button>
        <button class="px-2 py-1 bg-white border-1 border-gray-200  rounded-full text-xs">Shortlisted</button>
        <button class="px-2 py-1 bg-white border-1 border-gray-200  rounded-full text-xs">Rejected</button>
      </div>
    <div class="mt-3 space-y-3">
        {#each candidates as candidate}
          <button
            class="flex items-start gap-3 p-3 border rounded-lg shadow-sm hover:shadow-md transition cursor-pointer w-full
              {selectedCandidate?.name === candidate.name ? 'bg-green-100 border-green-600' : ''}"
            on:click={() => (selectedCandidate = candidate)}
          >
            <img src={candidate.avatar} alt={candidate.name} class="w-10 h-10 rounded-full object-cover" />
            <div class="flex-1 flex flex-col items-start">
              <h3 class="font-semibold text-gray-800 text-sm">{candidate.name}</h3>
              <p class="text-start text-xs text-gray-800">{candidate.role}</p>
              <p class="text-start text-[10px] text-gray-500">Applied for: {candidate.applied}</p>
              <p class="flex text-start text-[10px] text-gray-500 gap-0.5">
                <GraduationCap class="h-3 w-3 m-0.5" />{candidate.university}
                <Briefcase class="w-3 h-3 m-0.5" />{candidate.year}
                <Clock class="w-3 h-3 m-0.5" />{candidate.time}</p>
              <span
                class="inline-block mt-1 text-xs px-2 py-0.5 rounded-lg font-semibold
                  {candidate.status === 'New' ? 'bg-blue-100 text-blue-700 border-1 border-blue-700' : ''}
                  {candidate.status === 'Shortlisted' ? 'bg-purple-100 text-purple-700 border-1 border-purple-700' : ''}
                  {candidate.status === 'Rejected' ? 'bg-red-100 text-red-700 border-1 border-red-700' : ''}
                  {candidate.status === 'Reviewed' ? 'bg-yellow-100 text-yellow-700 border-1 border-yellow-700' : ''}"
              >
                {candidate.status}
              </span>
            </div>
          </button>
        {/each}
      </div>
    </div>

    <div class="col-span-2 bg-white rounded-lg">
      {#if selectedCandidate}
        <div class="flex items-start gap-3 p-3 rounded-lg mt-3">
          <img src={selectedCandidate.avatar} alt={selectedCandidate.name} class="w-20 h-20 rounded-full object-cover" />
          <div class="flex-1 flex flex-col items-start">
            <div class="flex items-center justify-between w-full">
              <h3 class="font-semibold text-gray-800 text-lg">{selectedCandidate.name}</h3>
              <div class="flex gap-2 font-semibold">
                <button class="flex px-4 py-1 text-xs bg-gray-100 rounded border border-gray-300 hover:bg-gray-500"><StickyNote class="w-4 h-4 mx-1" />Notes</button>
                <button class="flex px-4 py-1 text-xs bg-green-500 text-white rounded border border-green-300 hover:bg-green-700"><Check class="w-4 h-4 mx-1" />Accept</button>
                <button class="flex px-4 py-1 text-xs bg-red-500 text-white rounded border border-red-300 hover:bg-red-700"><X class="w-4 h-4 mx-1" />Reject</button>
              </div>
            </div>
            <p class="text-sm text-gray-800">{selectedCandidate.role}</p>
            <p class="text-sm text-gray-500">Applied for: {selectedCandidate.applied}</p>
            <p class="flex text-sm text-gray-500 gap-1 mt-1">
              <GraduationCap class="h-4 w-4 m-0.5" />{selectedCandidate.university}
              <Briefcase class="w-4 h-4 m-0.5" />{selectedCandidate.year}
              <Clock class="w-4 h-4 m-0.5" />{selectedCandidate.time}
            </p>
            <span
              class="inline-block mt-2 text-xs px-2 py-0.5 rounded-lg font-semibold
                {selectedCandidate.status === 'New' ? 'bg-blue-100 text-blue-700 border-1 border-blue-700' : ''}
                {selectedCandidate.status === 'Shortlisted' ? 'bg-purple-100 text-purple-700 border-1 border-purple-700' : ''}
                {selectedCandidate.status === 'Rejected' ? 'bg-red-100 text-red-700 border-1 border-red-700' : ''}
                {selectedCandidate.status === 'Reviewed' ? 'bg-yellow-100 text-yellow-700 border-1 border-yellow-700' : ''}"
            >
              {selectedCandidate.status}
            </span>
          </div>
        </div>
        <div class="p-4">       
          <hr class="w-full my-4 border-gray-300" />
            <h3 class="flex gap-1 font-semibold text-gray-800 text-lg my-4"><User class="mt-0.5" />Contact Information</h3>
            <div class="grid grid-cols-2 gap-4 w-full text-sm text-gray-700">
              <div class="flex flex-col">
                <span class="font-semibold text-gray-500 text-xs">Email</span>
                <span class="text-gray-800">{selectedCandidate.email || "alice.johnson@gmail.com"}</span>
              </div>
              <div class="flex flex-col">
                <span class="font-semibold text-gray-500 text-xs">Phone</span>
                <span class="text-gray-800">{selectedCandidate.phone || "+66 123 4567"}</span>
              </div>
              <div class="flex flex-col">
                <span class="font-semibold text-gray-500 text-xs">Address</span>
                <span class="text-gray-800">{selectedCandidate.address || "Bangkok, Thailand"}</span>
              </div>
              <div class="flex flex-col">
                <span class="font-semibold text-gray-500 text-xs">LinkedIn</span>
                <span class="text-gray-800">{selectedCandidate.linkedin || "None"}</span>
              </div>
            </div>
        </div>
        <div class="px-4 py-2">
          <h3 class="flex gap-1 font-semibold text-gray-800 text-lg mt-2"><GraduationCap class="mt-0.5" />Education</h3>
          <div class="grid grid-cols-1 gap-2 w-full text-sm mt-3">
            {#if selectedCandidate.education}
              {#each selectedCandidate.education as edu}
                <div class="flex flex-col">
                  <span class="font-semibold text-sm">{edu.degree}</span>
                  <span class="text-gray-800 mt-1">{edu.university} {edu.period}</span>
                  <span class="text-gray-800 mt-1">GPA {edu.gpa}</span>
                </div>
              {/each}
            {:else}
              <span class="text-gray-500 text-xs">No education details provided</span>
            {/if}
          </div>
        </div>
        <div class="px-4 py-2">
          <h3 class="flex gap-1 font-semibold text-gray-800 text-lg mt-2">
            <CodeXml class="mt-0.5" />Skills & Technologies
          </h3>
          <div class="flex flex-wrap gap-2 mt-3">
            {#if selectedCandidate.skills && selectedCandidate.skills.length > 0}
              {#each selectedCandidate.skills as skill}
                <span class="px-2 py-1 text-xs bg-gray-200 text-gray-800 rounded-lg border-1 border-gray-300">{skill}</span>
              {/each}
            {:else}
              <span class="text-gray-500 text-xs">No skills provided</span>
            {/if}
          </div>
        </div>
        <div class="px-4 py-2 mb-2">
          <h3 class="flex gap-1 font-semibold text-gray-800 text-lg mt-2">
            <FileText class="mt-0.5" />Documents & Portfolio
          </h3>
          <div class="flex flex-col gap-4 mt-3">
            {#if selectedCandidate.documents && selectedCandidate.documents.length > 0}
              {#each selectedCandidate.documents as doc}
                <a href={doc.name.startsWith('http') ? doc.name : '#'} target="_blank" class="flex items-center gap-4 bg-gray-100 border border-gray-300 rounded p-3 hover:bg-gray-200 transition">
                  {#await Promise.resolve(getDocIcon(doc.name)) then Icon}
                    <Icon class="w-8 h-8 text-gray-600" />
                  {/await}
                  <div class="flex flex-col">
                    <span class="text-sm font-semibold text-gray-800">{doc.name}</span>
                    <span class="text-xs text-gray-500">{doc.description}</span>
                  </div>
                </a>
              {/each}
            {:else}
              <span class="text-gray-500 text-xs">No documents uploaded</span>
            {/if}
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>

<script>
  import { Search, GraduationCap, Briefcase, Clock, StickyNote, Check, X, User, CodeXml, FileText, Github, Link } from 'lucide-svelte';

  let candidates = [
    {
      name: "Alice Johnson",
      role: "Frontend Developer",
      applied: "Senior Frontend Developer",
      status: "New",
      avatar: "https://randomuser.me/api/portraits/women/44.jpg",
      time: "2 hours ago",
      university: "KU",
      year: "2 years",
      email: "alice.johnson@gmail.com",
      phone: "+66 123 4567",
      address: "Bangkok, Thailand",
      linkedin: "linkedin.com/in/alicejohnson",
      education: [
        {
          degree: "Bachelor Degree of Computer Science",
          university: "Kasetsart University",
          period: "2015-2019",
          gpa: "4.00/4.00"
        }
      ],
      skills: [
        "HTML", "CSS", "JavaScript", "TypeScript", "React", 
        "Vue.js", "Svelte", "Node.js", "Express", "TailwindCSS",
        "Git", "Webpack", "REST API", "GraphQL", "Testing (Jest)"
      ],
      documents: [
        { name: "Resume_Alice_Johnson.pdf", description: "236kB" },
        { name: "Portfolio Website", description: "Personal projects" },
        { name: "Github Profile", description: "GitHub repositories" }
      ]
    },
    {
      name: "Michael Chen",
      role: "AI Researcher",
      applied: "Lead AI Scientist",
      status: "Shortlisted",
      avatar: "https://randomuser.me/api/portraits/men/32.jpg",
      time: "5 hours ago",
      university: "CU",
      year: "3 years",
      email: "michael.chen@gmail.com",
      phone: "+66 234 5678",
      address: "Chiang Mai, Thailand",
      linkedin: "linkedin.com/in/michaelchen",
      education: [
        {
          degree: "Master Degree of AI Engineering",
          university: "Chulalongkorn University",
          period: "2016-2018",
          gpa: "3.90/4.00"
        }
      ],
      skills: ["Python", "TensorFlow", "PyTorch", "Machine Learning", "Deep Learning"],
      documents: [
        { name: "Resume_Michael_Chen.pdf", description: "412kB" },
        { name: "Research_Papers_Michael_Chen.pdf", description: "2.3MB" }
      ]
    },
    {
      name: "Sophia Lee",
      role: "QA Tester",
      applied: "Automation QA Engineer",
      status: "Rejected",
      avatar: "https://randomuser.me/api/portraits/women/68.jpg",
      time: "1 day ago",
      university: "CMU",
      year: "1 year",
      email: "sophia.lee@gmail.com",
      phone: "+66 345 6789",
      address: "Phuket, Thailand",
      linkedin: "linkedin.com/in/sophialee",
      education: [
        {
          degree: "Bachelor Degree of Software Engineering",
          university: "Chiang Mai University",
          period: "2017-2021",
          gpa: "3.80/4.00"
        }
      ],
      skills: ["Selenium", "Cypress", "Jest", "Postman", "Bug Tracking"],
      documents: []
    },
    {
      name: "David Smith",
      role: "Backend Developer",
      applied: "Senior Backend Developer",
      status: "Reviewed",
      avatar: "https://randomuser.me/api/portraits/men/41.jpg",
      time: "3 days ago",
      university: "TU",
      year: "4 years",
      email: "david.smith@gmail.com",
      phone: "+66 456 7890",
      address: "Bangkok, Thailand",
      linkedin: "linkedin.com/in/davidsmith",
      education: [
        {
          degree: "Bachelor Degree of Computer Science",
          university: "Thammasat University",
          period: "2014-2018",
          gpa: "3.95/4.00"
        }
      ],
      skills: ["Java", "Spring Boot", "REST API", "SQL", "NoSQL"],
      documents: []
    }
  ];

  function getDocIcon(doc) {
    if (doc.toLowerCase().endsWith('.pdf')) return FileText;
    if (doc.toLowerCase().includes('github')) return Github;
    return Link; // fallback for any other link
  }
  let selectedCandidate = candidates[0];
</script>
