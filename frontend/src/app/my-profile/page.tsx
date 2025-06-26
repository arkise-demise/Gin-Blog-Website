// app/my-profile/page.tsx (Conceptual)
"use client";

import React, { useState, useEffect } from 'react';
import { useAuth } from '@/app/context/AuthContext';
import api from '../../../utils/axios-config';
import { useRouter } from 'next/navigation';

interface UserProfile {
    first_name?: string;
    last_name?: string;
    email: string;
    phone?: string;
    bio?: string;
    profile_picture_url?: string;
    location?: string;
    website?: string;
}

export default function MyProfilePage() {
    const { user, isAuthenticated, loading: authLoading } = useAuth();
    const router = useRouter();
    const [profile, setProfile] = useState<UserProfile | null>(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState<string | null>(null);
    const [successMessage, setSuccessMessage] = useState<string | null>(null);
    const [file, setFile] = useState<File | null>(null);
    const [previewUrl, setPreviewUrl] = useState<string | null>(null);

    useEffect(() => {
        if (authLoading) return;

        if (!isAuthenticated) {
            router.push('/login');
            return;
        }

        const fetchProfile = async () => {
            try {
                const response = await api.get<UserProfile>('/my-profile');
                setProfile(response.data);
                if (response.data.profile_picture_url) {
                    setPreviewUrl(response.data.profile_picture_url);
                }
            } catch (err: any) {
                setError(err.response?.data?.message || 'Failed to fetch profile.');
                console.error('Fetch profile error:', err);
            } finally {
                setLoading(false);
            }
        };

        fetchProfile();
    }, [isAuthenticated, authLoading, router]);

    const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
        const { name, value } = e.target;
        setProfile(prev => ({ ...prev!, [name]: value }));
    };

    const handleFileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
        if (e.target.files && e.target.files[0]) {
            const selectedFile = e.target.files[0];
            setFile(selectedFile);

            // Create a preview URL for the selected image
            const reader = new FileReader();
            reader.onloadend = () => {
                setPreviewUrl(reader.result as string);
            };
            reader.readAsDataURL(selectedFile);
        } else {
            setFile(null);
            // Don't clear previewUrl if it's from the existing profile
            // Only clear if the user explicitely wants to clear the pic
        }
    };

    const handleClearPicture = () => {
        setFile(null);
        setPreviewUrl(null);
        setProfile(prev => ({ ...prev!, profile_picture_url: undefined })); // Clear URL in state
    };


    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setError(null);
        setSuccessMessage(null);

        const formData = new FormData();
        // Append text fields
        formData.append('first_name', profile?.first_name || '');
        formData.append('last_name', profile?.last_name || '');
        formData.append('phone', profile?.phone || '');
        formData.append('bio', profile?.bio || '');
        formData.append('location', profile?.location || '');
        formData.append('website', profile?.website || '');

        // Append file if selected
        if (file) {
            formData.append('profile_picture', file);
        } else if (previewUrl === null && profile?.profile_picture_url) {
            // If previewUrl is null and profile already had a picture,
            // it means the user explicitly cleared it.
            formData.append('clear_profile_picture', 'true');
        }


        try {
            const response = await api.put('/my-profile', formData, {
                headers: {
                    'Content-Type': 'multipart/form-data', // IMPORTANT!
                },
            });
            const data = response.data as { message: string; user: UserProfile };
            setSuccessMessage(data.message);
            // Update the user context in AuthContext if needed
            // This would involve passing `response.data.user` back to login or a setUser function
            // (AuthContext currently only calls login during initial login, you might need a `updateUser` function)
            // Example: updateUser(response.data.user);
            setProfile(data.user); // Update local state with fresh data from backend
            setPreviewUrl(data.user.profile_picture_url || null); // Update preview with new URL
            setFile(null); // Clear file input after successful upload
        } catch (err: any) {
            setError(err.response?.data?.message || 'Failed to update profile.');
            console.error('Update profile error:', err);
        }
    };

    if (loading || authLoading) {
        return <div className="text-center py-10">Loading profile...</div>;
    }

    if (error) {
        return <div className="text-center py-10 text-red-500">Error: {error}</div>;
    }

    if (!profile) {
        return <div className="text-center py-10">No profile data available.</div>;
    }

    return (
        <div className="container mx-auto p-4 max-w-2xl bg-gray-700 text-white rounded-lg shadow-md mt-10">
            <h1 className="text-3xl font-bold mb-6 text-center text-blue-300">My Profile</h1>
            {successMessage && <div className="bg-green-500 text-white p-3 rounded mb-4 text-center">{successMessage}</div>}
            {error && <div className="bg-red-500 text-white p-3 rounded mb-4 text-center">{error}</div>}

            <form onSubmit={handleSubmit} className="space-y-4">
                {/* Profile Picture Section */}
                <div className="flex flex-col items-center mb-6">
                    {previewUrl ? (
                        <img
                            src={previewUrl}
                            alt="Profile Preview"
                            className="w-32 h-32 rounded-full object-cover mb-4 border-2 border-blue-400"
                        />
                    ) : (
                        <div className="w-32 h-32 rounded-full bg-gray-600 flex items-center justify-center text-gray-400 text-6xl mb-4">
                            ?
                        </div>
                    )}
                    <label htmlFor="profile_picture_input" className="cursor-pointer bg-blue-500 hover:bg-blue-600 text-white font-bold py-2 px-4 rounded-lg transition duration-300">
                        Upload Profile Picture
                        <input
                            type="file"
                            id="profile_picture_input"
                            name="profile_picture"
                            accept="image/*"
                            className="hidden"
                            onChange={handleFileChange}
                        />
                    </label>
                    {previewUrl && (
                        <button
                            type="button"
                            onClick={handleClearPicture}
                            className="mt-2 text-red-400 hover:text-red-500 text-sm"
                        >
                            Clear Picture
                        </button>
                    )}
                </div>

                <div>
                    <label htmlFor="first_name" className="block text-gray-300 text-sm font-bold mb-2">First Name:</label>
                    <input
                        type="text"
                        id="first_name"
                        name="first_name"
                        value={profile.first_name || ''}
                        onChange={handleChange}
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline bg-gray-800 border-gray-600"
                    />
                </div>
                <div>
                    <label htmlFor="last_name" className="block text-gray-300 text-sm font-bold mb-2">Last Name:</label>
                    <input
                        type="text"
                        id="last_name"
                        name="last_name"
                        value={profile.last_name || ''}
                        onChange={handleChange}
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline bg-gray-800 border-gray-600"
                    />
                </div>
                <div>
                    <label htmlFor="email" className="block text-gray-300 text-sm font-bold mb-2">Email:</label>
                    <input
                        type="email"
                        id="email"
                        name="email"
                        value={profile.email || ''}
                        readOnly // Email is usually not editable
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline bg-gray-800 border-gray-600 cursor-not-allowed"
                    />
                </div>
                <div>
                    <label htmlFor="phone" className="block text-gray-300 text-sm font-bold mb-2">Phone:</label>
                    <input
                        type="text"
                        id="phone"
                        name="phone"
                        value={profile.phone || ''}
                        onChange={handleChange}
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline bg-gray-800 border-gray-600"
                    />
                </div>
                <div>
                    <label htmlFor="bio" className="block text-gray-300 text-sm font-bold mb-2">Bio:</label>
                    <textarea
                        id="bio"
                        name="bio"
                        value={profile.bio || ''}
                        onChange={handleChange}
                        rows={4}
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline bg-gray-800 border-gray-600"
                    ></textarea>
                </div>
                <div>
                    <label htmlFor="location" className="block text-gray-300 text-sm font-bold mb-2">Location:</label>
                    <input
                        type="text"
                        id="location"
                        name="location"
                        value={profile.location || ''}
                        onChange={handleChange}
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline bg-gray-800 border-gray-600"
                    />
                </div>
                <div>
                    <label htmlFor="website" className="block text-gray-300 text-sm font-bold mb-2">Website:</label>
                    <input
                        type="url"
                        id="website"
                        name="website"
                        value={profile.website || ''}
                        onChange={handleChange}
                        className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline bg-gray-800 border-gray-600"
                    />
                </div>

                <div className="flex items-center justify-between">
                    <button
                        type="submit"
                        className="bg-blue-600 hover:bg-blue-700 text-white font-bold py-2 px-6 rounded-lg focus:outline-none focus:shadow-outline transition duration-300"
                    >
                        Update Profile
                    </button>
                </div>
            </form>
        </div>
    );
}