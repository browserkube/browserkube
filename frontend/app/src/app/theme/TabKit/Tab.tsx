import React, { type CSSProperties } from 'react';
import './Tabs.scss';

export interface TabProps {
  value: string;
  label: string;
  selected?: boolean;
  disabled?: boolean;
  onClick?: (e: React.SyntheticEvent) => void;
  style?: CSSProperties;
}

export const Tab: React.FC<TabProps> = ({ value, label, selected = false, disabled = false, onClick, style }) => {
  const handleClick = (e: React.SyntheticEvent) => {
    if (!disabled && onClick) {
      onClick(e);
    }
  };

  return (
    <div
      style={{
        cursor: disabled ? 'not-allowed' : 'pointer',
        borderBottom: selected ? '3px solid #00829B' : 'none',
        borderRadius: '1px',
        opacity: disabled ? 0.5 : 1,
        ...style,
      }}
      className={`tab ${selected ? 'active-tab' : ''}`}
      data-value={value}
      data-disabled={disabled}
      onClick={handleClick}>
      {label}
    </div>
  );
};
