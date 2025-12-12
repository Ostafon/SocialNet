import { useEffect, useState } from 'react';
import api from '../../api/api';
import { CommentForm } from './CommentForm';
import { CommentItem } from './CommentItem';
import './comments.scss';

export const CommentSection = ({ postId, onCommentsChange }) => {
  const [comments, setComments] = useState([]);
  const [loading, setLoading] = useState(false);

  const fetchComments = async () => {
    try {
      setLoading(true);
      const response = await api.get(`posts/${postId}/comments`);
      console.log('COMMENTS RAW →', response.data);
      setComments(response.data.comments || []);
    } catch (err) {
      console.error('Ошибка при загрузке комментариев:', err);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchComments();
  }, [postId]);

  const handleAddComment = comment => {
    setComments(prev => [comment, ...prev]);

    if (onCommentsChange) {
      onCommentsChange(); // сообщаем PostCard что комментарий добавился
    }
  };

  const handleDeleteComment = commentId => {
    setComments(prev => prev.filter(comment => comment.id !== commentId));
  };

  return (
    <div className='comment-section'>
      <CommentForm postId={postId} onAddComment={handleAddComment} />

      {loading ? (
        <p>Загрузка комментариев...</p>
      ) : (
        <div className='comment-section__list'>
          {comments.map(comment => (
            <CommentItem
              key={comment.id}
              comment={comment}
              onDelete={handleDeleteComment}
            />
          ))}
        </div>
      )}
    </div>
  );
};
