import { useEffect, useRef, useState } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { useChatStore } from '../../store/useChatStore';
import './chatPage.scss';

export const ChatPage = () => {
  const { id } = useParams();
  const navigate = useNavigate();
  const chatId = String(id);

  const messages = useChatStore(s => s.messages[chatId]) || [];
  const chats = useChatStore(s => s.chats);
  const currentChat = chats.find(c => String(c.id) === chatId);

  const messagesRef = useRef(null);
  const inputRef = useRef(null);
  const [text, setText] = useState('');
  const [isTyping, setIsTyping] = useState(false);

  useEffect(() => {
    if (!chatId) return;

    const state = useChatStore.getState();

    state.setActiveChat(chatId); // Включает стрим внутри useChatStore
    state.loadMessages(chatId); // Загружает историю

    return () => {
      state.stopAllStreams(); // Останавливает при выходе
    };
  }, [chatId]);

  useEffect(() => {
    if (messagesRef.current) {
      messagesRef.current.scrollTop = messagesRef.current.scrollHeight;
    }
  }, [messages]);

  const handleSend = () => {
    if (!text.trim()) return;

    const state = useChatStore.getState();
    state.sendMessage(chatId, text);
    setText('');
    inputRef.current?.focus();
  };

  const handleKeyDown = e => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault();
      handleSend();
    }
  };

  const formatTime = timestamp => {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    return date.toLocaleTimeString('ru-RU', {
      hour: '2-digit',
      minute: '2-digit'
    });
  };

  const formatDate = timestamp => {
    if (!timestamp) return '';
    const date = new Date(timestamp);
    const today = new Date();
    const yesterday = new Date(today);
    yesterday.setDate(yesterday.getDate() - 1);

    if (date.toDateString() === today.toDateString()) {
      return 'Сегодня';
    } else if (date.toDateString() === yesterday.toDateString()) {
      return 'Вчера';
    }
    return date.toLocaleDateString('ru-RU', {
      day: 'numeric',
      month: 'long'
    });
  };

  if (!chatId) {
    return <div className='tg-empty'>Выберите чат слева</div>;
  }

  return (
    <div className='tg-chat-page'>
      {/* HEADER */}
      <div className='tg-chat-header'>
        <button
          className='back-btn'
          onClick={() => navigate('/message')}
          title='Назад'
        >
          <svg width='24' height='24' viewBox='0 0 24 24' fill='none'>
            <path
              d='M15 18l-6-6 6-6'
              stroke='currentColor'
              strokeWidth='2'
              strokeLinecap='round'
              strokeLinejoin='round'
            />
          </svg>
        </button>

        <div className='header-left'>
          <div className='avatar-container'>
            <img
              className='tg-header-avatar'
              src={
                currentChat?.avatarUrl ||
                'https://cdn-icons-png.flaticon.com/512/3177/3177440.png'
              }
              alt=''
            />
            <span className='online-dot'></span>
          </div>
          <div className='header-info'>
            <div className='tg-header-name'>
              {currentChat?.name || `Чат #${chatId}`}
            </div>
            <div className='tg-header-status'>
              {isTyping ? (
                <span className='typing-indicator'>
                  <span className='dot'></span>
                  <span className='dot'></span>
                  <span className='dot'></span>
                  печатает...
                </span>
              ) : (
                'онлайн'
              )}
            </div>
          </div>
        </div>

        <div className='header-actions'>
          <button className='header-btn' title='Поиск'>
            <svg width='20' height='20' viewBox='0 0 24 24' fill='none'>
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
          </button>
          <button className='header-btn header-menu' title='Меню'>
            <svg width='20' height='20' viewBox='0 0 24 24' fill='none'>
              <circle cx='12' cy='5' r='2' fill='currentColor' />
              <circle cx='12' cy='12' r='2' fill='currentColor' />
              <circle cx='12' cy='19' r='2' fill='currentColor' />
            </svg>
          </button>
        </div>
      </div>

      {/* MESSAGES */}
      <div className='tg-messages' ref={messagesRef}>
        {messages.length === 0 ? (
          <div className='no-messages'>
            <div className='no-messages-icon'>💬</div>
            <p>Нет сообщений</p>
            <span>Начните переписку</span>
          </div>
        ) : (
          messages.map((msg, idx) => {
            const showDate =
              idx === 0 ||
              formatDate(msg.createdAt) !==
                formatDate(messages[idx - 1]?.createdAt);

            return (
              <div key={msg.id}>
                {showDate && (
                  <div className='date-divider'>
                    <span>{formatDate(msg.createdAt)}</span>
                  </div>
                )}
                <div
                  className={`tg-msg ${
                    msg.senderId === localStorage.getItem('userId')
                      ? 'mine'
                      : ''
                  }`}
                >
                  <div className='tg-msg-content'>
                    <div className='tg-msg-text'>{msg.content}</div>
                    <div className='tg-msg-meta'>
                      <span className='tg-msg-time'>
                        {formatTime(msg.createdAt)}
                      </span>
                      {msg.senderId === localStorage.getItem('userId') && (
                        <svg
                          className='msg-status'
                          width='16'
                          height='16'
                          viewBox='0 0 24 24'
                          fill='none'
                        >
                          <path
                            d='M20 6L9 17L4 12'
                            stroke='currentColor'
                            strokeWidth='2'
                            strokeLinecap='round'
                            strokeLinejoin='round'
                          />
                        </svg>
                      )}
                    </div>
                  </div>
                </div>
              </div>
            );
          })
        )}
      </div>

      {/* INPUT BOX */}
      <div className='tg-input-box'>
        <button
          className='action-btn attach-btn mobile-hidden'
          title='Прикрепить'
        >
          <svg width='24' height='24' viewBox='0 0 24 24' fill='none'>
            <path
              d='M21.44 11.05l-9.19 9.19a6 6 0 11-8.49-8.49l8.57-8.57A4 4 0 1 1 18 8.84l-8.59 8.57a2 2 0 0 1-2.83-2.83l8.49-8.48'
              stroke='currentColor'
              strokeWidth='2'
              strokeLinecap='round'
              strokeLinejoin='round'
            />
          </svg>
        </button>

        <div className='input-container'>
          <textarea
            ref={inputRef}
            value={text}
            onChange={e => setText(e.target.value)}
            onKeyDown={handleKeyDown}
            placeholder='Напишите сообщение...'
            rows={1}
            className='message-input'
          />

          <button className='emoji-btn' title='Эмодзи'>
            <svg width='24' height='24' viewBox='0 0 24 24' fill='none'>
              <circle
                cx='12'
                cy='12'
                r='10'
                stroke='currentColor'
                strokeWidth='2'
              />
              <path
                d='M8 14s1.5 2 4 2 4-2 4-2'
                stroke='currentColor'
                strokeWidth='2'
                strokeLinecap='round'
              />
              <line
                x1='9'
                y1='9'
                x2='9.01'
                y2='9'
                stroke='currentColor'
                strokeWidth='2'
                strokeLinecap='round'
              />
              <line
                x1='15'
                y1='9'
                x2='15.01'
                y2='9'
                stroke='currentColor'
                strokeWidth='2'
                strokeLinecap='round'
              />
            </svg>
          </button>
        </div>

        <button
          className={`send-btn ${text.trim() ? 'active' : ''}`}
          onClick={handleSend}
          disabled={!text.trim()}
          title='Отправить'
        >
          <svg width='24' height='24' viewBox='0 0 24 24' fill='none'>
            <path
              d='M22 2L11 13M22 2l-7 20-4-9-9-4 20-7z'
              fill='currentColor'
            />
          </svg>
        </button>
      </div>
    </div>
  );
};
