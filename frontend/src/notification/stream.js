import { NotificationServiceClient } from './proto/notification_grpc_web_pb';
import { StreamRequest } from './proto/notification_pb';

const HOST = process.env.REACT_APP_GRPC_URL || 'http://192.168.1.5:8081';
const client = new NotificationServiceClient(HOST);

export function createNotificationStream({ onMessage, onError }) {
  const req = new StreamRequest();
  req.setUserId(localStorage.getItem('userId') || '');

  const metadata = {
    Authorization: `Bearer ${localStorage.getItem('accessToken')}`,
    'x-grpc-web': '1'
  };

  console.log('📡 Connecting to:', HOST);

  const stream = client.streamNotifications(req, metadata);

  // ======================
  // ✔️ ЕДИНСТВЕННЫЙ ОБРАБОТЧИК STREAM DATA
  // ======================
  stream.on('data', notif => {
    try {
      console.log('🔥 STREAM RAW PROTOBUF:', notif);

      const msg = notif.toObject();
      console.log('🔥 STREAM PARSED:', msg);

      const normalized = {
        id: msg.id || '',
        userId: msg.userId || '',
        type: msg.type || 'unknown',
        referenceId: msg.referenceId || null,
        content: msg.content || '',
        read: msg.read ?? false,
        createdAt: msg.createdAt || new Date().toISOString()
      };

      console.log('🔥 STREAM NORMALIZED:', normalized);

      onMessage(normalized);
    } catch (err) {
      console.error('❌ ERROR PARSING STREAM MESSAGE:', err);
    }
  });

  // ======================
  // ❌ ОШИБКА
  // ======================
  stream.on('error', err => {
    console.error('❌ STREAM ERROR:', err);
    onError?.(err);
  });

  // ======================
  // ❌ ЗАКРЫТО
  // ======================
  stream.on('end', () => {
    console.log('ℹ️ STREAM CLOSED');
  });

  return () => {
    console.log('⛔ Stream cancelled manually');
    stream.cancel();
  };
}
