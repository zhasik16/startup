import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import './globals.css';
import Navigation from '@/components/Layout/Navigation';
import { SessionProviderWrapper } from '@/components/SessionProviderWrapper';

const inter = Inter({ subsets: ['latin'] });

export const metadata: Metadata = {
  title: 'Aegis AI - Automated Security & Compliance',
  description: 'AI-powered security analysis for your codebase',
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className={inter.className}>
        <SessionProviderWrapper>
          <Navigation />
          <main>{children}</main>
        </SessionProviderWrapper>
      </body>
    </html>
  );
}