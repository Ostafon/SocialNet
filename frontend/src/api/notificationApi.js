// src/api/notificationApi.js
import api from './api';

const NotificationAPI = {
  list: (limit = 50, offset = 0) =>
    api.get(`/notifications?limit=${limit}&offset=${offset}`).then(r => r.data),

  markAsRead: id => api.post(`/notifications/${id}/read`).then(r => r.data),

  markAllAsRead: () => api.post('/notifications/read-all').then(r => r.data),

  delete: id => api.delete(`/notifications/${id}`).then(r => r.data),

  clearAll: () => api.delete('/notifications/clear').then(r => r.data)
};

export default NotificationAPI;
