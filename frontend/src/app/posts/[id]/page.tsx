// app/posts/[id]/page.tsx
import api from '../../../../utils/axios-config';
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

export default async function PostPage({ params }: { params: { id: string } }) {
  const { id } = params;
  let post: BlogPost | null = null;
  let error: string | null = null;

  try {
    const response = await api.get(`/allpost/${id}`);
    const data = response.data as { data: BlogPost };
    post = data.data;
  } catch (err: any) {
    console.error(`Failed to fetch post with ID ${id}:`, err);
    error = err.response?.data?.message || 'Failed to load post.';
  }

  if (error) {
    return <div className="text-red-500 text-center text-lg mt-8">{error}</div>;
  }

  if (!post) {
    return <div className="text-gray-600 text-center text-lg mt-8">Post not found.</div>;
  }

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <div className="bg-white rounded-lg shadow-md overflow-hidden max-w-4xl mx-auto p-8">
        {post.image && (
          <img
            src={post.image}
            alt={post.title}
            className="w-full h-96 object-cover rounded-lg mb-6"
          />
        )}
        <h1 className="text-4xl font-bold text-gray-800 mb-4">{post.title}</h1>
        <p className="text-gray-700 text-lg mb-6 leading-relaxed">{post.description}</p>
        <p className="text-md text-gray-500">
          Posted by: {post.user?.first_name} {post.user?.last_name}
        </p>
      </div>
    </div>
  );
}