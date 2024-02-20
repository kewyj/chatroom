import React from 'react';
import LoginPage from './pages/Login';
import ChatPage from './pages/Chat';

import { BrowserRouter, Route, Routes } from 'react-router-dom';

export interface AppProps {}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path='/' element={<LoginPage />} />
        <Route path="chat" element={<ChatPage />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App;