import { useMemo } from 'react';
import { lang } from '@app/constants';

const imageContainer = {
  display: 'flex',
  flexDirection: 'column',
  width: '140px',
  height: '120px',
  marginBottom: '24px',
} as const;

const imageEtc = {
  borderRadius: '4px',
  backgroundColor: '#FFF',
  color: '#00829B',
  display: 'flex',
  padding: '0px 40px 0px 41px',
  justifyContent: 'center',
  alignItems: 'center',
  height: '96px',
  marginBottom: '24px',
  fontSize: '13px',
} as const;

const imageStyle = {
  margin: '3.73px 4px',
  maxWidth: '100%',
} as const;

const sreenshotTitle = {
  fontSize: '11px',
  textAlign: 'center',
  marginTop: '8px',
  color: '#A2AAB5',
} as const;

const screenshotContainer = {
  display: 'flex',
  flexWrap: 'wrap',
  gap: '16px',
  padding: '0 0 32px 24px',
} as const;

const { numberOfScreenshotsToShow } = lang.attachments;

export const ScreenshotsComponent = ({ screeshotArr }: { screeshotArr: string[] }) => {
  const imageUrls = useMemo(
    () => screeshotArr.map((base64Image) => `data:image/png;base64,${base64Image}`),
    [screeshotArr],
  );

  return (
    <div style={screenshotContainer}>
      {imageUrls.slice(0, numberOfScreenshotsToShow).map((imageUrl: string, index: number) => {
        if (index === 11) {
          return (
            <div key={`screenshot_img_${index}`} style={imageEtc}>
              {`+ ${imageUrls.length - 11} more`}
            </div>
          );
        }
        return (
          <div key={`screenshot_img_${index}`} style={imageContainer}>
            <img src={imageUrl} alt={`Image_screenshot${index + 1}`} style={imageStyle}></img>
            <div style={sreenshotTitle}>{`Image_screenshot${index + 1}`}</div>
          </div>
        );
      })}
    </div>
  );
};
