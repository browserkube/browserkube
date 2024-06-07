export const CloseChipIcon = (color?: string) => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" width="8" height="8" viewBox="0 0 8 8" fill="none">
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M0.196911 0.195262C-0.0656375 0.455611 -0.0656366 0.877723 0.196912 1.13807L3.08301 4L0.196912 6.86193C-0.0656366 7.12228 -0.0656375 7.54439 0.196911 7.80474C0.459459 8.06509 0.885134 8.06509 1.14768 7.80474L4.03378 4.94281L6.85232 7.73773C7.11486 7.99808 7.54054 7.99808 7.80309 7.73773C8.06564 7.47738 8.06563 7.05527 7.80309 6.79492L4.98456 4L7.80309 1.20508C8.06563 0.944727 8.06564 0.522616 7.80309 0.262267C7.54054 0.00191751 7.11486 0.00191734 6.85232 0.262267L4.03378 3.05719L1.14768 0.195262C0.885134 -0.0650874 0.459459 -0.0650873 0.196911 0.195262Z"
        fill={color ?? '#3F3F3F'}
      />
    </svg>
  );
};
