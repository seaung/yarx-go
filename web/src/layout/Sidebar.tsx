import React from 'react';
import { Layout, Menu } from 'antd';
import {
  AppstoreOutlined,
  DesktopOutlined,
  FileOutlined,
  TeamOutlined,
  UserOutlined,
  SettingOutlined,
  BarChartOutlined,
  CloudOutlined
} from '@ant-design/icons';
import './Sidebar.css';

const { Sider } = Layout;

interface SidebarProps {
  collapsed: boolean;
}

const Sidebar: React.FC<SidebarProps> = ({ collapsed }) => {
  return (
    <Sider trigger={null} collapsible collapsed={collapsed} className="sidebar">
      <div className="logo">
        <span className="logo-text">{collapsed ? 'X' : '洞鉴 YAR-X'}</span>
      </div>
      <Menu
        theme="dark"
        mode="inline"
        defaultSelectedKeys={['1']}
        items={[
          {
            key: '1',
            icon: <AppstoreOutlined />,
            label: '执行工作台',
          },
          {
            key: '2',
            icon: <DesktopOutlined />,
            label: '应用管理',
          },
          {
            key: '3',
            icon: <CloudOutlined />,
            label: '资源管理',
          },
          {
            key: '4',
            icon: <TeamOutlined />,
            label: '用户管理',
            children: [
              {
                key: '4-1',
                label: '用户列表',
              },
              {
                key: '4-2',
                label: '权限管理',
              },
            ],
          },
          {
            key: '5',
            icon: <BarChartOutlined />,
            label: '监控中心',
          },
          {
            key: '6',
            icon: <FileOutlined />,
            label: '日志中心',
          },
          {
            key: '7',
            icon: <SettingOutlined />,
            label: '系统设置',
          },
        ]}
      />
    </Sider>
  );
};

export default Sidebar;