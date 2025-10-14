import { NextRequest } from 'next/server'

export async function GET(request: NextRequest) {
  return Response.json({
    message: "Debug API is working",
    timestamp: new Date().toISOString(),
    environment: {
      NODE_ENV: process.env.NODE_ENV,
      BACKEND_URL: process.env.NEXT_PUBLIC_BACKEND_URL,
    }
  })
}