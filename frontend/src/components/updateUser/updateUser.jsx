import { useState } from 'react';
import {
  FaCalendar,
  FaCamera,
  FaFacebook,
  FaGithub,
  FaInstagram,
  FaLinkedin,
  FaSave,
  FaTelegram,
  FaTwitter,
  FaUser
} from 'react-icons/fa';
import { useNavigate } from 'react-router-dom';
import api from '../../api/api';
import './updateUser.scss';

export const UpdateUser = () => {
  const [avatarUrl, setAvatarUrl] = useState(null);
  const [preview, setPreview] = useState(null);
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [birthDate, setBirthDate] = useState('');
  const [bio, setBio] = useState('');
  const [facebook, setFacebook] = useState('');
  const [instagram, setInstagram] = useState('');
  const [twitter, setTwitter] = useState('');
  const [linkedin, setLinkedin] = useState('');
  const [github, setGithub] = useState('');
  const [telegram, setTelegram] = useState('');
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');
  const userId = localStorage.getItem('userId');
  const navigate = useNavigate();

  const handleAvatarChange = e => {
    const file = e.target.files[0];
    if (file) {
      setAvatarUrl(file);
      setPreview(URL.createObjectURL(file));
    }
  };

  const handleSubmit = async e => {
    e.preventDefault();

    if (!firstName || !lastName) {
      setError('Пожалуйста, заполните обязательные поля.');
      return;
    }

    try {
      // Собираем социальные сети
      const socialLinks = [];
      if (facebook) socialLinks.push(facebook);
      if (instagram) socialLinks.push(instagram);
      if (twitter) socialLinks.push(twitter);
      if (linkedin) socialLinks.push(linkedin);
      if (github) socialLinks.push(github);
      if (telegram) socialLinks.push(telegram);

      // Обновляем профиль
      const userData = {
        firstName,
        lastName,
        birthDate: birthDate
          ? new Date(birthDate).toISOString().split('T')[0]
          : null,
        bio,
        socialLinks: socialLinks.length > 0 ? socialLinks : null
      };

      await api.patch(`/users/${userId}`, userData);

      // 2️⃣ Если выбран аватар — отправляем отдельно
      if (avatarUrl) {
        const avatarBase64 = await new Promise((resolve, reject) => {
          const reader = new FileReader();
          reader.readAsDataURL(avatarUrl);
          reader.onload = () => resolve(reader.result);
          reader.onerror = err => reject(err);
        });

        const avatarData = {
          id: userId,
          avatar: avatarBase64.split(',')[1],
          filename: avatarUrl.name
        };

        await api.post(`/users/${userId}/avatar`, avatarData);
      }

      // 3️⃣ После успешного обновления убираем флаг нового пользователя
      localStorage.removeItem('isNewUser');

      // 4️⃣ Сообщаем пользователю об успехе
      setSuccess('Профиль успешно обновлён!');
      setError('');

      // 5️⃣ Переходим на страницу профиля
      navigate('/profile');
    } catch (err) {
      console.error('Ошибка при обновлении профиля:', err);
      setError('Не удалось обновить профиль.');
      setSuccess('');
    }
  };

  return (
    <div className='update-user-page'>
      <div className='update-user-container'>
        <div className='update-user-header'>
          <h1>Настройка профиля</h1>
          <p>Заполните информацию о себе</p>
        </div>

        <div className='avatar-upload-section'>
          <div className='avatar-preview-wrapper'>
            {preview ? (
              <img src={preview} alt='avatar' className='avatar-preview' />
            ) : (
              <div className='avatar-placeholder'>
                <FaUser size={50} />
              </div>
            )}
            <label htmlFor='avatar-upload' className='avatar-upload-label'>
              <FaCamera size={20} />
            </label>
            <input
              id='avatar-upload'
              type='file'
              accept='image/*'
              onChange={handleAvatarChange}
              style={{ display: 'none' }}
            />
          </div>
          <p className='avatar-hint'>Нажмите на камеру, чтобы загрузить фото</p>
        </div>

        <form onSubmit={handleSubmit} className='update-user-form'>
          <div className='form-row'>
            <div className='input-group'>
              <FaUser className='input-icon' />
              <input
                type='text'
                placeholder='Имя'
                value={firstName}
                onChange={e => setFirstName(e.target.value)}
                required
              />
            </div>
            <div className='input-group'>
              <FaUser className='input-icon' />
              <input
                type='text'
                placeholder='Фамилия'
                value={lastName}
                onChange={e => setLastName(e.target.value)}
                required
              />
            </div>
          </div>

          <div className='input-group date-input'>
            <FaCalendar className='input-icon' />
            <input
              type='date'
              value={birthDate}
              onChange={e => setBirthDate(e.target.value)}
              placeholder='Дата рождения'
            />
          </div>

          <div className='textarea-group'>
            <textarea
              placeholder='Расскажите о себе...'
              value={bio}
              onChange={e => setBio(e.target.value)}
              rows='4'
            />
          </div>

          <div className='social-section'>
            <h3 className='section-title'>Социальные сети</h3>
            <div className='social-inputs'>
              <div className='input-group social-input'>
                <FaFacebook className='input-icon social-icon facebook' />
                <input
                  type='url'
                  placeholder='Facebook URL'
                  value={facebook}
                  onChange={e => setFacebook(e.target.value)}
                />
              </div>

              <div className='input-group social-input'>
                <FaInstagram className='input-icon social-icon instagram' />
                <input
                  type='url'
                  placeholder='Instagram URL'
                  value={instagram}
                  onChange={e => setInstagram(e.target.value)}
                />
              </div>

              <div className='input-group social-input'>
                <FaTwitter className='input-icon social-icon twitter' />
                <input
                  type='url'
                  placeholder='Twitter URL'
                  value={twitter}
                  onChange={e => setTwitter(e.target.value)}
                />
              </div>

              <div className='input-group social-input'>
                <FaLinkedin className='input-icon social-icon linkedin' />
                <input
                  type='url'
                  placeholder='LinkedIn URL'
                  value={linkedin}
                  onChange={e => setLinkedin(e.target.value)}
                />
              </div>

              <div className='input-group social-input'>
                <FaGithub className='input-icon social-icon github' />
                <input
                  type='url'
                  placeholder='GitHub URL'
                  value={github}
                  onChange={e => setGithub(e.target.value)}
                />
              </div>

              <div className='input-group social-input'>
                <FaTelegram className='input-icon social-icon telegram' />
                <input
                  type='url'
                  placeholder='Telegram URL'
                  value={telegram}
                  onChange={e => setTelegram(e.target.value)}
                />
              </div>
            </div>
          </div>

          {error && <div className='message error-message'>{error}</div>}
          {success && <div className='message success-message'>{success}</div>}

          <button type='submit' className='submit-button'>
            <FaSave /> Сохранить профиль
          </button>
        </form>
      </div>
    </div>
  );
};
