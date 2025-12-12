import { useState } from 'react';
import { useNavigate, useSearchParams } from 'react-router-dom';
import api from '../../api/api'; // ✅ импортируем api.js
import './resetPassword.scss';

export const ResetPassword = ({ setUser }) => {
  const [searchParams] = useSearchParams();
  const token = searchParams.get('token');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleUpdate = async e => {
    e.preventDefault();
    if (!newPassword || !confirmPassword || newPassword !== confirmPassword) {
      setError('Пароли не совпадают');
      return;
    }

    try {
      await api.post('/auth/reset-password', {
        resetToken: token,
        newPassword,
        confirmPassword
      });

      // ⚡ Удаляем старый токен и пользователя
      localStorage.removeItem('accessToken');
      localStorage.removeItem('refreshToken');
      localStorage.removeItem('userId');
      setUser(null);

      setNewPassword('');
      setConfirmPassword('');
      setError('');
      navigate('/login'); // редирект на логин
    } catch (err) {
      setError(err.response?.data?.message || 'Ошибка при отправке запроса');
    }
  };

  return (
    <section className='update-password'>
      <div className='update-password__container'>
        <div className='update-password__header'>
          <h1 className='update-password__title'>Обновление пароля</h1>
          <form className='update-password__form' onSubmit={handleUpdate}>
            <input
              type='password'
              placeholder='Новый пароль'
              value={newPassword}
              onChange={e => setNewPassword(e.target.value)}
            />
            <input
              type='password'
              placeholder='Подтвердите новый пароль'
              value={confirmPassword}
              onChange={e => setConfirmPassword(e.target.value)}
            />
            {error && <div className='forgot-password__error'>{error}</div>}
            <button type='submit'>Обновить пароль</button>
          </form>
        </div>
      </div>
    </section>
  );
};
