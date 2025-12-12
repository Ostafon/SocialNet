import { useEffect, useState } from 'react';
import { useNavigate } from 'react-router-dom';
import api from '../../api/api';
import './friendList.scss';

export const FriendList = ({ onlyFollowing = false }) => {
  const [friends, setFriends] = useState([]);
  const [loading, setLoading] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    loadUser();
  }, []);

  const loadUser = async () => {
    setLoading(true);
    try {
      const response = await api.get('/users');

      console.log('USERS RESPONSE:', response.data);

      let list = [];

      // --- универсальная обработка ---
      if (Array.isArray(response.data)) {
        list = response.data;
      } else if (Array.isArray(response.data.users)) {
        list = response.data.users;
      } else if (Array.isArray(response.data.items)) {
        list = response.data.items;
      } else {
        console.warn('❗ API returned unexpected format');
      }

      if (onlyFollowing) {
        list = list.filter(user => user.isFollowing === true);
      }

      setFriends(list);
    } catch (error) {
      console.error('Error fetching users:', error);
    } finally {
      setLoading(false);
    }
  };

  if (loading) return <div>Загрузка...</div>;

  return (
    <div className='friend-list'>
      <h2 className='friend-list__title'>
        {onlyFollowing ? 'Мои подписки' : 'Все пользователи'}
      </h2>

      <ul className='friend-list__items'>
        {friends.map(friend => (
          <li key={friend.id} className='friend-list__item'>
            <img
              src={friend.avatarUrl || 'default-avatar.png'}
              alt='avatar'
              className='friend-list__avatar'
            />

            <div className='users-list__info'>
              <p className='users-list__name'>
                {friend.firstName} {friend.lastName}
              </p>

              <button
                className='users-list__profile-btn'
                onClick={() => navigate(`/profile/${friend.id}`)}
              >
                Перейти
              </button>
            </div>
          </li>
        ))}
      </ul>
    </div>
  );
};
