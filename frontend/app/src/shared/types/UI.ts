export enum MODAL_TYPE {
  CREATE_SESSION = 'createSession',
  GENERATE_CODE_SNIPPET = 'generateCodeSnippet',
}

export interface ModalProps {
  id: string;
  zIndex: number;
}

export interface Modal {
  id: string;
  zIndex: number;
  component: MODAL_TYPE;
}

export interface UIData {
  isSideBarOpen: boolean;
  modals: Modal[];
}

export interface AttachmentsTabProps {
  setTab: (newState: string) => void;
  isTerminated: boolean;
}
