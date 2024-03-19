import LoginPage from './pages/Login';
import ChatPage from './pages/Chat';
import RoomListPage from './pages/RoomList';
import { Provider } from 'react-redux';
import store from './store';
import { BrowserRouter as Router, Route, Routes } from 'react-router-dom';

function App() {
  return (
    <Provider store={store}>
        <Routes>
          <Route path="/" element={<LoginPage />} />
          <Route path="/chat" element={<ChatPage />} />
          <Route path="/rooms" element={<RoomListPage />} />
        </Routes>
    </Provider>
  );
}

export default App;