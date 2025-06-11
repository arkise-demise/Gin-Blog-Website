// utils/axios-config.ts
import axios from 'axios';

const api = axios.create({
  baseURL: 'http://localhost:8080/api', // Your Gin backend URL
  withCredentials: true, // Important for sending and receiving cookies (JWT)
});

export default api;