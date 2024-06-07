import React, { useEffect, useRef, useState, type ReactNode, type ReactElement } from 'react';

import './DropButton.scss';
import { ArrowDown } from '@shared/icons/arrowDown';
import { ArrowUp } from '@shared/icons/arrowUp';

interface DropDownKitProps {
  icon?: ReactNode;
  placeholder: string;
  children: Array<ReactElement<DropDownProps>>;
  style?: React.CSSProperties;
}

interface DropDownProps {
  onCustomClick: () => void;
  value: string;
  style?: React.CSSProperties;
  icon?: ReactNode;
}

export const DropDown: React.FC<DropDownProps> = ({ onCustomClick, value, style, icon }) => {
  const handleItemClick = () => {
    onCustomClick();
  };

  return (
    <div className="dropdown-item" style={style} onClick={handleItemClick}>
      {icon}
      {value}
    </div>
  );
};

export const DropDownKit: React.FC<DropDownKitProps> = ({ icon, placeholder, children, style }) => {
  const [isOpen, setIsOpen] = useState(false);
  const dropdownRef = useRef<HTMLDivElement | null>(null);

  useEffect(() => {
    const handleOutsideClick = (event: MouseEvent) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
        setIsOpen(false);
      }
    };

    window.addEventListener('click', handleOutsideClick);

    return () => {
      window.removeEventListener('click', handleOutsideClick);
    };
  }, []);

  const iconHandler = () => {
    event?.stopPropagation();
    setIsOpen((prevState) => !prevState);
  };

  const handleButtonClick = () => {
    setIsOpen(!isOpen);
  };

  const handleDropDownItemClick = () => {
    setIsOpen(false);
  };

  return (
    <div className="dropdown" ref={dropdownRef} style={style} onClick={handleButtonClick}>
      <div className="dropdown-button">
        {placeholder}
        <div className="arrow-icon" onClick={iconHandler}>
          {isOpen ? <ArrowUp color="white" /> : <ArrowDown color="white" />}
        </div>
      </div>
      {isOpen && (
        <div className="dropdown-list">
          {React.Children.map(children, (child) => {
            return React.cloneElement(child as React.ReactElement, {
              onClick: handleDropDownItemClick,
            });
          })}
        </div>
      )}
    </div>
  );
};
