import { Outlet, useParams } from 'react-router-dom';
import { MessageList } from '../../components/messageList/messageList';
import './message.scss';

export const Message = () => {
  const { id } = useParams();

  return (
    <div className='tg-layout'>
      <div className='tg-left'>
        <MessageList />
      </div>

      <div className='tg-right'>
        {id ? (
          <Outlet />
        ) : (
          <div className='tg-empty'>
            <div className='empty-icon'>💬</div>
            <h2>Выберите чат</h2>
            <p>Откройте беседу, чтобы начать общение</p>
          </div>
        )}
      </div>
    </div>
  );
};
