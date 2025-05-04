import React from 'react';
import { Table, Tag, Space, Button } from 'antd';
import { SyncOutlined } from '@ant-design/icons';
import './ServerTable.css';

interface ServerData {
  key: string;
  serverName: string;
  serverIP: string;
  runningStatus: string;
  serverVersion: string;
  region: string;
  cpuUsage: string;
  memoryUsage: string;
  runningTime: string;
  operations: string;
}

const ServerTable: React.FC = () => {
  // 模拟数据
  const data: ServerData[] = [];
  
  const columns = [
    {
      title: '服务器名称',
      dataIndex: 'serverName',
      key: 'serverName',
    },
    {
      title: '服务器 IP',
      dataIndex: 'serverIP',
      key: 'serverIP',
    },
    {
      title: '运行状态',
      dataIndex: 'runningStatus',
      key: 'runningStatus',
      render: (status: string) => (
        <Tag color={status === '运行中' ? 'green' : 'red'}>
          {status}
        </Tag>
      ),
    },
    {
      title: '服务器版本',
      dataIndex: 'serverVersion',
      key: 'serverVersion',
    },
    {
      title: '区域',
      dataIndex: 'region',
      key: 'region',
    },
    {
      title: 'CPU 使用率',
      dataIndex: 'cpuUsage',
      key: 'cpuUsage',
    },
    {
      title: '内存使用情况',
      dataIndex: 'memoryUsage',
      key: 'memoryUsage',
    },
    {
      title: '运行时间',
      dataIndex: 'runningTime',
      key: 'runningTime',
    },
    {
      title: '操作',
      key: 'operations',
      render: () => (
        <Space size="middle">
          <Button type="link" size="small">查看</Button>
          <Button type="link" size="small">编辑</Button>
          <Button type="link" size="small" danger>删除</Button>
        </Space>
      ),
    },
  ];

  return (
    <div className="server-table-container">
      <div className="table-header">
        <div className="table-filters">
          <Button type="primary" icon={<SyncOutlined />}>刷新列表</Button>
        </div>
        <div className="table-pagination">
          <span>10 条/页</span>
        </div>
      </div>
      <Table 
        columns={columns} 
        dataSource={data} 
        pagination={{ 
          position: ['bottomRight'],
          showSizeChanger: true,
          showQuickJumper: true,
          showTotal: (total) => `共 ${total} 条`
        }} 
        locale={{ emptyText: '暂无数据' }}
      />
    </div>
  );
};

export default ServerTable;