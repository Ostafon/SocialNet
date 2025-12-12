import { useEffect, useState } from 'react';
import {
  FaBirthdayCake,
  FaEdit,
  FaEnvelope,
  FaGlobe,
  FaInfoCircle,
  FaUserPlus
} from 'react-icons/fa';
import { useNavigate, useParams } from 'react-router-dom';
import api from '../../api/api';
import { PostCard } from '../../components/postCard/postCard';
import { Sidebar } from '../../components/sidebar/sidebar';
import { useChatStore } from '../../store/useChatStore';
import './profile.scss';

export const Profile = () => {
  const navigate = useNavigate();
  const { id } = useParams(); // ⬅️ ВАЖНО: получаем ID из URL

  const loggedId = String(localStorage.getItem('userId'));
  const profileId = id ?? loggedId; // если нет id → показываем свой профиль

  const [user, setUser] = useState(null);
  const [posts, setPosts] = useState([]);
  const [loading, setLoading] = useState(true);
  const [isFollowing, setIsFollowing] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);

  const isOwner = profileId === loggedId;

  // ==========================
  // 📌 Загрузка данных пользователя
  // ==========================
  useEffect(() => {
    const fetchUser = async () => {
      try {
        const res = await api.get(`/users/${profileId}`);
        setUser(res.data);

        if (!isOwner) {
          setIsFollowing(res.data?.isFollowing || false);
        }
      } catch (err) {
        console.error('Ошибка загрузки профиля:', err);
      } finally {
        setLoading(false);
      }
    };

    fetchUser();
  }, [profileId, isOwner]);

  // ==========================
  // 📌 Загрузка постов
  // ==========================
  useEffect(() => {
    const fetchPosts = async () => {
      try {
        const res = await api.get(`/users/${profileId}/posts`);

        const postsArray = Array.isArray(res.data)
          ? res.data
          : res.data?.posts || res.data?.data || [];

        const formatted = postsArray
          .map(p => ({
            id: p.id || p.ID,
            userId: p.userId || p.UserId,
            content: p.content || p.Content,
            imageUrl: p.imageUrl || p.ImageUrl,
            createdAt: p.createdAt || p.CreatedAt,
            likesCount: p.likesCount || p.LikesCount || 0,
            commentsCount: p.commentsCount || p.CommentsCount || 0
          }))
          .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt));

        setPosts(formatted);
      } catch (err) {
        console.error('Ошибка загрузки постов:', err);
      }
    };

    fetchPosts();
  }, [profileId]);

  // ==========================
  // 👥 Follow / Unfollow
  // ==========================
  const handleFollow = async () => {
    try {
      setFollowLoading(true);

      let res;

      if (isFollowing) {
        // 🔻 Unfollow
        res = await api.delete(`/users/${profileId}/follow`);
        console.log('UNFOLLOW RESPONSE:', res.data);
        setIsFollowing(false);
      } else {
        // 🔺 Follow
        res = await api.post(`/users/${profileId}/follow`, {});
        console.log('FOLLOW RESPONSE:', res.data);

        setIsFollowing(true);
      }
    } catch (err) {
      console.error('Ошибка follow/unfollow:', err);
    } finally {
      setFollowLoading(false);
    }
  };
  const handleOpenChat = async () => {
    try {
      const userId = localStorage.getItem('userId');

      // 1️⃣ загружаем список чатов
      const chats = useChatStore.getState().chats;

      // 2️⃣ ищем чат, где участники — вы и profileId
      const existing = chats.find(
        c =>
          c.participants.includes(userId) && c.participants.includes(profileId)
      );

      if (existing) {
        console.log('Чат уже существует → открываем:', existing.id);
        navigate(`/message/${existing.id}`);
        return;
      }

      // 3️⃣ если чата нет → создаём новый
      const res = await api.post('/chats', {
        participants: [profileId]
      });

      const chatId = res.data.id;

      // 4️⃣ переходим в чат
      navigate(`/message/${chatId}`);
    } catch (err) {
      console.error('Ошибка создания чата:', err);
    }
  };

  // ==========================
  // UI рендер
  // ==========================
  if (loading) {
    return (
      <div className='profile-container'>
        <Sidebar />
        <div className='loading-wrapper'>
          <div className='spinner'></div>
          <p>Загрузка профиля...</p>
        </div>
      </div>
    );
  }

  if (!user) {
    return (
      <div className='profile-container'>
        <Sidebar />
        <div className='error-wrapper'>
          <p>Пользователь не найден</p>
        </div>
      </div>
    );
  }

  return (
    <div className='profile-container'>
      <Sidebar />
      <div className='profile-main-content'>
        {/* Cover */}
        <div className='cover-section'>
          <div className='cover-gradient'></div>
        </div>

        {/* Header */}
        <div className='profile-header'>
          <div className='avatar-wrapper'>
            <div className='avatar'>
              {user.avatarUrl || user.avatar || user.base64Avatar ? (
                <img
                  src={
                    user.avatarUrl ||
                    user.avatar ||
                    `data:image/jpeg;base64,${user.base64Avatar}`
                  }
                  alt='avatar'
                />
              ) : (
                <div className='avatar-placeholder'>
                  {user.firstName?.[0]}
                  {user.lastName?.[0]}
                </div>
              )}
            </div>
            <div className='avatar-ring'></div>
          </div>

          <div className='user-info'>
            <h1 className='user-name'>
              {user.firstName} {user.lastName}
            </h1>
            <p className='user-username'>@{user.username}</p>
          </div>

          <div className='action-buttons'>
            {isOwner && (
              <button
                className='btn btn-primary'
                onClick={() => navigate('/update-user')}
              >
                <FaEdit /> Редактировать
              </button>
            )}

            {!isOwner && (
              <button
                className='btn btn-secondary'
                disabled={followLoading}
                onClick={handleFollow}
              >
                <FaUserPlus />
                {isFollowing ? 'Удалить из друзей' : 'Добавить в друзья'}
              </button>
            )}

            {!isOwner && (
              <button className='btn btn-secondary' onClick={handleOpenChat}>
                <FaEnvelope /> Сообщение
              </button>
            )}
          </div>
        </div>

        {/* Profile Content */}
        <div className='profile-grid'>
          {/* Info cards */}
          <div className='info-section'>
            <div className='card info-card'>
              <div className='card-header'>
                <FaInfoCircle className='card-icon' />
                <h3>Информация</h3>
              </div>

              <div className='card-body'>
                <div className='info-item'>
                  <FaBirthdayCake className='info-icon' />
                  <div>
                    <span className='info-label'>Дата рождения</span>
                    <p className='info-value'>{user.birthDate || '—'}</p>
                  </div>
                </div>

                <div className='info-item'>
                  <FaInfoCircle className='info-icon' />
                  <div>
                    <span className='info-label'>О себе</span>
                    <p className='info-value'>
                      {user.bio || 'Пока нет информации'}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            {/* Social links */}
            <div className='card info-card'>
              <div className='card-header'>
                <FaGlobe className='card-icon' />
                <h3>Соцсети</h3>
              </div>

              <div className='card-body'>
                {user.socialLinks?.length > 0 ? (
                  <ul className='social-links'>
                    {user.socialLinks.map((link, i) => (
                      <li key={i}>
                        <a
                          href={link}
                          target='_blank'
                          rel='noreferrer'
                          className='social-link'
                        >
                          {link}
                        </a>
                      </li>
                    ))}
                  </ul>
                ) : (
                  <p className='empty-state'>Нет соц. сетей</p>
                )}
              </div>
            </div>
          </div>

          {/* Posts section */}
          <div className='posts-section'>
            <div className='card posts-card'>
              <div className='card-header'>
                <h3>Посты</h3>
                <span className='badge'>{posts.length}</span>
              </div>

              <div className='card-body'>
                {posts.length === 0 ? (
                  <div className='empty-posts'>
                    <div className='empty-icon'>📝</div>
                    <p>Пока нет постов</p>
                  </div>
                ) : (
                  posts.map(post => (
                    <PostCard
                      key={post.id}
                      post={post}
                      currentUserId={loggedId}
                    />
                  ))
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
};
