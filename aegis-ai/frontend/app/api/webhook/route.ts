import { NextRequest, NextResponse } from 'next/server';
import { PullRequestEvent } from '@/types';

export async function POST(request: NextRequest) {
  try {
    const event: PullRequestEvent = await request.json();
    
    console.log('🔔 Webhook received:', {
      action: event.action,
      prNumber: event.number,
      repository: event.repository.name
    });

    // Forward to your Go backend
    const backendResponse = await fetch(`${process.env.BACKEND_URL}/webhook`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(event),
    });

    if (!backendResponse.ok) {
      throw new Error('Backend processing failed');
    }

    return NextResponse.json(
      { status: 'processing', message: 'Analysis started' },
      { status: 202 }
    );
  } catch (error) {
    console.error('Webhook error:', error);
    return NextResponse.json(
      { error: 'Webhook processing failed' },
      { status: 500 }
    );
  }
}