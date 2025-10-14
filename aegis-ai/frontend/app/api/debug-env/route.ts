// app/api/debug-env/route.ts
export async function GET() {
  return Response.json({
    hasGithubId: !!process.env.GITHUB_ID,
    hasGithubSecret: !!process.env.GITHUB_SECRET,
    hasNextAuthSecret: !!process.env.NEXTAUTH_SECRET,
    hasGithubClientId: !!process.env.GITHUB_CLIENT_ID,
    hasGithubClientSecret: !!process.env.GITHUB_CLIENT_SECRET,
  })
}