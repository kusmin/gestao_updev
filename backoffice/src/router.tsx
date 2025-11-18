import React from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router-dom';
import ProtectedRoute from './components/ProtectedRoute';
import DashboardPage from './pages/DashboardPage';
import TenantListPage from './pages/tenants/TenantListPage';
import TenantForm from './pages/tenants/TenantForm';
import UserListPage from './pages/users/UserListPage';
import UserForm from './pages/users/UserForm';
import ClientListPage from './pages/clients/ClientListPage';
import ClientForm from './pages/clients/ClientForm';
import ProductListPage from './pages/products/ProductListPage';
import ProductForm from './pages/products/ProductForm';
import ServiceListPage from './pages/services/ServiceListPage';
import ServiceForm from './pages/services/ServiceForm';
import AppointmentListPage from './pages/appointments/AppointmentListPage';
import AppointmentForm from './pages/appointments/AppointmentForm';
import SaleListPage from './pages/sales/SaleListPage';
import SaleForm from './pages/sales/SaleForm';
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
        element: <TenantForm />,
      },
      {
        path: '/tenants/edit/:id',
        element: <TenantForm />,
      },
      {
        path: '/users',
        element: <UserListPage />,
      },
      {
        path: '/users/new',
        element: <UserForm />,
      },
      {
        path: '/users/edit/:id',
        element: <UserForm />,
      },
      {
        path: '/clients',
        element: <ClientListPage />,
      },
      {
        path: '/clients/new',
        element: <ClientForm />,
      },
      {
        path: '/clients/edit/:id',
        element: <ClientForm />,
      },
      {
        path: '/products',
        element: <ProductListPage />,
      },
      {
        path: '/products/new',
        element: <ProductForm />,
      },
      {
        path: '/products/edit/:id',
        element: <ProductForm />,
      },
      {
        path: '/services',
        element: <ServiceListPage />,
      },
      {
        path: '/services/new',
        element: <ServiceForm />,
      },
      {
        path: '/services/edit/:id',
        element: <ServiceForm />,
      },
      {
        path: '/appointments',
        element: <AppointmentListPage />,
      },
      {
        path: '/appointments/new',
        element: <AppointmentForm />,
      },
      {
        path: '/appointments/edit/:id',
        element: <AppointmentForm />,
      },
      {
        path: '/sales',
        element: <SaleListPage />,
      },
      {
        path: '/sales/new',
        element: <SaleForm />,
      },
      {
        path: '/sales/edit/:id',
        element: <SaleForm />,
      },
    ],
  },
]);

const AppRouter: React.FC = () => {
  return <RouterProvider router={router} />;
};

export default AppRouter;
