import { FriendList } from '../../components/friendList/friendList';
import { PostList } from '../../components/postList/postList';
import { Sidebar } from '../../components/sidebar/sidebar';
import './home.scss';

export const Home = () => {
  const currentUser = JSON.parse(localStorage.getItem('currentUser'));

  return (
    <section className='home'>
      <Sidebar />

      <div className='home__content'>
        <h1 className='home__title'>Лента постов</h1>

        <div className='home__layout'>
          {/* Левая колонка – посты */}
          <div className='home__posts'>
            <PostList currentUser={currentUser} />
          </div>

          {/* Правая колонка – друзья */}
          <aside className='home__friends'>
            <FriendList onlyFollowing={true} />
          </aside>
        </div>
      </div>
    </section>
  );
};
