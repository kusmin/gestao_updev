import React from 'react';
import {
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
} from '@mui/material';
import {
  Dashboard,
  People,
  Business,
  Inventory,
  Receipt,
  Event,
  Group,
} from '@mui/icons-material';
import { Link } from 'react-router-dom';

const menuItems = [
  { text: 'Dashboard', icon: <Dashboard />, path: '/' },
  { text: 'Tenants', icon: <Business />, path: '/tenants' },
  { text: 'Users', icon: <Group />, path: '/users' },
  { text: 'Clients', icon: <People />, path: '/clients' },
  { text: 'Products', icon: <Inventory />, path: '/products' },
  { text: 'Services', icon: <Receipt />, path: '/services' },
  { text: 'Appointments', icon: <Event />, path: '/appointments' },
  { text: 'Sales', icon: <Receipt />, path: '/sales' },
];

const Sidebar: React.FC = () => {
  return (
    <List>
      {menuItems.map((item) => (
        <ListItem key={item.text} disablePadding>
          <ListItemButton component={Link} to={item.path}>
            <ListItemIcon>{item.icon}</ListItemIcon>
            <ListItemText primary={item.text} />
          </ListItemButton>
        </ListItem>
      ))}
    </List>
  );
};

export default Sidebar;
