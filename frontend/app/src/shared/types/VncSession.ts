export interface VncSessionType {
  sessionId: string;
  locked: boolean;
  swipeableOpened: boolean;
  swipeableAnchor: 'left' | 'top' | 'right' | 'bottom' | undefined;
  isInitialized: boolean;
}
