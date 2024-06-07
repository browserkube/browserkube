import React, { type ReactNode, useState, useRef } from 'react';
import styles from './ButtonDrop.module.scss';
import { ArrowRight } from './arrowRight';

interface DropValueProps {
  title: string;
  label: string;
  value?: string[];
  onMouseOver?: () => void;
  onClick?: (value: string) => void;
}

const tooltipCursor: React.CSSProperties = {
  content: '',
  position: 'absolute',
  left: '13%',
  top: '0',
  border: 'solid transparent',
  borderBottomColor: '#FFF',
  borderWidth: '8px',
  transform: 'translateY(-100%)',
} as const;

// DropValue.tsx
export const DropValue: React.FC<DropValueProps> = ({ title, label, value, onMouseOver, onClick }) => {
  const [isHovered, setIsHovered] = useState(false);

  const handleMouseLeave = () => {
    const timeoutID = setTimeout(() => {
      setIsHovered(false);
    }, 40);
    timeoutRef.current = timeoutID;
  };

  const handleListMouseEnter = () => {
    clearTimeout(timeoutRef.current);
  };

  const handleListMouseLeave = () => {
    handleMouseLeave();
  };

  const timeoutRef = useRef<NodeJS.Timeout | undefined>();

  const handleClick = () => {
    if (onClick) {
      if (!value) {
        onClick(label);
      }
    }
  };

  return (
    <div
      className={styles.dropValueContainer}
      onMouseOver={() => {
        if (onMouseOver) {
          onMouseOver();
        }
        setIsHovered(true);
        clearTimeout(timeoutRef.current);
      }}
      onMouseLeave={handleMouseLeave}
      onClick={handleClick}>
      <div className={styles.dropValueTitleContainer}>
        <div className={styles.title}>{title}</div>
        {value && value.length > 0 && <ArrowRight />}
      </div>
      {isHovered && value && (
        <div
          className={styles.dropValueListContainer}
          onMouseEnter={handleListMouseEnter}
          onMouseLeave={handleListMouseLeave}>
          {value.map((item) => (
            <div
              key={item}
              className={styles.dropValueItem}
              onClick={() => {
                if (onClick) {
                  onClick(item);
                }
              }}>
              {item}
            </div>
          ))}
        </div>
      )}
    </div>
  );
};

interface ButtonDropProps {
  title: string;
  onClick: (result: Record<string, string>) => void;
  children: ReactNode;
}

// ButtonDrop.tsx
export const ButtonDrop: React.FC<ButtonDropProps> = ({ title, onClick, children }) => {
  const [selectedValue, setSelectedValue] = useState<Record<string, string>>({});

  const handleMouseOver = () => {
    setSelectedValue({});
  };

  const handleChooseFilterValue = (label: string, value: string, isLabelSelect = false) => {
    const result = isLabelSelect ? { [label]: '' } : { [value.toLowerCase()]: '' };
    setSelectedValue(result);
    onClick(result);
  };

  return (
    <div className={styles.buttonDropContainer}>
      <div style={tooltipCursor} />
      <div className={styles.buttonDropTitle}>{title}</div>
      {React.Children.map(children, (child, index) => {
        if (React.isValidElement(child)) {
          const clonedChild = React.cloneElement(child as React.ReactElement, {
            key: index,
            onMouseOver: handleMouseOver,
            onClick: (value: string) => {
              if (child.props?.label && child.props?.value) {
                handleChooseFilterValue(child.props.label, value);
              } else if (child.props?.label) {
                handleChooseFilterValue(child.props.label, child.props.label, true);
              }
            },
          });
          return clonedChild;
        }
        return child;
      })}
    </div>
  );
};
