import axios from 'axios';

const api = axios.create({
  baseURL: 'http://192.168.1.5:8080/api/v1',
  headers: { 'ngrok-skip-browser-warning': 'true' }
});

// Interceptor для добавления accessToken в каждый запрос
api.interceptors.request.use(config => {
  const token = localStorage.getItem('accessToken');
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

// Interceptor для обработки 401 и обновления токена
api.interceptors.response.use(
  response => response,
  async error => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      originalRequest._retry = true;

      const refreshToken = localStorage.getItem('refreshToken');
      if (refreshToken) {
        try {
          const res = await axios.post(
            'http://192.168.1.5:8080/api/v1/auth/refresh',
            { refreshToken },
            { headers: { 'ngrok-skip-browser-warning': 'true' } }
          );

          const { accessToken: newAccessToken } = res.data;
          localStorage.setItem('accessToken', newAccessToken);

          originalRequest.headers.Authorization = `Bearer ${newAccessToken}`;
          return api(originalRequest);
        } catch (err) {
          console.error('Не удалось обновить accessToken', err);
        }
      }
    }

    return Promise.reject(error);
  }
);

export default api;
