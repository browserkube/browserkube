/* eslint-disable @typescript-eslint/no-explicit-any */
import React, { type ReactNode, type ChangeEvent, type CSSProperties } from 'react';
import { type TabProps } from './Tab';
import './Tabs.scss';

interface TabsProps {
  value: string;
  onChange?: (event: ChangeEvent<HTMLDivElement>, label: string) => void;
  children: ReactNode;
  style?: CSSProperties;
}

export const Tabs: React.FC<TabsProps> = ({ value, onChange, children, style }) => {
  const handleChange = (event: ChangeEvent<HTMLDivElement>, tabValue: string) => {
    if (onChange) {
      onChange(event, tabValue);
    }
  };

  return (
    <div className="tabs-container" style={style}>
      <div className="tabs">
        {React.Children.map(children, (child, index) => {
          if (React.isValidElement(child)) {
            return React.cloneElement(child as React.ReactElement<TabProps>, {
              key: index,
              onClick: (e: any) => {
                handleChange(e, child.props.value);
              },
              selected: value === child.props.value,
            });
          }
          return null;
        })}
      </div>
    </div>
  );
};
