// src/store/useChatStore.js
import { create } from 'zustand';
import {
  grpcListChats,
  grpcListMessages,
  grpcSendMessage
} from '../chat/chatClient';
import { createChatStream } from '../chat/stream';

// --- функция для защиты от дублей ---
function safeAddMessage(state, chatId, msg) {
  const arr = state.messages[chatId] || [];

  // если сообщение уже существует — не добавляем
  if (arr.some(m => m.id === msg.id)) return arr;

  return [...arr, msg];
}

export const useChatStore = create((set, get) => ({
  chats: [],
  messages: {}, // { [chatId]: Message[] }
  streams: {}, // { [chatId]: cancelFn }
  streamStatus: {}, // { [chatId]: 'connected' | 'connecting' | 'disconnected' }
  activeChat: null,
  loadingChats: false,
  loadingMessages: false,

  // ===== ЧАТЫ =====
  loadChats: async () => {
    try {
      set({ loadingChats: true });
      const chats = await grpcListChats();

      console.log('🔥 ЧАТЫ ПРИШЛИ С БЭКА:', chats);

      set({ chats, loadingChats: false });
    } catch (err) {
      console.error('❌ loadChats error:', err);
      if (err.message) console.error('Error message:', err.message);
      if (err.code) console.error('Error code:', err.code);
      set({ loadingChats: false, chats: [] });
    }
  },

  // ===== СООБЩЕНИЯ =====
  loadMessages: async chatId => {
    const id = String(chatId);
    try {
      set({ loadingMessages: true });

      const msgs = await grpcListMessages(id);

      set(state => ({
        messages: {
          ...state.messages,
          [id]: msgs
        },
        loadingMessages: false
      }));
    } catch (err) {
      console.error('❌ loadMessages error:', err);
      if (err.message) console.error('Error message:', err.message);
      set(state => ({
        messages: {
          ...state.messages,
          [id]: []
        },
        loadingMessages: false
      }));
    }
  },

  // ===== ОТПРАВКА =====
  sendMessage: async (chatId, content) => {
    const id = String(chatId);
    try {
      const msg = await grpcSendMessage(id, content);

      set(state => ({
        messages: {
          ...state.messages,
          [id]: safeAddMessage(state, id, msg)
        }
      }));

      return msg;
    } catch (err) {
      console.error('❌ sendMessage error:', err);
      if (err.message) console.error('Error message:', err.message);
      throw err;
    }
  },

  // ===== АКТИВНЫЙ ЧАТ =====
  setActiveChat: chatId => {
    const id = String(chatId);
    const { activeChat, stopAllStreams, startChatStream, loadMessages } = get();

    if (activeChat === id) return;

    // Останавливаем все старые стримы
    stopAllStreams();

    // Устанавливаем новый активный чат
    set({ activeChat: id });

    // Загружаем сообщения
    loadMessages(id);

    // Запускаем стрим для нового чата
    startChatStream(id);
  },

  // ===== СТРИМ =====
  startChatStream: chatId => {
    const id = String(chatId);
    const { streams, streamStatus } = get();

    if (streams[id]) {
      console.log('⚠️ Stream already exists for chat:', id);
      return;
    }

    console.log('🔔 Starting notification stream for chat:', id);

    // Устанавливаем статус подключения
    set(state => ({
      streamStatus: { ...state.streamStatus, [id]: 'connecting' }
    }));

    try {
      const cancel = createChatStream(id, {
        onMessage: msg => {
          console.log('📩 New message from stream:', msg);
          set(state => ({
            messages: {
              ...state.messages,
              [id]: safeAddMessage(state, id, msg)
            }
          }));
        },
        onConnect: () => {
          console.log('✅ Stream connected for chat:', id);
          set(state => ({
            streamStatus: { ...state.streamStatus, [id]: 'connected' }
          }));
        },
        onDisconnect: () => {
          console.log('⚠️ Stream disconnected for chat:', id);
          set(state => ({
            streamStatus: { ...state.streamStatus, [id]: 'disconnected' }
          }));
        },
        onReconnect: () => {
          console.log('🔄 Stream reconnecting for chat:', id);
          set(state => ({
            streamStatus: { ...state.streamStatus, [id]: 'connecting' }
          }));
        }
      });

      set(state => ({
        streams: { ...state.streams, [id]: cancel }
      }));
    } catch (err) {
      console.error('❌ startChatStream error:', err);
      set(state => ({
        streamStatus: { ...state.streamStatus, [id]: 'disconnected' }
      }));
    }
  },

  // ===== СТОП СТРИМА =====
  stopChatStream: chatId => {
    const id = String(chatId);
    const { streams } = get();

    if (streams[id] && typeof streams[id] === 'function') {
      console.log('⛔️ Stopping stream for chat:', id);
      streams[id]();

      set(state => {
        const newStreams = { ...state.streams };
        const newStatus = { ...state.streamStatus };
        delete newStreams[id];
        delete newStatus[id];

        return { streams: newStreams, streamStatus: newStatus };
      });
    }
  },

  // ===== СТОП ВСЕХ СТРИМОВ =====
  stopAllStreams: () => {
    const { streams } = get();
    console.log('⛔️ Stopping all streams...');

    Object.entries(streams).forEach(([chatId, fn]) => {
      if (typeof fn === 'function') {
        console.log('⛔️ Cancelling stream for chat:', chatId);
        fn();
      }
    });

    set({ streams: {}, streamStatus: {} });
  },

  // ===== ПОЛУЧИТЬ СТАТУС СТРИМА =====
  getStreamStatus: chatId => {
    const { streamStatus } = get();
    return streamStatus[String(chatId)] || 'disconnected';
  }
}));
