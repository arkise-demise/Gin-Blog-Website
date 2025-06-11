// app/posts/[id]/edit/page.tsx
'use client';

import { useState, useEffect } from 'react';
import api from '../../../../../utils/axios-config';
import { useRouter } from 'next/navigation';

interface BlogPost {
  id: number;
  title: string;
  description: string;
  image: string;
}

export default function EditPostPage({ params }: { params: { id: string } }) {
  const { id } = params;
  const router = useRouter();
  const [formData, setFormData] = useState<BlogPost>({
    id: parseInt(id),
    title: '',
    description: '',
    image: '',
  });
  const [loading, setLoading] = useState(true);
  const [message, setMessage] = useState('');
  const [error, setError] = useState('');
  const [imageFile, setImageFile] = useState<File | null>(null);


  useEffect(() => {
    const fetchPost = async () => {
      try {
        const response = await api.get(`/allpost/${id}`);
        const data = response.data as { data: BlogPost };
        setFormData(data.data);
      } catch (err: any) {
        console.error(`Failed to fetch post with ID ${id}:`, err);
        setError(err.response?.data?.message || 'Failed to load post for editing.');
        // Redirect to login if unauthorized or post not found
        if (err.response?.status === 401 || err.response?.status === 404) {
          router.push('/login');
        }
      } finally {
        setLoading(false);
      }
    };
    fetchPost();
  }, [id, router]);

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleImageChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    if (e.target.files && e.target.files[0]) {
      setImageFile(e.target.files[0]);
    }
  };

  const handleImageUpload = async () => {
    if (!imageFile) return formData.image; // Return existing image if no new file

    const data = new FormData();
    data.append('image', imageFile);

    try {
      const response = await api.post('/upload', data, {
        headers: {
          'Content-Type': 'multipart/form-data',
        },
      });
      return (response.data as { url: string }).url;
    } catch (err) {
      console.error('Image upload failed:', err);
      setError('Image upload failed. Please try again.');
      return '';
    }
  };


  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage('');
    setError('');

    let imageUrl = formData.image; // Start with the existing image
    if (imageFile) {
      imageUrl = await handleImageUpload();
      if (!imageUrl) return; // Stop if image upload failed
    }

    try {
      const updatedData = { ...formData, image: imageUrl };
      await api.put(`/updatepost/${id}`, updatedData);
      setMessage('Post updated successfully!');
      router.push(`/posts/${id}`); // Redirect to the post detail page
    } catch (err: any) {
      setError(err.response?.data?.message || 'Failed to update post.');
      if (err.response?.status === 401) {
        router.push('/login'); // Redirect to login if unauthorized
      }
    }
  };

  if (loading) {
    return <div className="text-center text-gray-600 mt-8">Loading post for editing...</div>;
  }

  if (error) {
    return <div className="text-red-500 text-center text-lg mt-8">{error}</div>;
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-100">
      <div className="bg-white p-8 rounded-lg shadow-md w-full max-w-lg">
        <h2 className="text-3xl font-bold text-center text-gray-800 mb-6">Edit Post</h2>
        <form onSubmit={handleSubmit} className="space-y-4">
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="title">
              Title:
            </label>
            <input
              type="text"
              id="title"
              name="title"
              value={formData.title}
              onChange={handleChange}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              required
            />
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="description">
              Description:
            </label>
            <textarea
              id="description"
              name="description"
              value={formData.description}
              onChange={handleChange}
              rows={5}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
              required
            ></textarea>
          </div>
          <div>
            <label className="block text-gray-700 text-sm font-bold mb-2" htmlFor="image">
              Image:
            </label>
            {formData.image && (
                <img src={formData.image} alt="Current Image" className="w-32 h-32 object-cover mb-2 rounded" />
            )}
            <input
              type="file"
              id="image"
              name="image"
              accept="image/*"
              onChange={handleImageChange}
              className="shadow appearance-none border rounded w-full py-2 px-3 text-gray-700 leading-tight focus:outline-none focus:shadow-outline"
            />
          </div>
          <button
            type="submit"
            className="bg-blue-500 hover:bg-blue-700 text-white font-bold py-2 px-4 rounded focus:outline-none focus:shadow-outline w-full"
          >
            Update Post
          </button>
        </form>
        {message && <p className="text-green-500 text-center mt-4">{message}</p>}
        {error && <p className="text-red-500 text-center mt-4">{error}</p>}
      </div>
    </div>
  );
}