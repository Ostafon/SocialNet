import { useState } from 'react';
import { FaEnvelope, FaLock, FaUser, FaUserPlus } from 'react-icons/fa';
import { useNavigate } from 'react-router-dom';
import api from '../../api/api';
import { logo } from '../../assets';
import './registration.scss';

export const Registration = ({ setUser }) => {
  const [email, setEmail] = useState('');
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const navigate = useNavigate();

  const handleRegister = async e => {
    e.preventDefault();
    try {
      const response = await api.post('/auth/register', {
        email,
        username,
        password
      });

      const { accessToken, refreshToken } = response.data;
      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', refreshToken);

      const profileResponse = await api.get('/auth/me');
      localStorage.setItem('userId', profileResponse.data.id);
      localStorage.setItem('isNewUser', 'true');
      setUser(profileResponse.data);

      setEmail('');
      setUsername('');
      setPassword('');

      navigate('/update-user');
    } catch (error) {
      console.error(
        'Ошибка при регистрации',
        error.response?.data || error.message
      );
      alert(error.response?.data?.message || 'Ошибка при регистрации');
    }
  };

  return (
    <section className='registration'>
      <div className='registration-container'>
        <div className='registration__left'>
          <img className='registration__image' src={logo} alt='logo' />
          <div className='registration__brand'>
            <h1>GS Social</h1>
            <p>Создайте аккаунт и начните общаться</p>
          </div>
        </div>

        <div className='registration__right'>
          <div className='registration__form-wrapper'>
            <div className='registration__header'>
              <h2>Регистрация</h2>
              <p>Это быстро и бесплатно</p>
            </div>

            <form className='registration__form' onSubmit={handleRegister}>
              <div className='input-group'>
                <FaEnvelope className='input-icon' />
                <input
                  type='email'
                  placeholder='Email'
                  value={email}
                  onChange={e => setEmail(e.target.value)}
                />
              </div>

              <div className='input-group'>
                <FaUser className='input-icon' />
                <input
                  type='text'
                  placeholder='Логин'
                  value={username}
                  onChange={e => setUsername(e.target.value)}
                />
              </div>

              <div className='input-group'>
                <FaLock className='input-icon' />
                <input
                  type='password'
                  placeholder='Пароль'
                  value={password}
                  onChange={e => setPassword(e.target.value)}
                />
              </div>

              <button type='submit' className='registration__button'>
                <FaUserPlus /> Зарегистрироваться
              </button>
            </form>
          </div>
        </div>
      </div>
    </section>
  );
};
