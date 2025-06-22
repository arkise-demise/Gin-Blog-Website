// components/Navbar.tsx
"use client";

import Link from 'next/link';
import { useAuth } from '../app/context/AuthContext';
import { useRouter } from 'next/navigation';

export default function Navbar() {
  const { isAuthenticated, isAdmin, logout } = useAuth();
  const router = useRouter();

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  // Reusable classes for smaller, blue-themed action links
  const actionLinkClasses = "text-white text-sm font-medium py-1.5 px-3 rounded-md transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-opacity-50";

  // Classes for Login/Register (kept original size as they are less common once logged in)
  const authLinkClasses = "bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50";


  return (
    <nav className="bg-gray-800 p-4 shadow-md">
      <div className="container mx-auto flex justify-between items-center">
        {/* Site Logo (Left Side) */}
        <Link href="/" className="flex items-center space-x-3">
          <img
            src="/gin-blog-logo.png"
            alt="Gin Blog Logo"
            className="h-30 w-30 object-contain"
          />
          <span className="text-white text-3xl font-bold italic hover:text-blue-400 transition duration-300">
            <span className="text-pink-400 font-extrabold italic">Gin Blog</span>
          </span>
        </Link>

        {/* Navigation Buttons (Right Side) */}
        <div className="flex space-x-4">
          {/* Conditional rendering based on isAuthenticated */}
          {!isAuthenticated && (
            // For guests (not logged in)
            <>
              <Link href="/register" className={authLinkClasses.replace('bg-blue-500', 'bg-blue-500')}> {/* Keep blue for Register */}
                Register
              </Link>
              <Link href="/login" className={authLinkClasses.replace('bg-green-500', 'bg-indigo-500')}> {/* Changed Login to indigo for slight variation but same theme */}
                Login
              </Link>
            </>
          )}

            {isAuthenticated && !isAdmin && (
            // For authenticated users who are NOT admins
            <>
              <Link href="/" className={`${actionLinkClasses} bg-blue-500 hover:bg-blue-600 focus:ring-blue-500`}>
                Home
              </Link>
              <Link href="/create-post" className={`${actionLinkClasses} bg-blue-500 hover:bg-blue-600 focus:ring-blue-500`}>
                Create Post
              </Link>
              <Link href="/my-posts" className={`${actionLinkClasses} bg-blue-500 hover:bg-blue-600 focus:ring-blue-500`}>
                My Posts
              </Link>
              {/* NEW: My Profile link for regular users */}
              <Link href="/my-profile" className={`${actionLinkClasses} bg-blue-500 hover:bg-blue-600 focus:ring-blue-500`}>
                My Profile
              </Link>
              <button
              onClick={handleLogout}
              className={`${actionLinkClasses} bg-blue-500 hover:bg-blue-600 focus:ring-blue-500`}
              >
              Logout
              </button>
            </>
            )}

            {isAuthenticated && isAdmin && (
            // This block is intentionally empty for authenticated admins.
            // Their navigation is handled by the sidebar.
            <></>
            )}
        </div>
      </div>
    </nav>
  );
}