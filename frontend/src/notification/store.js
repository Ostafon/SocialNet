import { create } from 'zustand';
import { createNotificationStream } from '../notification/stream';
import * as NotifAPI from './api';

export const useNotificationStore = create((set, get) => ({
  notifications: [],
  unreadCount: 0,

  stopFn: null,
  isStreaming: false,

  // =========================
  // 📌 Загрузка истории
  // =========================
  loadNotifications: async () => {
    try {
      const list = await NotifAPI.listNotifications();

      const normalized = list.map(n => ({
        id: n.id,
        userId: n.userId,
        type: n.type,
        referenceId: n.referenceId || null,
        content: n.content,
        read: n.read === true,
        createdAt: n.createdAt
      }));

      set({
        notifications: normalized,
        unreadCount: normalized.filter(n => !n.read).length
      });
    } catch (err) {
      console.error('Ошибка загрузки уведомлений:', err);
    }
  },

  // =========================
  // 🔔 Запуск gRPC стрима
  // =========================
  startStream: () => {
    if (get().isStreaming) return;

    console.log('🔔 Starting notification stream...');

    const stopFn = createNotificationStream({
      onMessage: notif => {
        console.log('🔔 STREAM MESSAGE:', notif);

        set(state => {
          // избегаем дубликатов
          if (state.notifications.some(n => n.id === notif.id)) {
            return state;
          }

          return {
            notifications: [notif, ...state.notifications],
            unreadCount: state.unreadCount + (notif.read ? 0 : 1)
          };
        });
      },

      onError: () => {
        console.warn('❌ Stream error → reconnecting...');
        set({ isStreaming: false, stopFn: null });
        setTimeout(() => get().startStream(), 3000);
      }
    });

    set({ stopFn, isStreaming: true });
  },

  // =========================
  // ⛔ Остановка стрима
  // =========================
  stopStream: () => {
    const stop = get().stopFn;
    if (stop) stop();

    set({ stopFn: null, isStreaming: false });
  },

  // =========================
  // ☑️ Отметить как прочитанное
  // =========================
  markAsRead: async id => {
    try {
      await NotifAPI.markAsRead(id);

      set(state => {
        const updated = state.notifications.map(n =>
          n.id === id ? { ...n, read: true } : n
        );

        return {
          notifications: updated,
          unreadCount: updated.filter(n => !n.read).length
        };
      });
    } catch (err) {
      console.error('Ошибка markAsRead:', err);
    }
  },

  // =========================
  // 🗑 Удалить уведомление
  // =========================
  deleteNotification: async id => {
    try {
      await NotifAPI.deleteNotification(id);

      set(state => {
        const updated = state.notifications.filter(n => n.id !== id);

        return {
          notifications: updated,
          unreadCount: updated.filter(n => !n.read).length
        };
      });
    } catch (err) {
      console.error('Ошибка deleteNotification:', err);
    }
  },

  // =========================
  // 🧹 Очистить все
  // =========================
  clearAll: async () => {
    try {
      await NotifAPI.clearAll();
      set({ notifications: [], unreadCount: 0 });
    } catch (err) {
      console.error('Ошибка clearAll:', err);
    }
  }
}));

// Для дебага через DevTools
if (typeof window !== 'undefined') {
  window.useNotificationStore = useNotificationStore;
}
