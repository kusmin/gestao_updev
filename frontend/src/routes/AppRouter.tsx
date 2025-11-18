import { Navigate, BrowserRouter, Routes, Route } from 'react-router-dom';
import AppLayout from '../components/layout/AppLayout';
import ClientListPage from '../pages/clients/ClientListPage';
import DashboardPage from '../pages/dashboard/DashboardPage';
import LoginPage from '../pages/auth/LoginPage';
import SignupPage from '../pages/auth/SignupPage';
import ProtectedRoute from '../components/routing/ProtectedRoute';
import AppointmentListPage from '../pages/appointments/AppointmentListPage';
import ProductListPage from '../pages/products/ProductListPage';
import SaleListPage from '../pages/sales/SaleListPage';
import ProductForm from '../pages/products/ProductForm';
import SaleForm from '../pages/sales/SaleForm';

const AppRouter: React.FC = () => (
  <BrowserRouter>
    <Routes>
      <Route path="/login" element={<LoginPage />} />
      <Route path="/signup" element={<SignupPage />} />
      <Route element={<ProtectedRoute />}>
        <Route path="/" element={<AppLayout />}>
          <Route index element={<DashboardPage />} />
          <Route path="clients" element={<ClientListPage />} />
          <Route path="appointments" element={<AppointmentListPage />} />
          <Route path="products" element={<ProductListPage />} />
          <Route path="products/new" element={<ProductForm />} />
          <Route path="products/edit/:id" element={<ProductForm />} />
          <Route path="sales" element={<SaleListPage />} />
          <Route path="sales/new" element={<SaleForm />} />
          <Route path="sales/edit/:id" element={<SaleForm />} />
        </Route>
      </Route>
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  </BrowserRouter>
);

export default AppRouter;
