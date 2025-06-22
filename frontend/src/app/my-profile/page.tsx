// app/my-profile/page.tsx
'use client';

import { useEffect, useState } from 'react';
import api from '../../../utils/axios-config';
import { useAuth } from '../../app/context/AuthContext'; // Adjust path
import { useRouter } from 'next/navigation';

interface UserProfile {
  id: number;
  first_name: string;
  last_name: string;
  bio: string;
  profile_picture_url: string;
  location: string;
  website: string;
  // Add other fields you fetch/display for my profile
}

export default function MyProfilePage() {
  const { isAuthenticated, user } = useAuth();
  const router = useRouter();
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [formData, setFormData] = useState<UserProfile | null>(null);
  const [updateMessage, setUpdateMessage] = useState<string | null>(null);

  useEffect(() => {
    if (!isAuthenticated) {
      router.push('/login'); // Redirect if not logged in
      return;
    }

    async function fetchMyProfile() {
      try {
        setLoading(true);
        const response = await api.get('/my-profile'); // Call the /my-profile endpoint
        const data = response.data as { data: UserProfile };
        setProfile(data.data);
        setFormData(data.data); // Initialize form data with current profile
      } catch (err: any) {
        console.error('Failed to fetch my profile:', err);
        setError(err.response?.data?.message || 'Failed to load my profile.');
      } finally {
        setLoading(false);
      }
    }
    fetchMyProfile();
  }, [isAuthenticated, router]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData({ ...formData!, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setUpdateMessage(null);
    try {
      const response = await api.put('/my-profile', formData); // Update my profile
      const data = response.data as { data: UserProfile; message?: string };
      setProfile(data.data); // Update displayed profile
      setUpdateMessage(data.message || 'Profile updated successfully!');
      setIsEditing(false); // Exit editing mode
    } catch (err: any) {
      console.error('Failed to update profile:', err);
      setUpdateMessage(err.response?.data?.message || 'Failed to update profile.');
    }
  };

  if (loading) return <p className="text-center text-gray-600">Loading profile...</p>;
  if (error) return <p className="text-red-500 text-center">{error}</p>;
  if (!profile) return <p className="text-gray-600 text-center">Profile not found.</p>;

  return (
    <div className="container mx-auto p-8 bg-white shadow-lg rounded-lg">
      <h1 className="text-4xl font-bold text-gray-800 mb-8 text-center">My Profile</h1>

      {updateMessage && (
        <p className={`text-center mb-4 ${updateMessage.includes('successfully') ? 'text-green-600' : 'text-red-500'}`}>
          {updateMessage}
        </p>
      )}

      {!isEditing ? (
        // Display Mode
        <div className="flex flex-col items-center mb-8">
          {profile.profile_picture_url ? (
            <img
              src={profile.profile_picture_url}
              alt={`${profile.first_name}'s profile`}
              className="w-32 h-32 rounded-full object-cover border-4 border-indigo-400"
            />
          ) : (
            <div className="w-32 h-32 rounded-full bg-gray-300 flex items-center justify-center text-gray-600 text-5xl font-bold">
              {profile.first_name ? profile.first_name.charAt(0) : 'U'}
            </div>
          )}
          <h2 className="text-4xl font-bold text-gray-800 mt-4">
            {profile.first_name} {profile.last_name}
          </h2>
          <p className="text-gray-600 text-lg mt-2">Email: {user?.email}</p> {/* Display email from auth context */}
          {profile.location && (
            <p className="text-gray-600 text-lg mt-2">
              <i className="fas fa-map-marker-alt mr-2"></i>{profile.location}
            </p>
          )}
          {profile.website && (
            <p className="text-blue-500 hover:underline mt-1">
              <a href={profile.website} target="_blank" rel="noopener noreferrer">
                {profile.website}
              </a>
            </p>
          )}

          {profile.bio && (
            <div className="mt-8 w-full">
              <h3 className="text-2xl font-semibold text-gray-800 mb-2 border-b pb-2">About Me</h3>
              <p className="text-gray-700 leading-relaxed">{profile.bio}</p>
            </div>
          )}

          <button
            onClick={() => setIsEditing(true)}
            className="mt-8 bg-indigo-500 hover:bg-indigo-600 text-white font-semibold py-2 px-6 rounded-lg transition duration-300 ease-in-out transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-opacity-50"
          >
            Edit Profile
          </button>
        </div>
      ) : (
        // Edit Mode
        <form onSubmit={handleSubmit} className="space-y-6">
          <div>
            <label htmlFor="first_name" className="block text-gray-700 text-sm font-bold mb-2">First Name:</label>
            <input
              type="text"
              id="first_name"
              name="first_name"
              value={formData?.first_name || ''}
              onChange={handleChange}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
          </div>
          <div>
            <label htmlFor="last_name" className="block text-gray-700 text-sm font-bold mb-2">Last Name:</label>
            <input
              type="text"
              id="last_name"
              name="last_name"
              value={formData?.last_name || ''}
              onChange={handleChange}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
          </div>
          <div>
            <label htmlFor="bio" className="block text-gray-700 text-sm font-bold mb-2">Bio:</label>
            <textarea
              id="bio"
              name="bio"
              value={formData?.bio || ''}
              onChange={handleChange}
              rows={4}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline resize-none"
            ></textarea>
          </div>
          <div>
            <label htmlFor="profile_picture_url" className="block text-gray-700 text-sm font-bold mb-2">Profile Picture URL:</label>
            <input
              type="url"
              id="profile_picture_url"
              name="profile_picture_url"
              value={formData?.profile_picture_url || ''}
              onChange={handleChange}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
          </div>
          <div>
            <label htmlFor="location" className="block text-gray-700 text-sm font-bold mb-2">Location:</label>
            <input
              type="text"
              id="location"
              name="location"
              value={formData?.location || ''}
              onChange={handleChange}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
          </div>
          <div>
            <label htmlFor="website" className="block text-gray-700 text-sm font-bold mb-2">Website:</label>
            <input
              type="url"
              id="website"
              name="website"
              value={formData?.website || ''}
              onChange={handleChange}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
          </div>
          <div className="flex justify-end space-x-4">
            <button
              type="button"
              onClick={() => setIsEditing(false)}
              className="bg-gray-400 hover:bg-gray-500 text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
            >
              Cancel
            </button>
            <button
              type="submit"
              className="bg-green-500 hover:bg-green-600 text-white font-semibold py-2 px-4 rounded-lg transition duration-300"
            >
              Save Changes
            </button>
          </div>
        </form>
      )}
    </div>
  );
}