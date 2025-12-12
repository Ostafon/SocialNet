import { useState } from 'react';
import { FaCamera, FaImage, FaSmile } from 'react-icons/fa';
import api from '../../api/api';
import './addPost.scss';

export const AddPost = ({ onPostAdded, currentUser }) => {
  const [content, setContent] = useState('');
  const [image, setImage] = useState(null);
  const [preview, setPreview] = useState(null);
  const [loading, setLoading] = useState(false);

  const handleImageChange = e => {
    const file = e.target.files[0];
    if (file) {
      setImage(file);
      setPreview(URL.createObjectURL(file));
    }
  };

  const clearImage = () => {
    setImage(null);
    setPreview(null);
  };

  const handleSubmit = async e => {
    e.preventDefault();
    if (!content.trim() && !image) return;

    try {
      setLoading(true);

      let imageBase64 = null;
      let fileName = null;

      // 🖼️ Конвертируем фото в base64, если выбрано
      if (image) {
        imageBase64 = await new Promise((resolve, reject) => {
          const reader = new FileReader();
          reader.readAsDataURL(image);
          reader.onload = () => resolve(reader.result.split(',')[1]); // убираем "data:image/png;base64,"
          reader.onerror = err => reject(err);
        });
        fileName = image.name;
      }

      // 📦 Отправляем всё в одном запросе
      const postData = {
        content,
        image: imageBase64,
        fileName
      };

      const response = await api.post('/posts', postData);
      const createdPost = response.data;

      if (onPostAdded) onPostAdded(createdPost);
      clearImage();
      setContent('');
    } catch (error) {
      console.error('Ошибка при добавлении поста:', error);
    } finally {
      setLoading(false);
    }
  };

  return (
    <form className='add-post' onSubmit={handleSubmit}>
      <div className='add-post__top'>
        <div className='add-post__avatar'>
          {currentUser?.avatarUrl || currentUser?.avatar ? (
            <img
              src={currentUser?.avatarUrl || currentUser?.avatar}
              alt='avatar'
            />
          ) : (
            <div className='avatar-placeholder'>
              {currentUser?.firstName?.[0] || 'U'}
            </div>
          )}
        </div>
        <textarea
          value={content}
          onChange={e => setContent(e.target.value)}
          placeholder='Что у вас нового?'
          className='add-post__textarea'
          rows={3}
        />
      </div>

      {preview && (
        <div className='add-post__preview'>
          <img src={preview} alt='preview' />
          <button
            type='button'
            className='add-post__remove'
            onClick={clearImage}
          >
            ✕
          </button>
        </div>
      )}

      <div className='add-post__actions'>
        <div className='add-post__options'>
          <label className='add-post__option'>
            <FaImage className='option-icon' />
            <span>Фото</span>
            <input type='file' accept='image/*' onChange={handleImageChange} />
            <FaCamera className='option-icon' />
            <FaSmile className='option-icon' />
          </label>
        </div>
        <button type='submit' className='add-post__submit' disabled={loading}>
          {loading ? 'Публикуется...' : 'Опубликовать'}
        </button>
      </div>
    </form>
  );
};
