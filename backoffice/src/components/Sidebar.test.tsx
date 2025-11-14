import { MemoryRouter } from 'react-router-dom';
import { render, screen } from '@testing-library/react';
import { describe, it, expect } from 'vitest';
import Sidebar from './Sidebar';

const MENU_ITEMS = [
  { label: 'Dashboard', path: '/' },
  { label: 'Tenants', path: '/tenants' },
  { label: 'Users', path: '/users' },
  { label: 'Clients', path: '/clients' },
  { label: 'Products', path: '/products' },
  { label: 'Services', path: '/services' },
  { label: 'Appointments', path: '/appointments' },
  { label: 'Sales', path: '/sales' },
];

describe('Sidebar', () => {
  const renderSidebar = () =>
    render(
      <MemoryRouter>
        <Sidebar />
      </MemoryRouter>,
    );

  it('renderiza todos os links configurados', () => {
    renderSidebar();

    MENU_ITEMS.forEach(({ label }) => {
      expect(screen.getByRole('link', { name: label })).toBeInTheDocument();
    });
  });

  it('aponta cada item para o caminho correto', () => {
    renderSidebar();

    MENU_ITEMS.forEach(({ label, path }) => {
      expect(screen.getByRole('link', { name: label })).toHaveAttribute(
        'href',
        path,
      );
    });
  });
});
