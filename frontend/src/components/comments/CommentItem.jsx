import { useEffect, useState } from 'react';
import { MdDeleteForever } from 'react-icons/md';
import api from '../../api/api';
import './comments.scss';

export const CommentItem = ({ comment, onDelete }) => {
  const [user, setUser] = useState(null);

  // Загружаем данные пользователя по userId
  useEffect(() => {
    const fetchUser = async () => {
      try {
        const res = await api.get(`/users/${comment.userId}`);
        setUser(res.data);
      } catch (err) {
        console.error('Ошибка загрузки пользователя комментария:', err);
      }
    };

    fetchUser();
  }, [comment.userId]);

  const handleDelete = async () => {
    try {
      await api.delete(`comments/${comment.id}`);
      onDelete(comment.id);
    } catch (err) {
      console.error('Ошибка при удалении комментария:', err);
    }
  };

  const avatar = user?.avatarUrl || '/default-avatar.png';

  return (
    <div className='comment'>
      <img src={avatar} alt='User avatar' className='comment__avatar' />

      <div className='comment__body'>
        <div className='comment__header'>
          <strong>
            {user ? `${user.firstName} ${user.lastName}` : 'Загрузка...'}
          </strong>
          <MdDeleteForever onClick={handleDelete} />
        </div>

        <p className='comment__text'>{comment.content}</p>

        <span className='comment__date'>
          {new Date(comment.createdAt).toLocaleDateString('ru-RU', {
            day: '2-digit',
            month: '2-digit',
            year: '2-digit',
            hour: '2-digit',
            minute: '2-digit'
          })}
        </span>
      </div>
    </div>
  );
};
