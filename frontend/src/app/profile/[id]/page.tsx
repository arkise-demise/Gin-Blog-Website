// app/profile/[id]/page.tsx
'use client';

import { useEffect, useState } from 'react';
import { useParams } from 'next/navigation';
import api from '../../../../utils/axios-config';
interface UserProfile {
  id: number;
  first_name: string;
  last_name: string;
  bio: string;
  profile_picture_url: string;
  location: string;
  website: string;
}

export default function UserProfilePage() {
  const { id } = useParams(); // Get the user ID from the URL
  const [profile, setProfile] = useState<UserProfile | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    async function fetchUserProfile() {
      try {
        setLoading(true);
        const response = await api.get(`/users/${id}/profile`);
        setProfile((response.data as { data: UserProfile }).data);
      } catch (err: any) {
        console.error('Failed to fetch user profile:', err);
        setError(err.response?.data?.message || 'Failed to load user profile.');
      } finally {
        setLoading(false);
      }
    }
    if (id) {
      fetchUserProfile();
    }
  }, [id]);

  if (loading) return <p className="text-center text-gray-600">Loading profile...</p>;
  if (error) return <p className="text-red-500 text-center">{error}</p>;
  if (!profile) return <p className="text-gray-600 text-center">Profile not found.</p>;

  return (
    <div className="container mx-auto p-8 bg-white shadow-lg rounded-lg">
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
        <h1 className="text-4xl font-bold text-gray-800 mt-4">
          {profile.first_name} {profile.last_name}
        </h1>
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
      </div>

      {profile.bio && (
        <div className="mb-8">
          <h2 className="text-2xl font-semibold text-gray-800 mb-2 border-b pb-2">About Me</h2>
          <p className="text-gray-700 leading-relaxed">{profile.bio}</p>
        </div>
      )}

      {/* You might want to display posts by this user here too */}
      {/* <h2 className="text-2xl font-semibold text-gray-800 mb-4 border-b pb-2">Posts by {profile.first_name}</h2> */}
      {/* Add logic to fetch and display posts by this user */}
    </div>
  );
}