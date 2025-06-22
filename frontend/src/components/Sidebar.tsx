// components/Sidebar.tsx
'use client';

import Link from 'next/link';
import { useAuth } from '../app/context/AuthContext';
import { usePathname } from 'next/navigation';

export default function Sidebar() {
  const { isAuthenticated, isAdmin, logout } = useAuth();
  const pathname = usePathname();

  // This sidebar should ONLY render if the user is authenticated AND is an admin
  if (!isAuthenticated || !isAdmin) {
    return null; // Don't render anything if not an authenticated admin
  }

  const isActive = (path: string) => pathname === path;
  const baseLinkClasses = "block py-2.5 px-4 rounded transition duration-200 hover:bg-gray-700 hover:text-white";
  const activeLinkClasses = "bg-gray-700 text-white";

  return (
    <aside className="w-64 bg-gray-800 text-gray-200 min-h-screen p-4 flex flex-col shadow-lg">
      <div className="text-2xl font-bold text-white mb-8 text-center">
        Admin Panel
      </div>
      <nav className="flex-1">
        <ul className="space-y-2">
          <li>
            <Link href="/" className={`${baseLinkClasses} ${isActive('/') ? activeLinkClasses : ''}`}>
              Home
            </Link>
          </li>
          <li>
            <Link href="/admin/dashboard" className={`${baseLinkClasses} ${isActive('/admin/dashboard') ? activeLinkClasses : ''}`}>
              Dashboard
            </Link>
          </li>
          <li>
            <Link href="/create-post" className={`${baseLinkClasses} ${isActive('/create-post') ? activeLinkClasses : ''}`}>
              Create Post
            </Link>
          </li>
          <li>
            <Link href="/my-posts" className={`${baseLinkClasses} ${isActive('/my-posts') ? activeLinkClasses : ''}`}>
              My Posts
            </Link>
          </li>
          {/* NEW: My Profile link for admins */}
          <li>
            <Link href="/my-profile" className={`${baseLinkClasses} ${isActive('/my-profile') ? activeLinkClasses : ''}`}>
              My Profile
            </Link>
          </li>

          <li>
            <button
              onClick={logout}
              className={`w-full text-left py-2.5 px-4 rounded transition duration-200 hover:bg-red-600 hover:text-white ${isActive('/logout') ? activeLinkClasses : ''} mt-4`}
         >
              Logout
            </button>
          </li>
        </ul>
      </nav>
      <div className="mt-auto text-sm text-gray-500 text-center">
        &copy; {new Date().getFullYear()}
      </div>
    </aside>
  );
}