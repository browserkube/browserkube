import { Dropdown, FieldText, Checkbox } from '@reportportal/ui-kit';
import { useSelector } from 'react-redux';
import { useState } from 'react';
import { createPortal } from 'react-dom';
import { type FormValues, type FormData } from '@shared/types/createSession';
import { getResolutionOptions } from '@shared/utils/getResolutionOptions';
import { getBrowserVersionOptions } from '@shared/utils/getBrowserVersionOptions';
import { getBrowserTypeOptions } from '@shared/utils/getBrowserTypeOptions';
import { getOSOptions } from '@shared/utils/getOSOptions';
import { getBrowserChip } from '@shared/utils/getBrowserChip';
import { LinuxIcon } from '@shared/chips/linux';
import { getBrowsers } from '@redux/browsers/browsersSelectors';
import { getFormData } from '@shared/utils/getFormData';
import styles from './CreateSessionModalForm.module.scss';

const tooltipContainer: React.CSSProperties = {
  position: 'relative',
} as const;

const tooltipStyle: React.CSSProperties = {
  display: 'block',
  position: 'absolute',
  borderRadius: '8px',
  background: '#3f3f3ff2',
  color: 'white',
  padding: '16px',
  transition: 'opacity 1s',
  width: 'max-content',
  fontSize: '13px',
  zIndex: '99',
  right: '60px',
  bottom: '40px',
} as const;

const tooltipCursor: React.CSSProperties = {
  content: '',
  position: 'absolute',
  left: '50%',
  border: 'solid transparent',
  borderTopColor: '#3f3f3ff2',
  borderWidth: '8px',
  transform: 'translateY(100%)',
  zIndex: '1000',
} as const;

interface CreateSessionModalFormProps {
  formValues: FormValues;
  onChange: (values: Partial<FormValues>) => void;
}

export const CreateSessionModalForm = ({ formValues, onChange }: CreateSessionModalFormProps) => {
  const browserData = useSelector(getBrowsers);
  const formData: FormData = getFormData(browserData);
  const { browserOptions, foundOS } = getBrowserTypeOptions(formData, formValues.platformName);
  const { browserVersionOptions, foundBrowserVersion } = getBrowserVersionOptions(foundOS, formValues.browserName);
  const resolutionOptions = getResolutionOptions(foundBrowserVersion, formValues.browserVersion);

  const [tooltip, isTooltip] = useState<boolean>(false);

  const saveValue = (fieldName: keyof FormValues, value: string | boolean) => {
    onChange({ [fieldName]: value });
  };

  return (
    <>
      <div // TODO: uncomment when resolve with z-index issue, now the tooltip is cut by the parent element
        style={tooltipContainer}
        onMouseEnter={() => {
          isTooltip(true);
        }}
        onMouseLeave={() => {
          isTooltip(true);
        }}>
        <div className={`${styles.field_container}`}>
          {/* TODO: replace it by label prop in Dropdown when added */}
          <div id={'platform_tooltip'}>Platform</div>
          <Dropdown
            icon={<LinuxIcon />}
            disabled={true} // TODO: remove, when other platforName would be supported
            className={`${styles.field} ${styles.dropdown} ${styles.icon}`}
            options={getOSOptions(formData)}
            value={formValues.platformName}
            onChange={(value) => {
              saveValue('platformName', String(value));
            }}
          />
          {tooltip &&
            createPortal(
              <div style={tooltipStyle}>
                We currently support only Linux.
                <div style={tooltipCursor} />
              </div>,
              document.getElementById('platform_tooltip') ?? document.body
            )}
        </div>
      </div>
      <div className={styles.field_container}>
        {/* TODO: replace it by label prop in Dropdown when added */}
        <div>Browser</div>
        <Dropdown
          icon={getBrowserChip(formValues.browserName)}
          className={`${styles.field} ${styles.dropdown} ${styles.icon}`}
          options={browserOptions}
          value={formValues.browserName}
          onChange={(value) => {
            saveValue('browserName', String(value));
          }}
        />
      </div>
      <div className={styles.field_container}>
        {/* TODO: replace it by label prop in Dropdown when added */}
        <div>Browser version</div>
        <Dropdown
          className={styles.field}
          options={browserVersionOptions}
          value={formValues.browserVersion}
          onChange={(value) => {
            saveValue('browserVersion', String(value));
          }}
        />
      </div>
      <div className={styles.field_container}>
        {/* TODO: replace it by label prop in Dropdown when added */}
        <div>Screen resolution</div>
        <Dropdown
          className={styles.field}
          options={resolutionOptions}
          value={formValues.screenResolution}
          onChange={(value) => {
            saveValue('screenResolution', String(value));
          }}
        />
      </div>
      <div className={styles.field_container}>
        <FieldText
          className={styles.field}
          label="Session Name"
          value={formValues.sessionName}
          onChange={(event) => {
            saveValue('sessionName', event.target.value);
          }}
        />
      </div>
      <div className={styles.field_container}>
        <Checkbox
          value={formValues.recordVideo}
          onChange={(event) => {
            saveValue('recordVideo', event.target.checked);
          }}>
          Record Video
        </Checkbox>
      </div>
    </>
  );
};
