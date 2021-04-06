import React, { useEffect, useState } from 'react'
import { history } from 'umi';
import { Button, Table, Modal, Input, message } from 'antd'
import {
  UploadOutlined, FileAddOutlined,
  ShareAltOutlined, DownloadOutlined
} from '@ant-design/icons'
import dayjs from 'dayjs'
import axios from 'axios'
import LeftMenu from '../../components/LeftMenu'
import styles from './index.css'

interface IData {
  key: number,
  name: string,
  size: number,
  date: number
}

export enum Type {
  all = 0,
  picture = 1,
  video = 2
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
  const [data, setData] = useState<IData[]>([])

  const columns = [
    {
      title: 'Name',
      dataIndex: 'name',
      width: 200,
      render: (text: string) => <a><img src="http://cdn.blogleeee.com/folder.png" style={{ marginRight: '6px' }} />{text}</a>
    },
    {
      title: 'Size',
      dataIndex: 'size',
      width: 200,
    },
    {
      title: 'Date',
      dataIndex: 'date',
      width: 600,
      render: (date: number) => (
        <div style={{ display: 'flex', justifyContent: "space-between" }}>
          <span>{dayjs(date).format('YYYY-MM-DD HH:mm:ss')}</span>
          <span style={{ marginRight: '60%' }} className={styles.iconSpan}>
            <ShareAltOutlined style={{ color: '#1890ff', fontSize: '16px', marginRight: '10px', cursor: "pointer" }} />
            <DownloadOutlined style={{ color: '#1890ff', fontSize: '16px', cursor: "pointer" }} />
          </span>
        </div>
      )
    },
    // {
    //   title: '',
    //   dataIndex: 'action',
    //   render: () => {
    //     return (
    //       <span>
    //         <ShareAltOutlined style={{ color: '#1890ff' }} />
    //         <DownloadOutlined style={{ color: '#1890ff' }} />
    //       </span>
    //     )
    //   }
    // }
  ];

  // useEffect(() => {
  //   axios.get('/files?type=all').then(res => {
  //     setData(res.data)
  //   }).catch(() => { })
  // }, [])

  const changeType = (type: Type) => {
    axios.get(`/files?type=${type}`).then(res => {
      setData(res.data)
    }).catch(() => { })
  }

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
      <div className={styles.contentLeft}>
        <LeftMenu changeType={changeType}/>
      </div>
      <div className={styles.contentRight}>
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
          <Table
            onRow={record => ({
              onMouseEnter(e: any) {
                e.target.parentNode.classList.add(styles.visible)
              },
              onMouseLeave(e: any) {
                e.target.parentNode.classList.remove(styles.visible)
              }
            })}
            rowSelection={{ selectedRowKeys, onChange: onSelectChange }}
            columns={columns}
            dataSource={data}
            tableLayout="fixed"
          />
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
      </div>
    </div>
    // <div className={styles.container}>

    // </div >
  )
}

export default DownLoad