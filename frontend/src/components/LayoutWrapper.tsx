// components/LayoutWrapper.tsx
'use client';

import { useAuth } from '../app/context/AuthContext';
import Sidebar from './Sidebar';

export default function LayoutWrapper({ children }: { children: React.ReactNode }) {
  const { isAdmin } = useAuth();

  return (
    // min-h-[calc(100vh-64px)] adjusts for a 64px tall Navbar. Adjust if your Navbar height differs
    <div className="flex min-h-[calc(100vh-64px)]">
      {isAdmin && <Sidebar />} {/* Sidebar only renders if isAdmin is true */}
      <main className={`flex-1 p-8 overflow-y-auto ${!isAdmin ? 'w-full' : ''}`}>
        {children}
      </main>
    </div>
  );
}