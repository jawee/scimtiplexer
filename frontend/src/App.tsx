import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom';
import './App.css'
import { ToastProvider } from './context/ToastContext';
import { AuthProvider } from './context/AuthContext';
import { ToastContainer } from './components/Toast';
import Protected from './components/Protected';
import LoginPage from './pages/LoginPage';

function App() {
  return (
    <ToastProvider>
      <AuthProvider>
        <BrowserRouter>
          <Routes>
            <Route path="/login" element={<LoginPage />} />
            {/* <Route path="/signup" element={<SignupPage />} /> */}
            {/* <Route path="/about" element={<AboutPage />} /> */}
            {/* <Route path="/privacy" element={<PrivacyPage />} /> */}
            {/* <Route path="/terms" element={<TermsPage />} /> */}

            <Route element={<Protected />}>
              {/* <Route path="/rooms" element={<RoomsPage />} /> */}
              {/* <Route path="/chat/:roomId" element={<ChatPage />} /> */}
              {/* <Route path="/profile" element={<ProfilePage />} /> */}
              {/* <Route index element={<Navigate to="/rooms" replace />} /> */}
            </Route>

            <Route path="*" element={<Navigate to="/login" replace />} />
          </Routes>
          <ToastContainer />
        </BrowserRouter>
      </AuthProvider>
    </ToastProvider>
  );
}

export default App
