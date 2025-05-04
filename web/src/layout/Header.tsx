import React from 'react';
import { Layout, Button, Space } from 'antd';
import { MenuFoldOutlined, MenuUnfoldOutlined, UserOutlined } from '@ant-design/icons';
import './Header.css';

const { Header } = Layout;

interface HeaderProps {
  collapsed: boolean;
  toggle: () => void;
}

const HeaderComponent: React.FC<HeaderProps> = ({ collapsed, toggle }) => {
  return (
    <Header className="site-header">
      <div className="header-left">
        <Button
          type="text"
          icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
          onClick={toggle}
          className="trigger-button"
        />
        <div className="header-tabs">
          <div className="active-tab">执行工作台 <Button type="text" size="small" className="close-tab">×</Button></div>
          <Button type="text" className="add-tab">+ 添加标签页</Button>
        </div>
      </div>
      <div className="header-right">
        <Space>
          <span className="user-info">
            <UserOutlined /> admin
          </span>
        </Space>
      </div>
    </Header>
  );
};

export default HeaderComponent;