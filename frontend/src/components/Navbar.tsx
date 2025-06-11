// components/Navbar.tsx
"use client";

import Link from 'next/link';
import { useAuth } from '../app/context/AuthContext'; // Import useAuth
import { useRouter } from 'next/navigation'; // For redirection

export default function Navbar() {
  const { isLoggedIn, logout } = useAuth(); // Use the context hook
  const router = useRouter();

  const handleLogout = () => {
    logout(); // Call the logout function from context
    router.push('/login'); // Redirect after logout
  };

  return (
    <nav className="bg-gray-800 p-4 shadow-md">
      <div className="container mx-auto flex justify-between items-center">
        {/* Site Logo (Left Side) */}
        <Link href="/" className="flex items-center space-x-3">
          <img
            src="/gin-blog-logo.png"
            alt="Gin Blog Logo"
            className="h-30 w-30 object-contain" // Increased from h-14 w-14 to h-20 w-20
          />
          <span className="text-white text-3xl font-bold italic hover:text-blue-400 transition duration-300">
            <span className="text-pink-400 font-extrabold italic">Gin Blog</span>
          </span>
        </Link>

        {/* Navigation Buttons (Right Side) */}
        <div className="flex space-x-4">
          {!isLoggedIn && (
            <>
             <Link href="/register" className="bg-blue-500 hover:bg-blue-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-opacity-50">
                Register
              </Link>
              <Link href="/login" className="bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-green-500 focus:ring-opacity-50">
                Login
              </Link>
            </>
          )}

          {isLoggedIn && (
            <>
              <Link href="/create-post" className="bg-purple-500 hover:bg-purple-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-opacity-50">
                Create Post
              </Link>
              <Link href="/my-posts" className="bg-yellow-500 hover:bg-yellow-600 text-gray-800 font-semibold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-yellow-500 focus:ring-opacity-50">
                My Posts
              </Link>
               <button // Basic Logout button example
                onClick={handleLogout} // Call the context's logout function
                className="bg-red-500 hover:bg-red-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-opacity-50"
              >
                Logout
              </button>
            </>
          )}
        </div>
      </div>
    </nav>
  );
}