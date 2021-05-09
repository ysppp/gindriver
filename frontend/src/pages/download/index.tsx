import React, { useEffect, useState } from 'react'
import { history } from 'umi';
import {
  Button, Table, Modal,
  Input, message, Form, Card,
  Upload, Dropdown, Popconfirm,
  Popover
} from 'antd'
import { FormInstance } from 'antd/lib/form'
import {
  UploadOutlined, FileAddOutlined,
  ShareAltOutlined, DownloadOutlined,
  InboxOutlined, MoreOutlined, FolderOutlined
} from '@ant-design/icons'
import dayjs from 'dayjs'
import clipBoard from 'clipboard'
import { cloneDeep } from 'lodash'
import axios from '../../utils/axios'
import { getUserInfo, invalidSessionJumpBack } from '../../utils/user';
import { webauthnLogout } from '../../utils/webauthn'
import LeftMenu from '../../components/LeftMenu'
import EditableCell from '../../components/editableCell'
import IconFont from '../../components/icon'
import styles from './index.css'

interface IData {
  ParentFolderId: number
  key: number
  FileId: number,
  FileName: string,
  Size: number,
  UploadTime: number
  FolderId: number,
  FolderName: string,
  time: number
}

interface IUploadData {
  user: string,
  uploadPerm: { disabled: boolean },
  uploadHint: React.ReactElement,
  jwt: null | string
}

export enum Type {
  all = 0,
  text = 1,
  picture = 2,
  video = 3,
  music
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
  const [fileList, setFileList] = useState<{ originFileObj: Blob }[]>([])
  const [currentFolderArrId, setCurrentFolderArrId] = useState<number[]>([0])
  const [editingKey, setEditingKey] = useState<number>();
  const [renameValue, setRenameValue] = useState<string>('')

  const formRef = React.createRef<FormInstance>();
  const columns = [
    {
      title: 'Name',
      dataIndex: 'FileName',
      width: 200,
      editable: true,
      onCell: (record: IData) => ({
        editing: record.key === editingKey,
        nameValue: renameValue,
        renameInputOnChange: renameInputOnChange,
        cancelRename: cancelRename,
        confirmRename: () => { confirmRename(record) }
      }),
      render: (text: string, record: IData) => {
        return (
          <a onClick={() => { getFileByFileFolder(record) }} style={{ display: 'flex', alignItems: 'center' }}>
            <IconFont type={record.FolderName ? "icon-folder" : 'icon-file'} style={{ fontSize: '25px', marginRight: '3px' }} />
            {text}
          </a>
        )
      }
    },
    {
      title: 'Size',
      dataIndex: 'Size',
      width: 200,
    },
    {
      title: 'Date',
      dataIndex: 'UploadTime',
      width: 600,
      render: (date: number, record: IData) => (
        <div style={{ display: 'flex', justifyContent: "space-between" }}>
          <span>{dayjs(date).format('YYYY-MM-DD HH:mm:ss')}</span>
          <span style={{ marginRight: '60%' }} className={styles.iconSpan}>
            {
              record.FileId ?
                <>
                  <Popconfirm
                    title={`是否要分享文件${record.FileName}？`}
                    onConfirm={() => {
                      shareFile(record)
                    }}
                    onCancel={() => { }}
                    okText="Yes"
                    cancelText="No"
                  >
                    <ShareAltOutlined
                      style={{ color: '#1890ff', fontSize: '16px', cursor: "pointer" }}
                    />
                  </Popconfirm>

                  <Popconfirm
                    title={`确定要下载${record.FileName}吗？`}
                    onConfirm={() => {
                      downloadFile(record)
                    }}
                    onCancel={() => { }}
                    okText="Yes"
                    cancelText="No"
                  >
                    <DownloadOutlined style={{ color: '#1890ff', fontSize: '16px', cursor: "pointer", marginLeft: '10px' }} />
                  </Popconfirm>
                  <Popover
                    content={
                      <div className={styles.popoverContent}>
                        <span>移动到</span>
                        <span onClick={() => renameFileOrFolder(record)} style={{ margin: '4px 0' }}>重命名</span>
                        <span onClick={() => { deleteFile(record) }}>删除</span>
                      </div>
                    }
                  >

                    <MoreOutlined style={{ color: '#1890ff', fontSize: '16px', cursor: "pointer", transform: "rotateZ(90deg)", marginLeft: '10px' }} />
                  </Popover>
                </>
                :
                <span style={{marginRight: '25px'}}>
                  <IconFont
                    type="icon-icon_rename"
                    style={{ color: '#1890ff', fontSize: '16px', cursor: "pointer", marginRight: '10px' }}
                    onClick={() => renameFileOrFolder(record)}
                  />
                  <Popconfirm
                    title={`确定要删除${record.FolderName}吗？`}
                    onConfirm={() => {
                      deleteFolder(record)
                    }}
                    onCancel={() => { }}
                    okText="Yes"
                    cancelText="No"
                  >
                    <IconFont type="icon-delete" style={{ color: '#1890ff', fontSize: '16px', cursor: "pointer" }}/>
                  </Popconfirm>
                </span>
            }
          </span>
        </div>
      )
    },
  ];

  const renameInputOnChange = (e: any) => {
    setRenameValue(e.target.value)
  }

  const renameFileOrFolder = (record: IData) => {
    setEditingKey(record.key)
    setRenameValue(record.FileName || record.FolderName)
  }

  const cancelRename = () => {
    setEditingKey(undefined)
  }

  // 这里需要判断一下是修改文件还是文件夹
  const confirmRename = (record: IData) => {
    if (renameValue === record.FileName) {
      setEditingKey(undefined)
      return
    }
    let url = ""
    let data = {}
    if (record.FolderId) {
      url = 'api/file/folder/update'
      data.FileFolderName = renameValue
      data.FolderId = record.FolderId
    } else {
      url = 'api/file/update'
      data.fId = record.FileId
      data.fileName = renameValue
    }
    axios({
      url,
      method: 'post',
      data
    }).then((res) => {
      setEditingKey(undefined)
      getFilesData(currentFolderArrId[currentFolderArrId.length - 1], () => { message.success('修改成功') })
    }).catch(() => { })

  }

  const deleteFile = (record: IData) => {
    const newData = data.filter((item: IData) => {
      return record.key !== item.key
    })
    axios({
      url: 'api/file/delete',
      method: 'post',
      data: JSON.stringify({
        fId: record.FileId,
        parentFolderId: record.ParentFolderId
      })
    }).then((res) => {
      setData(newData)
      message.success('删除成功')
    }).catch(() => { })
  }

  const deleteFolder = (record: IData) => {
    const newData = data.filter((item: IData) => {
      return record.key !== item.key
    })
    axios({
      url: 'api/file/folder/delete',
      method: 'post',
      data: {
        folderId: record.FolderId
      }
    }).then((res) => {
      setData(newData)
      message.success('删除成功')
    }).catch(() => { })
  }

  const getFileByFileFolder = (record: IData) => {
    if (record.FileId) return
    const newArr = [...currentFolderArrId]
    newArr.push(record.FolderId)
    setCurrentFolderArrId(newArr)
    getFilesData(record.FolderId)
  }

  const downloadFile = (record: IData) => {
    if (record.FolderId) return
    axios({
      url: '/api/file/download',
      method: 'post',
      responseType: 'blob',
      data: JSON.stringify({
        fId: record.FileId
      })
    }).then((res) => {
      const blob = res.data
      // FileReader主要用于将文件内容读入内存
      const reader = new FileReader();
      reader.readAsDataURL(blob);
      // onload当读取操作成功完成时调用
      reader.onload = function (e) {
        const a = document.createElement('a');
        // 获取文件名fileName
        const fileName = record.FileName
        a.download = fileName;
        a.href = e.target?.result as string;
        document.body.appendChild(a);
        a.click();
        document.body.removeChild(a);
      }
    }).catch((err) => {
      console.log(err, 'err')
    })
  }

  const shareFile = (record: IData) => {
    axios({
      url: 'api/file/share/add',
      method: 'post',
      data: JSON.stringify({
        fileId: record.FileId,
        url: 'http://www.bickik.com/share'
      })
    }).then(res => {
      const data = res.data
      let clipboard = new clipBoard('.btn', {
        text() {
          return `链接: ${data.url} 提取码: ${data.code}`
        }
      })
      clipboard.on('success', () => {
        message.success('复制成功')
      })
      clipboard.on('error', () => {
        message.error('复制失败，请手动复制')
      })
      Modal.info({
        title: '成功创建私密链接',
        content: (
          <>
            <Input value={`${data.url}`} style={{ marginBottom: '6px' }} />
            <span>提取码</span>
            <Input value={`${data.code}`} style={{ margin: '4px 0 10px' }} />
            <Button className="btn" type="primary">复制链接及提取码</Button>
          </>
        ),
        closable: true,
        okText: "知道了",
        onOk() {
        }
      })
    }).catch(() => { })
  }

  const getFilesData = (folderId?: number, callback?: () => void) => {
    let url = '/api/file/getAll'
    if (folderId) {
      url = url + `?fId=${folderId}`
    }
    axios({
      url,
      method: 'get'
    }).then(res => {
      if (callback) {
        callback()
      }
      const { fileFolders, files } = res.data
      // 将文件夹和文件整合到一个数组当中
      const newFileFolders = fileFolders.map((item: IData) => {
        item.key = item.FolderId
        item.FileName = item.FolderName
        item.Size = 0
        item.UploadTime = item.time
        return item
      })
      const formData = newFileFolders.concat(files.map((item: IData) => {
        item.key = item.FileId
        return item
      }))
      setData(formData)

    }).catch(() => { })
  }

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
    if (type === 0) {
      getFilesData()
      return
    }
    axios.get(`/api/files?type=${type}`).then(res => {
      const filesData = res.data.files
      const formData = filesData.map((item: IData) => {
        item.key = item.FileId
        return item
      })
      setData(formData)
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
    const data = {
      fileFolderName: fileName,
      parentFolderId: currentFolderArrId[currentFolderArrId.length - 1]
    }
    axios({
      url: 'https://www.bickik.com/api/file/folder/add',
      method: 'post',
      data
    }).then(() => {
      message.success('创建成功')
      setFileName('')
      setModalVisible(false)
      getFilesData(currentFolderArrId[currentFolderArrId.length - 1])
    }).catch(() => {

    })
  }

  const menu = (
    <div className={styles.dropDownContainer}>
      <div onClick={onLogoutAction}>注销</div>
    </div>
  )

  const fileOnchange = (info: { fileList: { originFileObj: Blob }[] }) => {
    setFileList(info.fileList)
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
    let url = 'https://www.bickik.com/api/file/upload'
    if (currentFolderArrId.length > 1) {
      url = url + `?fId=${currentFolderArrId[currentFolderArrId.length - 1]}`
    }
    const formdata = new FormData()
    const filesLength = fileList.length
    for (let i = 0; i < filesLength; i++) {
      formdata.append("files", fileList[i].originFileObj)
    }
    axios({
      url,
      method: 'POST',
      data: formdata
    }).then(() => {
      message.success('上传文件成功')
      setFileList([])
      setUploadModalVisible(false)
      getFilesData(currentFolderArrId[currentFolderArrId.length - 1])
    }).catch(() => { })
  }

  const fileBeforeUpload = () => {
    return false
  }

  const returnLastLevel = () => {
    const length = currentFolderArrId.length
    getFilesData(currentFolderArrId[length - 2])
    setCurrentFolderArrId(currentFolderArrId.slice(0, length - 1))
  }

  useEffect(() => {
    getFilesData()
  }, [])

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
            {
              currentFolderArrId.length > 1 &&
              <Button
                className={styles.buttonMargin}
                onClick={returnLastLevel}
              >
                返回上一级
              </Button>
            }
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
            components={data.length ? {
              body: {
                cell: EditableCell
              }
            } : {}}
            onRow={(record: IData) => ({
              onMouseEnter(e: any) {
                e.target.parentNode.classList.add(styles.visible)
              },
              onMouseLeave(e: any) {
                e.target.parentNode.classList.remove(styles.visible)
              }
            })}
            rowSelection={{ selectedRowKeys, onChange: onSelectChange }}
            columns={columns as any}
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
                    beforeUpload={fileBeforeUpload}
                    fileList={fileList as any}
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
  )
}

export default DownLoad
