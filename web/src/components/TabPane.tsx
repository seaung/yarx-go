import React, { useState } from 'react';
import { Tabs } from 'antd';
import ServerTable from './ServerTable';
import './TabPane.css';

const { TabPane } = Tabs;

interface TabItem {
  key: string;
  title: string;
  content: React.ReactNode;
  closable?: boolean;
}

const TabPaneComponent: React.FC = () => {
  // 初始标签页
  const initialTabs: TabItem[] = [
    {
      key: '1',
      title: '执行工作台',
      content: <ServerTable />,
      closable: false,
    },
  ];

  const [activeKey, setActiveKey] = useState('1');
  const [tabs, setTabs] = useState(initialTabs);

  const onChange = (key: string) => {
    setActiveKey(key);
  };

  const onEdit = (targetKey: React.MouseEvent | React.KeyboardEvent | string, action: 'add' | 'remove') => {
    if (action === 'add') {
      addTab();
    } else if (action === 'remove') {
      removeTab(targetKey as string);
    }
  };

  const addTab = () => {
    const newActiveKey = `newTab${tabs.length + 1}`;
    setTabs([
      ...tabs,
      {
        key: newActiveKey,
        title: '新标签页',
        content: <div>新标签页内容</div>,
        closable: true,
      },
    ]);
    setActiveKey(newActiveKey);
  };

  const removeTab = (targetKey: string) => {
    const targetIndex = tabs.findIndex(tab => tab.key === targetKey);
    const newTabs = tabs.filter(tab => tab.key !== targetKey);
    
    // 如果删除的是当前激活的标签页，需要设置新的激活标签页
    if (newTabs.length && activeKey === targetKey) {
      const { key } = newTabs[targetIndex === newTabs.length ? targetIndex - 1 : targetIndex];
      setActiveKey(key);
    }
    
    setTabs(newTabs);
  };

  return (
    <div className="tab-container">
      <Tabs
        type="editable-card"
        onChange={onChange}
        activeKey={activeKey}
        onEdit={onEdit}
        className="custom-tabs"
      >
        {tabs.map(tab => (
          <TabPane tab={tab.title} key={tab.key} closable={tab.closable}>
            {tab.content}
          </TabPane>
        ))}
      </Tabs>
    </div>
  );
};

export default TabPaneComponent;