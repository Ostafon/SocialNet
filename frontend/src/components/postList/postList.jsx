import { useEffect, useState } from 'react';
import api from '../../api/api';
import { AddPost } from '../addPost/addPost';
import { PostCard } from '../postCard/postCard';
import './postList.scss';

export const PostList = ({ currentUser }) => {
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchPosts = async () => {
    try {
      setLoading(true);
      const response = await api.get('/feed');
      console.log('response.data =', response.data); // 👈 Посмотри в консоль

      // безопасно достаём массив постов
      const postsArray = Array.isArray(response.data)
        ? response.data
        : response.data?.data || response.data?.posts || [];

      const sorted = postsArray
        .map(p => ({
          id: p.id || p.ID,
          userId: p.userId || p.UserId,
          content: p.content || p.Content,
          imageUrl: p.imageUrl || p.ImageUrl,
          createdAt: p.createdAt || p.CreatedAt,
          likesCount: p.likesCount || p.LikesCount || 0,
          commentsCount: p.commentsCount || p.CommentsCount || 0
        }))
        .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));

      setPosts(sorted);
    } catch (error) {
      console.error('Error fetching posts:', error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchPosts();
  }, []);

  const handlePostAdded = newPost => {
    setPosts(prev => [newPost, ...prev]);
  };

  const handleLikeToggle = (postId, updatedPost) => {
    setPosts(prev =>
      prev.map(post => (post.id === postId ? updatedPost : post))
    );
  };

  return (
    <section className='post-list'>
      <AddPost onPostAdded={handlePostAdded} currentUser={currentUser} />
      {loading ? (
        <p>Загрузка постов...</p>
      ) : posts.length > 0 ? (
        posts.map(post => (
          <PostCard
            key={post.id || post.ID}
            post={post}
            onLikeToggle={handleLikeToggle}
            currentUser={currentUser}
          />
        ))
      ) : (
        <p>Постов нет. Будьте первым, кто опубликует!</p>
      )}
    </section>
  );
};
