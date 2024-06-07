export const DividerType = {
  HORIZON: 'horizontal',
  VERTICAL: 'vertical',
} as const;

const vertical = {
  width: '1px',
  height: '100%',
  backgroundColor: '#E3E7EC',
} as const;

const horizontal = {
  height: '1px',
  width: '100%',
  backgroundColor: '#E3E7EC',
} as const;

export const Divider = ({ type }: { type: string }) => {
  const style = type === DividerType.HORIZON ? horizontal : vertical;
  return <div style={style}></div>;
};
