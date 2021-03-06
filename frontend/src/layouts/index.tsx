import React from 'react';
import LeftMenu from '../components/LeftMenu'
import styles from './index.css';

const BasicLayout: React.FC = props => {
  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>FileFree</h1>
      {
        //location.href.indexOf('login') !== -1 ?
        props.children
        // :
        // <div className={styles.container}>
        //   <div className={styles.contentLeft}>
        //     <LeftMenu />
        //   </div>
        //   <div className={styles.contentRight}>
        //     {props.children}
        //   </div>
        // </div>
      }
    </div>
  );
};

export default BasicLayout;
