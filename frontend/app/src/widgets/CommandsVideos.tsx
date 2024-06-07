import { IconButton } from '@mui/material';
import { useCallback, useEffect, useRef, useState } from 'react';
import { useSelector } from 'react-redux';
import ReactPlayer from 'react-player';
import { v4 as uuidv4 } from 'uuid';
import styles from '@app/styles/commandsVideos.module.scss';
import { type AttachmentsTabProps } from '@shared/types/UI';
import { getActiveSessionId, getSessionDetails, getSessionDetailsCommands } from '@redux/sessionDetails/selectors';
import { API_URL } from '@shared/api';
import { type CommandsParams, type Commands } from '@shared/types/sessions';

import { BASE_URL } from '@shared/lib';
import { useAppDispatch } from '@hooks/useAppDispatch';
import { fetchSessionDetailsCommands } from '@redux/sessionDetails/thunk';
import { ArrowDown } from '@shared/icons/arrowDown';
import { ArrowUp } from '@shared/icons/arrowUp';
import { DownloadIcon } from '@shared/icons/downloadIcon';
import { SuccessIcon } from '@shared/icons/successIcon';
import { FailureIcon } from '@shared/icons/failureIcon';
import { CommandPlay } from '@shared/icons/commandPlay';
import { CommandPause } from '@shared/icons/commandPause';
import { ShowMoreArrows } from '@shared/icons/showMoreArrows';
import { getCommandTimestamp } from '@shared/utils/getCommandTimestamp';
import { lang } from '@app/constants';
import { debounce } from '@shared/utils/debounce';
import { getSessionFileUrl } from '@shared/utils/getSessionFileUrl';
import { textStyle } from './Attachments';
import { NoVideoComponent } from './NoVideoComponent';
import { CommandComponent } from './CommandComponent';

const commandLine = {
  display: 'flex',
  alignItems: 'center',
  justifyContent: 'space-between',
  margin: '8px 0',
  backgroundColor: 'white',
  borderRadius: '4px',
  width: '100%',
};

const commandsSection = {
  width: '50%',
  marginRight: '24px',
  marginTop: '4px',
  height: '350px',
  maxHeight: '50%',
  overflowY: 'scroll',
} as const;

const onlyCommands = {
  width: '100%',
  height: '90vh',
  overflowY: 'scroll',
} as const;

const {
  commandsAndVideos: { bucketSize, scrollThrottle, fetchCommandsThrottle, fetchingErrorMsg, initialPageToken },
} = lang;

export const CommandVideos = (props: AttachmentsTabProps) => {
  const dispatch = useAppDispatch();
  const activeSessionId = useSelector(getActiveSessionId);
  const { state, videoRefAddr } = useSelector(getSessionDetails);
  const { data, newPageToken } = useSelector(getSessionDetailsCommands);

  const [playCommand, setPlayCommand] = useState<Record<string, boolean>>({});
  const [showCommand, setShowCommand] = useState<Record<string, boolean>>({});
  const [isCommandsSection, toggleIsCommandsSection] = useState(true);
  const [isFetching, setIsFetching] = useState<boolean>(false);
  const [isVideoPlaying, setIsVideoPlaying] = useState<boolean>(false);
  const playerContainerRef = useRef<HTMLDivElement | null>(null);
  const playerRef = useRef<ReactPlayer | null>(null);
  const targetRef = useRef<HTMLDivElement | null>(null);
  const commandsContainerRef = useRef<HTMLDivElement | null>(null);

  const commandsPayload: CommandsParams = {
    pageToken: newPageToken,
    pageSize: bucketSize,
  };

  const AVAILIABLE_EXTRA_COMMANDS = newPageToken && newPageToken !== initialPageToken;
  const COMMANDS_AVAILIABLE = data && !!data[0].sessionId;
  const isSessionActiveTerminated = activeSessionId && state === 'Terminated';
  const videoUrl = `${BASE_URL}${getSessionFileUrl(activeSessionId, videoRefAddr)}`;

  const handleCommands = () => {
    toggleIsCommandsSection(!isCommandsSection);
  };

  const handleCommand = (id: string) => {
    setShowCommand((prevState) => ({
      ...prevState,
      [id]: !prevState[id],
    }));
  };

  const playStopVideo = (time: number, commandId: string) => {
    if (playerRef.current) {
      const flagToPlay = !playCommand[commandId];
      setPlayCommand({ [commandId]: flagToPlay });
      if (flagToPlay) {
        playerRef.current.seekTo(time);
      }
      setIsVideoPlaying(!!flagToPlay);
    }
  };

  const stopVideo = () => {
    setIsVideoPlaying(false);
    setPlayCommand({});
  };

  const fetchCommands = useCallback(async () => {
    if (isFetching || !newPageToken) {
      return;
    }

    setIsFetching(true);
    try {
      const URL = `${API_URL.SESSIONS}${activeSessionId}/commands`;
      await dispatch(fetchSessionDetailsCommands({ url: URL, params: commandsPayload }));
      console.log('Fetching new commands...');
      await new Promise((resolve) => setTimeout(resolve, fetchCommandsThrottle));
    } catch (error) {
      console.error(fetchingErrorMsg, error);
    } finally {
      setIsFetching(false);
      if (commandsContainerRef.current) {
        const newScrollPosition = commandsContainerRef.current.scrollTop - 100;
        commandsContainerRef.current.scrollTo({
          top: newScrollPosition,
          behavior: 'smooth',
        });
      }
    }
  }, [isFetching, newPageToken, commandsPayload, dispatch]);

  useEffect(() => {
    if (isSessionActiveTerminated) {
      const fetchCommands = async () => {
        const URL = `${API_URL.SESSIONS}${activeSessionId}/commands`;
        await dispatch(fetchSessionDetailsCommands({ url: URL, params: commandsPayload }));
      };
      void fetchCommands();
    }
  }, [activeSessionId, dispatch, isSessionActiveTerminated]);

  useEffect(() => {
    const container = commandsContainerRef.current;
    const handleScroll = () => {
      if (!container) {
        return;
      }

      const { scrollTop, scrollHeight, clientHeight } = container;
      const isBottom = scrollTop + clientHeight >= scrollHeight - 10;

      if (isBottom && !isFetching) {
        void fetchCommands();
      }
    };
    const debouncedHandleScroll = debounce(handleScroll, scrollThrottle);
    container?.addEventListener('scroll', debouncedHandleScroll);
    return () => {
      container?.removeEventListener('scroll', debouncedHandleScroll);
    };
  }, [fetchCommands, isFetching]);

  return (
    <>
      <div className={styles.logs_section}>
        <div className={styles.details_section}>
          <IconButton onClick={handleCommands}>{(isCommandsSection && <ArrowDown />) || <ArrowUp />}</IconButton>
          <div style={textStyle}>Commands & Video</div>
        </div>
        <div className={styles.download} style={COMMANDS_AVAILIABLE ? {} : { opacity: '50%' }}>
          <div>Download</div>
          <IconButton disabled={true} onClick={() => null}>
            <DownloadIcon />
          </IconButton>
        </div>
      </div>
      <div className={styles.commands_container}>
        {isCommandsSection &&
          (COMMANDS_AVAILIABLE ? (
            <>
              <div
                ref={commandsContainerRef}
                className="commands_main_container"
                style={videoRefAddr ? commandsSection : onlyCommands}>
                {data.map(({ commandId, statusCode, command, request, response, timestamp, method }: Commands) => {
                  const { secondsToJump, commandTime } = getCommandTimestamp(timestamp, data);
                  return (
                    <div key={uuidv4()}>
                      <div style={showCommand[commandId] ? { ...commandLine, margin: '8px 0 0' } : commandLine}>
                        <div className={styles.left_command_line}>
                          <IconButton
                            onClick={() => {
                              handleCommand(commandId);
                            }}>
                            {(showCommand && <ArrowDown />) || <ArrowUp />}
                          </IconButton>
                          <div className={styles.status_icon}>
                            {String(statusCode).startsWith('2') ? <SuccessIcon /> : <FailureIcon />}
                          </div>
                          <div className={styles.command_name}>{`${method} ${command}`}</div>
                        </div>
                        {!!videoRefAddr && (
                          <div className={styles.timestamp}>
                            {commandTime}
                            <IconButton
                              onClick={() => {
                                playStopVideo(secondsToJump, commandId);
                              }}>
                              {!playCommand[commandId] ? <CommandPlay /> : <CommandPause />}
                            </IconButton>
                          </div>
                        )}
                      </div>
                      {showCommand[commandId] && (
                        <CommandComponent props={{ statusCode, command, request, response, method }} />
                      )}
                    </div>
                  );
                })}
                {AVAILIABLE_EXTRA_COMMANDS && (
                  <div className={styles.more_commands}>
                    <ShowMoreArrows />
                    <div ref={targetRef} className={styles.more_commands_text}>
                      show more commands
                    </div>
                  </div>
                )}
              </div>
              {!!videoRefAddr && (
                <div ref={playerContainerRef} className={styles.player}>
                  <ReactPlayer
                    onPause={stopVideo}
                    ref={playerRef}
                    url={videoUrl}
                    playing={isVideoPlaying}
                    controls={true}
                    width="100%"
                    height="100%"
                  />
                </div>
              )}
            </>
          ) : (
            !!activeSessionId ?? <NoVideoComponent {...props} />
          ))}
      </div>
    </>
  );
};
