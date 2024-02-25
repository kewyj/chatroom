import { BrowserRouter, Route, Routes } from 'react-router-dom';
import LoginPage from './pages/Login';
import ChatPage from './pages/Chat';
import { Provider } from 'react-redux';
import store from './store';

function App() {
  return (
    <Provider store={store}>
    <BrowserRouter>
      <Routes>
        {/* Define routes without using the catch-all route */}
        <Route path="/" element={<LoginPage />} />
        <Route path="/chat" element={<ChatPage />} />
      </Routes>
      </BrowserRouter>
    </Provider>
  )
}

export default App;