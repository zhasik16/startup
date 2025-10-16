import NextAuth, { NextAuthOptions } from "next-auth"
import GitHubProvider from "next-auth/providers/github"

export const authOptions: NextAuthOptions = {
  providers: [
    GitHubProvider({
      clientId: process.env.GITHUB_ID!,
      clientSecret: process.env.GITHUB_SECRET!,
      authorization: {
        params: {
          scope: "repo user", // Fixed scopes - this is the key!
        },
      },
    })
  ],
  callbacks: {
    async jwt({ token, account }: any) {
      console.log('🔑 JWT Callback - Account:', account ? 'Has account' : 'No account')
      if (account) {
        token.accessToken = account.access_token
        token.provider = account.provider
        console.log('🔑 JWT - Access token set:', account.access_token?.substring(0, 20) + '...')
      }
      return token
    },
    async session({ session, token }: any) {
      console.log('🔑 Session Callback - Token has accessToken:', !!token.accessToken)
      session.accessToken = token.accessToken
      session.provider = token.provider
      console.log('🔑 Session - Access token:', session.accessToken?.substring(0, 20) + '...')
      return session
    },
  },
  debug: true, // Add debug to see what's happening
}

const handler = NextAuth(authOptions)
export { handler as GET, handler as POST }