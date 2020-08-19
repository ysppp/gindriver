import React, { useState }  from 'react';
import {Card, Form, Input, Upload} from 'antd';
import {InboxOutlined} from '@ant-design/icons';
import {FormInstance} from 'antd/lib/form';
import {getUserInfo, invalidSessionJumpBack} from '../utils/user';


import styles from './user.css';

class UserForm extends React.Component {

  constructor(props: Readonly<{}>) {
    super(props);
    this.state = {
      user: "User",
      uploadPerm: {disabled: true},
      uploadHint: <p>You need to be <b>admin</b> to upload files</p>,
      jwt: null
    };
    this.onLoadUserFetch().then((username: string | null) => {
      if (username === null) {
        return invalidSessionJumpBack();
      }
      this.setState({user: `User: ${username}`});
      this.setState({jwt: localStorage.getItem("jwt")});
      if (username === "admin") {
        this.setState({
          uploadPerm: {disabled: false},
          uploadHint: <p>Click or drag a file to this area to upload</p>
        });
      }
    });
  }

  formRef = React.createRef<FormInstance>();

  normFile = (e: { fileList: any; }) => {
    console.log('Upload event:', e);
    if (Array.isArray(e)) {
      return e;
    }
    return e && e.fileList;
  };

  onLoadUserFetch = async () => {
    const user = await getUserInfo();
    if (user) {
      return user;
    }
    return null;
  }

  render() {
    return (
      // @ts-ignore
      <Card title={this.state.user} style={{
        width: "300px",
      }} onLoad={() => this.onLoadUserFetch()}>
        <Form ref={this.formRef}>
          <Form.Item name="dragger" valuePropName="fileList" getValueFromEvent={this.normFile} noStyle >
            <Upload.Dragger name="files" action="/api/user/file/upload"
                            multiple={false}
                            headers={{
                              Authorization: `Bearer ${this.state.jwt}`
                            }}
                            {...this.state.uploadPerm}>
              <p className="ant-upload-drag-icon">
                <InboxOutlined/>
              </p>
              <p className="ant-upload-text" style={{padding: "1px"}}>{this.state.uploadHint}</p>
            </Upload.Dragger>
          </Form.Item>
        </Form>
      </Card>
    );
  }
}

export default UserForm;
