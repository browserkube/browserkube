import { getSessionId } from '../getSessionId';

describe('getSessionId', () => {
  test('should return the last segment of the hash array', () => {
    const mockLastHashElement = '123';
    const mockHashArray = ['#', 'some', 'path', mockLastHashElement];
    Object.defineProperty(window, 'location', {
      value: {
        hash: mockHashArray.join('/'),
      },
      writable: true,
    });

    const sessionId = getSessionId();
    expect(sessionId).toBe(mockLastHashElement);
  });

  test('should return an empty string when the last segment is empty', () => {
    const mockLastHashElement = '';
    const mockHashArray = ['#', 'some', 'path', mockLastHashElement];
    Object.defineProperty(window, 'location', {
      value: {
        hash: mockHashArray.join('/'),
      },
      writable: true,
    });

    const sessionId = getSessionId();
    expect(sessionId).toBe(mockLastHashElement);
  });

  test('should return the last segment when it contains special characters', () => {
    const mockLastHashElement = '!@#$%^&*()';
    const mockHashArray = ['#', 'some', 'path', mockLastHashElement];
    Object.defineProperty(window, 'location', {
      value: {
        hash: mockHashArray.join('/'),
      },
      writable: true,
    });

    const sessionId = getSessionId();
    expect(sessionId).toBe(mockLastHashElement);
  });
});
