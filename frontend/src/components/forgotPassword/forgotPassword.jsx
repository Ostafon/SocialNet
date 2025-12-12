import { useState } from 'react';
import api from '../../api/api'; // ✅ импортируем api.js
import './forgotPassword.scss';

export const ForgotPassword = () => {
  const [email, setEmail] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  const handleSubmit = async e => {
    e.preventDefault();
    if (!email) {
      setError('Пожалуйста, введите корректный email');
      return;
    }
    try {
      const response = await api.post('/auth/forgot-password', { email });

      console.log('Успех', response.data);
      setEmail('');
      setError('');
      setSuccess('Письмо для восстановления пароля отправлено!');
    } catch (error) {
      setError(error.response?.data?.message || 'Ошибка при отправке запроса');
      setSuccess('');
    }
  };

  return (
    <section className='forgot-password'>
      <div className='forgot-password__container'>
        <div className='forgot-password__header'>
          <h1 className='forgot-password__title'>Восстановление пароля</h1>
          <form className='forgot-password__form' onSubmit={handleSubmit}>
            <input
              type='text'
              placeholder='Email'
              value={email}
              onChange={e => setEmail(e.target.value)}
            />
            {error && <div className='forgot-password__error'>{error}</div>}
            {success && (
              <div className='forgot-password__success'>{success}</div>
            )}
            <button type='submit'>Отправить</button>
          </form>
        </div>
      </div>
    </section>
  );
};
