import { IconButton } from '@mui/material';
import { Checkbox } from '@reportportal/ui-kit';
import { useMemo, useState } from 'react';
import styles from '@app/styles/archive.module.scss';
import { lang } from '@app/constants';
import { Divider, DividerType } from 'widgets/Divider';
import { CloseChipIcon } from '@shared/chips/closeChipIcon';
import { CloseModal } from '@shared/icons/closeModal';
import { ArchivePic } from '@shared/icons/archivePic';
import { DownloadIcon } from '@shared/icons/downloadIcon';

const { archive } = lang;

const HARDCODED_DATA = [
  {
    name: 'vinodel.zip',
    memory: 235,
    id: '23',
  },
  {
    name: 'vinodel2.zip',
    memory: 245,
    id: '32',
  },
  {
    name: 'vinodel3.zip',
    memory: 999,
    id: '44',
  },
  {
    name: 'vinodel4.zip',
    memory: 235,
    id: '9d9d',
  },
  {
    name: 'vinodel5.zip',
    memory: 235,
    id: '8ss777',
  },
  {
    name: 'vinodel6.zip',
    memory: 235,
    id: '678vv89',
  },
  {
    name: 'vinodel4.zip',
    memory: 235,
    id: '99g',
  },
  {
    name: 'vinodel5.zip',
    memory: 235,
    id: '8777a',
  },
  {
    name: 'vinodel6.zip',
    memory: 235,
    id: '67889ss',
  },
  {
    name: 'vinodel4.zip',
    memory: 235,
    id: '991',
  },
  {
    name: 'vinodel5.zip',
    memory: 235,
    id: '8777ff',
  },
  {
    name: 'vinodel6.zip',
    memory: 235,
    id: '67889',
  },
];

type SelectionState = Record<string, boolean>;

interface Archive {
  name: string;
  memory: number;
  id: string;
}

export type OnCloseModal = () => void;

interface ArchiveModalProps {
  onClose: OnCloseModal;
}

export const ArchiveModal = ({ onClose }: ArchiveModalProps) => {
  const [selectSession, setSelectSession] = useState<SelectionState>({});
  const [selectedArchive, setSelectedArchive] = useState<Archive[]>([]);

  const clearAllSelect = () => {
    setSelectSession({});
    setSelectedArchive([]);
  };

  const selectArchive = (id: string) => {
    setSelectSession((prevState: SelectionState) => ({
      ...prevState,
      [id]: !prevState[id],
    }));

    const acrhieveIndex = selectedArchive.findIndex((archive) => archive.id === id);
    if (acrhieveIndex !== -1) {
      setSelectedArchive((prevState) => {
        const newState = [...prevState];
        newState.splice(acrhieveIndex, 1);
        return newState;
      });
    } else {
      setSelectedArchive((prevState) => [...prevState, HARDCODED_DATA.find((archive) => archive.id === id)!]);
    }
  };

  const { totalMemory, amountSelected } = useMemo(() => {
    const totalMemory = selectedArchive.reduce((total, archive) => total + archive.memory, 0);
    const amountSelected = selectedArchive.length;

    return { totalMemory, amountSelected };
  }, [selectArchive]);

  return (
    <div className={styles.background_container}>
      <div className={styles.container}>
        <div className={styles.header}>
          <div className={styles.title}>Archive</div>
          <IconButton onClick={onClose}>
            <CloseModal />
          </IconButton>
        </div>
        <div className={styles.cautionMessage}>{archive.cautionMessage}</div>
        <div className={styles.list_container}>
          {HARDCODED_DATA.map((session, index) => {
            const isChecked = selectSession[session.id] ?? false;
            return (
              <div key={`${index}_${session.id}`} className={styles.archive_container}>
                <Checkbox
                  value={isChecked}
                  onChange={() => {
                    selectArchive(session.id);
                  }}></Checkbox>
                <div className={styles.archive_line}>
                  <div className={styles.left_bar}>
                    <ArchivePic />
                    <div className={styles.session_title}>{session.name}</div>
                  </div>
                  <div className={styles.right_bar}>
                    <div className={styles.memory}>{`${session.memory}MB`}</div>
                    <IconButton onClick={() => null}>
                      <DownloadIcon />
                    </IconButton>
                  </div>
                </div>
              </div>
            );
          })}
        </div>
        <div className={styles.footer}>
          <div className={styles.selection_container}>
            <div className={styles.selected_amount}>{`${amountSelected} items selected`}</div>
            <Divider type={DividerType.VERTICAL} />
            <div className={styles.deselect_btn} onClick={clearAllSelect}>
              Deselect all
            </div>
            {CloseChipIcon('#00829B')}
          </div>
          <div className={styles.download_container}>
            <div className={styles.memory_amount}>{`${totalMemory} MB`}</div>
            <div className={styles.download_btn}>Download</div>
          </div>
        </div>
      </div>
    </div>
  );
};
