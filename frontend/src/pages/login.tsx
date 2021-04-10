import React from 'react';
import 'antd/dist/antd.css';
import './login.css';
import { Button, Card, Form, Input } from 'antd';
import { FormInstance } from 'antd/lib/form';
import { UserOutlined } from '@ant-design/icons';
import { webauthnLogin, webauthnReg, browserSupport } from '../utils/webauthn';


class IndexForm extends React.Component {
  formRef = React.createRef<FormInstance>();

  onFocusCheck = () => {
    browserSupport();
  }

  getUsername = () => {
    if (this.formRef.current != null) {
      let username = this.formRef.current.getFieldValue("username");
      if (username !== "") {
        return username;
      }
      return null;
    }
    return null;
  }

  onLoginAction = () => {
    let username = this.getUsername();
    if (username === null) {
      console.log("err!");
      return;
    }
    console.log(username);
    webauthnLogin(username);
  }

  onRegAction = () => {
    let username = this.getUsername();
    if (username === null) {
      console.log("err!");
      return;
    }
    console.log(username);
    webauthnReg(username);
  }

  render() {
    return (
      <div style={{height: '100%', width: '100%', display: 'flex', justifyContent: 'center', alignItems: 'flex-start'}}>
        <Card title={"Login"} style={{
          width: "300px",
          marginTop: '100px'
        }}>
          <Form
            name="index-form"
            style={{ maxWidth: "400px" }}
            ref={this.formRef}
          >
            <Form.Item
              name="username"
              rules={[
                {
                  required: true,
                  message: 'Please input your Username!',
                },
              ]}
            >
              <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="Username" onFocus={this.onFocusCheck} />
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" onClick={this.onRegAction}
                style={{ width: "100%" }}>
                Register
            </Button>
            </Form.Item>

            <Form.Item>
              <Button htmlType="submit" style={{ width: "100%" }} onClick={this.onLoginAction}>
                Login
            </Button>
            </Form.Item>
          </Form>
        </Card>
      </div>

    );
  }
};

export default IndexForm;
