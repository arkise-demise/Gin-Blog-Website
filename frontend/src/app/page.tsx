// app/page.tsx
// (No need for Link here if all navigation is in Navbar)
// import Link from 'next/link'; // You might still need Link for "Read More"
import api from '../../utils/axios-config';
import Link from 'next/link'; // Keep Link for the individual post "Read More"

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

export default async function Home() {
  let posts: BlogPost[] = [];
  let error: string | null = null;

  try {
    const response = await api.get('/allpost');
    posts = (response.data as { data: BlogPost[] }).data;
  } catch (err: any) {
    console.error('Failed to fetch posts:', err);
    error = err.response?.data?.message || 'Failed to load posts.';
  }

  return (
    <div className="min-h-screen bg-gray-100 p-4">
      <h1 className="text-4xl font-bold text-center text-indigo-600 mb-8 italic drop-shadow-lg">
        Welcome to the <span className="text-pink-500">Gin Blog</span>
      </h1>

      {error && <p className="text-red-500 text-center text-lg">{error}</p>}

      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
        {posts.length > 0 ? (
          posts.map((post) => (
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
                <p className="text-sm text-gray-500 mb-4">
                  By: {post.user?.first_name} {post.user?.last_name}
                </p>
                <Link href={`/posts/${post.id}`} className="text-blue-500 hover:underline">
                    Read More
                </Link>
              </div>
            </div>
          ))
        ) : (
          <p className="text-center text-gray-600 col-span-full">No posts found yet.</p>
        )}
      </div>
    </div>
  );
}