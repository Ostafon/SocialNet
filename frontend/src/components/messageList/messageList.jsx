import { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useChatStore } from '../../store/useChatStore';
import './messageList.scss';

export const MessageList = () => {
  const navigate = useNavigate();
  const { id: activeChatId } = useParams();

  const chats = useChatStore(s => s.chats);
  const loadingChats = useChatStore(s => s.loadingChats);
  const [searchQuery, setSearchQuery] = useState('');

  const hasLoadedRef = useRef(false);

  useEffect(() => {
    if (hasLoadedRef.current) return;
    hasLoadedRef.current = true;

    const loadData = async () => {
      try {
        await useChatStore.getState().loadChats();
      } catch (err) {
        console.error('Failed to load chats:', err);
      }
    };

    loadData();
  }, []);

  const filteredChats = chats.filter(chat =>
    chat.name?.toLowerCase().includes(searchQuery.toLowerCase())
  );

  if (loadingChats) {
    return (
      <div className='tg-chat-list'>
        <div className='tg-header'>
          <h2>Сообщения</h2>
        </div>
        <div className='tg-loading'>
          <div className='spinner'></div>
          <p>Загрузка чатов...</p>
        </div>
      </div>
    );
  }

  if (!chats || chats.length === 0) {
    return (
      <div className='tg-chat-list'>
        <div className='tg-header'>
          <h2>Сообщения</h2>
        </div>
        <div className='tg-empty-list'>
          <div className='empty-icon'>📭</div>
          <p>Нет чатов</p>
          <span>Начните новую беседу</span>
        </div>
      </div>
    );
  }

  return (
    <div className='tg-chat-list'>
      <div className='tg-header'>
        <h2>Сообщения</h2>
        <div className='tg-search'>
          <svg
            className='search-icon'
            width='20'
            height='20'
            viewBox='0 0 24 24'
            fill='none'
          >
            <circle
              cx='11'
              cy='11'
              r='8'
              stroke='currentColor'
              strokeWidth='2'
            />
            <path
              d='M21 21L16.65 16.65'
              stroke='currentColor'
              strokeWidth='2'
              strokeLinecap='round'
            />
          </svg>
          <input
            type='text'
            placeholder='Поиск чатов...'
            value={searchQuery}
            onChange={e => setSearchQuery(e.target.value)}
          />
        </div>
      </div>

      <div className='tg-chats'>
        {filteredChats.map(chat => (
          <div
            key={chat.id}
            className={`tg-chat-item ${
              String(activeChatId) === String(chat.id) ? 'active' : ''
            }`}
            onClick={() => navigate(`/message/${chat.id}`)}
          >
            <div className='avatar-wrapper'>
              <img
                className='tg-avatar'
                src={
                  chat.avatarUrl ||
                  'https://cdn-icons-png.flaticon.com/512/3177/3177440.png'
                }
                alt=''
              />
              <span className='online-indicator'></span>
            </div>

            <div className='tg-chat-info'>
              <div className='tg-chat-top'>
                <div className='tg-chat-name'>
                  {chat.name || `Чат #${chat.id}`}
                </div>
                {chat.lastMessageTime && (
                  <div className='tg-chat-time'>
                    {new Date(chat.lastMessageTime).toLocaleTimeString(
                      'ru-RU',
                      {
                        hour: '2-digit',
                        minute: '2-digit'
                      }
                    )}
                  </div>
                )}
              </div>
              <div className='tg-chat-last'>
                {chat.lastMessage || 'Нет сообщений'}
              </div>
            </div>

            {chat.unreadCount > 0 && (
              <div className='unread-badge'>{chat.unreadCount}</div>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};
