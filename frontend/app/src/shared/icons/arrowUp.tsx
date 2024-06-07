export const ArrowUp = ({ color = '#8D95A1' }: { color?: string | undefined }) => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" width="10" height="6" viewBox="0 0 10 6" fill="none">
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M5 1.23978e-05L9.65 4.65001C9.8433 4.84331 9.8433 5.15671 9.65 5.35001C9.4567 5.54331 9.1433 5.54331 8.95 5.35001L5 1.40001L1.05 5.35001C0.856701 5.54331 0.5433 5.54331 0.349999 5.35001C0.1567 5.15671 0.1567 4.84331 0.35 4.65001L5 1.23978e-05Z"
        fill={color}
      />
    </svg>
  );
};
