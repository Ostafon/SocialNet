import { useState } from 'react';
import { FaEnvelope, FaLock, FaSignInAlt } from 'react-icons/fa';
import { Link, useNavigate } from 'react-router-dom';
import api from '../../api/api';
import { logo } from '../../assets';
import './login.scss';

export const Login = ({ setUser }) => {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const navigate = useNavigate();

  const handleLogin = async e => {
    e.preventDefault();

    if (!email || !password) {
      setError('Пожалуйста, заполните все поля');
      return;
    }

    try {
      // 1️⃣ Логин через api.js
      const response = await api.post('/auth/login', { email, password });
      const { accessToken, refreshToken } = response.data;

      // 2️⃣ Сохраняем токены
      localStorage.setItem('accessToken', accessToken);
      localStorage.setItem('refreshToken', refreshToken);

      // 3️⃣ Получаем профиль
      const profileResponse = await api.get('/auth/me');
      setUser(profileResponse.data);
      localStorage.setItem('userId', profileResponse.data.id);

      // 4️⃣ Очистка формы и переход
      setEmail('');
      setPassword('');
      setError('');
      navigate('/');
    } catch (err) {
      console.error(err);
      setError(err.response?.data?.message || 'Неверный логин или пароль');
    }
  };

  return (
    <section className='login'>
      <div className='login-container'>
        <div className='login__left'>
          <img className='login__image' src={logo} alt='logo' />
          <div className='login__brand'>
            <h1>GS Social</h1>
            <p>Подключайтесь к миру</p>
          </div>
        </div>

        <div className='login__right'>
          <div className='login__form-wrapper'>
            <div className='login__header'>
              <h2>Вход в аккаунт</h2>
              <p>Добро пожаловать назад!</p>
            </div>

            <form className='login__form' onSubmit={handleLogin}>
              <div className='input-group'>
                <FaEnvelope className='input-icon' />
                <input
                  type='text'
                  placeholder='Email'
                  value={email}
                  onChange={e => setEmail(e.target.value)}
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

              {error && <div className='login__error'>{error}</div>}

              <Link className='login__forgot-password' to='/forgot-password'>
                Забыли пароль?
              </Link>

              <button type='submit' className='login__button'>
                <FaSignInAlt /> Войти
              </button>

              <div className='login__divider'>
                <span>ИЛИ</span>
              </div>

              <p className='login__register-prompt'>
                Нет аккаунта? <Link to='/registration'>Зарегистрироваться</Link>
              </p>
            </form>
          </div>
        </div>
      </div>
    </section>
  );
};
