import React from 'react'
import { Link } from 'umi'
import { Button, Menu } from 'antd';
import {
  FileImageOutlined,
  DesktopOutlined,
  UserOutlined,
  BugOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined
} from '@ant-design/icons';

import './index.css'

interface IProps {

}

interface IState {
  collapsed: boolean,
  selectedKeys: string[]
}

class LeftMenu extends React.Component<IProps, IState> {
  constructor(props: IProps) {
    super(props)
    this.state = {
      collapsed: true,
      selectedKeys: [],
    };
  }

  toggleCollapsed = () => {

    this.setState({
      collapsed: !this.state.collapsed,
    });
  };

  handleMenuClick = ({ key }: any) => {
    this.setState({
      selectedKeys: [key]
    })
  }

  componentDidMount() {
    // this.props.history.listen(
    //   () => {
    //     setTimeout(() => {
    //       const path = this.props.location.pathname.slice(1)
    //       this.setState({
    //         selectedKeys: [path]
    //       })
    //       switch (path) {
    //         case 'AddArticle':
    //           this.props.setCurrentPath('写作台')
    //           break
    //         case 'ArticleList':
    //           this.props.setCurrentPath('个人中心')
    //           break
    //         case 'Problem':
    //           this.props.setCurrentPath('问题反馈')
    //           break
    //         default:
    //           this.props.setCurrentPath('图片管理')
    //       }
    //     }, 0)
    //   }
    // )
    // this.setState({
    //   selectedKeys: [this.props.location.pathname.slice(1)]
    // })
  }


  render() {
    return (
      <div>
        <Menu
          selectedKeys={this.state.selectedKeys}
          theme="dark"
          mode="inline"
          inlineCollapsed={this.state.collapsed}
          className="menu"
          onClick={this.handleMenuClick}
        >
          <Menu.Item key="uploadDocument">
            <DesktopOutlined />
            <Link to="/uploadDocument">上传文件</Link>
          </Menu.Item>
          <Menu.Item key="downloadDocument">
            <FileImageOutlined />
            <Link to="/downloadDocument">下载文件</Link>
          </Menu.Item>
          <Menu.Item key="shareDocument">
            <UserOutlined />
            <Link to="/shareDocument">分享文件</Link>
          </Menu.Item>
          {/* <Menu.Item key="Problem">
            <BugOutlined />
            <Link to="/Problem"></Link>
          </Menu.Item> */}
          <div onClick={this.toggleCollapsed} className="footer">
            {this.state.collapsed ? <MenuUnfoldOutlined style={{ color: 'white' }} /> : <MenuFoldOutlined style={{ color: 'white' }} />}
          </div>
        </Menu>
      </div>
    );
  }
}

export default LeftMenu