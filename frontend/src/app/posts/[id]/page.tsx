// app/posts/[id]/page.tsx
// This remains a Server Component for initial data fetching

import api from '../../../../utils/axios-config';
import Image from 'next/image';
import Link from 'next/link';
import CommentSection from '../../../components/CommentSection';

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
}

interface Comment {
  id: number;
  content: string;
  user_id: number;
  blog_id: number;
  created_at: string;
  updated_at: string;
  user: User;
}

interface PostPageProps {
  params: { id: string };
}

export default async function PostPage({ params }: PostPageProps) {
  // --- FIX 2: Await params and access id safely ---
  // Although params usually comes synchronously, Next.js recommends awaiting.
  // Destructure and access properties on a resolved object.
  const { id: postId } = await params;
  // --- END FIX 2 ---

  let post: BlogPost | null = null;
  let comments: Comment[] = [];
  let error: string | null = null;

  try {
    // --- FIX 1: Remove the redundant /api/ prefix from api.get calls ---
    // The baseURL in axios-config.ts already provides '/api'
    const postResponse = await api.get(`/allpost/${postId}`); // Changed from /api/allpost
    post = (postResponse.data as { data: BlogPost }).data;

    const commentsResponse = await api.get(`/posts/${postId}/comments`); // Changed from /api/posts
    comments = (commentsResponse.data as { data: Comment[] }).data;
    // --- END FIX 1 ---

  } catch (err: any) {
    console.error(`Failed to fetch post ${postId} or comments:`, err);
    if (err.response?.status === 404) {
      error = "Post not found.";
    } else {
      error = err.response?.data?.message || 'Failed to load post due to an unexpected error.';
    }
  }

  if (error) {
    return (
      <div className="min-h-screen bg-gray-100 p-4 flex items-center justify-center">
        <p className="text-red-500 text-lg">{error}</p>
      </div>
    );
  }

  if (!post) {
    return (
      <div className="min-h-screen bg-gray-100 p-4 flex items-center justify-center">
        <p className="text-gray-600 text-lg">Loading post data or post not available.</p>
      </div>
    );
  }

  const postDate = post.created_at ? new Date(post.created_at).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
  }) : 'Date N/A';

  return (
    <div className="min-h-screen bg-gray-100 p-8 flex justify-center">
      <div className="max-w-4xl w-full bg-white rounded-lg shadow-lg overflow-hidden">
        {post.image && (
          <div className="relative w-full h-80">
            <Image
              src={post.image}
              alt={post.title}
              fill
              style={{ objectFit: 'cover' }}
              className="rounded-t-lg"
              priority
            />
          </div>
        )}
        <div className="p-8">
          <h1 className="text-4xl font-extrabold text-gray-900 mb-4 leading-tight">{post.title}</h1>
          <p className="text-gray-600 text-sm mb-6">
            By{' '}
            <span className="font-semibold text-blue-600">
              {post.user?.first_name} {post.user?.last_name}
            </span>{' '}
            on {postDate}
          </p>
          <div className="prose prose-lg max-w-none text-gray-800 leading-relaxed break-words whitespace-pre-wrap">
            {post.description}
          </div>

          <hr className="my-8 border-gray-300" />

          <CommentSection initialComments={comments} postId={parseInt(postId)} />

        </div>
      </div>
    </div>
  );
}