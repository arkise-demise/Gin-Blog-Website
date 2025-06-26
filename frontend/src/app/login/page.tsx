// app/login/page.tsx
'use client';

import { useState } from 'react';
import api from '../../../utils/axios-config';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '../context/AuthContext';

export default function Login() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const router = useRouter();
  const { login } = useAuth();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage('');
    setError('');

    try {
      const response = await api.post('/login', formData);
      const responseData = response.data as {
        message: string;
        user: { id: number; email: string; role: 'user' | 'admin'; first_name?: string; last_name?: string; profile_picture_url?: string; }; // Added profile_picture_url
      };

      setMessage(responseData.message);
      login(responseData.user);

      if (responseData.user.role === 'admin') {
        router.push('/admin/dashboard');
      } else {
        router.push('/');
      }

    } catch (err: any) {
      setError(err.response?.data?.message || 'Login failed. Please check your credentials.'); // More user-friendly error
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-blue-500 to-purple-600 p-4"> {/* Enhanced background */}
      <div className="bg-white p-8 sm:p-10 rounded-xl shadow-2xl w-full max-w-md border border-gray-200"> {/* Sharper shadow, rounded corners, subtle border */}
        <h2 className="text-4xl font-extrabold text-center text-gray-900 mb-8">Welcome Back!</h2> {/* Larger, bolder title */}
        <p className="text-center text-gray-600 mb-6">Sign in to access your blog.</p> {/* Subtitle */}
        <form onSubmit={handleSubmit} className="space-y-6"> {/* More vertical space */}
          <div>
            <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="email"> {/* Bolder label */}
              Email Address
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-200" // Modern input style
              placeholder="you@example.com"
              required
            />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="password">
              Password
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-200"
              placeholder="••••••••"
              required
            />
          </div>
          <button
            type="submit"
            className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-3 px-6 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 w-full transition duration-300 transform hover:scale-105" // More prominent button
          >
            Login
          </button>
        </form>
        {message && <p className="text-green-600 text-center mt-6 font-medium">{message}</p>} {/* Stronger color */}
        {error && <p className="text-red-600 text-center mt-6 font-medium">{error}</p>} {/* Stronger color */}
        <p className="text-center text-gray-600 text-sm mt-8"> {/* More margin */}
          Don't have an account?{' '}
          <Link href="/register" className="text-blue-600 hover:text-blue-800 hover:underline font-semibold"> {/* Bolder link, stronger hover */}
            Register here
          </Link>
        </p>
      </div>
    </div>
  );
}