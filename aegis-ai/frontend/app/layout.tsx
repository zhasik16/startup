import type { Metadata } from 'next';
import { Inter } from 'next/font/google';
import Navigation from '@/components/Layout/Navigation';
import './globals.css';

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
        <div className="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
          <Navigation />
          <main>{children}</main>
        </div>
      </body>
    </html>
  );
}