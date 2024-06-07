/**
 *
 * @param mode type of session
 * @returns chip of auto/manual sessionMode
 */

export const getSessionMode = (mode: boolean) => {
  const iconText = mode ? 'Manual' : 'Auto';

  const modeStyles = {
    display: 'flex',
    alignItems: 'center',
    padding: '0px 6px 0 6px',
    borderRadius: '3px',
    background: mode ? '#3E7BE6' : '#3AA76D', // green for Auto and Blue for Manual
    color: 'white',
    height: '20px',
    width: '49px',
    fontSize: '11px',
    justifyContent: 'center',
  };

  return <div style={modeStyles}>{iconText}</div>;
};
