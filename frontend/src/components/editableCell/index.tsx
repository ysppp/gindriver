import React from 'react'
import { Input } from 'antd';
import { CheckCircleOutlined, CloseCircleOutlined } from '@ant-design/icons'

interface IProps {
  nameValue: string
  editing: boolean,
  renameInputOnChange: () => void
  children: React.ReactElement
  cancelRename: () => void
  confirmRename: () => void
}

const EditableCell: React.FC<IProps> = ({
  editing,
  children,
  renameInputOnChange,
  nameValue,
  cancelRename,
  confirmRename
}) => {
  return (
    <td>
      {editing ? (
        <div style={{display: 'flex', alignItems: 'center'}}>
          <Input value={nameValue} onChange={renameInputOnChange} style={{width: '70%'}}/>
          <CheckCircleOutlined
            onClick={confirmRename}
            style={{color: '#1890ff', fontSize: '20px', marginRight: '8px', marginLeft: '4px', cursor: 'pointer'}}
          />
          <CloseCircleOutlined
            style={{color: '#1890ff', fontSize: '20px', cursor: 'pointer'}}
            onClick={() => cancelRename()}
          />
        </div>
      ) : (
        children
      )}
    </td>
  );
};

export default EditableCell