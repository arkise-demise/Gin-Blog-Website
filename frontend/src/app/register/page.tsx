// app/register/page.tsx
'use client';

import { useState } from 'react';
import api from '../../../utils/axios-config';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

export default function Register() {
  const [formData, setFormData] = useState({
    first_name: '',
    last_name: '',
    email: '',
    phone: '',
    password: '',
  });
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const router = useRouter();

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage('');
    setError('');

    try {
      const response = await api.post('/register', formData);
      setMessage((response.data as { message: string }).message);
      // Optional: Clear form after successful registration
      setFormData({
        first_name: '',
        last_name: '',
        email: '',
        phone: '',
        password: '',
      });
      router.push('/login?registered=true'); // Redirect to login with a query param
    } catch (err: any) {
      setError(err.response?.data?.message || 'Registration failed. Please try again.'); // More user-friendly error
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-blue-500 to-purple-600 p-4"> {/* Consistent background */}
      <div className="bg-white p-8 sm:p-10 rounded-xl shadow-2xl w-full max-w-md border border-gray-200"> {/* Consistent card styling */}
        <h2 className="text-4xl font-extrabold text-center text-gray-900 mb-8">Join Gin Blog!</h2> {/* Catchy title */}
        <p className="text-center text-gray-600 mb-6">Create your account to start blogging.</p>
        <form onSubmit={handleSubmit} className="space-y-6"> {/* Consistent spacing */}
          <div>
            <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="first_name">
              First Name
            </label>
            <input
              type="text"
              id="first_name"
              name="first_name"
              value={formData.first_name}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-200"
              placeholder="John"
              required
            />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="last_name">
              Last Name
            </label>
            <input
              type="text"
              id="last_name"
              name="last_name"
              value={formData.last_name}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-200"
              placeholder="Doe"
              required
            />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="email">
              Email Address
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus://focus:ring-blue-500 focus:border-transparent transition duration-200"
              placeholder="you@example.com"
              required
            />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-semibold mb-2" htmlFor="phone">
              Phone (Optional)
            </label>
            <input
              type="text"
              id="phone"
              name="phone"
              value={formData.phone}
              onChange={handleChange}
              className="w-full px-4 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent transition duration-200"
              placeholder="e.g., +251912345678"
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
          {/* Distinct button color */}
          <button
            type="submit"
            className="bg-purple-600 hover:bg-purple-700 text-white font-bold py-3 px-6 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:ring-offset-2 w-full transition duration-300 transform hover:scale-105"
          >
            Register
          </button>
        </form>
        {message && <p className="text-green-600 text-center mt-6 font-medium">{message}</p>}
        {error && <p className="text-red-600 text-center mt-6 font-medium">{error}</p>}
        <p className="text-center text-gray-600 text-sm mt-8">
          Already have an account?{' '}
          <Link href="/login" className="text-purple-600 hover:text-purple-800 hover:underline font-semibold"> {/* Distinct link color */}
            Login here
          </Link>
        </p>
      </div>
    </div>
  );
}