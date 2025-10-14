// app/api/debug/auth/route.ts
import { getServerSession } from "next-auth/next"
import { authOptions } from "../../auth/[...nextauth]/route"

export async function GET() {
  const session = await getServerSession(authOptions)
  
  return Response.json({
    session: session ? {
      user: session.user,
      hasAccessToken: !!session.accessToken,
      accessToken: session.accessToken ? 'PRESENT' : 'MISSING',
      provider: session.provider,
    } : 'NO_SESSION',
    authOptions: {
      hasGitHubId: !!process.env.GITHUB_ID,
      hasGitHubSecret: !!process.env.GITHUB_SECRET,
    }
  })
}