import React, { useState } from 'react';
import { Layout } from 'antd';
import Sidebar from './Sidebar';
import HeaderComponent from './Header';
import './MainLayout.css';

const { Content } = Layout;

interface MainLayoutProps {
  children: React.ReactNode;
}

const MainLayout: React.FC<MainLayoutProps> = ({ children }) => {
  const [collapsed, setCollapsed] = useState(false);

  const toggleSidebar = () => {
    setCollapsed(!collapsed);
  };

  return (
    <Layout className="main-layout">
      <Sidebar collapsed={collapsed} />
      <Layout className="site-layout">
        <HeaderComponent collapsed={collapsed} toggle={toggleSidebar} />
        <Content className="site-content">
          {children}
        </Content>
      </Layout>
    </Layout>
  );
};

export default MainLayout;