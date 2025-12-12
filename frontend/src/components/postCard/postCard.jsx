import { useEffect, useState } from 'react';
import { AiFillHeart, AiOutlineHeart } from 'react-icons/ai';
import api from '../../api/api';
import { CommentSection } from '../comments/CommentSection';
import './postCard.scss';

export const PostCard = ({ post, onLikeToggle, currentUser }) => {
  const [isLiking, setIsLiking] = useState(false);
  const [author, setAuthor] = useState(null);
  const [showComments, setShowComments] = useState(false);

  // 🧠 Подгружаем данные пользователя по его ID
  useEffect(() => {
    const fetchAuthor = async () => {
      try {
        const res = await api.get(`/users/${post.userId}`);
        setAuthor(res.data);
      } catch (err) {
        console.error('Ошибка при загрузке автора поста:', err);
      }
    };

    if (post.userId) fetchAuthor();
  }, [post.userId]);

  const handleLike = async () => {
    try {
      setIsLiking(true);
      const response = await api.post(`/posts/${post.id}/like`);
      onLikeToggle(post.id, {
        ...post,
        likesCount: response.data.likesCount,
        isLiked: response.data.isLiked
      });
    } catch (err) {
      console.error('Ошибка при лайке:', err);
    } finally {
      setIsLiking(false);
    }
  };

  const avatarUrl =
    author?.avatarUrl ||
    author?.image_url ||
    (author?.base64Avatar
      ? `data:image/jpeg;base64,${author.base64Avatar}`
      : null) ||
    currentUser?.avatarUrl ||
    '/default-avatar.png';

  return (
    <article className='post-card'>
      <div className='post-card__header'>
        <img src={avatarUrl} alt='userAvatar' className='post-card__avatar' />
        <div className='post-card__info'>
          <h4 className='post-card__author'>
            {author?.firstName
              ? `${author.firstName} ${author.lastName || ''}`
              : 'Аноним'}
          </h4>
          <p className='post-card__date'>
            {new Date(post.createdAt).toLocaleDateString('ru-RU', {
              day: '2-digit',
              month: 'short',
              year: 'numeric',
              hour: '2-digit',
              minute: '2-digit'
            })}
          </p>
        </div>
      </div>

      {post.content && <p className='post-card__text'>{post.content}</p>}

      {post.imageUrl && (
        <div className='post-card__image-wrapper'>
          <img src={post.imageUrl} alt='post' className='post-card__image' />
        </div>
      )}

      <div className='post-card__actions'>
        <button
          className={`post-card__like-button ${post.isLiked ? 'liked' : ''}`}
          onClick={handleLike}
          disabled={isLiking}
        >
          {post.isLiked ? <AiFillHeart /> : <AiOutlineHeart />}
          {post.likesCount ?? 0}
        </button>
        <button
          className='post-card__comment-button'
          onClick={() => setShowComments(prev => !prev)}
        >
          💬 {post.commentsCount ?? 0}
        </button>
        {showComments && (
          <CommentSection
            postId={post.id}
            onCommentsChange={() =>
              onLikeToggle(post.id, {
                ...post,
                commentsCount: (post.commentsCount ?? 0) + 1
              })
            }
          />
        )}
      </div>
    </article>
  );
};
