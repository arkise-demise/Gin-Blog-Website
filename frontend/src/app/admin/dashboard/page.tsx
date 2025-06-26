// app/admin/dashboard/page.tsx
'use client';

import { useEffect, useState } from 'react';
import api from '../../../../utils/axios-config';
import { useRouter } from 'next/navigation';
import Link from 'next/link';
import { useAuth } from '../../context/AuthContext';

interface User {
  id: number;
  first_name: string;
  last_name: string;
}

interface BlogPost {
  id: number;
  title: string;
  description: string;
  image: string; 
  user: User; 
  created_at: string; 
  is_approved: boolean; 
}

export default function AdminDashboard() {
  const [posts, setPosts] = useState<BlogPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();
  const { isAdmin, loading: authLoading } = useAuth();

  useEffect(() => {
    if (authLoading) {
      return;
    }

    if (!isAdmin) {
      router.push('/login');
      return;
    }

    const fetchAllPosts = async () => {
      try {
        const response = await api.get('/admin/posts');
        const data = response.data as { data: BlogPost[] };
        setPosts(data.data);
        setError(null);
      } catch (err: any) {
        console.error('Failed to fetch all posts for admin dashboard:', err);
        setError(err.response?.data?.message || 'Failed to load all posts due to an unexpected error.');
        if (err.response?.status === 401 || err.response?.status === 403) {
          router.push('/login');
        }
      } finally {
        setLoading(false);
      }
    };

    if (isAdmin) {
      fetchAllPosts();
    }
  }, [isAdmin, authLoading, router]);

  const handleApprove = async (postId: number) => {
    if (!confirm('Are you sure you want to approve this post?')) {
      return;
    }
    try {
      await api.put(`/admin/posts/${postId}/approve`);
      // Update the state to reflect the change immediately in the UI
      setPosts((prevPosts) =>
        prevPosts.map((post) =>
          post.id === postId ? { ...post, is_approved: true } : post
        )
      );
      alert('Post approved successfully!');
    } catch (err: any) {
      console.error('Failed to approve post:', err);
      alert(err.response?.data?.message || 'Failed to approve post.');
      if (err.response?.status === 401 || err.response?.status === 403) {
        router.push('/login');
      }
    }
  };

  const handleReject = async (postId: number) => {
    if (!confirm('Are you sure you want to reject and delete this post? This action cannot be undone.')) {
      return;
    }
    try {
      // Use PUT for reject as per your backend route definition
      await api.put(`/admin/posts/${postId}/reject`);
      // Filter out the rejected/deleted post from the list
      setPosts((prevPosts) => prevPosts.filter((post) => post.id !== postId));
      alert('Post rejected and deleted successfully!');
    } catch (err: any) {
      console.error('Failed to reject post:', err);
      alert(err.response?.data?.message || 'Failed to reject post.');
      if (err.response?.status === 401 || err.response?.status === 403) {
        router.push('/login');
      }
    }
  };

  const handleDelete = async (postId: number) => {
    if (!confirm('Are you absolutely sure you want to delete this post? This action cannot be undone.')) {
      return;
    }
    try {
      // The backend route for admin deleting any post is /api/admin/posts/:id
      await api.delete(`/admin/posts/${postId}`);
      setPosts((prevPosts) => prevPosts.filter((post) => post.id !== postId));
      alert('Post deleted successfully!');
    } catch (err: any) {
      console.error('Failed to delete post:', err);
      alert(err.response?.data?.message || 'Failed to delete post.');
      if (err.response?.status === 401 || err.response?.status === 403) {
        router.push('/login');
      }
    }
  };

  if (authLoading || loading) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-100">
        <div className="text-center text-gray-600 text-lg">Loading admin dashboard...</div>
      </div>
    );
  }

  if (!isAdmin) {
    return (
      <div className="flex justify-center items-center min-h-screen bg-gray-100">
        <div className="text-center text-red-500 text-lg">Access Denied. You are not authorized to view this page.</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="text-red-500 text-center text-lg mt-8 min-h-screen bg-gray-100 p-4">
        {error}
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <h1 className="text-4xl font-bold text-center text-gray-800 mb-8">Admin Dashboard - All Posts</h1>

      {posts.length === 0 ? (
        <p className="text-center text-gray-600 text-lg">No posts found in the system. Perhaps you need to create some!</p>
      ) : (
        <div className="max-w-7xl mx-auto overflow-x-auto bg-white rounded-lg shadow-md p-6">
          <table className="min-w-full table-auto">
            <thead>
              {/* REMOVED WHITESPACE: All <th> tags are now on the same line as <tr> */}
              <tr className="bg-gray-200 text-gray-700 uppercase text-sm leading-normal"><th className="py-3 px-6 text-left">ID</th><th className="py-3 px-6 text-left">Title</th><th className="py-3 px-6 text-left">Description</th><th className="py-3 px-6 text-left">Author</th><th className="py-3 px-6 text-left">Created At</th><th className="py-3 px-6 text-left">Status</th><th className="py-3 px-6 text-center">Actions</th></tr>
            </thead>
            <tbody className="text-gray-600 text-sm font-light">
              {posts.map((post) => (
                // REMOVED WHITESPACE: All <td> tags are now on the same line as <tr>
                <tr key={post.id} className="border-b border-gray-200 hover:bg-gray-100"><td className="py-3 px-6 text-left whitespace-nowrap">{post.id}</td><td className="py-3 px-6 text-left font-medium max-w-xs overflow-hidden text-ellipsis whitespace-nowrap">
                    {post.title}
                  </td><td className="py-3 px-6 text-left max-w-xs overflow-hidden text-ellipsis whitespace-nowrap">
                    {post.description}
                  </td><td className="py-3 px-6 text-left">
                    {post.user?.first_name} {post.user?.last_name}
                  </td><td className="py-3 px-6 text-left">
                    {post.created_at ? new Date(post.created_at).toLocaleDateString() : 'N/A'}
                  </td><td className="py-3 px-6 text-left">
                    <span className={`px-2 py-1 rounded-full text-xs font-semibold ${
                      post.is_approved ? 'bg-green-200 text-green-800' : 'bg-yellow-200 text-yellow-800'
                    }`}>
                      {post.is_approved ? 'Approved' : 'Pending'}
                    </span>
                  </td><td className="py-3 px-6 text-center">
                    <div className="flex items-center justify-center space-x-3">
                      {!post.is_approved && ( // Only show Approve/Reject if the post is pending
                        <>
                          <button
                            onClick={() => handleApprove(post.id)}
                            className="bg-green-500 text-white px-3 py-1 rounded-md hover:bg-green-600 transition duration-300"
                          >
                            Approve
                          </button>
                          <button
                            onClick={() => handleReject(post.id)}
                            className="bg-orange-500 text-white px-3 py-1 rounded-md hover:bg-orange-600 transition duration-300"
                          >
                            Reject
                          </button>
                        </>
                      )}
                      <Link href={`/posts/${post.id}/edit`} passHref>
                        <button className="bg-blue-500 text-white px-3 py-1 rounded-md hover:bg-blue-600 transition duration-300">
                          Edit
                        </button>
                      </Link>
                      <button
                        onClick={() => handleDelete(post.id)}
                        className="bg-red-500 text-white px-3 py-1 rounded-md hover:bg-red-600 transition duration-300"
                      >
                        Delete
                      </button>
                    </div>
                  </td></tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}