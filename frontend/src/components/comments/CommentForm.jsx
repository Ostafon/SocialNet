import { useState } from 'react';
import api from '../../api/api';
import './comments.scss';

export const CommentForm = ({ postId, onAddComment }) => {
  const [text, setText] = useState('');
  const [sending, setSending] = useState(false);

  const sendComment = async e => {
    e.preventDefault();

    if (!text.trim()) return;

    try {
      setSending(true);

      const response = await api.post(`posts/${postId}/comments`, {
        post_id: postId,
        content: text
      });
      onAddComment(response.data);

      setText('');
    } catch (err) {
      console.error('Ошибка при отправке комментария:', err);
    } finally {
      setSending(false);
    }
  };

  return (
    <form className='comment-form' onSubmit={sendComment}>
      <textarea
        className='comment-form__textarea'
        value={text}
        onChange={e => setText(e.target.value)}
        placeholder='Напишите комментарий...'
      />

      <button
        type='submit'
        className='comment-form__submit'
        disabled={sending || !text.trim()}
      >
        {sending ? 'Отправка...' : 'Отправить'}
      </button>
    </form>
  );
};
