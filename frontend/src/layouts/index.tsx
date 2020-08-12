import React from 'react';
import styles from './index.css';

const BasicLayout: React.FC = props => {
  return (
    <div className={styles.normal}>
      <h1 className={styles.title}>GinDriver</h1>
      <div style={{
        display: "flex",
        justifyContent: "center",
        alignItems: "center",
        alignContent: "center",
        paddingTop: "30px",
      }}>
      {props.children}
      </div>
    </div>
  );
};

export default BasicLayout;
