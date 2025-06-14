// components/CommentSection.tsx
'use client'; 

import { useState } from 'react';
import api from '../../utils/axios-config'; 
import { useRouter } from 'next/navigation'; 

interface User {
  id: number;
  first_name: string;
  last_name: string;
}

interface Comment {
  id: number;
  content: string;
  user_id: number;
  blog_id: number;
  created_at: string;
  updated_at: string;
  user: User; // User object will be preloaded by backend
}

interface CommentSectionProps {
  initialComments: Comment[];
  postId: number;
}

export default function CommentSection({ initialComments, postId }: CommentSectionProps) {
  const [comments, setComments] = useState<Comment[]>(initialComments);
  const [newCommentContent, setNewCommentContent] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [commentError, setCommentError] = useState<string | null>(null);

  const router = useRouter();

  const handleCommentSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setCommentError(null); // Clear previous errors
    setIsSubmitting(true);

    if (!newCommentContent.trim()) {
      setCommentError("Comment cannot be empty.");
      setIsSubmitting(false);
      return;
    }

    try {
      const response = await api.post(`/posts/${postId}/comments`, {
        content: newCommentContent,
      });

      const addedComment = (response.data as { comment: Comment }).comment; 

      setComments(prevComments => [...prevComments, addedComment]);
      setNewCommentContent('');
      alert('Comment added successfully!');

    } catch (err: any) {
      console.error('Failed to add comment:', err);
      if (err.response?.status === 401) {
        // User is not authenticated, redirect to login
        alert('You must be logged in to comment.');
        router.push('/login');
      } else {
        setCommentError(err.response?.data?.message || 'Failed to add comment.');
      }
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <div className="mt-8">
      <h2 className="text-2xl font-bold text-gray-900 mb-4">Comments</h2>

      {/* Display Existing Comments */}
      {comments.length === 0 ? (
        <p className="text-gray-600 mb-4">No comments yet. Be the first to comment!</p>
      ) : (
        <div className="space-y-4 mb-6">
          {comments.map((comment) => {
            const commentDate = new Date(comment.created_at).toLocaleDateString('en-US', {
              year: 'numeric',
              month: 'long',
              day: 'numeric',
              hour: '2-digit',
              minute: '2-digit',
            });
            return (
              <div key={comment.id} className="bg-gray-50 p-4 rounded-lg shadow-sm border border-gray-200">
                <p className="text-gray-800 leading-relaxed">{comment.content}</p>
                <p className="text-gray-500 text-sm mt-2">
                  — <span className="font-semibold">{comment.user.first_name} {comment.user.last_name}</span> on {commentDate}
                </p>
              </div>
            );
          })}
        </div>
      )}

      <hr className="my-6 border-gray-200" />

      {/* Comment Submission Form */}
      <h3 className="text-xl font-semibold text-gray-900 mb-4">Add a Comment</h3>
      <form onSubmit={handleCommentSubmit} className="space-y-4">
        <div>
          <label htmlFor="commentContent" className="sr-only">Your Comment</label>
          <textarea
            id="commentContent"
            className="w-full p-3 border border-gray-300 rounded-lg focus:ring-blue-500 focus:border-blue-500 shadow-sm"
            rows={4}
            placeholder="Write your comment here..."
            value={newCommentContent}
            onChange={(e) => setNewCommentContent(e.target.value)}
            required
            disabled={isSubmitting} 
          ></textarea>
        </div>
        {commentError && <p className="text-red-500 text-sm">{commentError}</p>}
        <button
          type="submit"
          className="bg-blue-600 text-white px-6 py-3 rounded-lg hover:bg-blue-700 transition duration-300 shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          disabled={isSubmitting} 
        >
          {isSubmitting ? 'Posting Comment...' : 'Post Comment'}
        </button>
      </form>
      <p className="text-gray-600 text-sm mt-2">
        You must be logged in to post a comment. If you are not logged in, you will be redirected.
      </p>
    </div>
  );
}