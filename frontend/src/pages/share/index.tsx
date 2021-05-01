import { Button, Input, message } from 'antd'
import React, { useEffect, useState } from 'react'
import { Modal } from 'antd'
import {
  DownloadOutlined, FileAddOutlined
} from '@ant-design/icons'
import axios from '@/utils/axios'
import { queryURLParams } from '../../utils/helper'
import styles from './index.css'


interface IFileInfo {
  id: number
  username: string
  fileType: string
  filename: string
  hash: string
}

const Share: React.FC = () => {
  const [modalVisible, setModalVisible] = useState<boolean>(false)
  const [inputValue, setInputValue] = useState<string>('')
  const [fileInfo, setFileInfo] = useState<IFileInfo>()

  const getFileInfo = () => {
    const hashObject = queryURLParams(location.href)
    axios({
      url: 'api/file/share/show',
      method: 'post',
      data: JSON.stringify({
        hash: hashObject.f
      })
    }).then((res) => {
      setFileInfo(res.data)
    }).catch(() => {

    })
  }

  const downloadShareFile = (callback: () => void) => {
    const hashObject = queryURLParams(location.href)
    axios({
      url: 'api/file/share/download',
      method: 'post',
      responseType: 'blob',
      data: {
        fileId: fileInfo?.id,
        code: inputValue,
        hash: hashObject.f
      }
    }).then((res) => {
      const blob = res.data
      // FileReader主要用于将文件内容读入内存
      const reader = new FileReader();
      reader.readAsDataURL(blob);
      // onload当读取操作成功完成时调用
      reader.onload = function (e) {
        const a = document.createElement('a');
        // 获取文件名fileName
        const fileName = fileInfo?.filename
        a.download = fileName!;
        a.href = e.target?.result as string;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
        callback()
      }
    }).catch((err) => {
      message.error('提取码错误')
    })
  }

  const modalOnCancle = () => {
    setModalVisible(false)
    setInputValue('')
  }

  const modalOnOk = () => {
    downloadShareFile(() => {
      message.success('正在下载文件')
      setModalVisible(false)
    })
  }

  const inputOnChange = (e: any) => {
    setInputValue(e.target.value)
  }

  useEffect(() => {
    getFileInfo()
  }, [])

  return (
    <div className={styles.shareBackground}>
      <div className={styles.shareContainer}>
        <div className={styles.header}>
          <span>{fileInfo?.filename}</span>
          <Button
            icon={<DownloadOutlined />}
            onClick={() => setModalVisible(true)}
          >
            下载文件
          </Button>
        </div>
        <div className={styles.body}>
          <div style={{ textAlign: 'center' }}>
            <FileAddOutlined style={{ fontSize: '70px' }} />
            <span style={{ marginTop: '8px', display: 'block' }}>{fileInfo?.filename}</span>
          </div>

        </div>
      </div>
      <Modal
        maskClosable={false}
        width={300}
        visible={modalVisible}
        onCancel={modalOnCancle}
        onOk={modalOnOk}
      >
        请输入提取码
        <Input
          value={inputValue}
          onChange={inputOnChange}
          style={{ marginTop: '6px' }}
        />
      </Modal>
    </div>
  )
}

export default Share
