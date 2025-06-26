// components/Navbar.tsx
"use client";

import Link from 'next/link';
import { useAuth } from '../app/context/AuthContext';
import { useRouter } from 'next/navigation';
import { useState, useEffect } from 'react';

export default function Navbar() {
  const { isAuthenticated, isAdmin, user, loading, logout } = useAuth();
  const router = useRouter();

  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false);
  const [isClient, setIsClient] = useState(false);

  useEffect(() => {
    setIsClient(true);
  }, []);

  const handleLogout = () => {
    logout();
    router.push('/login');
  };

  // Reusable classes for non-colored action links (text-only)
  const actionLinkClasses = "text-white text-sm font-medium py-1.5 px-3 rounded-md transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none border-2 border-transparent hover:border-blue-400";

  // Classes for Login/Register (kept original size with their colors)
  const authLinkClasses = "bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50";

  // Admin Profile Picture specific classes (no transform, no padding, just the image)
  const adminProfilePicClasses = "w-12 h-12 rounded-full object-cover border-2 border-blue-400";


  return (
    <nav className="bg-gray-800 p-4 shadow-md">
      <div className="container mx-auto flex justify-between items-center">
        {/* Site Logo (Left Side) - Always visible */}
        <Link href="/" className="flex items-center space-x-3">
          <img
            src="/gin-blog-logo.png"
            alt="Gin Blog Logo"
            className="h-12 w-auto object-contain"
          />
          <span className="text-white text-3xl font-bold italic hover:text-blue-400 transition duration-300">
            <span className="text-pink-400 font-extrabold italic">Gin Blog</span>
          </span>
        </Link>

        {/* Conditional rendering for the Right Side of the Navbar */}
        {loading ? (
          // Loading state placeholder for all scenarios
          <div className="flex space-x-4 items-center">
              <div className="w-24 h-6 bg-gray-700 rounded animate-pulse"></div>
              <div className="w-24 h-6 bg-gray-700 rounded animate-pulse"></div>
              <div className="w-12 h-12 rounded-full bg-gray-700 animate-pulse"></div>
          </div>
        ) : !isAuthenticated ? (
          // Guest navigation (Login/Register) - always visible desktop & mobile
          <div className="flex space-x-4 items-center">
            <Link href="/register" className={authLinkClasses}>
              Register
            </Link>
            <Link href="/login" className={authLinkClasses.replace('bg-blue-500', 'bg-indigo-500')}>
              Login
            </Link>
          </div>
        ) : isAdmin ? (
          // Admin navigation (only profile pic and logout on the right)
          <div className="flex space-x-4 items-center">
            {/* Admin Profile Picture */}
            <Link href="/my-profile" className="bg-transparent !px-0 !py-0 !rounded-full"> {/* No extra padding/margin on link */}
                {user?.profile_picture_url ? (
                    <img
                        src={user.profile_picture_url}
                        alt="Profile"
                        className={adminProfilePicClasses}
                    />
                ) : (
                    <div className={`${adminProfilePicClasses} bg-blue-500 flex items-center justify-center text-md font-semibold text-white`}>
                        {user?.first_name ? user.first_name[0].toUpperCase() : 'U'}
                    </div>
                )}
            </Link>
          </div>
        ) : (
          // Regular Authenticated User navigation (with hamburger for mobile)
          <>
            {/* Hamburger Icon for Mobile (Right Side) - ONLY RENDER ON CLIENT for regular users */}
            {isClient && (
              <div className="md:hidden">
                <button
                  onClick={() => setIsMobileMenuOpen(!isMobileMenuOpen)}
                  className="text-white focus:outline-none focus:ring-2 focus:ring-blue-500 p-2 rounded"
                  aria-label="Toggle mobile menu"
                >
                  <svg className="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                    {isMobileMenuOpen ? (
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M6 18L18 6M6 6l12 12"></path>
                    ) : (
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth="2" d="M4 6h16M4 12h16M4 18h16"></path>
                    )}
                  </svg>
                </button>
              </div>
            )}

            {/* Desktop Navigation Buttons (Right Side - Regular User) */}
            <div className="hidden md:flex space-x-4 items-center">
              <Link href="/" className={actionLinkClasses}>
                Home
              </Link>
              <Link href="/create-post" className={actionLinkClasses}>
                Create Post
              </Link>
              <Link href="/my-posts" className={actionLinkClasses}>
                My Posts
              </Link>

              {/* Logout Button */}
              <button
                onClick={handleLogout}
                className={actionLinkClasses}
              >
                Logout
              </button>

              {/* Profile Picture */}
              <Link href="/my-profile" className={`${actionLinkClasses} bg-transparent !px-0 !py-0 !rounded-full`}>
                  {user?.profile_picture_url ? (
                      <img
                          src={user.profile_picture_url}
                          alt="Profile"
                          className="w-12 h-12 rounded-full object-cover border-2 border-blue-400 transition transform hover:scale-110"
                      />
                  ) : (
                      <div className="w-12 h-12 rounded-full bg-blue-500 flex items-center justify-center text-md font-semibold text-white border-2 border-blue-400">
                          {user?.first_name ? user.first_name[0].toUpperCase() : 'U'}
                      </div>
                  )}
              </Link>
            </div>
          </>
        )}
      </div>

      {/* Mobile Menu (Collapsible) - ONLY RENDER ON CLIENT for regular users */}
      {isClient && isAuthenticated && !isAdmin && ( // Only show mobile menu for regular users
        <div className={`md:hidden ${isMobileMenuOpen ? 'block' : 'hidden'} mt-4`}>
          <div className="flex flex-col space-y-2 px-2 pb-3 sm:px-3">
            {loading ? (
              <div className="flex flex-col space-y-2">
                <div className="w-full h-8 bg-gray-700 rounded animate-pulse"></div>
                <div className="w-full h-8 bg-gray-700 rounded animate-pulse"></div>
                <div className="w-full h-8 bg-gray-700 rounded animate-pulse"></div>
              </div>
            ) : (
              <>
                <Link href="/" className={`${actionLinkClasses} bg-gray-700 hover:bg-gray-600 w-full text-center`} onClick={() => setIsMobileMenuOpen(false)}>Home</Link>
                <Link href="/create-post" className={`${actionLinkClasses} bg-gray-700 hover:bg-gray-600 w-full text-center`} onClick={() => setIsMobileMenuOpen(false)}>Create Post</Link>
                <Link href="/my-posts" className={`${actionLinkClasses} bg-gray-700 hover:bg-gray-600 w-full text-center`} onClick={() => setIsMobileMenuOpen(false)}>My Posts</Link>
                {/* Profile and Logout for mobile menu */}
                <Link href="/my-profile" className={`${actionLinkClasses} bg-gray-700 hover:bg-gray-600 w-full text-center`} onClick={() => setIsMobileMenuOpen(false)}>My Profile</Link>
                <button
                  onClick={() => { handleLogout(); setIsMobileMenuOpen(false); }}
                  className={`${actionLinkClasses} bg-gray-700 hover:bg-gray-600 w-full text-center`}
                >
                  Logout
                </button>
              </>
            )}
          </div>
        </div>
      )}
    </nav>
  );
}