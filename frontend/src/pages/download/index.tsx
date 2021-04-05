import React, { useEffect, useState } from 'react'
import { history } from 'umi';
import { Button, Table, Modal, Input, message } from 'antd'
import { UploadOutlined, FileAddOutlined, FileOutlined } from '@ant-design/icons'
import dayjs from 'dayjs'
import styles from './index.css'

interface IData {
  key: number,
  name: string,
  size: number,
  date: number
}

const dataSource: IData[] = [];
for (let i = 0; i < 4; i++) {
  dataSource.push({
    key: i,
    name: `文件${i}`,
    size: 0,
    date: new Date().getTime(),
  });
}


const DownLoad: React.FC = () => {
  const [selectedRowKeys, setSelectedRowKeys] = useState<number[]>([])
  const [modalVisible, setModalVisible] = useState<boolean>(false)
  const [fileName, setFileName] = useState<string>('')
  const [data, setData] = useState<IData[]>(dataSource)

  useEffect(() => {

  })
  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      render: (text: string) => <a><img src="http://cdn.blogleeee.com/folder.png" style={{marginRight: '6px'}}/>{text}</a>
    },
    {
      title: 'Size',
      dataIndex: 'size',
    },
    {
      title: 'Date',
      dataIndex: 'date',
      render: (date: number) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')
    },
  ];

  const onSelectChange = (selectedRowKeys: any) => {
    console.log('selectedRowKeys changed: ', selectedRowKeys);
    setSelectedRowKeys(selectedRowKeys)
  };

  const modalCancel = () => {
    setModalVisible(false)
    setFileName('')
  }

  const handleOk = () => {
    const newData = [...data]
    newData.unshift({
      key: 20,
      name: fileName,
      size: 0,
      date: new Date().getTime(),
    });
    setFileName('')
    setModalVisible(false)
    setData(newData)
    message.success('创建成功')
  }
  return (
    <div className={styles.container}>
      <div className={styles.headButton}>
        <Button
          type="primary"
          icon={<UploadOutlined />}
          onClick={() => history.push('/uploadDocument')}
        >
          上传
        </Button>
        <Button
          icon={<FileAddOutlined />}
          className={styles.buttonMargin}
          onClick={() => setModalVisible(true)}
        >
          新建文件夹
        </Button>
      </div>
      <div className={styles.contentBody}>
        <Table rowSelection={{ selectedRowKeys, onChange: onSelectChange }} columns={columns} dataSource={data} />
      </div>
      <Modal
        width={300}
        maskClosable={false}
        onCancel={modalCancel}
        visible={modalVisible}
        onOk={handleOk}
      >
        <div style={{ display: 'flex', alignItems: 'center', paddingTop: '26px' }}>
          <span style={{ flexBasis: '60px' }}>文件名</span>
          <Input value={fileName} onChange={e => setFileName(e.target.value)}></Input>
        </div>
      </Modal>
    </div >
  )
}

export default DownLoad