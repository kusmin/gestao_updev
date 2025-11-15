import { render, screen } from '@testing-library/react';
import { MemoryRouter, Route, Routes } from 'react-router-dom';
import { beforeEach, describe, expect, it, vi } from 'vitest';

const useAuthMock = vi.hoisted(() => vi.fn());

vi.mock('../contexts/AuthContext', () => ({
  useAuth: useAuthMock,
}));

import ProtectedRoute from './ProtectedRoute';

describe('ProtectedRoute', () => {
  beforeEach(() => {
    useAuthMock.mockReset();
  });

  const renderWithRoutes = () =>
    render(
      <MemoryRouter initialEntries={['/dashboard']}>
        <Routes>
          <Route path="/login" element={<div>Login Page</div>} />
          <Route element={<ProtectedRoute />}>
            <Route path="/dashboard" element={<div>Protected Content</div>} />
          </Route>
        </Routes>
      </MemoryRouter>,
    );

  it('redirects to the login page when no token is available', () => {
    useAuthMock.mockReturnValue({ token: null });

    renderWithRoutes();

    expect(screen.getByText('Login Page')).toBeInTheDocument();
  });

  it('renders the protected content when a token exists', () => {
    useAuthMock.mockReturnValue({ token: 'token-value' });

    renderWithRoutes();

    expect(screen.getByText('Protected Content')).toBeInTheDocument();
    expect(screen.getByText(/gestão updev - backoffice/i)).toBeInTheDocument();
  });
});
