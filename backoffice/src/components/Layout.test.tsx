import { render, screen } from '@testing-library/react';
import { MemoryRouter } from 'react-router-dom';
import { type ReactNode } from 'react';
import Layout from './Layout';

describe('Layout', () => {
  const renderLayout = (children: ReactNode) =>
    render(
      <MemoryRouter>
        <Layout>{children}</Layout>
      </MemoryRouter>,
    );

  it('renders the application header', () => {
    renderLayout(<div />);

    expect(
      screen.getByText(/gestão updev - backoffice/i),
    ).toBeInTheDocument();
  });

  it('renders provided children inside the main content area', () => {
    renderLayout(<p>Content goes here</p>);

    expect(screen.getByText('Content goes here')).toBeInTheDocument();
  });

  it('includes the navigation menu items', () => {
    renderLayout(<div />);

    expect(screen.getByRole('link', { name: /dashboard/i })).toBeInTheDocument();
    expect(screen.getByRole('link', { name: /tenants/i })).toBeInTheDocument();
  });
});
