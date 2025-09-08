<script>

</script>

<div class="min-h-screen flex flex-col lg:flex-row">
  <!-- Left Side - Sign up section -->
  <div class="w-full lg:w-3/6 flex items-center justify-center p-8">
    <div class="w-full max-w-lg">
      <div class="text-center mb-8">
        <h2 class="text-2xl font-semibold text-gray-800 mb-2">Sign up</h2>
        <div class="flex items-baseline justify-center gap-1 mb-3">
          <h1 class="text-2xl font-semibold text-gray-900">Job Applier</h1>
          <h1 class="text-2xl font-semibold text-green-700">3000</h1>
        </div>
        <p class="text-sm text-gray-500">Are you looking for a job or hiring talented graduates?</p>
      </div>

      <div class="space-y-3 mb-8">
        <!-- Option 1 - Student/Alumni -->
        <a href="/auth/signup/student" class="block border border-gray-200 rounded-lg p-4 cursor-pointer hover:border-green-600 transition-colors">
          <div class="flex items-start space-x-3">
            <div class="flex-shrink-0">
              <div class="flex -space-x-2">
                <div class="w-8 h-8 bg-purple-100 rounded-full flex items-center justify-center border-2 border-white">
                </div>
                <div class="w-8 h-8 bg-blue-100 rounded-full flex items-center justify-center border-2 border-white">
                </div>
                <div class="w-8 h-8 bg-green-100 rounded-full flex items-center justify-center border-2 border-white">
                </div>
              </div>
            </div>
            <div class="flex-1">
              <h3 class="font-medium text-gray-900 text-sm">I'm a student/alumni</h3>
              <p class="text-xs text-gray-500 mt-1">I'm a student/alumni looking for jobs.</p>
            </div>
          </div>
        </a>

        <!-- Option 2 - Recruiter -->
        <a href="/auth/signup/company" class="block border border-gray-200 rounded-lg p-4 cursor-pointer hover:border-green-600 transition-colors">
          <div class="flex items-start space-x-3">
            <div class="flex-shrink-0">
              <div class="flex -space-x-2">
                <div class="w-8 h-8 bg-red-100 rounded-full flex items-center justify-center border-2 border-white">
                </div>
                <div class="w-8 h-8 bg-yellow-100 rounded-full flex items-center justify-center border-2 border-white">
                </div>
                <div class="w-8 h-8 bg-indigo-100 rounded-full flex items-center justify-center border-2 border-white">
                </div>
              </div>
            </div>
            <div class="flex-1">
              <h3 class="font-medium text-gray-900 text-sm">I want to hire talents</h3>
              <p class="text-xs text-gray-500 mt-1">I'm a recruiter hiring talented graduates.</p>
            </div>
          </div>
        </a>
      </div>

      <div class="text-center">
        <p class="text-sm text-gray-500">
          Already have an account? 
          <a href="/auth/login" class="text-green-600 hover:text-green-700 font-medium">Log in</a>
        </p>
      </div>
    </div>
  </div>

  <!-- Right Side -->
  <div class="w-full lg:w-3/5 bg-gradient-to-br from-green-50 to-blue-50 relative overflow-hidden">
    <div class="h-full flex flex-col justify-between p-12">
      <div class="max-w-2xl pt-20 px-20">
        <p class="text-gray-700 text-2xl leading-relaxed mb-6">
          Job Applier 3000 connects CPE students, alumni, and companies in a platform where students find opportunities and recruiters reach verified talent
        </p>
        <p class="text-2xl font-semibold text-gray-900">Join now to get started!</p>
      </div>

      <div class="flex justify-end">
        <div class="w-full max-w-4xl -mr-18 -mb-18">
          <div class="bg-white rounded-3xl border-5 border-black overflow-hidden">
            <!-- Image placeholder -->
            <div class="w-full h-160 bg-gray-200 flex items-center justify-center">
              <div class="text-center">
                <p class="text-gray-500 text-lg">Dashboard Image</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</div>
