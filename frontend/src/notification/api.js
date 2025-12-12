// src/notifications/api.js
import NotificationAPI from '../api/notificationApi';

// ==============================
//  Получить список уведомлений
// ==============================
export async function listNotifications(limit = 50, offset = 0) {
  const res = await NotificationAPI.list(limit, offset);

  // backend возвращает:
  // { notifications: [...] }
  return res.notifications || [];
}

// ==============================
//  Отметить одно уведомление
// ==============================
export async function markAsRead(id) {
  return NotificationAPI.markAsRead(id);
}

// ==============================
//  Отметить все уведомления
// ==============================
export async function markAllAsRead() {
  return NotificationAPI.markAllAsRead();
}

// ==============================
//  Удалить одно уведомление
// ==============================
export async function deleteNotification(id) {
  return NotificationAPI.delete(id);
}

// ==============================
//  Удалить все уведомления
// ==============================
export async function clearAll() {
  return NotificationAPI.clearAll();
}
