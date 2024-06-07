export enum REDUCER_STATUS {
  IDLE = 'idle',
  PENDING = 'pending',
  FULFILLED = 'fulfilled',
  REJECTED = 'rejected',
}

interface DefaultStatusData {
  status: REDUCER_STATUS;
  error: string | null;
}

export type ReducerType<Data, StatusData = DefaultStatusData> = StatusData & {
  data: Data;
};
