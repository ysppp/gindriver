import React from 'react'
import { Link } from 'umi'
import { withRouter } from 'react-router'
import { Menu } from 'antd';
import {
  FileImageOutlined,
  DesktopOutlined,
  UserOutlined,
  MenuUnfoldOutlined,
  MenuFoldOutlined
} from '@ant-design/icons';
import { Type } from '../../pages/download'
import './index.css'

interface IProps {
  history: any,
  location: any,
  changeType: (type: Type) => void
}

interface IState {
  collapsed: boolean,
  selectedKeys: string[]
}

class LeftMenu extends React.Component<IProps, IState> {
  constructor(props: IProps) {
    super(props)
    this.state = {
      collapsed: false,
      selectedKeys: ['uploadDocument'],
    };
  }

  toggleCollapsed = () => {
    if (this.state.collapsed) {
      document.documentElement.style.setProperty('--basicWidth', '200px')
    } else {
      document.documentElement.style.setProperty('--basicWidth', '0')
    }
    this.setState({
      collapsed: !this.state.collapsed,
    });
  };

  handleMenuClick = ({ key }: any) => {
    this.setState({
      selectedKeys: [key]
    })
  }

  // componentDidMount() {

  //   this.props.history.listen(
  //     () => {
  //       setTimeout(() => {
  //         const path = this.props.location.pathname.slice(1)
  //         this.setState({
  //           selectedKeys: [path]
  //         })
  //         // switch (path) {
  //         //   case 'AddArticle':
  //         //     this.props.setCurrentPath('写作台')
  //         //     break
  //         //   case 'ArticleList':
  //         //     this.props.setCurrentPath('个人中心')
  //         //     break
  //         //   case 'Problem':
  //         //     this.props.setCurrentPath('问题反馈')
  //         //     break
  //         //   default:
  //         //     this.props.setCurrentPath('图片管理')
  //         // }
  //       }, 0)
  //     }
  //   )
  //   const path = this.props.location.pathname.slice(1)
  //   if (path) {
  //     this.setState({
  //       selectedKeys: [path]
  //     })
  //   }
  // }


  render() {
    return (
      <div>
        <Menu
          selectedKeys={this.state.selectedKeys}
          defaultActiveFirst={true}
          theme="light"
          mode="inline"
          inlineCollapsed={this.state.collapsed}
          className="menu"
          onClick={this.handleMenuClick}
        >
          <Menu.Item key="uploadDocument">
            <DesktopOutlined />
            {/* <Link to="/uploadDocument">quanbu</Link> */}
            <a onClick={() => this.props.changeType(Type.all)}>全部</a>
          </Menu.Item>
          <Menu.Item key="downloadDocument">
            <FileImageOutlined />
            {/* <Link to="/downloadDocument">下载文件</Link> */}
            <a onClick={() => this.props.changeType(Type.picture)}>图片</a>
          </Menu.Item>
          <Menu.Item key="shareDocument">
            <UserOutlined />
            {/* <Link to="/shareDocument">分享文件</Link> */}
            <a onClick={() => this.props.changeType(Type.video)}>视频</a>
          </Menu.Item>

          {/* <Menu.Item key="Problem">
            <BugOutlined />
            <Link to="/Problem"></Link>
          </Menu.Item> */}
          <div onClick={this.toggleCollapsed} className="footer">
            {this.state.collapsed ? <MenuUnfoldOutlined style={{ color: 'balck' }} /> : <MenuFoldOutlined style={{ color: 'black' }} />}
          </div>
        </Menu>
      </div>
    );
  }
}

export default withRouter(LeftMenu as any)
