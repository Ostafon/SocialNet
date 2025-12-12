// src/chat/chatClient.js
import { ChatServiceClient } from './proto/chat_grpc_web_pb';
import {
  EmptyRequest,
  ListMessagesRequest,
  SendMessageRequest
} from './proto/chat_pb';

const HOST = process.env.REACT_APP_GRPC_URL || 'http://192.168.1.5:8082';

const options = {
  format: 'text'
};

export const client = new ChatServiceClient(HOST, null, options);

// ===== МЕТАДАННЫЕ =====
export function getMetadata() {
  const token = localStorage.getItem('accessToken');
  const userId = localStorage.getItem('userId');

  const meta = {};
  if (token) meta['authorization'] = `Bearer ${token}`;
  if (userId) meta['user-id'] = userId;

  return meta;
}

// ===== УНИВЕРСАЛЬНЫЙ PROMISE =====
function promisifyUnary(method, request) {
  const metadata = getMetadata();

  return new Promise((resolve, reject) => {
    try {
      method.call(client, request, metadata, (err, resp) => {
        if (err) {
          console.error('❌ gRPC Error:', {
            code: err.code,
            message: err.message
          });
          return reject(err);
        }

        if (!resp) {
          console.error('❌ Empty response from gRPC');
          return reject(new Error('Empty response'));
        }

        resolve(resp);
      });
    } catch (err) {
      console.error('❌ Exception in promisifyUnary:', err);
      reject(err);
    }
  });
}

// ===== УНИВЕРСАЛЬНАЯ ФУНКЦИЯ ДЛЯ ЧТЕНИЯ СПИСКА =====
function extractArray(response, possibleMethods) {
  for (const method of possibleMethods) {
    if (typeof response[method] === 'function') {
      try {
        const result = response[method]();
        if (Array.isArray(result)) {
          console.log(
            `✅ Found array via ${method}() - ${result.length} items`
          );
          return result;
        }
        // Если вернулась обёртка - пробуем вызвать методы на ней
        if (result && typeof result === 'object') {
          const wrapperMethods = Object.getOwnPropertyNames(
            Object.getPrototypeOf(result)
          ).filter(
            m =>
              typeof result[m] === 'function' &&
              m.toLowerCase().includes('list')
          );

          for (const wrapperMethod of wrapperMethods) {
            try {
              const innerResult = result[wrapperMethod]();
              if (Array.isArray(innerResult)) {
                console.log(
                  `✅ Found array via ${method}().${wrapperMethod}() - ${innerResult.length} items`
                );
                return innerResult;
              }
            } catch (e) {}
          }
        }
      } catch (e) {
        console.log(`⚠️  ${method}() failed:`, e.message);
      }
    }
  }
  return null;
}

// ===== API =====

export async function grpcListChats() {
  try {
    console.log('📡 Calling grpcListChats...');
    const req = new EmptyRequest();
    const res = await promisifyUnary(client.listChats, req);

    console.log('📦 Response type:', res.constructor.name);

    // Пробуем все возможные методы для получения списка чатов
    const chatsArray = extractArray(res, [
      'getChatsList',
      'getChats',
      'toObject',
      'getChatslist',
      'chats'
    ]);

    if (!chatsArray) {
      console.error('❌ Could not find chats array in response');
      console.log(
        'Available methods:',
        Object.getOwnPropertyNames(Object.getPrototypeOf(res))
      );
      return [];
    }

    console.log(`✅ Processing ${chatsArray.length} chats`);

    return chatsArray.map(c => {
      const getId = c.getId || c.getid || c.id;
      const getName = c.getName || c.getname || c.name;
      const getIsGroup = c.getIsgroup || c.getIsGroup || c.is_group;
      const getLastMessage =
        c.getLastmessage || c.getLastMessage || c.last_message;
      const getParticipants =
        c.getParticipantsList || c.getparticipantsList || c.participants;
      const getCreatedAt =
        c.getCreatedat || (c.getCreatedAt ? c.getCreatedAt() : c.created_at);
      const getUpdatedAt =
        c.getUpdatedat || (c.getUpdatedAt ? c.getUpdatedAt() : c.updated_at);

      return {
        id: typeof getId === 'function' ? getId.call(c) : getId || '',
        name: typeof getName === 'function' ? getName.call(c) : getName || '',
        isGroup:
          typeof getIsGroup === 'function'
            ? getIsGroup.call(c)
            : getIsGroup || false,
        lastMessage:
          typeof getLastMessage === 'function'
            ? getLastMessage.call(c)
            : getLastMessage || null,
        participants:
          typeof getParticipants === 'function'
            ? getParticipants.call(c)
            : getParticipants || [],
        createdAt:
          typeof getCreatedAt === 'function'
            ? getCreatedAt.call(c)
            : getCreatedAt || '',
        updatedAt:
          typeof getUpdatedAt === 'function'
            ? getUpdatedAt.call(c)
            : getUpdatedAt || ''
      };
    });
  } catch (err) {
    console.error('❌ grpcListChats failed:', err);
    throw err;
  }
}

export async function grpcListMessages(chatId, limit = 50, offset = 0) {
  try {
    console.log(`📡 Calling grpcListMessages for chat: ${chatId}`);
    const req = new ListMessagesRequest();

    // Пробуем оба варианта
    if (typeof req.setChatid === 'function') {
      req.setChatid(String(chatId));
    } else {
      req.setChatId(String(chatId));
    }
    req.setLimit(limit);
    req.setOffset(offset);

    const res = await promisifyUnary(client.listMessages, req);

    console.log('📦 Response type:', res.constructor.name);

    // Пробуем все возможные методы для получения списка сообщений
    const messagesArray = extractArray(res, [
      'getMessagesList',
      'getMessages',
      'toObject',
      'getMessageslist',
      'messages'
    ]);

    if (!messagesArray) {
      console.error('❌ Could not find messages array in response');
      console.log(
        'Available methods:',
        Object.getOwnPropertyNames(Object.getPrototypeOf(res))
      );
      return [];
    }

    console.log(`✅ Processing ${messagesArray.length} messages`);

    return messagesArray.map(m => {
      const getId = m.getId || m.getid || m.id;
      const getChatId = m.getChatid || m.getChatId || m.chat_id;
      const getSenderId = m.getSenderid || m.getSenderId || m.sender_id;
      const getContent = m.getContent || m.getcontent || m.content;
      const getContentType =
        m.getContenttype || m.getContentType || m.content_type;
      const getMediaUrl = m.getMediaurl || m.getMediaUrl || m.media_url;
      const getCreatedAt = m.getCreatedat || m.getCreatedAt || m.created_at;
      const getRead = m.getRead || m.read;

      return {
        id: typeof getId === 'function' ? getId.call(m) : getId || '',
        chatId:
          typeof getChatId === 'function' ? getChatId.call(m) : getChatId || '',
        senderId:
          typeof getSenderId === 'function'
            ? getSenderId.call(m)
            : getSenderId || '',
        content:
          typeof getContent === 'function'
            ? getContent.call(m)
            : getContent || '',
        contentType:
          typeof getContentType === 'function'
            ? getContentType.call(m)
            : getContentType || 'text',
        mediaUrl:
          typeof getMediaUrl === 'function'
            ? getMediaUrl.call(m)
            : getMediaUrl || '',
        createdAt:
          typeof getCreatedAt === 'function'
            ? getCreatedAt.call(m)
            : getCreatedAt || '',
        read: typeof getRead === 'function' ? getRead.call(m) : getRead || false
      };
    });
  } catch (err) {
    console.error('❌ grpcListMessages failed:', err);
    throw err;
  }
}

export async function grpcSendMessage(chatId, content) {
  try {
    console.log(`📡 Calling grpcSendMessage for chat: ${chatId}`);
    const req = new SendMessageRequest();

    // Пробуем оба варианта
    if (typeof req.setChatid === 'function') {
      req.setChatid(String(chatId));
    } else {
      req.setChatId(String(chatId));
    }
    req.setContent(content);

    if (typeof req.setContenttype === 'function') {
      req.setContenttype('text');
    } else {
      req.setContentType('text');
    }

    const res = await promisifyUnary(client.sendMessage, req);

    console.log('📦 SendMessage response type:', res.constructor.name);

    const getId = res.getId || res.getid || res.id;
    const getChatId = res.getChatid || res.getChatId || res.chat_id;
    const getSenderId = res.getSenderid || res.getSenderId || res.sender_id;
    const getContent = res.getContent || res.getcontent || res.content;
    const getContentType =
      res.getContenttype || res.getContentType || res.content_type;
    const getMediaUrl = res.getMediaurl || res.getMediaUrl || res.media_url;
    const getCreatedAt = res.getCreatedat || res.getCreatedAt || res.created_at;

    return {
      id: typeof getId === 'function' ? getId.call(res) : getId || '',
      chatId:
        typeof getChatId === 'function' ? getChatId.call(res) : getChatId || '',
      senderId:
        typeof getSenderId === 'function'
          ? getSenderId.call(res)
          : getSenderId || '',
      content:
        typeof getContent === 'function'
          ? getContent.call(res)
          : getContent || '',
      contentType:
        typeof getContentType === 'function'
          ? getContentType.call(res)
          : getContentType || 'text',
      mediaUrl:
        typeof getMediaUrl === 'function'
          ? getMediaUrl.call(res)
          : getMediaUrl || '',
      createdAt:
        typeof getCreatedAt === 'function'
          ? getCreatedAt.call(res)
          : getCreatedAt || new Date().toISOString(),
      read: false
    };
  } catch (err) {
    console.error('❌ grpcSendMessage failed:', err);
    throw err;
  }
}
