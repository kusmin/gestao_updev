import React from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import ProtectedRoute from './components/ProtectedRoute';
import DashboardPage from './pages/DashboardPage';
import TenantListPage from './pages/tenants/TenantListPage';
import UserListPage from './pages/users/UserListPage';
import ClientListPage from './pages/clients/ClientListPage';
import ProductListPage from './pages/products/ProductListPage';
import ServiceListPage from './pages/services/ServiceListPage';
import AppointmentListPage from './pages/appointments/AppointmentListPage';
import SaleListPage from './pages/sales/SaleListPage';
import LoginPage from './pages/LoginPage';

const router = createBrowserRouter([
  {
    path: '/login',
    element: <LoginPage />,
  },
  {
    path: '/',
    element: <ProtectedRoute />,
    children: [
      {
        path: '/',
        element: <DashboardPage />,
      },
      {
        path: '/tenants',
        element: <TenantListPage />,
      },
      {
        path: '/users',
        element: <UserListPage />,
      },
      {
        path: '/clients',
        element: <ClientListPage />,
      },
      {
        path: '/products',
        element: <ProductListPage />,
      },
      {
        path: '/services',
        element: <ServiceListPage />,
      },
      {
        path: '/appointments',
        element: <AppointmentListPage />,
      },
      {
        path: '/sales',
        element: <SaleListPage />,
      },
    ],
  },
]);

const AppRouter: React.FC = () => {
  return <RouterProvider router={router} />;
};

export default AppRouter;
