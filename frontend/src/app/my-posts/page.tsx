// app/my-posts/page.tsx
'use client';

import { useEffect, useState } from 'react';
import api from '../../../utils/axios-config';
import { useRouter } from 'next/navigation';
import Link from 'next/link';

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
}

export default function MyPosts() {
  const [posts, setPosts] = useState<BlogPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const router = useRouter();

  useEffect(() => {
    const fetchMyPosts = async () => {
      try {
        const response = await api.get('/myposts'); // Update to the correct endpoint
        const data = response.data as { data: BlogPost[] };
        setPosts(data.data);
      } catch (err: any) {
        console.error('Failed to fetch user posts:', err);
        setError(err.response?.data?.message || 'Failed to load your posts.');
        if (err.response?.status === 401) {
          router.push('/login'); // Redirect to login if unauthorized
        }
      } finally {
        setLoading(false);
      }
    };

    fetchMyPosts();
  }, [router]);

  const handleDelete = async (postId: number) => {
    if (!confirm('Are you sure you want to delete this post?')) {
      return;
    }
    try {
      await api.delete(`/deletepost/${postId}`);
      setPosts(posts.filter((post) => post.id !== postId));
      alert('Post deleted successfully!');
    } catch (err: any) {
      console.error('Failed to delete post:', err);
      alert(err.response?.data?.message || 'Failed to delete post.');
    }
  };

  if (loading) {
    return <div className="text-center text-gray-600 mt-8">Loading your posts...</div>;
  }

  if (error) {
    return <div className="text-red-500 text-center text-lg mt-8">{error}</div>;
  }

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <h1 className="text-4xl font-bold text-center text-gray-800 mb-8">Your Posts</h1>

      {posts.length === 0 ? (
        <p className="text-center text-gray-600 text-lg">You haven't created any posts yet.</p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          {posts.map((post) => (
            <div key={post.id} className="bg-white rounded-lg shadow-md overflow-hidden transform transition duration-300 hover:scale-105">
              {post.image && (
                <img
                  src={post.image}
                  alt={post.title}
                  className="w-full h-48 object-cover"
                />
              )}
              <div className="p-6">
                <h2 className="text-2xl font-semibold text-gray-800 mb-2">
                  {post.title}
                </h2>
                <p className="text-gray-600 mb-4 line-clamp-3">
                  {post.description}
                </p>
                <div className="flex justify-between items-center mt-4">
                  <Link href={`/posts/${post.id}`} className="text-blue-500 hover:underline">
                      Read More
                  </Link>
                  <button
                    onClick={() => handleDelete(post.id)}
                    className="bg-red-500 text-white px-4 py-2 rounded-lg hover:bg-red-600 transition duration-300"
                  >
                    Delete
                  </button>
                </div>
              </div>
            </div>
          ))}
        </div>
      )}
    </div>
  );
}