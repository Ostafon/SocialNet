import api from './api';
export const SearchAPI = {
  users: (query, limit = 20) =>
    api.get(`/search/users?query=${query}&limit=${limit}`),

  posts: (query, limit = 20) =>
    api.get(`/search/posts?query=${query}&limit=${limit}`),

  all: (query, limit = 20) => api.get(`/search?query=${query}&limit=${limit}`)
};
