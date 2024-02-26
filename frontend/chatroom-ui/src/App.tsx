import LoginPage from './pages/Login';
import ChatPage from './pages/Chat';
import { Provider } from 'react-redux';
import store from './store';
import { HashRouter as Router, Route, Routes } from 'react-router-dom';

function App() {
  return (
    <Provider store={store}>
      <Router>
        <Routes>
          <Route path="/" element={<LoginPage />} />
          <Route path="/chat" element={<ChatPage />} />
        </Routes>
      </Router>
    </Provider>
  );
}

export default App;