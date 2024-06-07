export const ArrowDown = ({ color = '#8D95A1' }: { color?: string | undefined }) => {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" width="10" height="6" viewBox="0 0 10 6" fill="none">
      <path
        fillRule="evenodd"
        clipRule="evenodd"
        d="M5 5.7L0.350001 1.05C0.156701 0.856701 0.1567 0.5433 0.35 0.35C0.5433 0.1567 0.8567 0.1567 1.05 0.35L5 4.3L8.95 0.35C9.1433 0.156701 9.4567 0.1567 9.65 0.35C9.8433 0.5433 9.8433 0.8567 9.65 1.05L5 5.7Z"
        fill={color}
      />
    </svg>
  );
};
