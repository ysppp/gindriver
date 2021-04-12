import React, { useEffect, useState } from 'react'
import { history } from 'umi';
import {
  Button, Table, Modal,
  Input, message, Form, Card,
  Upload, Dropdown, Menu
} from 'antd'
import { FormInstance } from 'antd/lib/form'
import {
  UploadOutlined, FileAddOutlined,
  ShareAltOutlined, DownloadOutlined,
  InboxOutlined
} from '@ant-design/icons'
import dayjs from 'dayjs'
import { cloneDeep } from 'lodash'
import axios from 'axios'
import { getUserInfo, invalidSessionJumpBack } from '../../utils/user';
import { webauthnLogout } from '../../utils/webauthn'
import LeftMenu from '../../components/LeftMenu'
import styles from './index.css'

interface IData {
  key: number,
  name: string,
  size: number,
  date: number
}

interface IUploadData {
  user: string,
  uploadPerm: { disabled: boolean },
  uploadHint: React.ReactElement,
  jwt: null | string
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
  const [uploadData, setUploadData] = useState<IUploadData>({
    user: 'User',
    uploadPerm: { disabled: false },
    uploadHint: <p>Click or drag a file to this area to upload</p>,
    jwt: null
  })
  const [uploadModalVisible, setUploadModalVisible] = useState<boolean>(false)
  const [fileList, setFileList] = useState([])

  const formRef = React.createRef<FormInstance>();
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

  useEffect(() => {
    onLoadUserFetch().then((username: string | null) => {
      if (username === null) {
        return invalidSessionJumpBack();
      }
      const newUploadData = cloneDeep(uploadData)
      newUploadData.user = `User: ${username}`
      newUploadData.jwt = localStorage.getItem("jwt")
      newUploadData.uploadPerm = { disabled: false }
      newUploadData.uploadHint = <p>Click or drag a file to this area to upload</p>
      setUploadData(newUploadData)
    });
  }, [])

  const normFile = (e: { fileList: any; }) => {
    if (Array.isArray(e)) {
      return e;
    }
    return e && e.fileList;
  };

  const onLoadUserFetch = async () => {
    const user = await getUserInfo();
    if (user) {
      return user;
    }
    return null;
  }

  const onLogoutAction = () => {
    let username = uploadData.user;
    if (username == null) {
      console.log("err!");
      return;
    }
    console.log(username);
    webauthnLogout(username);
  }

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

  const menu = (
    <div className={styles.dropDownContainer}>
      <div onClick={onLogoutAction}>注销</div>
    </div>
  )

  const fileOnchange = (info) => {
    if (info.file?.status === 'done' || info.file?.status === 'removed') {
      setFileList(info.fileList)
    }
  }

  const uploadOnCancle = () => {
    setUploadModalVisible(false)
    setFileList([])
  }

  const uploadOnOk = () => {
    if (fileList.length === 0) {
      message.error('请上传文件')
      return
    }
    const data = {
      files: fileList,
      userName: uploadData.user
    }
    axios({
      url: '/api/user/file/upload',
      method: 'POST',
      data
    }).then(() => {
      message.success('上传文件成功')
      setFileList([])
      setModalVisible(false)
    }).catch(() => {})
  }

  return (
    <div className={styles.container}>
      <div className={styles.contentLeft}>
        <LeftMenu changeType={changeType} />
      </div>
      <div className={styles.contentRight}>
        <div className={styles.headButton}>
          <div>
            <Button
              type="primary"
              icon={<UploadOutlined />}
              onClick={() => setUploadModalVisible(true)}
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
          <div className={styles.headButtonRight}>
            <Dropdown
              overlay={menu}
              placement="bottomCenter"
              trigger={['hover']}
            >
              <div>
                <img src="http://cdn.blogleeee.com/wtp4ln2hccunw9g" />
                <span>{uploadData.user}</span>
              </div>
            </Dropdown>

          </div>
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
        <Modal
          maskClosable={false}
          visible={uploadModalVisible}
          onCancel={uploadOnCancle}
          onOk={uploadOnOk}
        >
          <div style={{ display: 'flex', justifyContent: 'center' }}>
            <Card
              style={{
                width: "300px",
              }}
              onLoad={() => onLoadUserFetch()}
            >
              {/* <Button htmlType='submit' style={{ width: "100%" }} onClick={onLogoutAction}>
                注销
              </Button> */}
              <Form ref={formRef}>

                <Form.Item name="dragger" valuePropName="fileList" getValueFromEvent={normFile} noStyle >
                  <Upload.Dragger
                    name="files"
                    multiple={true}
                    headers={{
                      Authorization: `Bearer ${uploadData.jwt}`
                    }}
                    fileList={fileList}
                    onChange={fileOnchange}
                    {...uploadData.uploadPerm}
                  >
                    <p className="ant-upload-drag-icon">
                      <InboxOutlined />
                    </p>
                    <p className="ant-upload-text" style={{ padding: "1px" }}>{uploadData.uploadHint}</p>
                  </Upload.Dragger>
                </Form.Item>
              </Form>
            </Card>
          </div>
        </Modal>
      </div>
    </div>
    // <div className={styles.container}>

    // </div >
  )
}

export default DownLoad
