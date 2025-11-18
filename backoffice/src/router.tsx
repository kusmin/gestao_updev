import React from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import ProtectedRoute from './components/ProtectedRoute';
import DashboardPage from './pages/DashboardPage';
import TenantListPage from './pages/tenants/TenantListPage';
import TenantFormWrapper from './pages/tenants/TenantFormWrapper';
import UserListPage from './pages/users/UserListPage';
import UserFormWrapper from './pages/users/UserFormWrapper';
import ClientListPage from './pages/clients/ClientListPage';
import ClientFormWrapper from './pages/clients/ClientFormWrapper';
import ProductListPage from './pages/products/ProductListPage';
import ProductFormWrapper from './pages/products/ProductFormWrapper';
import ServiceListPage from './pages/services/ServiceListPage';
import ServiceFormWrapper from './pages/services/ServiceFormWrapper';
import AppointmentListPage from './pages/appointments/AppointmentListPage';
import AppointmentFormWrapper from './pages/appointments/AppointmentFormWrapper';
import SaleListPage from './pages/sales/SaleListPage';
import SaleFormWrapper from './pages/sales/SaleFormWrapper';
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
        path: '/tenants/new',
        element: <TenantFormWrapper />,
      },
      {
        path: '/tenants/edit/:id',
        element: <TenantFormWrapper />,
      },
      {
        path: '/users',
        element: <UserListPage />,
      },
      {
        path: '/users/new',
        element: <UserFormWrapper />,
      },
      {
        path: '/users/edit/:id',
        element: <UserFormWrapper />,
      },
      {
        path: '/clients',
        element: <ClientListPage />,
      },
      {
        path: '/clients/new',
        element: <ClientFormWrapper />,
      },
      {
        path: '/clients/edit/:id',
        element: <ClientFormWrapper />,
      },
      {
        path: '/products',
        element: <ProductListPage />,
      },
      {
        path: '/products/new',
        element: <ProductFormWrapper />,
      },
      {
        path: '/products/edit/:id',
        element: <ProductFormWrapper />,
      },
      {
        path: '/services',
        element: <ServiceListPage />,
      },
      {
        path: '/services/new',
        element: <ServiceFormWrapper />,
      },
      {
        path: '/services/edit/:id',
        element: <ServiceFormWrapper />,
      },
      {
        path: '/appointments',
        element: <AppointmentListPage />,
      },
      {
        path: '/appointments/new',
        element: <AppointmentFormWrapper />,
      },
      {
        path: '/appointments/edit/:id',
        element: <AppointmentFormWrapper />,
      },
      {
        path: '/sales',
        element: <SaleListPage />,
      },
      {
        path: '/sales/new',
        element: <SaleFormWrapper />,
      },
      {
        path: '/sales/edit/:id',
        element: <SaleFormWrapper />,
      },
    ],
  },
]);

const AppRouter: React.FC = () => {
  return <RouterProvider router={router} />;
};

export default AppRouter;
