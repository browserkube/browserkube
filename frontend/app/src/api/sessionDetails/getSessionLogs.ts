import { type CommandsResponse, type CommandsParams } from '@shared/types/sessions';
import { api } from '@shared/api';

export const getSessionLogs = async (url: string): Promise<string> => {
  return await api.get<string>({ url, allowErrorHandling: false });
};

export const getSessionCommands = async (url: string, params: CommandsParams): Promise<CommandsResponse> => {
  return await api.get<CommandsResponse>({ url, params, allowErrorHandling: false });
};
