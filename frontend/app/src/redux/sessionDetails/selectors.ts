import { type SessionDetailsCommands, type SessionDetails, type SessionDetailsLogs } from '@shared/types/sessions';
import { type REDUCER_STATUS } from '@shared/types/reducerType';
import { type AppState } from '../store';

export const getActiveSessionId = (state: AppState): string => state.sessionDetails.data.activeSessionId;

export const getActiveSessionState = (state: AppState): string => state.sessionDetails.data.sessionData.state;

export const getSessionDetails = (state: AppState): SessionDetails => state.sessionDetails.data.sessionData;

export const getSessionDetailsCommands = (state: AppState): SessionDetailsCommands =>
  state.sessionDetails.data.commands;

export const getSessionDetailsLogs = (state: AppState): SessionDetailsLogs => state.sessionDetails.data.logs;

export const getSessionStatus = (state: AppState): REDUCER_STATUS => state.sessionDetails.status;
