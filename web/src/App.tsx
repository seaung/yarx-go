import React from 'react';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/lib/locale/zh_CN';
import MainLayout from './layout/MainLayout';
import TabPaneComponent from './components/TabPane';
import './App.css';

function App() {
  return (
    <ConfigProvider locale={zhCN}>
      <MainLayout>
        <TabPaneComponent />
      </MainLayout>
    </ConfigProvider>
  );
}

export default App;
