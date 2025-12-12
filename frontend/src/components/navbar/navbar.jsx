import { useEffect, useState } from 'react';
import { BiHome, BiMessageRounded } from 'react-icons/bi';
import { CgProfile } from 'react-icons/cg';
import { FiLogOut } from 'react-icons/fi';
import { IoIosNotifications, IoIosSearch } from 'react-icons/io';
import { Link, useNavigate } from 'react-router-dom';

import { SearchAPI } from '../../api/searchApi';
import { logo } from '../../assets';

import { useNotificationStore } from '../../notification/store';

import './navbar.scss';

export const Navbar = ({ setUser }) => {
  const navigate = useNavigate();

  // ================================
  // 🔥 SEARCH STATE
  // ================================
  const [query, setQuery] = useState('');
  const [searchResults, setSearchResults] = useState({ users: [], posts: [] });
  const [isOpen, setIsOpen] = useState(false);

  // ================================
  // 🔥 NOTIFICATION STORE (Zustand)
  // ================================
  const notifications = useNotificationStore(state => state.notifications);
  const unreadCount = useNotificationStore(state => state.unreadCount);

  const loadNotifications = useNotificationStore(
    state => state.loadNotifications
  );
  const startStream = useNotificationStore(state => state.startStream);
  const stopStream = useNotificationStore(state => state.stopStream);

  const markAsRead = useNotificationStore(state => state.markAsRead);
  const deleteNotification = useNotificationStore(
    state => state.deleteNotification
  );
  const markAllAsRead = useNotificationStore(state => state.markAllAsRead);
  const clearAll = useNotificationStore(state => state.clearAll);

  const [notifOpen, setNotifOpen] = useState(false);

  // ================================
  // 🔍 SEARCH LOGIC
  // ================================
  useEffect(() => {
    if (!query.trim()) {
      setSearchResults({ users: [], posts: [] });
      setIsOpen(false);
      return;
    }

    const timer = setTimeout(async () => {
      try {
        const res = await SearchAPI.all(query, 20);
        setSearchResults({
          users: res.data.users || [],
          posts: res.data.posts || []
        });
        setIsOpen(true);
      } catch (err) {
        console.error('Ошибка поиска:', err);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [query]);

  // ================================
  // 🔔 NOTIFICATIONS: LOAD + STREAM
  // ================================

  useEffect(() => {
    loadNotifications();
    startStream();

    return () => {
      stopStream();
    };
  }, []);

  const toggleNotifications = () => {
    setNotifOpen(prev => !prev);
  };

  // ================================
  // 🚪 LOGOUT
  // ================================
  const handleLogout = () => {
    stopStream();
    localStorage.removeItem('accessToken');
    localStorage.removeItem('refreshToken');
    localStorage.removeItem('userId');
    setUser(null);
    navigate('/login');
  };

  // ================================
  // RENDER
  // ================================
  return (
    <nav className='modern-navbar'>
      <div className='navbar-container'>
        {/* LOGO */}
        <Link to='/' className='navbar-brand'>
          <img className='navbar-logo' src={logo} alt='GS Social' />
          <span className='navbar-title'>GS Social</span>
        </Link>

        {/* ============================== */}
        {/* SEARCH */}
        {/* ============================== */}
        <div className='navbar-search'>
          <IoIosSearch className='search-icon' />

          <input
            type='text'
            placeholder='Поиск...'
            className='search-input'
            value={query}
            onChange={e => setQuery(e.target.value)}
            onFocus={() => {
              if (searchResults.users.length || searchResults.posts.length) {
                setIsOpen(true);
              }
            }}
          />

          {isOpen && (
            <div className='search-dropdown animated-dropdown'>
              <div className='dropdown-header'>Результаты поиска</div>

              {/* USERS */}
              {searchResults.users.length > 0 && (
                <div className='dropdown-section'>
                  <div className='section-title'>Пользователи</div>
                  {searchResults.users.map(user => (
                    <Link
                      key={user.id}
                      to={`/profile/${user.id}`}
                      className='dropdown-item'
                      onClick={() => setIsOpen(false)}
                    >
                      <img
                        src={
                          user.avatar ||
                          'https://cdn-icons-png.flaticon.com/512/3177/3177440.png'
                        }
                        className='item-avatar'
                        alt=''
                      />
                      <div className='item-info'>
                        <span className='item-title'>
                          {user.firstName} {user.lastName}
                        </span>
                        <span className='item-sub'>Профиль</span>
                      </div>
                    </Link>
                  ))}
                </div>
              )}

              {/* POSTS */}
              {searchResults.posts.length > 0 && (
                <div className='dropdown-section'>
                  <div className='section-title'>Посты</div>
                  {searchResults.posts.map(post => (
                    <Link
                      key={post.id}
                      to={`/posts/${post.id}`}
                      className='dropdown-item'
                      onClick={() => setIsOpen(false)}
                    >
                      <img
                        src={post.imageUrl}
                        className='item-avatar post-thumb'
                        alt=''
                      />
                      <div className='item-info'>
                        <span className='item-title'>{post.title}</span>
                        <span className='item-sub'>Открыть пост</span>
                      </div>
                    </Link>
                  ))}
                </div>
              )}

              {/* EMPTY */}
              {searchResults.users.length === 0 &&
                searchResults.posts.length === 0 && (
                  <div className='dropdown-empty'>Ничего не найдено</div>
                )}
            </div>
          )}
        </div>

        {/* ============================== */}
        {/* RIGHT SIDE ICONS */}
        {/* ============================== */}
        <div className='navbar-actions'>
          <Link to='/' className='nav-icon-link'>
            <BiHome className='nav-icon' />
          </Link>

          <Link to='/message' className='nav-icon-link'>
            <BiMessageRounded className='nav-icon' />
          </Link>

          {/* ============================== */}
          {/* 🔔 NOTIFICATION BUTTON */}
          {/* ============================== */}
          <div
            className='nav-icon-link notifications'
            onClick={toggleNotifications}
          >
            <IoIosNotifications className='nav-icon' />
            {unreadCount > 0 && (
              <span className='notification-badge'>{unreadCount}</span>
            )}
          </div>

          {/* ============================== */}
          {/* 🔔 NOTIFICATION DROPDOWN */}
          {/* ============================== */}
          {notifOpen && (
            <div className='notif-dropdown'>
              {notifications.length === 0 && (
                <div className='notif-empty'>Нет уведомлений</div>
              )}

              {notifications.length > 0 && (
                <div className='notif-controls'>
                  <button onClick={markAllAsRead}>Прочитать все</button>
                  <button className='delete' onClick={clearAll}>
                    Очистить все
                  </button>
                </div>
              )}

              {notifications.map(n => (
                <div
                  key={n.id}
                  className={`notif-item ${!n.isRead ? 'unread' : ''}`}
                >
                  <div className='notif-text'>{n.content}</div>

                  <div className='notif-actions'>
                    {!n.read && (
                      <button onClick={() => markAsRead(n.id)}>
                        Прочитать
                      </button>
                    )}
                    <button
                      className='delete'
                      onClick={() => deleteNotification(n.id)}
                    >
                      Удалить
                    </button>
                  </div>
                </div>
              ))}
            </div>
          )}

          <Link to='/profile' className='nav-icon-link'>
            <CgProfile className='nav-icon' />
          </Link>

          <button className='logout-btn' onClick={handleLogout}>
            <FiLogOut className='nav-icon' />
          </button>
        </div>
      </div>
    </nav>
  );
};
