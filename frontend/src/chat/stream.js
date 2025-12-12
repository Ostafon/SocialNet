import { getMetadata } from './chatClient';
import { ChatServiceClient } from './proto/chat_grpc_web_pb';
import { SubscribeRequest } from './proto/chat_pb';

const HOST = process.env.REACT_APP_GRPC_URL || 'http://192.168.1.5:8082';

const options = {
  format: 'text'
};

const client = new ChatServiceClient(HOST, null, options);

export function createChatStream(
  chatId,
  { onMessage, onConnect, onDisconnect, onReconnect }
) {
  console.log('📡 Connecting to:', HOST);
  console.log('🔔 Starting notification stream for chat:', chatId);

  let isCancelled = false;
  let reconnectTimeout;
  let currentStream = null;
  let reconnectAttempts = 0;
  const MAX_RECONNECT_ATTEMPTS = 10;

  function connect() {
    if (isCancelled) return;

    // Увеличиваем счетчик попыток переподключения
    if (reconnectAttempts > 0) {
      console.log(
        `🔄 Reconnect attempt ${reconnectAttempts}/${MAX_RECONNECT_ATTEMPTS}`
      );
      onReconnect?.();
    }

    const req = new SubscribeRequest();
    req.setChatIdsList([String(chatId)]);

    const metadata = getMetadata();

    currentStream = client.subscribeMessages(req, metadata);

    currentStream.on('data', msg => {
      if (isCancelled) return;

      // Сбрасываем счетчик попыток при успешном получении данных
      reconnectAttempts = 0;

      try {
        const parsedMsg = {
          id: msg.getId(),
          chatId: msg.getChatId(),
          senderId: msg.getSenderId(),
          content: msg.getContent() || '',
          contentType: msg.getContentType() || 'text',
          mediaUrl: msg.getMediaUrl() || '',
          createdAt: msg.getCreatedAt() || ''
        };

        console.log('📩 New message from stream:', parsedMsg);
        onMessage(parsedMsg);
      } catch (err) {
        console.error('❌ Failed to parse stream message:', err);
      }
    });

    currentStream.on('error', err => {
      if (isCancelled) return;

      // Игнорируем нормальное закрытие
      if (err.code === 0 || err.code === 1) {
        console.log('✅ Stream closed normally');
        return;
      }

      console.error('❌ Stream error:', {
        code: err.code,
        message: err.message,
        metadata: err.metadata
      });

      onDisconnect?.();

      // Автоматическое переподключение
      if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
        reconnectAttempts++;
        const delay = Math.min(
          1000 * Math.pow(2, reconnectAttempts - 1),
          30000
        ); // Exponential backoff
        console.log(`⏳ Reconnecting in ${delay}ms...`);

        reconnectTimeout = setTimeout(() => {
          connect();
        }, delay);
      } else {
        console.error('❌ Max reconnection attempts reached. Stream stopped.');
      }
    });

    currentStream.on('end', () => {
      if (isCancelled) return;

      console.log('📡 Stream ended');
      onDisconnect?.();

      // Переподключаемся при неожиданном завершении
      if (reconnectAttempts < MAX_RECONNECT_ATTEMPTS) {
        reconnectAttempts++;
        const delay = 2000;
        console.log(`⏳ Stream ended, reconnecting in ${delay}ms...`);

        reconnectTimeout = setTimeout(() => {
          connect();
        }, delay);
      }
    });

    currentStream.on('status', status => {
      console.log('📡 Stream status:', status);

      // Если статус OK - сообщаем о подключении
      if (status.code === 0) {
        onConnect?.();
      }
    });
  }

  // Начинаем подключение
  connect();

  // Возвращаем функцию отмены
  return () => {
    if (isCancelled) return;

    isCancelled = true;
    reconnectAttempts = MAX_RECONNECT_ATTEMPTS; // Предотвращаем дальнейшие попытки переподключения

    console.log('⛔️ Stream cancelled manually for chat:', chatId);

    // Очищаем таймаут переподключения
    if (reconnectTimeout) {
      clearTimeout(reconnectTimeout);
      reconnectTimeout = null;
    } // Отменяем текущий стрим
    if (currentStream) {
      try {
        currentStream.cancel();
      } catch (err) {
        console.error('Error cancelling stream:', err);
      }
      currentStream = null;
    }
  };
}
