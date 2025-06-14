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
  user: User; // Assuming user is always present after preload
}

export default function MyPosts() {
  const [posts, setPosts] = useState<BlogPost[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  // New state for editing functionality
  const [editingPostId, setEditingPostId] = useState<number | null>(null);
  const [editedPostData, setEditedPostData] = useState<Partial<BlogPost> | null>(null);

  const router = useRouter();

  useEffect(() => {
    const fetchMyPosts = async () => {
      try {
        const response = await api.get('/uniquepost');
        const data = response.data as { data: BlogPost[] };
        setPosts(data.data);

        // Scenario 1: Backend returns 200 OK, but the data array is empty
        if (data.data.length === 0) {
          setError("You haven't created any posts yet."); // Set a specific message
        }

      } catch (err: any) {
        console.error('Failed to fetch user posts:', err);

        // Scenario 2: Backend returns 404 specifically when there are no posts for the user
        if (err.response?.status === 404) {
          setPosts([]); // Explicitly set posts to an empty array
          setError("You haven't created any posts yet."); // Set a specific message for this case
        } else if (err.response?.status === 401) {
          router.push('/login'); // Redirect to login if unauthorized
        } else {
          // General error handling for other API errors
          setError(err.response?.data?.message || 'Failed to load your posts due to an unexpected error.');
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
      // Filter out the deleted post
      const updatedPosts = posts.filter((post) => post.id !== postId);
      setPosts(updatedPosts);
      alert('Post deleted successfully!');
      // If no posts are left after deletion, update the error state to show the prompt
      if (updatedPosts.length === 0) {
         setError("You haven't created any posts yet.");
      }
    } catch (err: any) {
      console.error('Failed to delete post:', err);
      alert(err.response?.data?.message || 'Failed to delete post.');
    }
  };

  // Function to handle clicking the "Edit" button
  const handleEditClick = (post: BlogPost) => {
    setEditingPostId(post.id);
    // Initialize edited data with current post data
    setEditedPostData({
      title: post.title,
      description: post.description,
      image: post.image,
    });
  };

  // Function to handle clicking the "Cancel" button in edit mode
  const handleCancelEdit = () => {
    setEditingPostId(null); // Exit edit mode
    setEditedPostData(null); // Clear edited data
  };

  // Function to handle changes in the edit form inputs
  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    const { name, value } = e.target;
    setEditedPostData(prev => ({
      ...prev,
      [name]: value
    }));
  };

  // Function to handle submitting the updated post data
  const handleUpdate = async () => {
    if (!editingPostId || !editedPostData) return; // Should not happen if edit mode is active

    if (!confirm('Are you sure you want to update this post?')) {
        return;
    }

    try {
      // Send the PUT request to your backend's UpdatePostById handler
      const response = await api.put(`/updatepost/${editingPostId}`, editedPostData);
      // Assert the type of response.data to avoid 'unknown' error
      const data = response.data as { post: BlogPost };
      const updatedPost = data.post;

      // Update the posts state to reflect the changes immediately
      setPosts(prevPosts =>
        prevPosts.map(post => (post.id === editingPostId ? updatedPost : post))
      );
      alert('Post updated successfully!');
      setEditingPostId(null); // Exit edit mode
      setEditedPostData(null); // Clear edited data
    } catch (err: any) {
      console.error('Failed to update post:', err);
      // Display backend's error message or a generic one
      alert(err.response?.data?.message || 'Failed to update post.');
    }
  };


  if (loading) {
    return (
      <div className="flex justify-center items-center min-h-screen">
        <div className="text-center text-gray-600 text-lg">Loading your posts...</div>
      </div>
    );
  }

  // --- Render Logic for No Posts or Specific Errors ---
  // This block handles both 404 from backend and 200 with empty array
  if (error && posts.length === 0) {
    return (
      <div className="min-h-screen bg-gray-100 p-4 flex flex-col items-center justify-center">
        <p className="text-gray-600 text-lg mb-4">{error}</p>
        <Link href="/create-post" passHref>
          <button className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition duration-300 shadow-md">
            Create Your First Post
          </button>
        </Link>
      </div>
    );
  }

  // General error display (for errors other than "no posts found")
  if (error) {
     return (
        <div className="text-red-500 text-center text-lg mt-8 min-h-screen">
          {error}
        </div>
     );
  }

  // --- Main Render Logic for Displaying Posts ---
  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <h1 className="text-4xl font-bold text-center text-gray-800 mb-8">Your Posts</h1>

      {posts.length === 0 ? (
        // This message will only show if 'error' is null but posts are still empty (e.g., initial load, then logout)
        // However, given the 'if (error && posts.length === 0)' block above, this might rarely be hit.
        <p className="text-center text-gray-600 text-lg">
            No posts found. Create your first post!
        </p>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8 max-w-7xl mx-auto">
          {posts.map((post) => (
            <div key={post.id} className="bg-white rounded-lg shadow-md overflow-hidden transform transition duration-300 hover:scale-105">
              {editingPostId === post.id ? (
                // --- RENDER EDIT FORM WHEN IN EDIT MODE ---
                <div className="p-6">
                  <h2 className="text-2xl font-semibold text-gray-800 mb-4">Edit Post</h2>
                  <div className="mb-4">
                    <label htmlFor={`title-${post.id}`} className="block text-gray-700 text-sm font-bold mb-2">Title:</label>
                    <input
                      type="text"
                      id={`title-${post.id}`}
                      name="title"
                      value={editedPostData?.title || ''}
                      onChange={handleChange}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                    />
                  </div>
                  <div className="mb-4">
                    <label htmlFor={`description-${post.id}`} className="block text-gray-700 text-sm font-bold mb-2">Description:</label>
                    <textarea
                      id={`description-${post.id}`}
                      name="description"
                      value={editedPostData?.description || ''}
                      onChange={handleChange}
                      rows={5}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                    ></textarea>
                  </div>
                  <div className="mb-4">
                    <label htmlFor={`image-${post.id}`} className="block text-gray-700 text-sm font-bold mb-2">Image URL:</label>
                    <input
                      type="text"
                      id={`image-${post.id}`}
                      name="image"
                      value={editedPostData?.image || ''}
                      onChange={handleChange}
                      className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
                    />
                  </div>
                  <div className="flex justify-end gap-2">
                    <button
                      onClick={handleCancelEdit}
                      className="bg-gray-500 text-white px-4 py-2 rounded-lg hover:bg-gray-600 transition duration-300"
                    >
                      Cancel
                    </button>
                    <button
                      onClick={handleUpdate}
                      className="bg-green-500 text-white px-4 py-2 rounded-lg hover:bg-green-600 transition duration-300"
                    >
                      Save Changes
                    </button>
                  </div>
                </div>
              ) : (
                // --- RENDER POST DISPLAY CARD IN NORMAL MODE ---
                <>
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
                      <div className="flex gap-2"> {/* Group buttons */}
                        <button
                          onClick={() => handleEditClick(post)}
                          className="bg-yellow-500 text-white px-4 py-2 rounded-lg hover:bg-yellow-600 transition duration-300"
                        >
                          Edit
                        </button>
                        <button
                          onClick={() => handleDelete(post.id)}
                          className="bg-red-500 text-white px-4 py-2 rounded-lg hover:bg-red-600 transition duration-300"
                        >
                          Delete
                        </button>
                      </div>
                    </div>
                  </div>
                </>
              )}
            </div>
          ))}
        </div>
      )}
    </div>
  );
}