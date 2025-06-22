// app/login/page.tsx
'use client';

import { useState } from 'react';
import api from '../../../utils/axios-config';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '../context/AuthContext'; // Import the useAuth hook from your context

export default function Login() {
  const [formData, setFormData] = useState({
    email: '',
    password: '',
  });
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const router = useRouter();
  const { login } = useAuth(); // Destructure the login function from AuthContext

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage('');
    setError('');

    try {
      const response = await api.post('/login', formData);
      // **IMPORTANT**: Your backend must now return a 'user' object within the response.data.
      // This 'user' object should contain at least 'id', 'email', and crucially, 'role'.
      // For example: { message: "Login successful", user: { id: 1, email: "admin@example.com", role: "admin", first_name: "Admin", last_name: "User" } }
      const responseData = response.data as {
        message: string;
        user: { id: number; email: string; role: 'user' | 'admin'; first_name?: string; last_name?: string; };
      };

      setMessage(responseData.message);

      // Call the login function from AuthContext and pass the user data received from the backend.
      // This updates the global authentication state with the user's details and role.
      login(responseData.user);

      // Redirect the user based on their role for a tailored experience.
      if (responseData.user.role === 'admin') {
        router.push('/admin/dashboard'); // Send admins to their dashboard
      } else {
        router.push('/'); // Regular users go to the homepage
      }

    } catch (err: any) {
      setError(err.response?.data?.message || 'Login failed.');
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-md">
        <h2 className="text-3xl font-bold text-center text-gray-800 mb-6">Login</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="email">
              Email:
            </label>
            <input
              type="email"
              id="email"
              name="email"
              value={formData.email}
              onChange={handleChange}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              required
            />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="password">
              Password:
            </label>
            <input
              type="password"
              id="password"
              name="password"
              value={formData.password}
              onChange={handleChange}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              required
            />
          </div>
          <button
            type="submit"
            className="bg-green-500 hover:bg-green-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
          >
            Login
          </button>
        </form>
        {message && <p className="text-green-500 text-center mt-4">{message}</p>}
        {error && <p className="text-red-500 text-center mt-4">{error}</p>}
        <p className="text-center text-gray-600 text-sm mt-4">
          Don't have an account?{' '}
          <Link href="/register" className="text-green-500 hover:underline">
            Register here
          </Link>
        </p>
      </div>
    </div>
  );
}