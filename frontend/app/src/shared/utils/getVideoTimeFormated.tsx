export const formatTime = (time: number): string => {
  const hours = Math.floor(time / (60 * 60 * 1000))
    .toString()
    .padStart(2, '0');
  const minutes = Math.floor((time % (60 * 60 * 1000)) / (60 * 1000))
    .toString()
    .padStart(2, '0');
  const seconds = Math.floor((time % (60 * 1000)) / 1000)
    .toString()
    .padStart(2, '0');
  return `${hours}:${minutes}:${seconds}`;
};
